package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"admin-service/data"
)

// Register handles the registration of new admin
func (app *Config) Register(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		AdminName string `json:"admin_name"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	_, err = app.Models.Admin.GetAdminByEmail(requestPayload.Email)
	if err == nil {
		app.errorJSON(w, fmt.Errorf("admin with email %s already exists", requestPayload.Email), http.StatusConflict)
		return
	}

	if err != sql.ErrNoRows {
		app.errorJSON(w, err)
		return
	}

	newAdmin := data.Admin{
		Email:        requestPayload.Email,
		AdminName:    requestPayload.AdminName,
		PasswordHash: requestPayload.Password,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	adminID, err := app.Models.Admin.InsertAdmin(newAdmin)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.logRequest("registration", fmt.Sprintf("Admin %s registered", newAdmin.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Admin %s registered successfully", newAdmin.Email),
		Data: map[string]interface{}{
			"admin_id": adminID,
		},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// GetAll retrieves all admins from the database and sends them as a JSON response.
func (app *Config) GetAll(w http.ResponseWriter, r *http.Request) {
	// Call the GetAllAdmins method to fetch the admins from the database.
	admins, err := app.Models.Admin.GetAllAdmins()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Admins retrieved successfully",
		Data: map[string]interface{}{
			"admins": admins,
		},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// ResetPassword handles the request to reset an admin's password by generating a token and sending an email.
func (app *Config) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	admin, err := app.Models.Admin.GetAdminByEmail(requestPayload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			app.errorJSON(w, fmt.Errorf("admin with email %s does not exist", requestPayload.Email), http.StatusNotFound)
			return
		}
		app.errorJSON(w, err)
		return
	}

	token, err := app.Models.Token.GenerateAndSavePasswordResetToken(requestPayload.Email, int(admin.ID))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if err := app.SendResetPasswordEmail(requestPayload.Email, token.PlainText); err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Password reset email sent successfully",
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// SendResetPasswordEmail sends a reset password link with a token to the admin's email
func (app *Config) SendResetPasswordEmail(email, token string) error {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	resetLink := fmt.Sprintf("https://fit_new_password.com/reset-password?token=%s", token)

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

// UpdatePassword handles the process of changing the admin's password after token verification.
func (app *Config) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token    string `json:"token"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	isValidToken, err := app.Models.Token.ValidateResetToken(requestPayload.Email, requestPayload.Token)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if !isValidToken {
		app.errorJSON(w, fmt.Errorf("invalid or expired token"), http.StatusUnauthorized)
		return
	}

	err = app.Models.Admin.UpdateAdminPassword(requestPayload.Email, requestPayload.Password)
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

// UpdateAdmin handles the update of an admin's information based on their ID passed in the URL.
func (app *Config) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	idStr := vars.Get("admin_id")

	adminID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || adminID < 1 {
		app.errorJSON(w, fmt.Errorf("invalid admin ID"), http.StatusBadRequest)
		return
	}

	var requestPayload struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	err = app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	updatedAdmin := data.Admin{
		ID:        adminID,
		Email:     requestPayload.Email,
		AdminName: requestPayload.Username,
		UpdatedAt: time.Now(),
	}

	err = app.Models.Admin.UpdateAdmin(updatedAdmin)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = app.logRequest("update_admin", fmt.Sprintf("Admin with ID %d updated", adminID))
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Admin with ID %d updated successfully", adminID),
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// DeleteAdmin handles the deletion of an admin based on their ID passed in the URL.
func (app *Config) DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	idStr := vars.Get("admin_id")

	adminID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || adminID < 1 {
		app.errorJSON(w, fmt.Errorf("invalid admin ID"), http.StatusBadRequest)
		return
	}

	err = app.Models.Admin.DeleteAdminByID(adminID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.logRequest("delete_admin", fmt.Sprintf("Admin with ID %d deleted", adminID))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Admin with ID %d deleted successfully", adminID),
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// InsertNewAdminHandler handles the insertion of a new admin email.
func (app *Config) AddNewAdmin(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.NewAdmin.InsertNewAdmin(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("New admin email %s inserted successfully", requestPayload.Email),
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// GetAllNewAdminsHandler retrieves all new admin emails and returns them as JSON.
func (app *Config) GetAllNewAdmins(w http.ResponseWriter, r *http.Request) {
	newAdmins, err := app.Models.NewAdmin.GetAllNewAdmins()
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "New admins retrieved successfully",
		Data: map[string]interface{}{
			"new_admins": newAdmins,
		},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// DeleteNewAdminHandler handles the deletion of a specific new admin email.
func (app *Config) DeleteNewAdmin(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.NewAdmin.DeleteNewAdmin(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Admin email %s deleted successfully", requestPayload.Email),
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
