package authHelpers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Ignacio-J-Maylin/arithmetic-calculator/middlewares"
)

func GetUserIDFromToken(r *http.Request) (int64, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return 0, errors.New("missing or invalid token")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := middlewares.ValidateJWT(tokenString)
	if err != nil || claims.UserID == 0 {
		return 0, errors.New("unauthorized")
	}

	return claims.UserID, nil
}
