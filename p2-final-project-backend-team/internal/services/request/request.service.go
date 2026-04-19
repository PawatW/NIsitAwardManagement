package request

import (
	"context"
	"fmt"
	"time"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/requests"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/responses"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/storage"
	academicTermRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/academic_term"
	requestRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/request"
	userRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/user"
	"github.com/google/uuid"
)

type Service interface {
	// student ops
	CreateRequest(ctx context.Context, req *requests.CreateRequestDTO) (*responses.RequestResponse, error)
	GetMyRequests(ctx context.Context, studentID uuid.UUID) ([]*responses.RequestListResponse, error)

	// ✅ changed: must pass viewer context for access control
	GetRequestDetails(ctx context.Context, requestID uuid.UUID, viewerID uuid.UUID, role enums.UserRole) (*responses.RequestResponse, error)

	// approve/reject
	ApproveRequest(ctx context.Context, requestID uuid.UUID, approverID uuid.UUID, role enums.UserRole, remark string) error
	RejectRequest(ctx context.Context, requestID uuid.UUID, rejectorID uuid.UUID, remark string, role enums.UserRole) error

	// award category change
	ChangeAwardCategory(ctx context.Context, requestID uuid.UUID, changerID uuid.UUID, newCategoryID int, additionalData models.AdditionalData, remark string) error

	// ✅ changed: must pass viewer context for access control
	GetRequestsByStatus(ctx context.Context, viewerID uuid.UUID, role enums.UserRole, status string, page, limit int) (*responses.PaginatedRequestResponse, error)
	GetRequestsWithFilters(ctx context.Context, viewerID uuid.UUID, role enums.UserRole, query *requests.RequestQueryDTO) (*responses.PaginatedRequestResponse, error)

	// document ops
	UploadDocument(ctx context.Context, requestID uuid.UUID, bucket string, objectKey string, originalName string, contentType string, size int64) error
	GetRequestDocuments(ctx context.Context, requestID uuid.UUID) ([]*responses.RequestDocumentResponse, error)
}

type service struct {
	requestRepo      requestRepo.Repository
	userRepo         userRepo.Repository
	academicTermRepo academicTermRepo.Repository
	storageClient    storage.StorageClient
}

func NewService(requestRepo requestRepo.Repository, userRepo userRepo.Repository, academicTermRepo academicTermRepo.Repository, storageClient storage.StorageClient) Service {
	return &service{
		requestRepo:      requestRepo,
		userRepo:         userRepo,
		academicTermRepo: academicTermRepo,
		storageClient:    storageClient,
	}
}

// ---------- scope helpers ----------

func (s *service) getScopeForRole(ctx context.Context, viewerID uuid.UUID, role enums.UserRole) (facultyID *uuid.UUID, departmentID *uuid.UUID, _ error) {
	// ✅ policy A: Admin sees all
	if role == enums.Admin {
		return nil, nil, nil
	}

	// CommitteeChair: ไม่ได้ขอให้จำกัด → ปล่อยผ่านแบบ global (ปรับได้ทีหลัง)
	if role == enums.CommitteeChair {
		return nil, nil, nil
	}

	// Student: scope จะถูกบังคับเป็น student_id = viewerID ใน list/query อยู่แล้ว
	if role == enums.Student {
		return nil, nil, nil
	}

	u, err := s.userRepo.FindByIDWithRelations(ctx, viewerID)
	if err != nil {
		return nil, nil, err
	}
	if u == nil {
		return nil, nil, exceptions.ErrViewerNotFound
	}

	switch role {
	case enums.HeadOfDepartment:
		// จำกัดภาค
		if u.DepartmentID == nil {
			return nil, nil, exceptions.ErrForbidden
		}
		return nil, u.DepartmentID, nil
	case enums.ViceDean, enums.Dean:
		// จำกัดคณะ
		if u.FacultyID == nil {
			return nil, nil, exceptions.ErrForbidden
		}
		return u.FacultyID, nil, nil
	default:
		return nil, nil, nil
	}
}

