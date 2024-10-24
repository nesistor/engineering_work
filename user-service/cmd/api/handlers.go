package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strconv"

	"user-service/data"
)

// Register handles the registration of new user
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

	payload := jsonResponse{
		Error:   false,
		Message: "Users retrieved successfully",
		Data: map[string]interface{}{
			"users": users,
		},
	}

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

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	isValidToken, err := app.Models.Token.ValidateResetToken(token, "ka")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if !isValidToken {
		app.errorJSON(w, fmt.Errorf("invalid or expired token"), http.StatusUnauthorized)
		return
	}

	err = app.Models.User.UpdateUserPassword(requestPayload.Email, requestPayload.Password)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

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


// DeleteUser handles the deletion of a user based on their ID passed in the URL.
func (app *Config) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the URL path
	vars := r.URL.Query()
	idStr := vars.Get("user_id")

	// Convert the user ID from string to int64
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || userID < 1 {
		app.errorJSON(w, fmt.Errorf("invalid user ID"), http.StatusBadRequest)
		return
	}

	// Call the DeleteUserByID method from the user model
	err = app.Models.User.DeleteUserByID(userID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Log the user deletion request
	err = app.logRequest("delete_user", fmt.Sprintf("User with ID %d deleted", userID))
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("User with ID %d deleted successfully", userID),
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
