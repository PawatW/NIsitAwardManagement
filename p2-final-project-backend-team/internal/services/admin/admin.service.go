package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/requests"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/user"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/campus"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/department"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/faculty"
)

type Service interface {
	CreateUser(ctx context.Context, req requests.CreateStaffRequest) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUserStatus(ctx context.Context, userID string, req requests.UpdateUserStatusRequest) (*models.User, error)
	CreateCampus(ctx context.Context, name string) error
	CreateFaculty(ctx context.Context, name string, campusID string) error
	CreateDepartment(ctx context.Context, name string, facultyID string) error
}

type service struct {
	userRepo          user.Repository
	campusService     campus.Service
	facultyService    faculty.Service
	departmentService department.Service
}

func NewAdminService(userRepo user.Repository, campusService campus.Service, facultyService faculty.Service, departmentService department.Service) Service {
	return &service{
		userRepo:          userRepo,
		campusService:     campusService,
		facultyService:    facultyService,
		departmentService: departmentService,
	}
}

func (s *service) CreateUser(ctx context.Context, req requests.CreateStaffRequest) error {
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check email existence")
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		log.Warn().Msg("User creation attempt with existing email")
		return exceptions.ErrEmailAlreadyExists
	}

	ptrCampusID, err := s.campusService.ValidateID(ctx, req.CampusID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to validate campus ID")
		return fmt.Errorf("failed to validate campus ID: %w", err)
	}

	ptrFacultyID, err := s.facultyService.ValidateID(ctx, req.FacultyID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to validate faculty ID")
		return fmt.Errorf("failed to validate faculty ID: %w", err)
	}

	ptrDepartmentID, err := s.departmentService.ValidateID(ctx, req.DepartmentID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to validate department ID")
		return fmt.Errorf("failed to validate department ID: %w", err)
	}
	if ptrCampusID != nil && ptrFacultyID != nil {
		if err := s.facultyService.ValidateBelongsToCampus(ctx, *ptrFacultyID, *ptrCampusID); err != nil {
			log.Error().Err(err).Msg("Faculty does not belong to campus")
			return err
		}
	}
	if ptrFacultyID != nil && ptrDepartmentID != nil {
		if err := s.departmentService.ValidateBelongsToFaculty(ctx, *ptrDepartmentID, *ptrFacultyID); err != nil {
			log.Error().Err(err).Msg("Department does not belong to faculty")
			return err
		}
	}

	var phoneNumber *string
	if req.PhoneNumber != "" {
		phoneNumber = &req.PhoneNumber
	}

	u := &models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PhoneNumber:  phoneNumber,
		Role:         req.Role,
		NisitID:      nil,
		CampusID:     ptrCampusID,
		FacultyID:    ptrFacultyID,
		DepartmentID: ptrDepartmentID,
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		log.Error().Err(err).Msg("Failed to create user record")
		return fmt.Errorf("failed to create user: %w", err)
	}

	log.Info().Str("role", string(req.Role)).Msg("User created successfully by admin")
	return nil
}

func (s *service) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	users, err := s.userRepo.GetAllWithRelations(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all users")
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

func (s *service) UpdateUserStatus(ctx context.Context, userID string, req requests.UpdateUserStatusRequest) (*models.User, error) {
	// Parse and validate user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Error().Err(err).Msgf("Invalid user ID format: %s", userID)
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Validate that IsActive is provided
	if req.IsActive == nil {
		log.Warn().Msg("IsActive field is required but not provided")
		return nil, fmt.Errorf("is_active field is required")
	}

	// Check if user exists
	existingUser, err := s.userRepo.FindByIDWithRelations(ctx, userUUID)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to find user with ID %s", userID)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if existingUser == nil {
		log.Warn().Msgf("User not found with ID %s", userID)
		return nil, exceptions.ErrUserNotFound
	}

	// Update user status
	if err := s.userRepo.UpdateUserStatus(ctx, userUUID, *req.IsActive); err != nil {
		log.Error().Err(err).Msgf("Failed to update user status for ID %s", userID)
		return nil, fmt.Errorf("failed to update user status: %w", err)
	}

	// Fetch updated user with relations
	updatedUser, err := s.userRepo.FindByIDWithRelations(ctx, userUUID)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to fetch updated user with ID %s", userID)
		return nil, fmt.Errorf("failed to fetch updated user: %w", err)
	}

	log.Info().Str("userID", userID).Bool("isActive", *req.IsActive).Msg("User status updated successfully")
	return updatedUser, nil
}

func (s *service) CreateCampus(ctx context.Context, name string) error {
	if err := s.campusService.Create(ctx, name); err != nil {
		log.Error().Err(err).Msg("Failed to create campus")
		return fmt.Errorf("failed to create campus: %w", err)
	}
	log.Info().Str("name", name).Msg("Campus created successfully")
	return nil
}

func (s *service) CreateFaculty(ctx context.Context, name string, campusID string) error {
	ptrCampusID, err := s.campusService.ValidateID(ctx, campusID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to validate campus ID for faculty")
		return fmt.Errorf("failed to validate campus ID: %w", err)
	}
	if ptrCampusID == nil {
		return exceptions.ErrCampusNotFound
	}
	if err := s.facultyService.Create(ctx, name, *ptrCampusID); err != nil {
		log.Error().Err(err).Msg("Failed to create faculty")
		return fmt.Errorf("failed to create faculty: %w", err)
	}
	log.Info().Str("name", name).Msg("Faculty created successfully")
	return nil
}

func (s *service) CreateDepartment(ctx context.Context, name string, facultyID string) error {
	ptrFacultyID, err := s.facultyService.ValidateID(ctx, facultyID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to validate faculty ID for department")
		return fmt.Errorf("failed to validate faculty ID: %w", err)
	}
	if ptrFacultyID == nil {
		return exceptions.ErrFacultyNotFound
	}
	if err := s.departmentService.Create(ctx, name, *ptrFacultyID); err != nil {
		log.Error().Err(err).Msg("Failed to create department")
		return fmt.Errorf("failed to create department: %w", err)
	}
	log.Info().Str("name", name).Msg("Department created successfully")
	return nil
}
