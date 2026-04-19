package responses

import (
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/google/uuid"
)

type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	User        UserInfo `json:"user"`
}

type UserInfo struct {
	UserID     uuid.UUID      `json:"user_id"`
	Email      string         `json:"email"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	PhoneNumber *string        `json:"phone_number,omitempty"`
	ProfileURL *string        `json:"profile_url,omitempty"`
	Role       enums.UserRole `json:"role"`
}

type AuthURLResponse struct {
	AuthURL string `json:"auth_url"`
}

type UserNotFoundResponse struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	RegisterToken string `json:"register_token"`
	Email         string `json:"email"`
}
