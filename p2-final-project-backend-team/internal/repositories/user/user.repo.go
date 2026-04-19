package user

import (
	"context"
	"errors"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindByIDWithRelations(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindByNisitID(ctx context.Context, nisitID string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindOrCreateByEmail(ctx context.Context, email, firstName, lastName, profileURL, authProvider string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	GetAllWithRelations(ctx context.Context) ([]*models.User, error)

	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	UpdateUserStatus(ctx context.Context, userID uuid.UUID, isActive bool) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msgf("user with ID %s not found", userID)
			return nil, nil
		}
		log.Error().Err(err).Msgf("Error to find user with ID %s", userID)
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByIDWithRelations(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).
		Preload("Campus").
		Preload("Faculty").
		Preload("Department").
		Where("user_id = ?", userID).
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msgf("user with ID %s not found", userID)
			return nil, nil
		}
		log.Error().Err(err).Msgf("Error to find user with relations ID %s", userID)
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByNisitID(ctx context.Context, nisitID string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).Where("nisit_id = ?", nisitID).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msgf("user with NisitID %s not found", nisitID)
			return nil, nil
		}
		log.Error().Err(err).Msgf("Error to find user with NisitID %s", nisitID)
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msgf("user with email %s not found", email)
			return nil, nil
		}
		log.Error().Err(err).Msgf("Error to find user with email %s", email)
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		log.Error().Err(err).Msg("failed to get all users")
		return nil, err
	}
	return users, nil
}

func (r *repository) GetAllWithRelations(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.WithContext(ctx).
		Preload("Campus").
		Preload("Faculty").
		Preload("Department").
		Find(&users).Error; err != nil {
		log.Error().Err(err).Msg("failed to get all users with relations")
		return nil, err
	}
	return users, nil
}

func (r *repository) Create(ctx context.Context, user *models.User) error {
	user.UserID = uuid.New()
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return err
	}
	log.Info().Msgf("user with ID %s created successfully", user.UserID)
	return nil
}

func (r *repository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("user_id = ?", user.UserID).
		Updates(user).Error; err != nil {
		log.Error().Err(err).Msgf("failed to update user with ID %s", user.UserID)
		return err
	}
	log.Info().Msgf("user with ID %s updated successfully", user.UserID)
	return nil
}

func (r *repository) UpdateUserStatus(ctx context.Context, userID uuid.UUID, isActive bool) error {
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("is_active", isActive).Error; err != nil {
		log.Error().Err(err).Msgf("failed to update user status for ID %s", userID)
		return err
	}
	log.Info().Msgf("user status updated successfully for ID %s to %v", userID, isActive)
	return nil
}

func (r *repository) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.User{}).Error; err != nil {
		log.Error().Err(err).Msgf("failed to delete user with ID %s", userID)
		return err
	}
	log.Info().Msgf("user with ID %s deleted successfully", userID)
	return nil
}

// FindOrCreateByEmail finds a user by email or creates a new one if not found
func (r *repository) FindOrCreateByEmail(ctx context.Context, email, firstName, lastName, profileURL, authProvider string) (*models.User, error) {
	existingUser, err := r.FindByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Msgf("error finding user by email %s", email)
		return nil, err
	}

	// If user exists, update profile fields
	if existingUser != nil {
		existingUser.FirstName = firstName
		existingUser.LastName = lastName
		existingUser.AuthProvider = authProvider
		if profileURL != "" {
			existingUser.ProfileURL = &profileURL
		}

		if err := r.Update(ctx, existingUser); err != nil {
			log.Error().Err(err).Msgf("failed to update existing user %s", email)
			return nil, err
		}
		return existingUser, nil
	}

	// Create new user
	newUser := &models.User{
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		AuthProvider: authProvider,
		Role:         enums.Student, // Default role for new users is Student but will be updated later in the handler
	}
	if profileURL != "" {
		newUser.ProfileURL = &profileURL
	}

	if err := r.Create(ctx, newUser); err != nil {
		log.Error().Err(err).Msgf("failed to create new user %s", email)
		return nil, err
	}
	log.Info().Msgf("new user with data %+v created successfully", newUser)

	return newUser, nil
}
