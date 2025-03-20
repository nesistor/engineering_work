package main

import (
	"auth/data"
	"auth/users"
	"auth/admin"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	userServiceAddress  = "user-service:50001"
	adminServiceAddress = "admin-service:50002"
	authTimeout         = 5 * time.Second
)

type (
	// AuthHandlerProvider definiuje kontrakt dla handlerów autentykacji
	AuthHandlerProvider interface {
		AuthenticateUser(http.ResponseWriter, *http.Request)
		AuthenticateAdmin(http.ResponseWriter, *http.Request)
		RefreshToken(http.ResponseWriter, *http.Request)
		RevokeToken(http.ResponseWriter, *http.Request)
		ValidateUserToken(http.ResponseWriter, *http.Request)
		ValidateAdminToken(http.ResponseWriter, *http.Request)
		Logout(http.ResponseWriter, *http.Request)
	}

	// TokenManager zarządza operacjami na tokenach
	TokenManager interface {
		GenerateToken(ctx context.Context, userID int, role data.Role, duration time.Duration, scope data.Scope, keyID string) (string, error)
		RefreshAccessToken(ctx context.Context, refreshToken string, keyID string) (string, error)
		InsertDeactivatedToken(ctx context.Context, token string, ttl time.Duration) error
	}

	// GRPCClientProvider dostarcza klientów gRPC
	GRPCClientProvider interface {
		NewUserClient(conn *grpc.ClientConn) users.UserServiceClient
		NewAdminClient(conn *grpc.ClientConn) admin.AdminServiceClient
	}

	authHandler struct {
		config          *Config
		tokenManager    TokenManager
		grpcClient      GRPCClientProvider
		logger          Logger
		userServiceConn *grpc.ClientConn
		adminServiceConn *grpc.ClientConn
	}
)

