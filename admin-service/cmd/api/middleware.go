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
