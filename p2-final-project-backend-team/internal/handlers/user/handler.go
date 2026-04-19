package user

import (
	"net/http"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/requests"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/auth"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/user"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/utils"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Register(c echo.Context) error
	GetCurrentUser(c echo.Context) error
}

type handler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) Handler {
	return &handler{
		userService: userService,
	}
}

func (h *handler) Register(c echo.Context) error {
	var req requests.StudentRegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation failed",
			"errors":  utils.AddUserFormatValidationError(err),
		})
	}

	if err := h.userService.Register(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to register student",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Student registered successfully. Please login with Google to continue.",
	})
}

// GetCurrentUser returns the current authenticated user's information from database
func (h *handler) GetCurrentUser(c echo.Context) error {
	// Get claims from context (set by auth middleware)
	claims, ok := c.Get("user").(*auth.JWTClaims)
	if !ok || claims == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "unauthorized",
		})
	}

	// Query database for full user data
	userData, err := h.userService.GetByID(c.Request().Context(), claims.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "user not found",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user retrieved successfully",
		"data":    userData,
	})
}
