package dev

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/auth"
	userRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/user"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GenerateTokenRequest struct {
	UserID string         `json:"user_id"`
	Email  string         `json:"email" validate:"required,email"`
	Role   enums.UserRole `json:"role" validate:"required"`

	CreateUser bool   `json:"create_user"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	NisitID    string `json:"nisit_id"`

	CampusID     string `json:"campus_id"`
	FacultyID    string `json:"faculty_id"`
	DepartmentID string `json:"department_id"`

	ExpireHours int `json:"expire_hours"`
}

type GenerateTokenResponse struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}

type DevHandler interface {
	GenerateToken(c echo.Context) error
}

type handler struct {
	config       *configs.Config
	tokenManager auth.TokenManager
	userRepo     userRepo.Repository
}

func NewDevHandler(config *configs.Config, tokenManager auth.TokenManager, userRepo userRepo.Repository) DevHandler {
	return &handler{config: config, tokenManager: tokenManager, userRepo: userRepo}
}

func (h *handler) GenerateToken(c echo.Context) error {
	if h.config.AppEnv == "production" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "This endpoint is only available in development mode"})
	}

	var req GenerateTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	req.Email = normalizeEmail(req.Email)
	if !isValidRole(req.Role) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid role. Valid roles: ADMIN, STUDENT, HEAD_OF_DEPARTMENT, VICE_DEAN, DEAN, COMMITTEE_CHAIR",
		})
	}

	userID := uuid.Nil
	if req.CreateUser {
		u, err := h.upsertUserLikeRealSystem(c.Request().Context(), &req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create/update user", "details": err.Error()})
		}
		userID = u.UserID
	} else {
		if req.UserID == "" {
			userID = uuid.New()
		} else {
			parsed, err := uuid.Parse(req.UserID)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id format"})
			}
			userID = parsed
		}
	}

	token, err := h.tokenManager.GenerateAccessToken(userID, req.Email, req.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	expireHours := h.config.JWT.ExpiryHours
	if req.ExpireHours > 0 {
		expireHours = req.ExpireHours
	}
	expiresAt := time.Now().Add(time.Hour * time.Duration(expireHours))

	return c.JSON(http.StatusOK, GenerateTokenResponse{
		Token:     token,
		UserID:    userID,
		Email:     req.Email,
		Role:      string(req.Role),
		ExpiresAt: expiresAt,
	})
}

func (h *handler) upsertUserLikeRealSystem(ctx context.Context, req *GenerateTokenRequest) (*modelsUser, error) {
	first := strings.TrimSpace(req.FirstName)
	last := strings.TrimSpace(req.LastName)
	if first == "" {
		first = "Test"
	}
	if last == "" {
		last = string(req.Role)
	}

	u, err := h.userRepo.FindOrCreateByEmail(ctx, req.Email, first, last, "", "test")
	if err != nil {
		return nil, err
	}

	// ✅ role เป็น enum ตรง model
	u.Role = req.Role

	if strings.TrimSpace(req.NisitID) != "" {
		nisit := strings.TrimSpace(req.NisitID)
		u.NisitID = &nisit
	}

	if req.CampusID != "" {
		if v, err := uuid.Parse(req.CampusID); err == nil {
			u.CampusID = &v
		}
	}
	if req.FacultyID != "" {
		if v, err := uuid.Parse(req.FacultyID); err == nil {
			u.FacultyID = &v
		}
	}
	if req.DepartmentID != "" {
		if v, err := uuid.Parse(req.DepartmentID); err == nil {
			u.DepartmentID = &v
		}
	}

	if err := h.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	return &modelsUser{UserID: u.UserID}, nil
}

type modelsUser struct {
	UserID uuid.UUID
}

func normalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func isValidRole(role enums.UserRole) bool {
	switch role {
	case enums.Admin,
		enums.Student,
		enums.HeadOfDepartment,
		enums.ViceDean,
		enums.Dean,
		enums.CommitteeChair:
		return true
	default:
		return false
	}
}