package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

// ContextKeyUserID is key to store userID
const ContextKeyUserID = contextKey("userID")

// AuthUserMiddleware checks Authentication Header for users and verifies their roles
func (app *Config) AuthUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

		// Parse the token without a key to inspect claims
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Here we return nil since we want to parse the token without validating the signature
			return nil, nil
		})

		if err != nil {
			app.errorJSON(w, fmt.Errorf("invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// Ensure claims are of the correct type
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			app.errorJSON(w, fmt.Errorf("invalid token claims"), http.StatusUnauthorized)
			return
		}

		// Check the role of the user
		role, ok := claims["role"].(string)
		if !ok || role != "user" {
			app.errorJSON(w, fmt.Errorf("unauthorized role: %s", role), http.StatusUnauthorized)
			return
		}

		// Retrieve the user ID from the claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			app.errorJSON(w, fmt.Errorf("invalid user ID in token"), http.StatusUnauthorized)
			return
		}

		// Store the user ID in the context
		ctx := context.WithValue(r.Context(), ContextKeyUserID, int64(userID))

		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthAdminMiddleware checks if the user is an admin and verifies the JWT token
func (app *Config) AuthAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// Parse the token without verifying the signature initially
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Only accept RSA signing method
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Always return the public key for verification
			publicKey, err := app.KeyManager.GetPublicKey() // Używamy jednego publicznego klucza
			if err != nil {
				return nil, fmt.Errorf("failed to get public key: %v", err)
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			app.errorJSON(w, fmt.Errorf("invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			app.errorJSON(w, fmt.Errorf("invalid token claims"), http.StatusUnauthorized)
			return
		}

		// Sprawdzanie roli
		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			app.errorJSON(w, fmt.Errorf("unauthorized: user is not an admin"), http.StatusForbidden)
			return
		}

		// Pobieranie user_id z claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			app.errorJSON(w, fmt.Errorf("invalid user ID in token"), http.StatusUnauthorized)
			return
		}

		// Ustawienie userID w kontekście
		ctx := context.WithValue(r.Context(), ContextKeyUserID, int64(userID))

		// Przekazanie nowego kontekstu do następnego handlera
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

