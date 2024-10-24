package main

import (
	"auth/data"
	"auth/users"
	"auth/admin"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	userServiceAddress = "user-service:50001"
	adminServiceAddress = "admin-service:50002"
)

// AuthenticateUser handles user authentication by validating email and password.
// It assigns the "user" role if authentication is successful.
func (app *Config) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Używamy user-service dla autoryzacji użytkownika
	conn, err := grpc.Dial(userServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := users.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.ValidateUser(ctx, &users.ValidateUserRequest{
		Email:    requestPayload.Email,
		Password: requestPayload.Password,
	})
	if err != nil || !response.IsValid {
		app.errorJSON(w, fmt.Errorf("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Generowanie tokenu z rolą "user"
	accessToken, err := app.Models.Token.GenerateToken(ctx, int(response.UserId), data.RoleUser, 15*time.Minute, data.ScopeAuthentication, "user-key")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshToken, err := app.Models.Token.GenerateToken(ctx, int(response.UserId), data.RoleUser, 7*24*time.Hour, data.ScopeRefresh, "user-key")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Logowanie operacji
	err = app.logRequest("AuthenticateUser", fmt.Sprintf("User %s authenticated successfully", requestPayload.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("User %s authenticated successfully", requestPayload.Email),
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

// RefreshToken handles the refresh token request and generates a new access token.
func (app *Config) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	accessToken, err := app.Models.Token.RefreshAccessToken(context.Background(), requestPayload.RefreshToken, "user-key")
	if err != nil {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	err = app.logRequest("RefreshToken", "Access token refreshed successfully")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Access token refreshed successfully",
		Status:  http.StatusOK,
		Data: map[string]interface{}{
			"access_token": accessToken,
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

// validateTokenWithRole validates a token, ensuring the correct role and scope.
func (app *Config) validateTokenWithRole(w http.ResponseWriter, r *http.Request, requiredRole string) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		app.errorJSON(w, fmt.Errorf("missing Authorization header"), http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		app.errorJSON(w, fmt.Errorf("invalid Authorization header format"), http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	userID, role, err := app.Models.Token.GetUserIDForToken(context.Background(), token, data.ScopeAuthentication)
	if err != nil || role != requiredRole {
		app.errorJSON(w, fmt.Errorf("unauthorized or invalid token"), http.StatusUnauthorized)
		return
	}

	// Logowanie operacji
	err = app.logRequest("ValidateToken", fmt.Sprintf("Token for user %d with role %s validated", userID, role))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Token is valid for %s", requiredRole),
		Status:  http.StatusOK,
		Data: map[string]interface{}{
			"user_id": userID,
			"role":    role,
		},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
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