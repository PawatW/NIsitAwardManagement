package admin

import (
	"net/http"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/requests"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/admin"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/utils"
	"github.com/labstack/echo/v4"
)

type AdminHandler interface {
	CreateUser(c echo.Context) error
	GetAllUsers(c echo.Context) error
	UpdateUserStatus(c echo.Context) error
	CreateCampus(c echo.Context) error
	CreateFaculty(c echo.Context) error
	CreateDepartment(c echo.Context) error
}

type handler struct {
	adminService admin.Service
}

func NewAdminHandler(
	adminService admin.Service,
) AdminHandler {
	return &handler{
		adminService: adminService,
	}
}

// CreateUser creates a new staff/faculty user
func (h *handler) CreateUser(c echo.Context) error {
	var req requests.CreateStaffRequest
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

	if err := h.adminService.CreateUser(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to create user",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "user created successfully",
	})
}

func (h *handler) GetAllUsers(c echo.Context) error {
	users, err := h.adminService.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get all users",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "successfully retrieved all users",
		"data":    users,
	})
}

func (h *handler) UpdateUserStatus(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "user ID is required",
		})
	}

	var req requests.UpdateUserStatusRequest
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

	updatedUser, err := h.adminService.UpdateUserStatus(c.Request().Context(), userID, req)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "user not found",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to update user status",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user status updated successfully",
		"data":    updatedUser,
	})
}

func (h *handler) CreateCampus(c echo.Context) error {
	var req requests.CreateCampusRequest
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

	if err := h.adminService.CreateCampus(c.Request().Context(), req.Name); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to create campus",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "campus created successfully",
	})
}

func (h *handler) CreateFaculty(c echo.Context) error {
	var req requests.CreateFacultyRequest
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

	if err := h.adminService.CreateFaculty(c.Request().Context(), req.Name, req.CampusID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to create faculty",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "faculty created successfully",
	})
}

func (h *handler) CreateDepartment(c echo.Context) error {
	var req requests.CreateDepartmentRequest
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

	if err := h.adminService.CreateDepartment(c.Request().Context(), req.Name, req.FacultyID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to create department",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "department created successfully",
	})
}