func (s *service) assertRequestScope(ctx context.Context, viewerID uuid.UUID, role enums.UserRole, req *models.Request) error {
	// Admin sees all
	if role == enums.Admin {
		return nil
	}

	// Student sees own only
	if role == enums.Student {
		if req.StudentID != viewerID {
			return exceptions.ErrForbidden
		}
		return nil
	}

	// CommitteeChair: ปล่อยผ่าน (ปรับ policy ได้ทีหลัง)
	if role == enums.CommitteeChair {
		if !canRoleViewStatus(role, req.CurrentStatus) {
			return exceptions.ErrForbidden
		}
		return nil
	}

	// For HOD/ViceDean/Dean: compare org ids using viewer's scope vs student's scope
	facultyID, departmentID, err := s.getScopeForRole(ctx, viewerID, role)
	if err != nil {
		return err
	}

	student := req.Student
	// ต้อง preload Student มาก่อน
	if role == enums.HeadOfDepartment {
		if departmentID == nil || student.DepartmentID == nil || *departmentID != *student.DepartmentID {
			return exceptions.ErrForbidden
		}
	}
	if role == enums.ViceDean || role == enums.Dean {
		if facultyID == nil || student.FacultyID == nil || *facultyID != *student.FacultyID {
			return exceptions.ErrForbidden
		}
	}
	if !canRoleViewStatus(role, req.CurrentStatus) {
		return exceptions.ErrForbidden
	}

	return nil
}

// ---------- core features ----------

