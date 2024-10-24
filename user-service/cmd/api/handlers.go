package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"user-service/data"
)

// Register handles the registration of new users. It processes the request payload,
// checks for existing users with the same email, hashes the provided password, and
// inserts the new user into the database.
// .
func (app *Config) Register(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	_, err = app.Models.User.GetUserByEmail(requestPayload.Email)
	if err == nil {
		app.errorJSON(w, fmt.Errorf("user with email %s already exists"), http.StatusConflict)
		return
	}

	if err != sql.ErrNoRows {
		app.errorJSON(w, err)
		return
	}

	newUser := data.User{
		Email:        requestPayload.Email,
		UserName:     requestPayload.Username,
		PasswordHash: requestPayload.Password,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userID, err := app.Models.User.InsertUser(newUser)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.logRequest("registration", fmt.Sprintf("User %s registered", newUser.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("User %s registered successfully", newUser.Email),
		Data: map[string]interface{}{
			"user_id": userID,
		},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// GetAll retrieves all users from the database and sends them as a JSON response.
func (app *Config) GetAll(w http.ResponseWriter, r *http.Request) {
	// Call the GetAllUsers method to fetch the users from the database.
	users, err := app.Models.User.GetAllUsers()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Prepare the response payload
	payload := jsonResponse{
		Error:   false,
		Message: "Users retrieved successfully",
		Data: map[string]interface{}{
			"users": users,
		},
	}

	// Write the JSON response
	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// CheckEmail handles checking if a user with the provided email exists in the database.
func (app *Config) CheckEmail(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	exists, err := app.Models.User.EmailExists(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var message string
	if exists {
		message = fmt.Sprintf("User with email %s exists", requestPayload.Email)
	} else {
		message = fmt.Sprintf("User with email %s does not exist", requestPayload.Email)
	}

	payload := jsonResponse{
		Error:   false,
		Message: message,
		Data: map[string]interface{}{
			"exists": exists,
		},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// ResetPassword handles changing a user's password. It extracts a token from the URL
// and a new password from the request payload, verifies the token, and updates the password.
func (app *Config) ResetPassword(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the URL
	token := r.URL.Query().Get("token")
	if token == "" {
		app.errorJSON(w, fmt.Errorf("token is required"), http.StatusBadRequest)
		return
	}

	// Create a struct to hold the new password
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Read the JSON payload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Verify the token. Check both return values.
	isValidToken, err := app.Models.User.VerifyToken(token)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if !isValidToken {
		app.errorJSON(w, fmt.Errorf("invalid or expired token"), http.StatusUnauthorized)
		return
	}

	// Update the password in the database
	err = app.Models.User.UpdateUserPassword(requestPayload.Email, requestPayload.Password)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Respond with a success message
	payload := jsonResponse{
		Error:   false,
		Message: "Password changed successfully",
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// SendResetPasswordEmail sends a reset password link with a token to the user's email
func (app *Config) SendResetPasswordEmail(email, token string) error {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	resetLink := fmt.Sprintf("https://fit_for_example.com/reset-password?token=%s", token)

	mailPayload := mailMessage{
		From:    "no-reply@your-domain.com",
		To:      email,
		Subject: "Password Reset Request",
		Message: fmt.Sprintf("Please use the following link to reset your password: %s", resetLink),
	}

	jsonData, err := json.Marshal(mailPayload)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	mailServiceURL := "http://mail-service/send"
	req, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send email, status code: %d", resp.StatusCode)
	}

	return nil
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
