package oauth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleProvider struct {
	config *oauth2.Config
}

// NewGoogleProvider creates a new Google OAuth provider
func NewGoogleProvider(cfg *configs.Config) Provider {
	return &googleProvider{
		config: &oauth2.Config{
			ClientID:     cfg.OAuth.GoogleClientID,
			ClientSecret: cfg.OAuth.GoogleClientSecret,
			RedirectURL:  cfg.OAuth.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// GetAuthURL generates the OAuth authorization URL
func (g *googleProvider) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode exchanges authorization code for access token
func (g *googleProvider) ExchangeCode(ctx context.Context, code string) (string, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}
	return token.AccessToken, nil
}

// GetUserInfo retrieves user information from Google
func (g *googleProvider) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	resp, err := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)).Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var googleUser struct {
		Email      string `json:"email"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture    string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &UserInfo{
		Email:      googleUser.Email,
		FirstName:  googleUser.GivenName,
		LastName:   googleUser.FamilyName,
		ProfileURL: googleUser.Picture,
	}, nil
}
