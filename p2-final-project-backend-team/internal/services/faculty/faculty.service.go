package faculty

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/faculty"
)

type Service interface {
	ValidateID(ctx context.Context, idStr string) (*uuid.UUID, error)
	Create(ctx context.Context, name string, campusID uuid.UUID) error
	GetByCampusID(ctx context.Context, campusID uuid.UUID) ([]*models.Faculty, error)
	ValidateBelongsToCampus(ctx context.Context, facultyID uuid.UUID, campusID uuid.UUID) error
}

type service struct {
	repo faculty.Repository
}

func NewFacultyService(repo faculty.Repository) Service {
	return &service{
		repo: repo,
	}
}


// ValidateID validates and parses faculty ID, then checks if it exists
func (s *service) ValidateID(ctx context.Context, idStr string) (*uuid.UUID, error) {
	if idStr == "" {
		return nil, nil
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("facultyID", idStr).Msg("Invalid faculty ID format")
		return nil, exceptions.ErrInvalidFacultyID
	}

	record, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("facultyID", id.String()).Msg("Failed to find faculty")
		return nil, fmt.Errorf("failed to find faculty: %w", err)
	}

	if record == nil {
		log.Warn().Str("facultyID", id.String()).Msg("Faculty not found")
		return nil, exceptions.ErrFacultyNotFound
	}

	return &id, nil
}

// Create creates a new faculty
func (s *service) Create(ctx context.Context, name string, campusID uuid.UUID) error {
	faculty := &models.Faculty{
		Name: name,
		CampusID: campusID,
	}

	if err := s.repo.Create(ctx, faculty); err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to create faculty in database")
		return fmt.Errorf("failed to create faculty: %w", err)
	}

	return nil
}

func (s *service) GetByCampusID(ctx context.Context, campusID uuid.UUID) ([]*models.Faculty, error) {
	faculties, err := s.repo.FindByCampusID(ctx, campusID)
	if err != nil {
		log.Error().Err(err).Str("campusID", campusID.String()).Msg("Failed to get faculties by campus")
		return nil, fmt.Errorf("failed to get faculties by campus: %w", err)
	}
	return faculties, nil
}

func (s *service) ValidateBelongsToCampus(ctx context.Context, facultyID uuid.UUID, campusID uuid.UUID) error {
	record, err := s.repo.FindByID(ctx, facultyID)
	if err != nil {
		log.Error().Err(err).Str("facultyID", facultyID.String()).Msg("Failed to validate faculty campus relationship")
		return fmt.Errorf("failed to validate faculty campus relationship: %w", err)
	}
	if record == nil {
		log.Warn().Str("facultyID", facultyID.String()).Msg("Faculty not found when validating campus")
		return exceptions.ErrFacultyNotFound
	}
	if record.CampusID != campusID {
		log.Warn().Str("facultyID", facultyID.String()).Str("campusID", campusID.String()).Msg("Faculty does not belong to campus")
		return exceptions.ErrFacultyNotInCampus
	}
	return nil
}