// NewAuthHandler tworzy nową instancję handlera autentykacji
func NewAuthHandler(
	cfg *Config,
	tm TokenManager,
	gc GRPCClientProvider,
	logger Logger,
) (AuthHandlerProvider, error) {
	userConn, err := grpc.Dial(userServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	adminConn, err := grpc.Dial(adminServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to admin service: %w", err)
	}

	return &authHandler{
		config:          cfg,
		tokenManager:    tm,
		grpcClient:      gc,
		logger:          logger,
		userServiceConn: userConn,
		adminServiceConn: adminConn,
	}, nil
}

// AuthenticateUser obsługuje autentykację użytkownika
func (h *authHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := h.config.readJSON(w, r, &req); err != nil {
		h.logger.Error("Invalid request payload", "error", err)
		h.config.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), authTimeout)
	defer cancel()

	userClient := h.grpcClient.NewUserClient(h.userServiceConn)
	resp, err := userClient.ValidateUser(ctx, &users.ValidateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil || !resp.IsValid {
		h.logger.Warn("Authentication failed", "email", req.Email, "error", err)
		h.config.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	if err := h.issueTokens(w, r, int(resp.UserId), data.RoleUser); err != nil {
		h.logger.Error("Token generation failed", "error", err)
		h.config.errorJSON(w, err)
		return
	}
}

// issueTokens generuje tokeny dostępu i refresh
func (h *authHandler) issueTokens(w http.ResponseWriter, r *http.Request, userID int, role data.Role) error {
	ctx, cancel := context.WithTimeout(r.Context(), authTimeout)
	defer cancel()

	accessToken, err := h.tokenManager.GenerateToken(
		ctx,
		userID,
		role,
		15*time.Minute,
		data.ScopeAuthentication,
		"user-key",
	)
	if err != nil {
		return fmt.Errorf("access token generation failed: %w", err)
	}

	refreshToken, err := h.tokenManager.GenerateToken(
		ctx,
		userID,
		role,
		7*24*time.Hour,
		data.ScopeRefresh,
		"user-key",
	)
	if err != nil {
		return fmt.Errorf("refresh token generation failed: %w", err)
	}

	if err := h.logger.LogRequest(r.Context(), "Authentication", fmt.Sprintf("%s authenticated", role)); err != nil {
		h.logger.Warn("Failed to log authentication", "error", err)
	}

	return h.config.writeJSON(w, http.StatusOK, jsonResponse{
		Error:   false,
		Message: "Authentication successful",
		Data: map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// RefreshToken obsługuje odświeżanie tokenu dostępu
func (h *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := h.config.readJSON(w, r, &req); err != nil {
		h.config.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), authTimeout)
	defer cancel()

	accessToken, err := h.tokenManager.RefreshAccessToken(ctx, req.RefreshToken, "user-key")
	if err != nil {
		h.logger.Warn("Token refresh failed", "error", err)
		h.config.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	if err := h.logger.LogRequest(r.Context(), "TokenRefresh", "Access token refreshed"); err != nil {
		h.logger.Warn("Failed to log token refresh", "error", err)
	}

	h.config.writeJSON(w, http.StatusOK, jsonResponse{
		Error:   false,
		Message: "Token refreshed",
		Data:    map[string]string{"access_token": accessToken},
	})
}

// validateTokenWithRole sprawdza poprawność tokenu i roli
func (h *authHandler) validateTokenWithRole(w http.ResponseWriter, r *http.Request, requiredRole data.Role) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.config.errorJSON(w, errors.New("authorization header required"), http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		h.config.errorJSON(w, errors.New("invalid authorization header format"), http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), authTimeout)
	defer cancel()

	userID, role, err := h.config.Models.Token.GetUserIDForToken(ctx, token, data.ScopeAuthentication)
	if err != nil || data.Role(role) != requiredRole {
		h.logger.Warn("Unauthorized access attempt", 
			"required_role", requiredRole,
			"actual_role", role,
			"error", err,
		)
		h.config.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	if err := h.logger.LogRequest(r.Context(), "TokenValidation", 
		fmt.Sprintf("User %d with role %s authorized", userID, role),
	); err != nil {
		h.logger.Warn("Failed to log token validation", "error", err)
	}

	h.config.writeJSON(w, http.StatusOK, jsonResponse{
		Error:   false,
		Message: "Access granted",
		Data:    map[string]interface{}{"user_id": userID, "role": role},
	})
}

// AuthenticateAdmin handles admin authentication by validating credentials.
func (app *Config) AuthenticateAdmin(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(adminServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := admin.NewAdminServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.ValidateAdmin(ctx, &admin.ValidateAdminRequest{
		Email:    requestPayload.Email,
		Password: requestPayload.Password,
	})
	if err != nil || !response.IsValid {
		app.errorJSON(w, fmt.Errorf("invalid admin credentials"), http.StatusUnauthorized)
		return
	}

	accessToken, err := app.Models.Token.GenerateToken(ctx, int(response.AdminId), data.RoleAdmin, 15*time.Minute, data.ScopeAuthentication, "user-key")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshToken, err := app.Models.Token.GenerateToken(ctx, int(response.AdminId), data.RoleAdmin, 7*24*time.Hour, data.ScopeRefresh, "user-key")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.logRequest("AuthenticateAdmin", fmt.Sprintf("Admin %s authenticated successfully", requestPayload.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Admin %s authenticated successfully", requestPayload.Email),
		Status:  http.StatusOK,
		Data: map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// ValidateUserToken validates a token for a regular user role.
func (app *Config) ValidateUserToken(w http.ResponseWriter, r *http.Request) {
	app.validateTokenWithRole(w, r, data.RoleUser)
}

// ValidateAdminToken validates a token for the admin role.
func (app *Config) ValidateAdminToken(w http.ResponseWriter, r *http.Request) {
	app.validateTokenWithRole(w, r, data.RoleAdmin)
}

// RevokeToken handles token revocation (blacklisting in Redis).
func (app *Config) RevokeToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Dezaktywowanie tokenu poprzez dodanie go do blacklisty (Redis)
	err = app.Models.Token.InsertDeactivatedToken(context.Background(), requestPayload.Token, 7*24*time.Hour) // 7 dni TTL
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Logowanie operacji
	err = app.logRequest("RevokeToken", "Token revoked successfully")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Token revoked successfully",
		Status:  http.StatusOK,
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// Logout handles the process of logging out a user or admin by revoking their access token.
func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        app.errorJSON(w, fmt.Errorf("missing Authorization header"), http.StatusUnauthorized)
        return
    }

    if !strings.HasPrefix(authHeader, "Bearer ") {
        app.errorJSON(w, fmt.Errorf("invalid Authorization header format"), http.StatusUnauthorized)
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")

    err := app.Models.Token.InsertDeactivatedToken(context.Background(), tokenString, 15*time.Minute) // TTL np. 15 minut na dezaktywację
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.logRequest("Logout", "User or Admin logged out successfully")
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := jsonResponse{
        Error:   false,
        Message: "Logged out successfully",
        Status:  http.StatusOK,
    }

    err = app.writeJSON(w, http.StatusOK, payload)
    if err != nil {
        app.errorJSON(w, err)
    }
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}