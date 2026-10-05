package services

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type AuthService struct {
	oauthConfig *oauth2.Config
}

func NewAuthService(clientID, clientSecret, redirectURL string) *AuthService {
	return &AuthService{
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL, // e.g. "http://localhost:8080/auth/google/callback"
			Scopes:       []string{sheets.SpreadsheetsReadonlyScope},
			Endpoint:     google.Endpoint,
		},
	}
}

func (s *AuthService) GetAuthURL() string {
	return s.oauthConfig.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce, // Forces Google to issue a new refresh token every time
	)
}

func (s *AuthService) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	tok, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	return tok, nil
}