// Student submit
func (s *service) CreateRequest(ctx context.Context, req *requests.CreateRequestDTO) (*responses.RequestResponse, error) {
	// Check if academic term exists and is open
	academicTerm, err := s.academicTermRepo.FindByID(ctx, req.AcademicTermID)
	if err != nil {
		return nil, err
	}
	if academicTerm == nil {
		return nil, exceptions.ErrAcademicTermNotFound
	}
	if !academicTerm.IsOpen {
		return nil, exceptions.ErrAcademicTermNotOpen
	}

	exists, err := s.requestRepo.ExistsByStudentAndTerm(ctx, req.StudentID, req.AcademicTermID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, exceptions.ErrDuplicateRequestInTerm
	}

	request := &models.Request{
		ID:                  uuid.New(),
		StudentID:           req.StudentID,
		AwardID:             req.AwardID,
		AcademicTermID:      req.AcademicTermID,
		AdditionalData:      req.AdditionalData,
		SnapshotStudyYear:   req.SnapshotStudyYear,
		SnapshotGPA:         req.SnapshotGPA,
		SnapshotAdvisor:     req.SnapshotAdvisor,
		SnapshotDateOfBirth: req.SnapshotDateOfBirth,
		SnapshotPhone:       req.SnapshotPhone,
		SnapshotAddress:     req.SnapshotAddress,
		CurrentStatus:       string(enums.PendingHOD),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	history := &models.RequestStatusHistory{
		RequestID: request.ID,
		Status:    string(enums.PendingHOD),
		UpdatedBy: req.StudentID,
		Remark:    "Request submitted",
		UpdatedAt: time.Now(),
	}

	if err := s.requestRepo.CreateWithHistory(ctx, request, history); err != nil {
		return nil, err
	}

	return s.toResponse(request), nil
}

type approvalStep struct {
	validFrom  []enums.RequestStatus
	nextStatus enums.RequestStatus
	label      string
}

var approvalSteps = map[string]approvalStep{
	"HOD":       {validFrom: []enums.RequestStatus{enums.PendingHOD}, nextStatus: enums.PendingViceDean, label: "HOD"},
	"ViceDean":  {validFrom: []enums.RequestStatus{enums.PendingViceDean}, nextStatus: enums.PendingDean, label: "Vice Dean"},
	"Dean":      {validFrom: []enums.RequestStatus{enums.PendingDean, enums.AwardChanged}, nextStatus: enums.PendingCommittee, label: "Dean"},
	"Committee": {validFrom: []enums.RequestStatus{enums.PendingCommittee}, nextStatus: enums.Approved, label: "Committee"},
}

var visibleStatusesByRole = map[enums.UserRole]map[enums.RequestStatus]struct{}{
	enums.HeadOfDepartment: {
		enums.PendingHOD:    {},
		enums.RejectedByHOD: {},
	},
	enums.ViceDean: {
		enums.PendingViceDean:    {},
		enums.RejectedByViceDean: {},
	},
	enums.Dean: {
		enums.PendingDean:    {},
		enums.RejectedByDean: {},
		enums.AwardChanged:   {},
	},
	enums.CommitteeChair: {
		enums.PendingCommittee:    {},
		enums.Approved:            {},
		enums.RejectedByCommittee: {},
	},
}

func canRoleViewStatus(role enums.UserRole, status string) bool {
	allowed, ok := visibleStatusesByRole[role]
	if !ok {
		return true
	}

	_, exists := allowed[enums.RequestStatus(status)]
	return exists
}

func allowedStatusesForRole(role enums.UserRole) []string {
	allowed, ok := visibleStatusesByRole[role]
	if !ok {
		return nil
	}

	statuses := make([]string, 0, len(allowed))
	for status := range allowed {
		statuses = append(statuses, string(status))
	}
	return statuses
}

func approvalStepForRole(role enums.UserRole) (approvalStep, error) {
	switch role {
	case enums.HeadOfDepartment:
		return approvalSteps["HOD"], nil
	case enums.ViceDean:
		return approvalSteps["ViceDean"], nil
	case enums.Dean:
		return approvalSteps["Dean"], nil
	case enums.CommitteeChair:
		return approvalSteps["Committee"], nil
	case enums.Admin:
		return approvalStep{}, exceptions.ErrRoleNotAllowedToApprove
	default:
		return approvalStep{}, exceptions.ErrRoleNotAllowedToApprove
	}
}

// approveRequest uses atomic update
func (s *service) approveRequest(ctx context.Context, requestID uuid.UUID, approverID uuid.UUID, remark string, step approvalStep) error {
	request, err := s.requestRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}

	valid := false
	for _, status := range step.validFrom {
		if request.CurrentStatus == string(status) {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("%w: %s", exceptions.ErrInvalidApprovalTransition, step.label)
	}

	request.CurrentStatus = string(step.nextStatus)
	request.UpdatedAt = time.Now()

	history := &models.RequestStatusHistory{
		RequestID: requestID,
		Status:    string(step.nextStatus),
		UpdatedBy: approverID,
		Remark:    remark,
		UpdatedAt: time.Now(),
	}

	return s.requestRepo.UpdateWithHistory(ctx, request, history)
}

func (s *service) ApproveRequest(ctx context.Context, requestID uuid.UUID, approverID uuid.UUID, role enums.UserRole, remark string) error {
	// ✅ scope check: preload Student for org compare
	reqWithStudent, err := s.requestRepo.FindByIDWithStudent(ctx, requestID)
	if err != nil {
		return err
	}
	if err := s.assertRequestScope(ctx, approverID, role, reqWithStudent); err != nil {
		return err
	}

	step, err := approvalStepForRole(role)
	if err != nil {
		return err
	}
	if len(step.validFrom) == 0 {
		return exceptions.ErrApprovalStepNotConfigured
	}

	return s.approveRequest(ctx, requestID, approverID, remark, step)
}

func (s *service) RejectRequest(ctx context.Context, requestID uuid.UUID, rejectorID uuid.UUID, remark string, role enums.UserRole) error {
	// ✅ scope check: preload Student for org compare
	reqWithStudent, err := s.requestRepo.FindByIDWithStudent(ctx, requestID)
	if err != nil {
		return err
	}
	if err := s.assertRequestScope(ctx, rejectorID, role, reqWithStudent); err != nil {
		return err
	}

	request, err := s.requestRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}

	step, stepErr := approvalStepForRole(role)
	if stepErr != nil {
		if stepErr == exceptions.ErrRoleNotAllowedToApprove {
			return exceptions.ErrRoleNotAllowedToReject
		}
		return stepErr
	}

	allowed := false
	for _, status := range step.validFrom {
		if request.CurrentStatus == string(status) {
			allowed = true
			break
		}
	}
	if !allowed {
		return exceptions.ErrInvalidRejectionTransition
	}

	var newStatus enums.RequestStatus
	switch role {
	case enums.HeadOfDepartment:
		newStatus = enums.RejectedByHOD
	case enums.ViceDean:
		newStatus = enums.RejectedByViceDean
	case enums.Dean:
		newStatus = enums.RejectedByDean
	case enums.CommitteeChair:
		newStatus = enums.RejectedByCommittee
	case enums.Admin:
		return exceptions.ErrRoleNotAllowedToReject
	default:
		return exceptions.ErrInvalidRejectionRole
	}

	request.CurrentStatus = string(newStatus)
	request.UpdatedAt = time.Now()

	history := &models.RequestStatusHistory{
		RequestID: requestID,
		Status:    string(newStatus),
		UpdatedBy: rejectorID,
		Remark:    remark,
		UpdatedAt: time.Now(),
	}

	return s.requestRepo.UpdateWithHistory(ctx, request, history)
}

