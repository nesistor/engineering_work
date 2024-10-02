package main

import (
	"auth-service/data"
	"auth-service/users"
	
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

// Authenticate handles user authentication by validating the provided email and password.
// It reads the request payload, establishes a gRPC connection to the user service. 
// If the credentials are valid, it generates access and refresh tokens.
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial("user-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if !response.IsValid {
		payload := jsonResponse{
			Error:   true,
			Message: response.Message,
			Status:  http.StatusUnauthorized,
		}
		err = app.writeJSON(w, http.StatusUnauthorized, payload)
		if err != nil {
			app.errorJSON(w, err)
		}
		return
	}

	accessToken, err := app.Models.Token.GenerateToken(int(response.UserId), 15*time.Minute, data.ScopeAuthentication)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshToken, err := app.Models.Token.GenerateToken(int(response.UserId), 7*24*time.Hour, data.ScopeRefresh)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", requestPayload.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", requestPayload.Email),
		Status:  http.StatusOK,
		Data: map[string]interface{}{
			"user_id":       response.UserId,
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

	// Odświeżenie access tokena na podstawie refresh tokena JWT
	accessToken, err := app.Models.Token.RefreshAccessToken(requestPayload.RefreshToken)
	if err != nil {
		app.errorJSON(w, err, http.StatusUnauthorized)
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

// RevokeToken ...
func (app *Config) RevokeToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Dodaj token do blacklisty w Redis
	err = app.Models.Token.DeleteToken(requestPayload.Token) 
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

// ValidateUserToken ...
func (app *Config) ValidateUserToken(w http.ResponseWriter, r *http.Request) {

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

	valid, err := app.Models.Token.IsTokenValid(token, data.ScopeAuthentication)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if !valid {
		payload := jsonResponse{
			Error:   true,
			Message: "Invalid or expired token",
			Status:  http.StatusUnauthorized,
		}
		err = app.writeJSON(w, http.StatusUnauthorized, payload)
		if err != nil {
			app.errorJSON(w, err)
		}
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Token is valid",
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
