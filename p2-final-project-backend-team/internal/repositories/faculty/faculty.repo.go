package faculty

import (
	"context"
	"errors"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	FindByID(ctx context.Context, facultyID uuid.UUID) (*models.Faculty, error)
	FindByCampusID(ctx context.Context, campusID uuid.UUID) ([]*models.Faculty, error)
	Create(ctx context.Context, faculty *models.Faculty) error
}

type repository struct {
	db *gorm.DB
}

func NewFacultyRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}


func (r *repository) FindByID(ctx context.Context, facultyID uuid.UUID) (*models.Faculty, error) {
	var f models.Faculty
	if err := r.db.WithContext(ctx).Where("id = ?", facultyID).First(&f).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msgf("faculty with ID %s not found", facultyID)
			return nil, nil
		}
		log.Error().Err(err).Msgf("Error to find faculty with ID %s", facultyID)
		return nil, err
	}
	return &f, nil
}

func (r *repository) FindByCampusID(ctx context.Context, campusID uuid.UUID) ([]*models.Faculty, error) {
	var faculties []*models.Faculty
	if err := r.db.WithContext(ctx).Where("campus_id = ?", campusID).Find(&faculties).Error; err != nil {
		log.Error().Err(err).Msgf("Error to find faculties with campus ID %s", campusID)
		return nil, err
	}
	return faculties, nil
}

func (r *repository) Create(ctx context.Context, faculty *models.Faculty) error {
	faculty.ID = uuid.New()
	if err := r.db.WithContext(ctx).Create(faculty).Error; err != nil {
		log.Error().Err(err).Msg("failed to create faculty")
		return err
	}
	log.Info().Msgf("faculty with ID %s created successfully", faculty.ID)
	return nil
}