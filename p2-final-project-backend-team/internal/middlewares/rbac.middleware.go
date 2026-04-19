package middlewares

import (
	"net/http"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/labstack/echo/v4"
)

type RBACMiddleware interface {
	RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc
}

type rbacMiddleware struct{}

func NewRBACMiddleware() RBACMiddleware {
	return &rbacMiddleware{}
}

// RequireAdmin middleware checks if the authenticated user has admin role
func (r *rbacMiddleware) RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get user role from context (set by auth Middleware)
		roleValue := c.Get(string(enums.UserRoleContextKey))
		if roleValue == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "unauthorized - no role found",
			})
		}

		role, ok := roleValue.(enums.UserRole)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "unauthorized - invalid role type",
			})
		}

		// Check if user is admin
		if role != enums.Admin {
			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "forbidden - admin access required",
			})
		}

		return next(c)
	}
}
