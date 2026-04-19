package request

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/requests"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/storage"
	requestService "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/request"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type RequestHandler interface {
	// Student endpoints
	CreateRequest(c echo.Context) error
	GetMyRequests(c echo.Context) error
	GetRequestDetails(c echo.Context) error
	UploadDocument(c echo.Context) error

	// Unified approval/rejection (role-aware)
	ApproveRequest(c echo.Context) error
	RejectRequest(c echo.Context) error

	// Award category change
	ChangeAwardCategory(c echo.Context) error

	// Query endpoints
	GetRequestsByStatus(c echo.Context) error
	GetAllRequests(c echo.Context) error
}

type handler struct {
	requestService requestService.Service
	storageClient  storage.StorageClient
}

func NewHandler(requestService requestService.Service, storageClient storage.StorageClient) RequestHandler {
	return &handler{
		requestService: requestService,
		storageClient:  storageClient,
	}
}

// Helper function
func getUserIDAndRole(c echo.Context) (uuid.UUID, enums.UserRole, error) {
	userIDVal := c.Get(string(enums.UserIDContextKey))
	if userIDVal == nil {
		return uuid.Nil, "", exceptions.ErrUnauthenticated
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return uuid.Nil, "", exceptions.ErrInvalidUserIDType
	}

	roleVal := c.Get(string(enums.UserRoleContextKey))
	if roleVal == nil {
		return uuid.Nil, "", exceptions.ErrRoleNotFound
	}
	role, ok := roleVal.(enums.UserRole)
	if !ok {
		return uuid.Nil, "", exceptions.ErrInvalidUserRoleType
	}

	return userID, role, nil
}

func mapRequestError(err error) (int, string) {
	switch {
	case errors.Is(err, exceptions.ErrInvalidRequestID),
		errors.Is(err, exceptions.ErrInvalidRequestBody),
		errors.Is(err, exceptions.ErrInvalidQueryParameters),
		errors.Is(err, exceptions.ErrFileRequired),
		errors.Is(err, exceptions.ErrFileTooLarge),
		errors.Is(err, exceptions.ErrInvalidFileType):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, exceptions.ErrUnauthenticated),
		errors.Is(err, exceptions.ErrInvalidUserIDType),
		errors.Is(err, exceptions.ErrRoleNotFound),
		errors.Is(err, exceptions.ErrInvalidUserRoleType),
		errors.Is(err, exceptions.ErrViewerNotFound):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, exceptions.ErrForbidden),
		errors.Is(err, exceptions.ErrRoleNotAllowedToApprove),
		errors.Is(err, exceptions.ErrRoleNotAllowedToReject),
		errors.Is(err, exceptions.ErrInvalidRejectionRole):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, exceptions.ErrRequestNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, exceptions.ErrDuplicateRequestInTerm),
		errors.Is(err, exceptions.ErrInvalidApprovalTransition),
		errors.Is(err, exceptions.ErrInvalidRejectionTransition):
		return http.StatusConflict, err.Error()
	case errors.Is(err, exceptions.ErrApprovalStepNotConfigured),
		errors.Is(err, exceptions.ErrFailedToOpenFile),
		errors.Is(err, exceptions.ErrFailedToUploadFile):
		return http.StatusInternalServerError, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func writeError(c echo.Context, err error) error {
	status, message := mapRequestError(err)
	return c.JSON(status, map[string]string{
		"message": message,
	})
}

// CreateRequest - Student creates new request
// POST /api/v1/requests
func (h *handler) CreateRequest(c echo.Context) error {
	var req requests.CreateRequestDTO
	if err := c.Bind(&req); err != nil {
		return writeError(c, exceptions.ErrInvalidRequestBody)
	}

	userID, _, err := getUserIDAndRole(c)
	if err != nil {
		return writeError(c, err)
	}

	req.StudentID = userID

	if err := utils.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	response, err := h.requestService.CreateRequest(c.Request().Context(), &req)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusCreated, response)
}

// GetMyRequests - Student gets their own requests
// GET /api/v1/requests/my
func (h *handler) GetMyRequests(c echo.Context) error {
	userID, _, err := getUserIDAndRole(c)
	if err != nil {
		return writeError(c, err)
	}

	requests, err := h.requestService.GetMyRequests(c.Request().Context(), userID)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, requests)
}

// GetRequestDetails - Get full request details
// GET /api/v1/requests/:id
func (h *handler) GetRequestDetails(c echo.Context) error {
	idStr := c.Param("id")
	requestID, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid request ID")
		return writeError(c, exceptions.ErrInvalidRequestID)
	}

	userID, role, err := getUserIDAndRole(c)
	if err != nil {
		return writeError(c, err)
	}

	response, err := h.requestService.GetRequestDetails(c.Request().Context(), requestID, userID, role)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get request details")
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

// ApproveRequest - Unified approval endpoint (role-aware)
// POST /api/v1/requests/:id/approve
func (h *handler) ApproveRequest(c echo.Context) error {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return writeError(c, exceptions.ErrInvalidRequestID)
	}

	var dto requests.ApproveRequestDTO
	if err := c.Bind(&dto); err != nil {
		return writeError(c, exceptions.ErrInvalidRequestBody)
	}

	approverID, role, err := getUserIDAndRole(c)
	if err != nil {
		return writeError(c, err)
	}

	if err := h.requestService.ApproveRequest(c.Request().Context(), requestID, approverID, role, dto.Remark); err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Request approved successfully",
	})
}

