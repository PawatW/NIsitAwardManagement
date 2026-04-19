package middlewares

import (
	"fmt"
	"net/http"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/auth"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/utils/tokenutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type authMiddleware struct {
	configs *configs.Config
}

type AuthMiddleware interface {
	Middleware(next echo.HandlerFunc) echo.HandlerFunc
}

func NewAuthMiddleware(configs *configs.Config) AuthMiddleware {
	return &authMiddleware{
		configs: configs,
	}
}

func (a *authMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString, err := tokenutil.GetTokenFromEchoHeader(c)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaims{}, func(token *jwt.Token) (any, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.configs.JWT.Secret), nil
		})
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
		}

		// Validate claims
		claims, ok := token.Claims.(*auth.JWTClaims)
		if !ok || !token.Valid {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "invalid claim"})
		}

		// Set claims and user ID to context
		c.Set("user", claims) // Set full claims as "user"
		c.Set(string(enums.UserIDContextKey), claims.UserID)
		c.Set(string(enums.UserEmailContextKey), claims.Email)
		c.Set(string(enums.UserRoleContextKey), claims.Role)

		return next(c)
	}
}
