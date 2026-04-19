package award

import (
	"net/http"
	"strconv"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/requests"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/award"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/utils"
	"github.com/labstack/echo/v4"
)

type AwardHandler interface {
	CreateAwardCategory(c echo.Context) error
	GetAllAwardCategories(c echo.Context) error
	GetAwardCategory(c echo.Context) error
	UpdateAwardCategory(c echo.Context) error
	DeleteAwardCategory(c echo.Context) error
	GetHonorRoll(c echo.Context) error
}

type handler struct {
	awardService award.Service
}

func NewHandler(awardService award.Service) AwardHandler {
	return &handler{
		awardService: awardService,
	}
}

func (h *handler) CreateAwardCategory(c echo.Context) error {
	req := new(requests.CreateAwardCategoryRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	response, err := h.awardService.CreateAwardCategory(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *handler) GetAwardCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID",
		})
	}

	response, err := h.awardService.GetAwardCategoryByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handler) GetAllAwardCategories(c echo.Context) error {
	responses, err := h.awardService.GetAllAwardCategories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *handler) UpdateAwardCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID",
		})
	}

	req := new(requests.UpdateAwardCategoryRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	response, err := h.awardService.UpdateAwardCategory(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handler) DeleteAwardCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID",
		})
	}

	if err := h.awardService.DeleteAwardCategory(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetHonorRoll returns approved award winners grouped by academic year and category
// GET /api/v1/honor-roll
// Query Parameters:
//   - academic_year: (optional) Filter by academic year (e.g., 2567)
//   - semester: (optional) Filter by semester ("first" or "second")
//   - campus_id: (optional) Filter by campus UUID
//   - award_id: (optional) Filter by specific award category ID
func (h *handler) GetHonorRoll(c echo.Context) error {
	var query requests.HonorRollQueryDTO
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid query parameters",
		})
	}

	response, err := h.awardService.GetHonorRoll(c.Request().Context(), &query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}