// RejectRequest - Unified rejection endpoint (role-aware)
// POST /api/v1/requests/:id/reject
func (h *handler) RejectRequest(c echo.Context) error {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return writeError(c, exceptions.ErrInvalidRequestID)
	}

	var dto requests.RejectRequestDTO
	if err := c.Bind(&dto); err != nil {
		return writeError(c, exceptions.ErrInvalidRequestBody)
	}

	if err := utils.ValidateStruct(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	rejectorID, role, err := getUserIDAndRole(c)
	if err != nil {
		return writeError(c, err)
	}

	if err := h.requestService.RejectRequest(c.Request().Context(), requestID, rejectorID, dto.Remark, role); err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Request rejected successfully",
	})
}

// ChangeAwardCategory - Change award category (Student Development Division)
// POST /api/v1/requests/:id/change-category
func (h *handler) ChangeAwardCategory(c echo.Context) error {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return writeError(c, exceptions.ErrInvalidRequestID)
	}

	var dto requests.ChangeAwardCategoryDTO
	if err := c.Bind(&dto); err != nil {
		return writeError(c, exceptions.ErrInvalidRequestBody)
	}

	if err := utils.ValidateStruct(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	// Get changer ID from auth context
	changerID, _, err := getUserIDAndRole(c)
	if err != nil {
		return writeError(c, err)
	}

	if err := h.requestService.ChangeAwardCategory(c.Request().Context(), requestID, changerID, dto.NewCategoryID, dto.AdditionalData, dto.Remark); err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Award category changed",
	})
}

// UploadDocument - Upload document to request
// POST /api/v1/requests/:id/documents
func (h *handler) UploadDocument(c echo.Context) error {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return writeError(c, exceptions.ErrInvalidRequestID)
	}

	userID, role, err := getUserIDAndRole(c)
	if err != nil {
		return writeError(c, err)
	}

	// Verify request ownership - only the request owner can upload documents
	requestDetails, err := h.requestService.GetRequestDetails(c.Request().Context(), requestID, userID, role)
	if err != nil {
		return writeError(c, err)
	}
	if requestDetails.StudentID != userID {
		return writeError(c, exceptions.ErrForbidden)
	}

	// Get the file from form data
	file, err := c.FormFile("file")
	if err != nil {
		return writeError(c, exceptions.ErrFileRequired)
	}

	// Validate file size (max 10 MB)
	const maxFileSize = 10 * 1024 * 1024 // 10 MB in bytes
	if file.Size > maxFileSize {
		return writeError(c, exceptions.ErrFileTooLarge)
	}

	// Validate file type (only .png, .jpeg, .jpg, .pdf)
	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".png":  true,
		".jpeg": true,
		".jpg":  true,
		".pdf":  true,
	}
	if !allowedExtensions[fileExt] {
		return writeError(c, exceptions.ErrInvalidFileType)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return writeError(c, exceptions.ErrFailedToOpenFile)
	}
	defer src.Close()

	// Generate unique object key for MinIO
	// Format: documents/{requestID}/{uuid}.{ext}
	objectKey := fmt.Sprintf("documents/%s/%s%s", requestID.String(), uuid.New().String(), fileExt)

	// Set correct content type based on validated file extension
	contentType := ""
	switch fileExt {
	case ".png":
		contentType = "image/png"
	case ".jpeg", ".jpg":
		contentType = "image/jpeg"
	case ".pdf":
		contentType = "application/pdf"
	}

	// Upload to MinIO
	bucket := h.storageClient.GetBucketName()
	if err := h.storageClient.UploadFile(c.Request().Context(), objectKey, src, file.Size, contentType); err != nil {
		return writeError(c, exceptions.ErrFailedToUploadFile)
	}

	// Save metadata to database
	if err := h.requestService.UploadDocument(
		c.Request().Context(),
		requestID,
		bucket,
		objectKey,
		file.Filename,
		contentType,
		file.Size,
	); err != nil {
		// If database save fails, try to delete the uploaded file from MinIO
		_ = h.storageClient.DeleteFile(c.Request().Context(), objectKey)
		return writeError(c, err)
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message":      "Document uploaded successfully",
		"bucket":       bucket,
		"object_key":   objectKey,
		"file_name":    file.Filename,
		"size":         fmt.Sprintf("%d", file.Size),
		"content_type": contentType,
	})
}

// GetRequestsByStatus - Get requests by status (paginated)
// GET /api/v1/requests/status/:status
// Query Parameters:
//   - status: (path param) Request status to filter
//   - page: (optional) Page number (default: 1)
//   - limit: (optional) Items per page (default: 10)
func (h *handler) GetRequestsByStatus(c echo.Context) error {
	status := c.Param("status")

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	viewerID, role, err := getUserIDAndRole(c)
	if err != nil {
		return writeError(c, err)
	}

	response, err := h.requestService.GetRequestsByStatus(c.Request().Context(), viewerID, role, status, page, limit)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

// GetAllRequests - Get all requests with filters (paginated)
// GET /api/v1/requests
// Query Parameters:
//   - student_id: (optional) Filter by student ID
//   - award_id: (optional) Filter by award category ID
//   - academic_term_id: (optional) Filter by academic term ID
//   - current_status: (optional) Filter by current request status
//   - page: (optional) Page number (min: 1)
//   - limit: (optional) Items per page (min: 1, max: 100)
func (h *handler) GetAllRequests(c echo.Context) error {
	var query requests.RequestQueryDTO
	if err := c.Bind(&query); err != nil {
		return writeError(c, exceptions.ErrInvalidQueryParameters)
	}

	viewerID, role, err := getUserIDAndRole(c)
	if err != nil {
		return writeError(c, err)
	}

	response, err := h.requestService.GetRequestsWithFilters(c.Request().Context(), viewerID, role, &query)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, response)
}
