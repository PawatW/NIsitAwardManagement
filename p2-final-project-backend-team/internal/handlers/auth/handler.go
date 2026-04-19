package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/auth"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AuthHandler interface {
	GetGoogleLoginURL(c echo.Context) error
	GoogleCallback(c echo.Context) error
}

type handler struct {
	authService auth.Service
	cfg         *configs.Config
}

func NewAuthHandler(authService auth.Service, cfg *configs.Config) AuthHandler {
	return &handler{
		authService: authService,
		cfg:         cfg,
	}
}

// GetGoogleLoginURL returns the Google OAuth authorization URL
func (h *handler) GetGoogleLoginURL(c echo.Context) error {
	authURL, err := h.authService.GetGoogleAuthURL()
	if err != nil {
		log.Error().Err(err).Msg("failed to get Google auth URL")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to generate authorization URL",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"auth_url": authURL,
	})
}

// GoogleCallback handles the OAuth callback from Google
func (h *handler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		redirectURL := h.redirectWithMsg("Missing authorization code")
		return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}

	state := c.QueryParam("state")
	if state == "" {
		redirectURL := h.redirectWithMsg("Missing state parameter")
		return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}

	loginResponse, userNotFoundResponse, err := h.authService.HandleGoogleCallback(c.Request().Context(), code, state)
	if err != nil {
		// Handle inactive user error - redirect to frontend with error message
		if errors.Is(err, exceptions.ErrUserInactive) {
			log.Warn().Err(err).Msg("inactive user attempted to login")
			redirectURL := h.redirectWithMsg("Your account has been deactivated. Please contact an administrator.")
			return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		}
		log.Error().Err(err).Msg("failed to handle Google callback")
		redirectURL := h.redirectWithMsg("Authentication failed")
		return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}

	// Check if user not found (needs registration)
	if userNotFoundResponse != nil {
		log.Info().Msgf("user not found, redirecting to registration page")
		registrationURL := h.redirectToRegistration(userNotFoundResponse.RegisterToken, userNotFoundResponse.Email)
		return c.Redirect(http.StatusTemporaryRedirect, registrationURL)
	}

	// User exists, redirect to frontend with access token
	successURL := h.redirectWithTokens(loginResponse.AccessToken)
	log.Info().Msg("redirecting to frontend with tokens")
	return c.Redirect(http.StatusTemporaryRedirect, successURL)
}

func (h *handler) redirectWithMsg(message string) string {
	frontendURL := fmt.Sprintf("%s/oauth/callback", h.cfg.AllowOrigins[0])
	fmt.Print(frontendURL)

	resURL := fmt.Sprintf(
		"%s?msg=%s",
		frontendURL,
		url.QueryEscape(message), // Using message as query parameter for now
	)

	return resURL
}

func (h *handler) redirectWithTokens(accessToken string) string {
	frontendURL := fmt.Sprintf("%s/oauth/callback", h.cfg.AllowOrigins[0])
	return fmt.Sprintf(
		"%s?accessToken=%s&refreshToken=%s",
		frontendURL,
		url.QueryEscape(accessToken),
		url.QueryEscape(accessToken), // Using accessToken as refreshToken for now
	)
}

func (h *handler) redirectToRegistration(registerToken, email string) string {
	frontendURL := fmt.Sprintf("%s/register", h.cfg.AllowOrigins[0])
	return fmt.Sprintf(
		"%s?register_token=%s&email=%s",
		frontendURL,
		url.QueryEscape(registerToken),
		url.QueryEscape(email),
	)
}
