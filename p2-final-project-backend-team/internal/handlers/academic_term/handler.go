package academic_term

import (
	"net/http"
	"strconv"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/requests"
	academic_term_service "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/academic_term"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/utils"
	"github.com/labstack/echo/v4"
)

type AcademicTermHandler interface {
	Create(c echo.Context) error
	GetAll(c echo.Context) error
	GetLatest(c echo.Context) error
	GetByID(c echo.Context) error
	UpdateIsOpen(c echo.Context) error
}

type handler struct {
	academicTermService academic_term_service.Service
}

func NewHandler(academicTermService academic_term_service.Service) AcademicTermHandler {
	return &handler{academicTermService: academicTermService}
}

func (h *handler) Create(c echo.Context) error {
	req := new(requests.CreateAcademicTermRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation failed",
			"errors":  err.Error(),
		})
	}

	if err := h.academicTermService.Create(c.Request().Context(), *req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to create academic term",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "academic term created successfully",
	})
}

func (h *handler) GetAll(c echo.Context) error {
	// Parse optional query parameter is_open
	var isOpenFilter *bool
	isOpenParam := c.QueryParam("is_open")
	if isOpenParam != "" {
		isOpen := isOpenParam == "true"
		isOpenFilter = &isOpen
	}

	terms, err := h.academicTermService.GetAll(c.Request().Context(), isOpenFilter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get academic terms",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "successfully retrieved all academic terms",
		"data":    terms,
	})
}

func (h *handler) GetLatest(c echo.Context) error {
	// Parse optional query parameter is_open
	var isOpenFilter *bool
	isOpenParam := c.QueryParam("is_open")
	if isOpenParam != "" {
		isOpen := isOpenParam == "true"
		isOpenFilter = &isOpen
	}

	term, err := h.academicTermService.GetLatest(c.Request().Context(), isOpenFilter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get latest academic term",
			"error":   err.Error(),
		})
	}
	if term == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "no academic term found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "successfully retrieved latest academic term",
		"data":    term,
	})
}

func (h *handler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid ID",
		})
	}

	term, err := h.academicTermService.FindByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get academic term",
			"error":   err.Error(),
		})
	}
	if term == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "academic term not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "successfully retrieved academic term",
		"data":    term,
	})
}

func (h *handler) UpdateIsOpen(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid ID",
		})
	}

	req := new(requests.UpdateAcademicTermRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	if err := h.academicTermService.UpdateIsOpen(c.Request().Context(), id, req.IsOpen); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to update academic term",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "academic term updated successfully",
	})
}
