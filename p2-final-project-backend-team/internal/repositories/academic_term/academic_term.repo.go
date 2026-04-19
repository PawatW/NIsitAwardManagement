package academic_term

import (
	"context"
	"errors"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context, isOpenFilter *bool) ([]*models.AcademicTerm, error)
	FindByID(ctx context.Context, id int) (*models.AcademicTerm, error)
	FindLatest(ctx context.Context, isOpenFilter *bool) (*models.AcademicTerm, error)
	Create(ctx context.Context, term *models.AcademicTerm) error
	FindByYearAndSemester(ctx context.Context, year int, semester string) (*models.AcademicTerm, error) // +
	FindOpenTerm(ctx context.Context) (*models.AcademicTerm, error)
	Update(ctx context.Context, term *models.AcademicTerm) error
}

type repository struct {
	db *gorm.DB
}

func NewAcademicTermRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context, isOpenFilter *bool) ([]*models.AcademicTerm, error) {
	var terms []*models.AcademicTerm
	query := r.db.WithContext(ctx)

	// Apply filter if provided
	if isOpenFilter != nil {
		query = query.Where("is_open = ?", *isOpenFilter)
	}

	if err := query.Order("academic_year DESC, semester DESC").Find(&terms).Error; err != nil {
		log.Error().Err(err).Msg("failed to get all academic terms")
		return nil, err
	}
	return terms, nil
}

func (r *repository) FindByID(ctx context.Context, id int) (*models.AcademicTerm, error) {
	var term models.AcademicTerm
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&term).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msgf("academic term with ID %d not found", id)
			return nil, nil
		}
		log.Error().Err(err).Msgf("failed to find academic term with ID %d", id)
		return nil, err
	}
	return &term, nil
}

func (r *repository) Create(ctx context.Context, term *models.AcademicTerm) error {
	if err := r.db.WithContext(ctx).Create(term).Error; err != nil {
		log.Error().Err(err).Msg("failed to create academic term")
		return err
	}
	log.Info().Msgf("academic term with ID %d created successfully", term.ID)
	return nil
}

func (r *repository) FindByYearAndSemester(ctx context.Context, year int, semester string) (*models.AcademicTerm, error) {
	var term models.AcademicTerm
	if err := r.db.WithContext(ctx).Where("academic_year = ? AND semester = ?", year, semester).First(&term).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(err).Msg("failed to find academic term by year and semester")
		return nil, err
	}
	return &term, nil
}

func (r *repository) FindLatest(ctx context.Context, isOpenFilter *bool) (*models.AcademicTerm, error) {
	var term models.AcademicTerm
	query := r.db.WithContext(ctx)

	// Apply filter if provided
	if isOpenFilter != nil {
		query = query.Where("is_open = ?", *isOpenFilter)
	}

	if err := query.Order("academic_year DESC, semester DESC").First(&term).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("no academic term found")
			return nil, nil
		}
		log.Error().Err(err).Msg("failed to find latest academic term")
		return nil, err
	}
	return &term, nil
}

func (r *repository) FindOpenTerm(ctx context.Context) (*models.AcademicTerm, error) {
	var term models.AcademicTerm
	if err := r.db.WithContext(ctx).Where("is_open = true").First(&term).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(err).Msg("failed to find open academic term")
		return nil, err
	}
	return &term, nil
}

func (r *repository) Update(ctx context.Context, term *models.AcademicTerm) error {
	if err := r.db.WithContext(ctx).Save(term).Error; err != nil {
		log.Error().Err(err).Msgf("failed to update academic term with ID %d", term.ID)
		return err
	}
	log.Info().Msgf("academic term with ID %d updated successfully", term.ID)
	return nil
}
