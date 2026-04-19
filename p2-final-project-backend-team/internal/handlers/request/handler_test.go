package request

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/requests"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/responses"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type mockRequestService struct {
	createRequestFn         func(ctx context.Context, req *requests.CreateRequestDTO) (*responses.RequestResponse, error)
	getMyRequestsFn         func(ctx context.Context, studentID uuid.UUID) ([]*responses.RequestListResponse, error)
	getRequestDetailsFn     func(ctx context.Context, requestID uuid.UUID, viewerID uuid.UUID, role enums.UserRole) (*responses.RequestResponse, error)
	approveRequestFn        func(ctx context.Context, requestID uuid.UUID, approverID uuid.UUID, role enums.UserRole, remark string) error
	rejectRequestFn         func(ctx context.Context, requestID uuid.UUID, rejectorID uuid.UUID, remark string, role enums.UserRole) error
	changeAwardCategoryFn   func(ctx context.Context, requestID uuid.UUID, changerID uuid.UUID, newCategoryID int, additionalData models.AdditionalData, remark string) error
	getRequestsByStatusFn   func(ctx context.Context, viewerID uuid.UUID, role enums.UserRole, status string, page, limit int) (*responses.PaginatedRequestResponse, error)
	getRequestsWithFilterFn func(ctx context.Context, viewerID uuid.UUID, role enums.UserRole, query *requests.RequestQueryDTO) (*responses.PaginatedRequestResponse, error)
	uploadDocumentFn        func(ctx context.Context, requestID uuid.UUID, bucket string, objectKey string, originalName string, contentType string, size int64) error
	getRequestDocumentsFn   func(ctx context.Context, requestID uuid.UUID) ([]*responses.RequestDocumentResponse, error)
}

func (m *mockRequestService) CreateRequest(ctx context.Context, req *requests.CreateRequestDTO) (*responses.RequestResponse, error) {
	return m.createRequestFn(ctx, req)
}

func (m *mockRequestService) GetMyRequests(ctx context.Context, studentID uuid.UUID) ([]*responses.RequestListResponse, error) {
	return m.getMyRequestsFn(ctx, studentID)
}

func (m *mockRequestService) GetRequestDetails(ctx context.Context, requestID uuid.UUID, viewerID uuid.UUID, role enums.UserRole) (*responses.RequestResponse, error) {
	return m.getRequestDetailsFn(ctx, requestID, viewerID, role)
}

func (m *mockRequestService) ApproveRequest(ctx context.Context, requestID uuid.UUID, approverID uuid.UUID, role enums.UserRole, remark string) error {
	return m.approveRequestFn(ctx, requestID, approverID, role, remark)
}

func (m *mockRequestService) RejectRequest(ctx context.Context, requestID uuid.UUID, rejectorID uuid.UUID, remark string, role enums.UserRole) error {
	return m.rejectRequestFn(ctx, requestID, rejectorID, remark, role)
}

func (m *mockRequestService) ChangeAwardCategory(ctx context.Context, requestID uuid.UUID, changerID uuid.UUID, newCategoryID int, additionalData models.AdditionalData, remark string) error {
	return m.changeAwardCategoryFn(ctx, requestID, changerID, newCategoryID, additionalData, remark)
}

func (m *mockRequestService) GetRequestsByStatus(ctx context.Context, viewerID uuid.UUID, role enums.UserRole, status string, page, limit int) (*responses.PaginatedRequestResponse, error) {
	return m.getRequestsByStatusFn(ctx, viewerID, role, status, page, limit)
}

func (m *mockRequestService) GetRequestsWithFilters(ctx context.Context, viewerID uuid.UUID, role enums.UserRole, query *requests.RequestQueryDTO) (*responses.PaginatedRequestResponse, error) {
	return m.getRequestsWithFilterFn(ctx, viewerID, role, query)
}

