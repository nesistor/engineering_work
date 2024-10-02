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

// Check Authentication Header is correct for Admin and User and return their ID name userID
func (app *Config) AuthMiddleware(next http.Handler) http.Handler {
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

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return nil, nil
		})

		if err != nil {
			app.errorJSON(w, fmt.Errorf("invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			app.errorJSON(w, fmt.Errorf("missing kid in token header"), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			app.errorJSON(w, fmt.Errorf("invalid token claims"), http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			app.errorJSON(w, fmt.Errorf("invalid role in token"), http.StatusUnauthorized)
			return
		}

		var publicKey interface{}

		if role == "admin" {
			publicKey, err = app.KeyManager.GetPublicAdminKey() 
			if err != nil {
				app.errorJSON(w, fmt.Errorf("failed to get admin public key: %v", err), http.StatusUnauthorized)
				return
			}
		} else {

			publicKey, err = app.KeyManager.GetPublicKey(kid)
			if err != nil {
				app.errorJSON(w, fmt.Errorf("invalid KID: %s", kid), http.StatusUnauthorized)
				return
			}
		}

		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			app.errorJSON(w, fmt.Errorf("invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(float64) 
		if !ok {
			app.errorJSON(w, fmt.Errorf("invalid user ID in token"), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, int64(userID))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
