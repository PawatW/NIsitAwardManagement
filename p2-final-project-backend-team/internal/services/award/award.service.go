package award

import (
	"context"
	"errors"
	"fmt"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/requests"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/responses"
	awardRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/award"
	requestRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	CreateAwardCategory(ctx context.Context, req *requests.CreateAwardCategoryRequest) (*responses.AwardCategoryResponse, error)
	GetAwardCategoryByID(ctx context.Context, id int) (*responses.AwardCategoryResponse, error)
	GetAllAwardCategories(ctx context.Context) ([]*responses.AwardCategoryResponse, error)
	UpdateAwardCategory(ctx context.Context, id int, req *requests.UpdateAwardCategoryRequest) (*responses.AwardCategoryResponse, error)
	DeleteAwardCategory(ctx context.Context, id int) error
	GetHonorRoll(ctx context.Context, query *requests.HonorRollQueryDTO) (*responses.HonorRollResponse, error)
}

type service struct {
	awardRepo   awardRepo.Repository
	requestRepo requestRepo.Repository
}

func NewService(awardRepo awardRepo.Repository, requestRepo requestRepo.Repository) Service {
	return &service{
		awardRepo:   awardRepo,
		requestRepo: requestRepo,
	}
}

func (s *service) CreateAwardCategory(ctx context.Context, req *requests.CreateAwardCategoryRequest) (*responses.AwardCategoryResponse, error) {

	//Validate form structure
	if err := s.validateFormStructure(req.FormStructure); err != nil {
		return nil, err
	}

	// Check if award category with same name already exists
	existing, err := s.awardRepo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, exceptions.ErrAwardCategoryAlreadyExists
	}

	award := &models.AwardCategory{
		Name:          req.Name,
		Description:   req.Description,
		FormStructure: req.FormStructure,
		IsActive:      true,
	}

	if req.IsActive != nil {
		award.IsActive = *req.IsActive
	}

	if err := s.awardRepo.Create(ctx, award); err != nil {
		return nil, err
	}

	return s.toResponse(award), nil
}

func (s *service) GetAwardCategoryByID(ctx context.Context, id int) (*responses.AwardCategoryResponse, error) {
	award, err := s.awardRepo.FindByID(ctx, id)
	if err != nil {
		return nil, awardRepo.HandleError(err)
	}

	return s.toResponse(award), nil
}

func (s *service) GetAllAwardCategories(ctx context.Context) ([]*responses.AwardCategoryResponse, error) {
	awards, err := s.awardRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*responses.AwardCategoryResponse, len(awards))
	for i, award := range awards {
		responses[i] = s.toResponse(award)
	}

	return responses, nil
}

func (s *service) UpdateAwardCategory(ctx context.Context, id int, req *requests.UpdateAwardCategoryRequest) (*responses.AwardCategoryResponse, error) {
	award, err := s.awardRepo.FindByID(ctx, id)
	if err != nil {
		return nil, awardRepo.HandleError(err)
	}

	if req.Name != nil {
		// Check name uniqueness
		existing, err := s.awardRepo.FindByName(ctx, *req.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, exceptions.ErrAwardCategoryAlreadyExists
		}
		award.Name = *req.Name
	}

	if req.Description != nil {
		award.Description = *req.Description
	}

	if req.FormStructure != nil {
		if err := s.validateFormStructure(*req.FormStructure); err != nil {
			return nil, err
		}
		award.FormStructure = *req.FormStructure
	}

	if req.IsActive != nil {
		award.IsActive = *req.IsActive
	}

	if err := s.awardRepo.Update(ctx, award); err != nil {
		return nil, err
	}

	return s.toResponse(award), nil
}

func (s *service) DeleteAwardCategory(ctx context.Context, id int) error {
	_, err := s.awardRepo.FindByID(ctx, id)
	if err != nil {
		return awardRepo.HandleError(err)
	}

	// GORM soft delete - automatically sets deleted_at
	return s.awardRepo.Delete(ctx, id)
}

// form validation method
func (s *service) validateFormStructure(fs models.FormStructure) error {
	if len(fs) == 0 {
		return errors.New("form structure must have at least one field")
	}

	fieldIDs := make(map[string]bool)
	validTypes := map[string]bool{
		"text": true, "number": true, "date": true, "textarea": true,
		"select": true, "radio": true, "checkbox": true, "email": true, "url": true,
	}

	for _, field := range fs {
		// Check duplicate IDs
		if fieldIDs[field.ID] {
			return errors.New("duplicate field ID: " + field.ID)
		}
		fieldIDs[field.ID] = true

		// Validate field type
		if !validTypes[field.Type] {
			return errors.New("invalid field type: " + field.Type)
		}

		// Validate select/radio have options
		if (field.Type == "select" || field.Type == "radio") && len(field.Options) == 0 {
			return errors.New("field " + field.ID + " requires options")
		}
	}

	return nil
}