func (m *mockRequestService) UploadDocument(ctx context.Context, requestID uuid.UUID, bucket string, objectKey string, originalName string, contentType string, size int64) error {
	return m.uploadDocumentFn(ctx, requestID, bucket, objectKey, originalName, contentType, size)
}

func (m *mockRequestService) GetRequestDocuments(ctx context.Context, requestID uuid.UUID) ([]*responses.RequestDocumentResponse, error) {
	return m.getRequestDocumentsFn(ctx, requestID)
}

type mockStorage struct{}

func (m *mockStorage) UploadFile(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) error {
	return nil
}
func (m *mockStorage) DownloadFile(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	return nil, nil
}
func (m *mockStorage) DeleteFile(ctx context.Context, objectKey string) error {
	return nil
}
func (m *mockStorage) GetPresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	return "", nil
}
func (m *mockStorage) GetPublicURL(objectKey string) string {
	return ""
}
func (m *mockStorage) EnsureBucketExists(ctx context.Context) error {
	return nil
}
func (m *mockStorage) GetBucketName() string {
	return "bucket"
}

func mustMessage(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	msg, ok := body["message"].(string)
	if !ok {
		t.Fatalf("message field not found: %v", body)
	}
	return msg
}

func TestCreateRequest_DuplicateTermReturns409(t *testing.T) {
	e := echo.New()
	userID := uuid.New()
	service := &mockRequestService{
		createRequestFn: func(ctx context.Context, req *requests.CreateRequestDTO) (*responses.RequestResponse, error) {
			return nil, exceptions.ErrDuplicateRequestInTerm
		},
	}
	h := NewHandler(service, &mockStorage{})

	body := `{"award_id":1,"academic_term_id":1,"additional_data":{"k":"v"},"snapshot_study_year":3,"snapshot_gpa":3.25,"snapshot_advisor":"advisor"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/requests", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(string(enums.UserIDContextKey), userID)
	c.Set(string(enums.UserRoleContextKey), enums.Student)

	if err := h.CreateRequest(c); err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", rec.Code)
	}
	if got := mustMessage(t, rec); got != exceptions.ErrDuplicateRequestInTerm.Error() {
		t.Fatalf("unexpected message: %s", got)
	}
}

func TestGetRequestDetails_NotFoundReturns404(t *testing.T) {
	e := echo.New()
	userID := uuid.New()
	requestID := uuid.New()
	service := &mockRequestService{
		getRequestDetailsFn: func(ctx context.Context, requestID uuid.UUID, viewerID uuid.UUID, role enums.UserRole) (*responses.RequestResponse, error) {
			return nil, exceptions.ErrRequestNotFound
		},
	}
	h := NewHandler(service, &mockStorage{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/requests/"+requestID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/requests/:id")
	c.SetParamNames("id")
	c.SetParamValues(requestID.String())
	c.Set(string(enums.UserIDContextKey), userID)
	c.Set(string(enums.UserRoleContextKey), enums.Student)

	if err := h.GetRequestDetails(c); err != nil {
		t.Fatalf("GetRequestDetails() error = %v", err)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
	_ = mustMessage(t, rec)
}

func TestApproveRequest_InvalidTransitionReturns409(t *testing.T) {
	e := echo.New()
	userID := uuid.New()
	requestID := uuid.New()
	service := &mockRequestService{
		approveRequestFn: func(ctx context.Context, requestID uuid.UUID, approverID uuid.UUID, role enums.UserRole, remark string) error {
			return errors.Join(exceptions.ErrInvalidApprovalTransition, errors.New("HOD"))
		},
	}
	h := NewHandler(service, &mockStorage{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/requests/"+requestID.String()+"/approve", strings.NewReader(`{"remark":"ok"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/requests/:id/approve")
	c.SetParamNames("id")
	c.SetParamValues(requestID.String())
	c.Set(string(enums.UserIDContextKey), userID)
	c.Set(string(enums.UserRoleContextKey), enums.HeadOfDepartment)

	if err := h.ApproveRequest(c); err != nil {
		t.Fatalf("ApproveRequest() error = %v", err)
	}
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", rec.Code)
	}
	_ = mustMessage(t, rec)
}

