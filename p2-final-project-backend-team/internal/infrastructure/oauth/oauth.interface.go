package oauth

import "context"

// UserInfo represents user information from OAuth provider
type UserInfo struct {
	Email      string
	FirstName  string
	LastName   string
	ProfileURL string
}

// Provider defines the interface for OAuth providers
type Provider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (string, error)
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)
}
