package tokenutil

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
)

// GetTokenFromEchoHeader extracts the JWT token from the Authorization header
func GetTokenFromEchoHeader(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	// Check if it's a Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}
