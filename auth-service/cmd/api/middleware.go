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

// AuthMiddleware checks Authentication Header for both users and admins
func (app *Config) AuthMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				// Always return the public key for verification
				// Używamy jednego publicznego klucza, ale należy podać właściwy identyfikator (np. "kid")
				publicKey, err := app.KeyManager.GetPublicKey("your-key-id") // Przekaż właściwy identyfikator klucza
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

			// Check the role
			role, ok := claims["role"].(string)
			if !ok || role != requiredRole {
				app.errorJSON(w, fmt.Errorf("unauthorized: user role is %s, required role is %s", role, requiredRole), http.StatusForbidden)
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
}
