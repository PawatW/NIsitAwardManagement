package campus


import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/campus"
)

type Service interface {
	ValidateID(ctx context.Context, idStr string) (*uuid.UUID, error)
	Create(ctx context.Context, name string) error
	GetAll(ctx context.Context) ([]*models.Campus, error)
}

type service struct {
	repo campus.Repository
}

func NewCampusService(repo campus.Repository) Service {
	return &service{
		repo: repo,
	}
}

// ValidateID validates and parses campus ID, then checks if it exists
func (s *service) ValidateID(ctx context.Context, idStr string) (*uuid.UUID, error) {
	if idStr == "" {
		return nil, nil
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("campusID", idStr).Msg("Invalid campus ID format")
		return nil, exceptions.ErrInvalidCampusID
	}

	record, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("campusID", id.String()).Msg("Failed to find campus")
		return nil, fmt.Errorf("failed to find campus: %w", err)
	}

	if record == nil {
		log.Warn().Str("campusID", id.String()).Msg("Campus not found")
		return nil, exceptions.ErrCampusNotFound
	}

	return &id, nil
}

// Create creates a new campus
func (s *service) Create(ctx context.Context, name string) error {
	campus := &models.Campus{
		Name: name,
	}

	if err := s.repo.Create(ctx, campus); err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to create campus in database")
		return fmt.Errorf("failed to create campus: %w", err)
	}

	return nil
}

func (s *service) GetAll(ctx context.Context) ([]*models.Campus, error) {
	campuses, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all campuses")
		return nil, fmt.Errorf("failed to get all campuses: %w", err)
	}
	return campuses, nil
}