func (s *service) toResponse(award *models.AwardCategory) *responses.AwardCategoryResponse {
	return &responses.AwardCategoryResponse{
		ID:            award.ID,
		Name:          award.Name,
		Description:   award.Description,
		FormStructure: award.FormStructure,
		IsActive:      award.IsActive,
		CreatedAt:     award.CreatedAt,
		UpdatedAt:     award.UpdatedAt,
	}
}

// GetHonorRoll retrieves approved award winners grouped by academic year and award category
func (s *service) GetHonorRoll(ctx context.Context, query *requests.HonorRollQueryDTO) (*responses.HonorRollResponse, error) {
	// Prepare filter parameters
	var academicYear *int
	if query.AcademicYear > 0 {
		academicYear = &query.AcademicYear
	}

	var semester *string
	if query.Semester != "" {
		semester = &query.Semester
	}

	var campusID *uuid.UUID
	if query.CampusID != uuid.Nil {
		campusID = &query.CampusID
	}

	var awardID *int
	if query.AwardID > 0 {
		awardID = &query.AwardID
	}

	// Fetch approved winners from repository
	requests, err := s.requestRepo.FindApprovedWinners(ctx, academicYear, semester, campusID, awardID)
	if err != nil {
		return nil, err
	}

	// Group by academic year and award category
	yearMap := make(map[string]*responses.AcademicYearGroupResponse)
	totalWinners := 0

	for _, req := range requests {
		// Skip if required data is missing
		if req.Student.UserID == uuid.Nil || req.Award.ID == 0 || req.AcademicTerm.ID == 0 {
			continue
		}

		totalWinners++

		// Create year key
		yearKey := fmt.Sprintf("%d-%s", req.AcademicTerm.AcademicYear, req.AcademicTerm.Semester)

		// Initialize year group if not exists
		if _, exists := yearMap[yearKey]; !exists {
			yearMap[yearKey] = &responses.AcademicYearGroupResponse{
				AcademicYear: req.AcademicTerm.AcademicYear,
				Semester:     req.AcademicTerm.Semester,
				Categories:   []responses.AwardCategoryGroupResponse{},
			}
		}

		yearGroup := yearMap[yearKey]

		// Find or create award category group
		var categoryGroup *responses.AwardCategoryGroupResponse
		for i := range yearGroup.Categories {
			if yearGroup.Categories[i].AwardID == req.Award.ID {
				categoryGroup = &yearGroup.Categories[i]
				break
			}
		}

		if categoryGroup == nil {
			yearGroup.Categories = append(yearGroup.Categories, responses.AwardCategoryGroupResponse{
				AwardID:   req.Award.ID,
				AwardName: req.Award.Name,
				Winners:   []responses.WinnerResponse{},
			})
			categoryGroup = &yearGroup.Categories[len(yearGroup.Categories)-1]
		}

		// Create winner response
		studentName := req.Student.FirstName + " " + req.Student.LastName
		nisitID := ""
		if req.Student.NisitID != nil {
			nisitID = *req.Student.NisitID
		}

		campusName := ""
		if req.Student.Campus != nil {
			campusName = req.Student.Campus.Name
		}

		facultyName := ""
		if req.Student.Faculty != nil {
			facultyName = req.Student.Faculty.Name
		}

		departmentName := ""
		if req.Student.Department != nil {
			departmentName = req.Student.Department.Name
		}

		winner := responses.WinnerResponse{
			StudentID:      req.Student.UserID,
			StudentName:    studentName,
			NisitID:        nisitID,
			FacultyName:    facultyName,
			DepartmentName: departmentName,
			CampusName:     campusName,
			GPA:            req.SnapshotGPA,
		}

		categoryGroup.Winners = append(categoryGroup.Winners, winner)
	}

	// Convert map to sorted slice
	var yearGroups []responses.AcademicYearGroupResponse
	for _, group := range yearMap {
		yearGroups = append(yearGroups, *group)
	}

	// Sort year groups by academic year descending (newest first)
	// Note: Winners are already sorted by GPA DESC from repository query
	for i := 0; i < len(yearGroups)-1; i++ {
		for j := i + 1; j < len(yearGroups); j++ {
			if yearGroups[i].AcademicYear < yearGroups[j].AcademicYear {
				yearGroups[i], yearGroups[j] = yearGroups[j], yearGroups[i]
			}
		}
	}

	return &responses.HonorRollResponse{
		Data:         yearGroups,
		TotalWinners: totalWinners,
	}, nil
}
