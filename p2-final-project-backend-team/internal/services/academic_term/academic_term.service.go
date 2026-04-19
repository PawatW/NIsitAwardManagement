package academic_term

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/requests"
	academic_term_repo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/academic_term"
)

type Service interface {
	Create(ctx context.Context, req requests.CreateAcademicTermRequest) error
	GetAll(ctx context.Context, isOpenFilter *bool) ([]*models.AcademicTerm, error)
	GetLatest(ctx context.Context, isOpenFilter *bool) (*models.AcademicTerm, error)
	FindByID(ctx context.Context, id int) (*models.AcademicTerm, error)
	UpdateIsOpen(ctx context.Context, id int, isOpen bool) error
}

type service struct {
	repo academic_term_repo.Repository
}

func NewAcademicTermService(repo academic_term_repo.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req requests.CreateAcademicTermRequest) error {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end_date format: %w", err)
	}

	if !endDate.After(startDate) {
		return fmt.Errorf("end_date must be after start_date")
	}

	// Check duplicate (academic_year + semester must be unique)
	existing, err := s.repo.FindByYearAndSemester(ctx, req.AcademicYear, req.Semester)
	if err != nil {
		return fmt.Errorf("failed to check duplicate academic term: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("academic term for %d %s semester already exists", req.AcademicYear, req.Semester)
	}

	// Check only one open term at a time
	if req.IsOpen {
		openTerm, err := s.repo.FindOpenTerm(ctx)
		if err != nil {
			return fmt.Errorf("failed to check open academic term: %w", err)
		}
		if openTerm != nil {
			return fmt.Errorf("there is already an open academic term (ID: %d, %d %s)", openTerm.ID, openTerm.AcademicYear, openTerm.Semester)
		}
	}

	term := &models.AcademicTerm{
		AcademicYear: req.AcademicYear,
		Semester:     req.Semester,
		StartDate:    startDate,
		EndDate:      endDate,
		IsOpen:       req.IsOpen,
	}

	if err := s.repo.Create(ctx, term); err != nil {
		log.Error().Err(err).Msg("failed to create academic term")
		return fmt.Errorf("failed to create academic term: %w", err)
	}

	return nil
}

func (s *service) GetAll(ctx context.Context, isOpenFilter *bool) ([]*models.AcademicTerm, error) {
	terms, err := s.repo.GetAll(ctx, isOpenFilter)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all academic terms")
		return nil, fmt.Errorf("failed to get all academic terms: %w", err)
	}
	return terms, nil
}

func (s *service) GetLatest(ctx context.Context, isOpenFilter *bool) (*models.AcademicTerm, error) {
	term, err := s.repo.FindLatest(ctx, isOpenFilter)
	if err != nil {
		log.Error().Err(err).Msg("failed to get latest academic term")
		return nil, fmt.Errorf("failed to get latest academic term: %w", err)
	}
	return term, nil
}

func (s *service) FindByID(ctx context.Context, id int) (*models.AcademicTerm, error) {
	term, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find academic term with ID %d", id)
		return nil, fmt.Errorf("failed to find academic term: %w", err)
	}
	return term, nil
}

func (s *service) UpdateIsOpen(ctx context.Context, id int, isOpen bool) error {
	term, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find academic term: %w", err)
	}
	if term == nil {
		return fmt.Errorf("academic term not found")
	}

	// if opening check no other term is already open
	if isOpen {
		openTerm, err := s.repo.FindOpenTerm(ctx)
		if err != nil {
			return fmt.Errorf("failed to check open academic term: %w", err)
		}
		if openTerm != nil && openTerm.ID != term.ID {
			return fmt.Errorf("there is already an open academic term (ID: %d, %d %s)", openTerm.ID, openTerm.AcademicYear, openTerm.Semester)
		}
	}

	term.IsOpen = isOpen
	if err := s.repo.Update(ctx, term); err != nil {
		return fmt.Errorf("failed to update academic term: %w", err)
	}
	return nil
}