func (s *service) ChangeAwardCategory(ctx context.Context, requestID uuid.UUID, changerID uuid.UUID, newCategoryID int, additionalData models.AdditionalData, remark string) error {
	request, err := s.requestRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}

	if request.CurrentStatus != string(enums.PendingCommittee) {
		return exceptions.ErrInvalidAwardCategoryChangeStatus
	}

	oldCategoryID := request.AwardID

	request.AwardID = newCategoryID
	request.AdditionalData = additionalData
	request.CurrentStatus = string(enums.AwardChanged)
	request.UpdatedAt = time.Now()

	change := &models.AwardCategoryChange{
		RequestID:     requestID,
		OldCategoryID: oldCategoryID,
		NewCategoryID: newCategoryID,
		ChangedBy:     changerID,
		ChangedAt:     time.Now(),
	}

	history := &models.RequestStatusHistory{
		RequestID: requestID,
		Status:    string(enums.AwardChanged),
		UpdatedBy: changerID,
		Remark:    remark,
		UpdatedAt: time.Now(),
	}

	return s.requestRepo.UpdateWithCategoryChange(ctx, request, change, history)
}

func (s *service) GetRequestDetails(ctx context.Context, requestID uuid.UUID, viewerID uuid.UUID, role enums.UserRole) (*responses.RequestResponse, error) {
	request, err := s.requestRepo.FindByIDWithRelations(ctx, requestID)
	if err != nil {
		return nil, err
	}

	// ✅ scope check (Admin ผ่าน, Student own-only, HOD/Dean by org)
	if err := s.assertRequestScope(ctx, viewerID, role, request); err != nil {
		return nil, err
	}

	return s.toDetailedResponse(request), nil
}

func (s *service) GetMyRequests(ctx context.Context, studentID uuid.UUID) ([]*responses.RequestListResponse, error) {
	requests, err := s.requestRepo.FindByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	return s.toListResponses(requests), nil
}

func (s *service) GetRequestsByStatus(ctx context.Context, viewerID uuid.UUID, role enums.UserRole, status string, page, limit int) (*responses.PaginatedRequestResponse, error) {
	return s.GetRequestsWithFilters(ctx, viewerID, role, &requests.RequestQueryDTO{
		CurrentStatus: status,
		Page:          page,
		Limit:         limit,
	})
}

