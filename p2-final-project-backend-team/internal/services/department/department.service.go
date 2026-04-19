package department

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/department"
)


type Service interface {
	ValidateID(ctx context.Context, idStr string) (*uuid.UUID, error)
	Create(ctx context.Context, name string, facultyID uuid.UUID) error
	GetByFacultyID(ctx context.Context, facultyID uuid.UUID) ([]*models.Department, error)
	ValidateBelongsToFaculty(ctx context.Context, departmentID uuid.UUID, facultyID uuid.UUID) error
}

type service struct {
	repo department.Repository
}

func NewDepartmentService(repo department.Repository) Service {
	return &service{
		repo: repo,
	}
}

// ValidateID validates and parses department ID, then checks if it exists
func (s *service) ValidateID(ctx context.Context, idStr string) (*uuid.UUID, error) {
	if idStr == "" {
		return nil, nil
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("departmentID", idStr).Msg("Invalid department ID format")
		return nil, exceptions.ErrInvalidDepartmentID
	}

	record, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("departmentID", id.String()).Msg("Failed to find department")
		return nil, fmt.Errorf("failed to find department: %w", err)
	}

	if record == nil {
		log.Warn().Str("departmentID", id.String()).Msg("Department not found")
		return nil, exceptions.ErrDepartmentNotFound
	}

	return &id, nil
}

// Create creates a new department
func (s *service) Create(ctx context.Context, name string, facultyID uuid.UUID) error {
	department := &models.Department{
		Name: name,
		FacultyID: facultyID,
	}

	if err := s.repo.Create(ctx, department); err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to create department in database")
		return fmt.Errorf("failed to create department: %w", err)
	}

	return nil
}

func (s *service) GetByFacultyID(ctx context.Context, facultyID uuid.UUID) ([]*models.Department, error) {
	departments, err := s.repo.FindByFacultyID(ctx, facultyID)
	if err != nil {
		log.Error().Err(err).Str("facultyID", facultyID.String()).Msg("Failed to get departments by faculty")
		return nil, fmt.Errorf("failed to get departments by faculty: %w", err)
	}
	return departments, nil
}

func (s *service) ValidateBelongsToFaculty(ctx context.Context, departmentID uuid.UUID, facultyID uuid.UUID) error {
	record, err := s.repo.FindByID(ctx, departmentID)
	if err != nil {
		log.Error().Err(err).Str("departmentID", departmentID.String()).Msg("Failed to validate department faculty relationship")
		return fmt.Errorf("failed to validate department faculty relationship: %w", err)
	}
	if record == nil {
		log.Warn().Str("departmentID", departmentID.String()).Msg("Department not found when validating faculty")
		return exceptions.ErrDepartmentNotFound
	}
	if record.FacultyID != facultyID {
		log.Warn().Str("departmentID", departmentID.String()).Str("facultyID", facultyID.String()).Msg("Department does not belong to faculty")
		return exceptions.ErrDepartmentNotInFaculty
	}
	return nil
}