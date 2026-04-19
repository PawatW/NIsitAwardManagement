package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/dto/responses"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/auth"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/oauth"
	userRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/user"
	"github.com/rs/zerolog/log"
)

type Service interface {
	GetGoogleAuthURL() (string, error)
	HandleGoogleCallback(ctx context.Context, code, state string) (*responses.LoginResponse, *responses.UserNotFoundResponse, error)
}

type service struct {
	oauthProvider oauth.Provider
	userRepo      userRepo.Repository
	tokenManager  auth.TokenManager
	stateStore    *stateStore
}

// stateStore manages OAuth state tokens with expiration
type stateStore struct {
	mu     sync.RWMutex
	states map[string]time.Time
}

// newStateStore creates a new state store
func newStateStore() *stateStore {
	store := &stateStore{
		states: make(map[string]time.Time),
	}
	// Start cleanup goroutine
	go store.cleanup()
	return store
}

// set stores a state with 5 minute expiration
func (s *stateStore) set(state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[state] = time.Now().Add(5 * time.Minute)
}

// validate checks if state exists and is not expired
func (s *stateStore) validate(state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	expiry, exists := s.states[state]
	if !exists {
		return false
	}
	if time.Now().After(expiry) {
		delete(s.states, state)
		return false
	}
	// Delete after validation (one-time use)
	delete(s.states, state)
	return true
}

// cleanup removes expired states every minute
func (s *stateStore) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for state, expiry := range s.states {
			if now.After(expiry) {
				delete(s.states, state)
			}
		}
		s.mu.Unlock()
	}
}

func NewAuthService(
	oauthProvider oauth.Provider,
	userRepo userRepo.Repository,
	tokenManager auth.TokenManager,
) Service {
	return &service{
		oauthProvider: oauthProvider,
		userRepo:      userRepo,
		tokenManager:  tokenManager,
		stateStore:    newStateStore(),
	}
}

// GetGoogleAuthURL generates and returns the Google OAuth authorization URL
func (s *service) GetGoogleAuthURL() (string, error) {
	// Generate a random state for CSRF protection
	state, err := generateState()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate state")
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	// Store state for validation
	s.stateStore.set(state)

	authURL := s.oauthProvider.GetAuthURL(state)
	log.Info().Str("state", state).Msg("generated Google OAuth URL with state")
	return authURL, nil
}

// HandleGoogleCallback handles the OAuth callback and returns login response or user not found response
func (s *service) HandleGoogleCallback(ctx context.Context, code, state string) (*responses.LoginResponse, *responses.UserNotFoundResponse, error) {
	// Validate state to prevent CSRF attacks
	if !s.stateStore.validate(state) {
		log.Error().Str("state", state).Msg("invalid or expired state")
		return nil, nil, fmt.Errorf("invalid or expired state token")
	}

	// Exchange authorization code for access token
	accessToken, err := s.oauthProvider.ExchangeCode(ctx, code)
	if err != nil {
		log.Error().Err(err).Msg("failed to exchange code")
		return nil, nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from Google
	userInfo, err := s.oauthProvider.GetUserInfo(ctx, accessToken)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user info")
		return nil, nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user exists in database
	existingUser, err := s.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		log.Error().Err(err).Msg("failed to find user")
		return nil, nil, fmt.Errorf("failed to find user: %w", err)
	}

	// If user not found, generate register token and return UserNotFoundResponse
	if existingUser == nil {
		log.Info().Msgf("user %s not found, generating register token", userInfo.Email)

		registerToken, err := s.tokenManager.GenerateRegisterToken(userInfo.Email)
		if err != nil {
			log.Error().Err(err).Msg("failed to generate register token")
			return nil, nil, fmt.Errorf("failed to generate register token: %w", err)
		}

		return nil, &responses.UserNotFoundResponse{
			Code:          "USER_NOT_FOUND",
			Message:       "User not found. Please complete registration.",
			RegisterToken: registerToken,
			Email:         userInfo.Email,
		}, nil
	}

	// Check if user is active
	if !existingUser.IsActive {
		log.Warn().Msgf("user %s attempted to login but account is deactivated", existingUser.Email)
		return nil, nil, exceptions.ErrUserInactive
	}

	// User exists, proceed with login
	// Update profile picture from Google every login
	if userInfo.ProfileURL != "" {
		existingUser.ProfileURL = &userInfo.ProfileURL
	}

	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		log.Error().Err(err).Msg("failed to update user profile picture")
		// Don't fail login if update fails, just log the error
	}

	// Generate JWT token
	jwtToken, err := s.tokenManager.GenerateAccessToken(existingUser.UserID, existingUser.Email, existingUser.Role)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate token")
		return nil, nil, fmt.Errorf("failed to generate token: %w", err)
	}

	log.Info().Msgf("user %s logged in successfully", existingUser.Email)

	return &responses.LoginResponse{
		AccessToken: jwtToken,
		User: responses.UserInfo{
			UserID:      existingUser.UserID,
			Email:       existingUser.Email,
			FirstName:   existingUser.FirstName,
			LastName:    existingUser.LastName,
			ProfileURL:  existingUser.ProfileURL,
			PhoneNumber: existingUser.PhoneNumber,
			Role:        existingUser.Role,
		},
	}, nil, nil
}

// generateState generates a random state for CSRF protection
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