func (s *service) GetRequestsWithFilters(ctx context.Context, viewerID uuid.UUID, role enums.UserRole, query *requests.RequestQueryDTO) (*responses.PaginatedRequestResponse, error) {
	// Parse query filters
	var studentID *uuid.UUID
	if query.StudentID != "" {
		parsed, err := uuid.Parse(query.StudentID)
		if err == nil {
			studentID = &parsed
		}
	}

	// ✅ enforce Student scope to own (กันรั่ว endpoint /requests)
	if role == enums.Student {
		own := viewerID
		studentID = &own
	}

	var awardID *int
	if query.AwardID > 0 {
		awardID = &query.AwardID
	}

	var termID *int
	if query.AcademicTermID > 0 {
		termID = &query.AcademicTermID
	}

	var status *string
	if query.CurrentStatus != "" {
		if !canRoleViewStatus(role, query.CurrentStatus) {
			return nil, exceptions.ErrForbidden
		}
		status = &query.CurrentStatus
	}

	page := query.Page
	if page < 1 {
		page = 1
	}

	limit := query.Limit
	if limit < 1 {
		limit = 10
	}

	// ✅ get org scope for HOD/ViceDean/Dean (Admin nil,nil)
	facultyID, departmentID, err := s.getScopeForRole(ctx, viewerID, role)
	if err != nil {
		return nil, err
	}

	requestsList, total, err := s.requestRepo.FindAllWithFiltersAndScope(
		ctx,
		studentID, awardID, termID, status,
		allowedStatusesForRole(role),
		facultyID, departmentID,
		page, limit,
	)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &responses.PaginatedRequestResponse{
		Data:       s.toListResponses(requestsList),
		TotalItems: total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *service) UploadDocument(ctx context.Context, requestID uuid.UUID, bucket string, objectKey string, originalName string, contentType string, size int64) error {
	document := &models.RequestDocument{
		RequestID:    requestID,
		Bucket:       bucket,
		ObjectKey:    objectKey,
		OriginalName: originalName,
		ContentType:  contentType,
		Size:         size,
	}
	return s.requestRepo.AddDocument(ctx, document)
}

func (s *service) GetRequestDocuments(ctx context.Context, requestID uuid.UUID) ([]*responses.RequestDocumentResponse, error) {
	documents, err := s.requestRepo.GetDocuments(ctx, requestID)
	if err != nil {
		return nil, err
	}

	var result []*responses.RequestDocumentResponse
	for _, doc := range documents {
		result = append(result, s.toDocumentResponse(doc))
	}

	return result, nil
}

func (s *service) toDocumentResponse(doc *models.RequestDocument) *responses.RequestDocumentResponse {
	fileURL := s.storageClient.GetPublicURL(doc.ObjectKey)
	if fileURL == "" {
		var err error
		fileURL, err = s.storageClient.GetPresignedURL(context.Background(), doc.ObjectKey, 24*time.Hour)
		if err != nil {
			fileURL = fmt.Sprintf("/api/v1/documents/%s/%s", doc.Bucket, doc.ObjectKey)
		}
	}

	return &responses.RequestDocumentResponse{
		ID:           doc.ID,
		RequestID:    doc.RequestID,
		Bucket:       doc.Bucket,
		ObjectKey:    doc.ObjectKey,
		OriginalName: doc.OriginalName,
		ContentType:  doc.ContentType,
		Size:         doc.Size,
		FileURL:      fileURL,
		UploadedAt:   doc.CreatedAt,
	}
}

func (s *service) toResponse(request *models.Request) *responses.RequestResponse {
	return &responses.RequestResponse{
		ID:                  request.ID,
		StudentID:           request.StudentID,
		AwardID:             request.AwardID,
		AcademicTermID:      request.AcademicTermID,
		AdditionalData:      request.AdditionalData,
		SnapshotStudyYear:   request.SnapshotStudyYear,
		SnapshotGPA:         request.SnapshotGPA,
		SnapshotAdvisor:     request.SnapshotAdvisor,
		SnapshotDateOfBirth: request.SnapshotDateOfBirth,
		SnapshotPhone:       request.SnapshotPhone,
		SnapshotAddress:     request.SnapshotAddress,
		CurrentStatus:       request.CurrentStatus,
		CreatedAt:           request.CreatedAt,
		UpdatedAt:           request.UpdatedAt,
	}
}

func (s *service) toDetailedResponse(request *models.Request) *responses.RequestResponse {
	resp := s.toResponse(request)

	if request.Student.UserID != uuid.Nil {
		resp.Student = s.toStudentBasicResponse(&request.Student)
	}

	if request.Award.ID > 0 {
		resp.Award = &responses.AwardCategoryResponse{
			ID:            request.Award.ID,
			Name:          request.Award.Name,
			Description:   request.Award.Description,
			FormStructure: request.Award.FormStructure,
			IsActive:      request.Award.IsActive,
			CreatedAt:     request.Award.CreatedAt,
			UpdatedAt:     request.Award.UpdatedAt,
		}
	}

	if len(request.Documents) > 0 {
		resp.Documents = make([]responses.RequestDocumentResponse, len(request.Documents))
		for i := range request.Documents {
			resp.Documents[i] = *s.toDocumentResponse(&request.Documents[i])
		}
	}

	return resp
}

func (s *service) toListResponses(requests []*models.Request) []*responses.RequestListResponse {
	var listResponses []*responses.RequestListResponse
	for _, req := range requests {
		student := (*responses.UserBasicResponse)(nil)
		if req.Student.UserID != uuid.Nil {
			student = s.toStudentBasicResponse(&req.Student)
		}
		awardName := ""
		if req.Award.ID > 0 {
			awardName = req.Award.Name
		}

		listResponses = append(listResponses, &responses.RequestListResponse{
			ID:            req.ID,
			Student:       student,
			AwardName:     awardName,
			CurrentStatus: req.CurrentStatus,
			SnapshotGPA:   req.SnapshotGPA,
			CreatedAt:     req.CreatedAt,
		})
	}
	return listResponses
}

func (s *service) toStudentBasicResponse(student *models.User) *responses.UserBasicResponse {
	if student == nil || student.UserID == uuid.Nil {
		return nil
	}

	resp := &responses.UserBasicResponse{
		ID:          student.UserID,
		FirstName:   student.FirstName,
		LastName:    student.LastName,
		PhoneNumber: student.PhoneNumber,
	}

	if student.NisitID != nil {
		resp.NisitID = *student.NisitID
	}
	if student.Campus != nil {
		resp.CampusName = student.Campus.Name
	}
	if student.Faculty != nil {
		resp.FacultyName = student.Faculty.Name
	}
	if student.Department != nil {
		resp.DepartmentName = student.Department.Name
	}

	return resp
}
