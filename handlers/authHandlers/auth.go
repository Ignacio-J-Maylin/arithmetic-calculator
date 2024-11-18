package authHandlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ignacio-J-Maylin/arithmetic-calculator/middlewares"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/models"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/service/userService"
)

func decodeCredentials(r *http.Request) (models.Credentials, error) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	return creds, err
}

func generateTokens(userID int64, username string) (string, string, error) {
	token, err := middlewares.GenerateJWT(userID, username)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := middlewares.GenerateRefreshToken(userID, username)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func sendJSONResponse(w http.ResponseWriter, status int, data map[string]string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func SignUp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds, err := decodeCredentials(r)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if !models.IsValidEmail(creds.Username) {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		err = userService.RegisterUser(db, creds.Username, creds.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		token, refreshToken, err := generateTokens(0, creds.Username)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, http.StatusCreated, map[string]string{
			"username":      creds.Username,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds, err := decodeCredentials(r)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		userID, isAuthenticated, err := userService.AuthenticateUser(db, creds.Username, creds.Password)
		if err != nil || !isAuthenticated {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, refreshToken, err := generateTokens(userID, creds.Username)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, http.StatusOK, map[string]string{
			"username":      creds.Username,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

func RefreshToken(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

		if refreshToken == "" {
			http.Error(w, "No refresh token provided", http.StatusUnauthorized)
			return
		}

		claims, err := middlewares.ValidateJWT(refreshToken)
		if err != nil {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		newToken, newRefreshToken, err := generateTokens(claims.UserID, claims.Username)
		if err != nil {
			http.Error(w, "Error generating new tokens", http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, http.StatusOK, map[string]string{
			"token":         newToken,
			"refresh_token": newRefreshToken,
		})
	}
}

func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendJSONResponse(w, http.StatusOK, map[string]string{"message": "Logout successful"})
	}
}
