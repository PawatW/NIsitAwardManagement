package award

import (
	"context"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, award *models.AwardCategory) error
	FindByID(ctx context.Context, id int) (*models.AwardCategory, error)
	FindAll(ctx context.Context) ([]*models.AwardCategory, error)
	FindAllActive(ctx context.Context) ([]*models.AwardCategory, error)
	Update(ctx context.Context, award *models.AwardCategory) error
	Delete(ctx context.Context, id int) error
	FindByName(ctx context.Context, name string) (*models.AwardCategory, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, award *models.AwardCategory) error {
	return r.db.WithContext(ctx).Create(award).Error
}

func (r *repository) FindByID(ctx context.Context, id int) (*models.AwardCategory, error) {
	var award models.AwardCategory
	err := r.db.WithContext(ctx).First(&award, id).Error
	if err != nil {
		return nil, err
	}
	return &award, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*models.AwardCategory, error) {
	var awards []*models.AwardCategory
	err := r.db.WithContext(ctx).Find(&awards).Error
	if err != nil {
		return nil, err
	}
	return awards, nil
}

func (r *repository) Update(ctx context.Context, award *models.AwardCategory) error {
	return r.db.WithContext(ctx).Save(award).Error
}

func (r *repository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.AwardCategory{}, id).Error
}

func (r *repository) FindByName(ctx context.Context, name string) (*models.AwardCategory, error) {
	var award models.AwardCategory
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&award).Error
	if err != nil {
		return nil, err
	}
	return &award, nil
}

func (r *repository) FindAllActive(ctx context.Context) ([]*models.AwardCategory, error) {
	var awards []*models.AwardCategory
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&awards).Error
	if err != nil {
		return nil, err
	}
	return awards, nil
}
