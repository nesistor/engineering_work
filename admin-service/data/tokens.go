package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

// TokenModel represents the model for working with tokens in the database.
type TokenModel struct {
	DB *sql.DB
}

// Token represents a token used for actions like password reset.
type Token struct {
	PlainText string    `json:"token"`
	UserID    int64     `json:"-"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

// GenerateToken generates a token with a specified time-to-live (TTL) for a user.
func (m *TokenModel) GenerateToken(userID int, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID: int64(userID),
		Expiry: time.Now().Add(ttl),
	}

	// Generate a random sequence of bytes for the token.
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Encode the random bytes to a base32 string (this is the token that will be sent to the user).
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Hash the token for secure storage in the database.
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]
	return token, nil
}

// SaveResetToken saves the password reset token and its expiration to the users table.
func (m *TokenModel) SaveResetToken(email string, token *Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Update the reset_token and reset_token_expiration fields in the users table.
	query := `UPDATE users SET reset_token = $1, reset_token_expiration = $2 WHERE email = $3`
	_, err := m.DB.ExecContext(ctx, query, token.PlainText, token.Expiry, email)
	if err != nil {
		return err
	}

	return nil
}

// GenerateAndSavePasswordResetToken generates a password reset token and saves it in the users table.
func (m *TokenModel) GenerateAndSavePasswordResetToken(email string, userID int) (*Token, error) {
	// Generate a token with a 1-hour TTL.
	token, err := m.GenerateToken(userID, 1*time.Hour)
	if err != nil {
		return nil, err
	}

	// Save the generated token in the users table.
	err = m.SaveResetToken(email, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// ValidateResetToken validates a provided password reset token by checking it against the stored token and expiration.
func (m *TokenModel) ValidateResetToken(email string, providedToken string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var storedToken string
	var expiry time.Time

	// Retrieve the stored token and its expiration from the users table.
	query := `SELECT reset_token, reset_token_expiration FROM users WHERE email = $1`
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&storedToken, &expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	// Check if the token has expired.
	if time.Now().After(expiry) {
		return false, nil
	}

	// Compare the provided token with the stored one.
	if providedToken != storedToken {
		return false, nil
	}

	return true, nil
}
