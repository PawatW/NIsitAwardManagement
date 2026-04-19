package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/requests"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/responses"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/auth"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/user"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/campus"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/department"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/faculty"
)

type Service interface {
	Register(ctx context.Context, req requests.StudentRegisterRequest) error
	GetByID(ctx context.Context, userID uuid.UUID) (*responses.FullUserResponse, error)
}

type service struct {
	userRepo          user.Repository
	campusService     campus.Service
	facultyService    faculty.Service
	departmentService department.Service
	tokenManager      auth.TokenManager
}

func NewUserService(userRepo user.Repository, campusService campus.Service, facultyService faculty.Service, departmentService department.Service, tokenManager auth.TokenManager) Service {
	return &service{
		userRepo:          userRepo,
		campusService:     campusService,
		facultyService:    facultyService,
		departmentService: departmentService,
		tokenManager:      tokenManager,
	}
}

func (s *service) Register(ctx context.Context, req requests.StudentRegisterRequest) error {
	// Validate register token
	tokenClaims, err := s.tokenManager.ValidateRegisterToken(req.RegisterToken)
	if err != nil {
		log.Error().Err(err).Msg("Invalid register token")
		return fmt.Errorf("invalid or expired register token: %w", err)
	}

	// Get email from token
	email := tokenClaims.Email

	// Check if user already exists (shouldn't happen, but double-check)
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check email existence")
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		log.Warn().Msg("Registration attempt with existing email")
		return exceptions.ErrEmailAlreadyExists
	}

	// Check if nisit ID already exists
	existingStudent, err := s.userRepo.FindByNisitID(ctx, req.NisitID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check nisit ID existence")
		return fmt.Errorf("failed to check nisit ID: %w", err)
	}
	if existingStudent != nil {
		log.Warn().Msg("Registration attempt with existing nisit ID")
		return exceptions.ErrNisitIDAlreadyExists
	}

	// Validate organizational IDs (all required for students)
	ptrCampusID, err := s.campusService.ValidateID(ctx, req.CampusID)
	if err != nil {
		return fmt.Errorf("failed to validate campus ID: %w", err)
	}
	if ptrCampusID == nil {
		log.Error().Msg("Campus ID is required for student registration")
		return exceptions.ErrCampusNotFound
	}

	ptrFacultyID, err := s.facultyService.ValidateID(ctx, req.FacultyID)
	if err != nil {
		return fmt.Errorf("failed to validate faculty ID: %w", err)
	}
	if ptrFacultyID == nil {
		log.Error().Msg("Faculty ID is required for student registration")
		return exceptions.ErrFacultyNotFound
	}

	ptrDepartmentID, err := s.departmentService.ValidateID(ctx, req.DepartmentID)
	if err != nil {
		return fmt.Errorf("failed to validate department ID: %w", err)
	}
	if ptrDepartmentID == nil {
		log.Error().Msg("Department ID is required for student registration")
		return exceptions.ErrDepartmentNotFound
	}

	// Validate organizational relationships
	if err := s.facultyService.ValidateBelongsToCampus(ctx, *ptrFacultyID, *ptrCampusID); err != nil {
		log.Error().Err(err).Msg("Faculty does not belong to campus")
		return err
	}

	if err := s.departmentService.ValidateBelongsToFaculty(ctx, *ptrDepartmentID, *ptrFacultyID); err != nil {
		log.Error().Err(err).Msg("Department does not belong to faculty")
		return err
	}

	// Create new student user (role is always STUDENT for self-registration)
	// ProfileURL will be set on first login via Google OAuth
	u := &models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        email,
		AuthProvider: "google",
		PhoneNumber:  &req.PhoneNumber,
		ProfileURL:   nil, // Will be updated on login
		Role:         enums.Student,
		NisitID:      &req.NisitID,
		CampusID:     ptrCampusID,
		FacultyID:    ptrFacultyID,
		DepartmentID: ptrDepartmentID,
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		log.Error().Err(err).Msg("Failed to create student record")
		return fmt.Errorf("failed to create student: %w", err)
	}

	log.Info().Msgf("Student %s registered successfully. Please login to continue.", email)
	return nil
}

// GetByID retrieves user by ID with all related data from database
func (s *service) GetByID(ctx context.Context, userID uuid.UUID) (*responses.FullUserResponse, error) {
	user, err := s.userRepo.FindByIDWithRelations(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get user with ID %s", userID)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		log.Warn().Msgf("User with ID %s not found", userID)
		return nil, exceptions.ErrUserNotFound
	}

	return responses.ToFullUserResponse(user), nil
}
