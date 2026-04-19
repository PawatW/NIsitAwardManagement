package organization

import (
	"net/http"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/campus"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/department"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/faculty"
	"github.com/labstack/echo/v4"
)

type OrganizationHandler interface {
	GetCampuses(c echo.Context) error
	GetFacultiesByCampus(c echo.Context) error
	GetDepartmentsByFaculty(c echo.Context) error
}

type handler struct {
	campusService     campus.Service
	facultyService    faculty.Service
	departmentService department.Service
}

func NewOrganizationHandler(
	campusService campus.Service,
	facultyService faculty.Service,
	departmentService department.Service,
) OrganizationHandler {
	return &handler{
		campusService:     campusService,
		facultyService:    facultyService,
		departmentService: departmentService,
	}
}

func (h *handler) GetCampuses(c echo.Context) error {
	campuses, err := h.campusService.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get campuses",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "successfully retrieved campuses",
		"data":    campuses,
	})
}

func (h *handler) GetFacultiesByCampus(c echo.Context) error {
	campusID := c.Param("campus_id")
	ptrCampusID, err := h.campusService.ValidateID(c.Request().Context(), campusID)
	if err != nil || ptrCampusID == nil {
		errorMessage := ""
		if err != nil {
			errorMessage = err.Error()
		}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid campus id",
			"error":   errorMessage,
		})
	}

	faculties, err := h.facultyService.GetByCampusID(c.Request().Context(), *ptrCampusID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get faculties",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "successfully retrieved faculties",
		"data":    faculties,
	})
}

func (h *handler) GetDepartmentsByFaculty(c echo.Context) error {
	facultyID := c.Param("faculty_id")
	ptrFacultyID, err := h.facultyService.ValidateID(c.Request().Context(), facultyID)
	if err != nil || ptrFacultyID == nil {
		errorMessage := ""
		if err != nil {
			errorMessage = err.Error()
		}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid faculty id",
			"error":   errorMessage,
		})
	}

	departments, err := h.departmentService.GetByFacultyID(c.Request().Context(), *ptrFacultyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get departments",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "successfully retrieved departments",
		"data":    departments,
	})
}