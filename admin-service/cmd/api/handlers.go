package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"admin-service/data"
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

	err = app.logRequest("registration", fmt.Sprintf("Admin %s registered", newUser.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Admin %s registered successfully", newUser.Email),
		Data: map[string]interface{}{
			"user_id": userID,
		},
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
