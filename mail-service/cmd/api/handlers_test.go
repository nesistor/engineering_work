package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockMailer struct {
	SendSMTPMessageFunc func(msg Message) error
}

func (m *MockMailer) SendSMTPMessage(msg Message) error {
	return m.SendSMTPMessageFunc(msg)
}

type MockConfig struct {
	Mailer        *MockMailer
	readJSONFunc  func(w http.ResponseWriter, r *http.Request, v interface{}) error
	errorJSONFunc func(w http.ResponseWriter, err error)
	writeJSONFunc func(w http.ResponseWriter, statusCode int, data interface{}) error
}

func (app *MockConfig) readJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return app.readJSONFunc(w, r, v)
}

func (app *MockConfig) errorJSON(w http.ResponseWriter, err error) {
	app.errorJSONFunc(w, err)
}

func (app *MockConfig) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	return app.writeJSONFunc(w, statusCode, data)
}

func (app *MockConfig) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func TestSendMail_Success(t *testing.T) {
	app := &MockConfig{
		Mailer: &MockMailer{
			SendSMTPMessageFunc: func(msg Message) error {
				return nil
			},
		},
		readJSONFunc: func(w http.ResponseWriter, r *http.Request, v interface{}) error {
			return json.NewDecoder(r.Body).Decode(v)
		},
		errorJSONFunc: func(w http.ResponseWriter, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
		writeJSONFunc: func(w http.ResponseWriter, statusCode int, data interface{}) error {
			w.WriteHeader(statusCode)
			return json.NewEncoder(w).Encode(data)
		},
	}

	payload := map[string]string{
		"from":    "test@example.com",
		"to":      "recipient@example.com",
		"subject": "Test Subject",
		"message": "Test Message",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/send-mail", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.SendMail)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}

	var jsonResponse jsonResponse
	err := json.Unmarshal(rr.Body.Bytes(), &jsonResponse)
	if err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	if jsonResponse.Error {
		t.Errorf("expected no error, got error: %v", jsonResponse.Message)
	}
	if jsonResponse.Message != "sent to recipient@example.com" {
		t.Errorf("unexpected message: %v", jsonResponse.Message)
	}
}

func TestSendMail_InvalidJSON(t *testing.T) {
	app := &MockConfig{
		Mailer: &MockMailer{
			SendSMTPMessageFunc: func(msg Message) error {
				return nil
			},
		},
		readJSONFunc: func(w http.ResponseWriter, r *http.Request, v interface{}) error {
			return fmt.Errorf("invalid JSON")
		},
		errorJSONFunc: func(w http.ResponseWriter, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		writeJSONFunc: func(w http.ResponseWriter, statusCode int, data interface{}) error {
			w.WriteHeader(statusCode)
			return json.NewEncoder(w).Encode(data)
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/send-mail", bytes.NewBuffer([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.SendMail)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestSendMail_SMTPError(t *testing.T) {
	app := &MockConfig{
		Mailer: &MockMailer{
			SendSMTPMessageFunc: func(msg Message) error {
				return fmt.Errorf("SMTP error")
			},
		},
		readJSONFunc: func(w http.ResponseWriter, r *http.Request, v interface{}) error {
			return json.NewDecoder(r.Body).Decode(v)
		},
		errorJSONFunc: func(w http.ResponseWriter, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
		writeJSONFunc: func(w http.ResponseWriter, statusCode int, data interface{}) error {
			w.WriteHeader(statusCode)
			return json.NewEncoder(w).Encode(data)
		},
	}

	payload := map[string]string{
		"from":    "test@example.com",
		"to":      "recipient@example.com",
		"subject": "Test Subject",
		"message": "Test Message",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/send-mail", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.SendMail)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	var jsonResponse jsonResponse
	err := json.Unmarshal(rr.Body.Bytes(), &jsonResponse)
	if err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	if !jsonResponse.Error {
		t.Errorf("expected error, got none")
	}
	if jsonResponse.Message != "SMTP error" {
		t.Errorf("unexpected message: %v", jsonResponse.Message)
	}
}
