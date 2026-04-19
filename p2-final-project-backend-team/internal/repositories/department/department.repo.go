package department

import (
	"context"
	"errors"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	FindByID(ctx context.Context, departmentID uuid.UUID) (*models.Department, error)
	FindByFacultyID(ctx context.Context, facultyID uuid.UUID) ([]*models.Department, error)
	Create(ctx context.Context, department *models.Department) error
}

type repository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByID(ctx context.Context, departmentID uuid.UUID) (*models.Department, error) {
	var d models.Department
	if err := r.db.WithContext(ctx).Where("id = ?", departmentID).First(&d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msgf("department with ID %s not found", departmentID)
			return nil, nil
		}
		log.Error().Err(err).Msgf("Error to find department with ID %s", departmentID)
		return nil, err
	}
	return &d, nil
}

func (r *repository) FindByFacultyID(ctx context.Context, facultyID uuid.UUID) ([]*models.Department, error) {
	var departments []*models.Department
	if err := r.db.WithContext(ctx).Where("faculty_id = ?", facultyID).Find(&departments).Error; err != nil {
		log.Error().Err(err).Msgf("Error to find departments with faculty ID %s", facultyID)
		return nil, err
	}
	return departments, nil
}

func (r *repository) Create(ctx context.Context, department *models.Department) error {
	department.ID = uuid.New()
	if err := r.db.WithContext(ctx).Create(department).Error; err != nil {
		log.Error().Err(err).Msg("failed to create department")
		return err
	}
	log.Info().Msgf("department with ID %s created successfully", department.ID)
	return nil
}