func TestRejectRequest_InvalidRoleReturns403(t *testing.T) {
	e := echo.New()
	userID := uuid.New()
	requestID := uuid.New()
	service := &mockRequestService{
		rejectRequestFn: func(ctx context.Context, requestID uuid.UUID, rejectorID uuid.UUID, remark string, role enums.UserRole) error {
			return exceptions.ErrInvalidRejectionRole
		},
	}
	h := NewHandler(service, &mockStorage{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/requests/"+requestID.String()+"/reject", strings.NewReader(`{"remark":"x"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/requests/:id/reject")
	c.SetParamNames("id")
	c.SetParamValues(requestID.String())
	c.Set(string(enums.UserIDContextKey), userID)
	c.Set(string(enums.UserRoleContextKey), enums.Admin)

	if err := h.RejectRequest(c); err != nil {
		t.Fatalf("RejectRequest() error = %v", err)
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
	_ = mustMessage(t, rec)
}

func TestUploadDocument_NoFileReturns400(t *testing.T) {
	e := echo.New()
	userID := uuid.New()
	requestID := uuid.New()
	service := &mockRequestService{
		getRequestDetailsFn: func(ctx context.Context, requestID uuid.UUID, viewerID uuid.UUID, role enums.UserRole) (*responses.RequestResponse, error) {
			return &responses.RequestResponse{StudentID: viewerID}, nil
		},
	}
	h := NewHandler(service, &mockStorage{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/requests/"+requestID.String()+"/documents", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/requests/:id/documents")
	c.SetParamNames("id")
	c.SetParamValues(requestID.String())
	c.Set(string(enums.UserIDContextKey), userID)
	c.Set(string(enums.UserRoleContextKey), enums.Student)

	if err := h.UploadDocument(c); err != nil {
		t.Fatalf("UploadDocument() error = %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	_ = mustMessage(t, rec)
}

func TestGetMyRequests_UnknownErrorReturns500(t *testing.T) {
	e := echo.New()
	userID := uuid.New()
	service := &mockRequestService{
		getMyRequestsFn: func(ctx context.Context, studentID uuid.UUID) ([]*responses.RequestListResponse, error) {
			return nil, errors.New("boom")
		},
	}
	h := NewHandler(service, &mockStorage{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/requests/my", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(string(enums.UserIDContextKey), userID)
	c.Set(string(enums.UserRoleContextKey), enums.Student)

	if err := h.GetMyRequests(c); err != nil {
		t.Fatalf("GetMyRequests() error = %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	if got := mustMessage(t, rec); got != "internal server error" {
		t.Fatalf("unexpected message: %s", got)
	}
}

func TestMapRequestError_StatusMapping(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{name: "invalid request id", err: exceptions.ErrInvalidRequestID, expected: http.StatusBadRequest},
		{name: "unauthenticated", err: exceptions.ErrUnauthenticated, expected: http.StatusUnauthorized},
		{name: "forbidden", err: exceptions.ErrForbidden, expected: http.StatusForbidden},
		{name: "request not found", err: exceptions.ErrRequestNotFound, expected: http.StatusNotFound},
		{name: "duplicate request", err: exceptions.ErrDuplicateRequestInTerm, expected: http.StatusConflict},
		{name: "invalid approval transition wrapped", err: errors.Join(exceptions.ErrInvalidApprovalTransition, errors.New("detail")), expected: http.StatusConflict},
		{name: "invalid rejection transition", err: exceptions.ErrInvalidRejectionTransition, expected: http.StatusConflict},
		{name: "upload fail", err: exceptions.ErrFailedToUploadFile, expected: http.StatusInternalServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status, _ := mapRequestError(tc.err)
			if status != tc.expected {
				t.Fatalf("expected %d, got %d", tc.expected, status)
			}
		})
	}
}
