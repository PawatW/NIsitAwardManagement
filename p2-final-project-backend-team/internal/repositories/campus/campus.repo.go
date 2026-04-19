package campus

import (
	"context"
	"errors"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)	

type Repository interface {
	FindByID(ctx context.Context, campusID uuid.UUID) (*models.Campus, error)
	GetAll(ctx context.Context) ([]*models.Campus, error)
	Create(ctx context.Context, campus *models.Campus) error
}

type repository struct {
	db *gorm.DB
}

func NewCampusRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByID(ctx context.Context, campusID uuid.UUID) (*models.Campus, error) {
	var c models.Campus
	if err := r.db.WithContext(ctx).Where("id = ?", campusID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msgf("campus with ID %s not found", campusID)
			return nil, nil
		}		
		log.Error().Err(err).Msgf("Error to find campus with ID %s", campusID)
		return nil, err

	}
	return &c, nil
}

func (r *repository) GetAll(ctx context.Context) ([]*models.Campus, error) {
    	var campuses []*models.Campus
	if err := r.db.WithContext(ctx).Find(&campuses).Error; err != nil {
		log.Error().Err(err).Msg("failed to get all campuses")
		return nil, err
	}
	return campuses, nil
}


func (r *repository) Create(ctx context.Context, campus *models.Campus) error {
	campus.ID = uuid.New()
	if err := r.db.WithContext(ctx).Create(campus).Error; err != nil {
		log.Error().Err(err).Msg("failed to create campus")
		return err
	}
	log.Info().Msgf("campus with ID %s created successfully", campus.ID)
	return nil
}