package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Service defines the interface for Google OAuth operations
type Service interface {
	GetAuthURL() string
	ExchangeCodeForToken(ctx context.Context, code string) (*GoogleUserInfo, error)
}

// ServiceImpl implements the Google OAuth service
type ServiceImpl struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

// GoogleUserInfo represents user information from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

// NewService creates a new Google OAuth service
func NewService(clientID, clientSecret, redirectURI string) Service {
	return &ServiceImpl{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

// GetAuthURL generates the Google OAuth authorization URL
func (s *ServiceImpl) GetAuthURL() string {
	params := url.Values{}
	params.Add("client_id", s.clientID)
	params.Add("redirect_uri", s.redirectURI)
	params.Add("scope", "openid email profile")
	params.Add("response_type", "code")
	params.Add("access_type", "offline")
	// Remove prompt: "consent" as it can cause issues

	authURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?%s", params.Encode())

	// Debug: log the generated URL (remove in production)
	fmt.Printf("Generated Google OAuth URL: %s\n", authURL)
	fmt.Printf("Redirect URI: %s\n", s.redirectURI)

	return authURL
}

// ExchangeCodeForToken exchanges authorization code for access token and user info
func (s *ServiceImpl) ExchangeCodeForToken(ctx context.Context, code string) (*GoogleUserInfo, error) {
	// Exchange code for access token
	token, err := s.exchangeCodeForAccessToken(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info using access token
	userInfo, err := s.getUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return userInfo, nil
}

// exchangeCodeForAccessToken exchanges authorization code for access token
func (s *ServiceImpl) exchangeCodeForAccessToken(ctx context.Context, code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", s.redirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		// Try to parse error response for better error messages
		var errorResp struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if json.Unmarshal(body, &errorResp) == nil {
			return "", fmt.Errorf("token exchange failed: %s - %s", errorResp.Error, errorResp.ErrorDescription)
		}
		return "", fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

// getUserInfo retrieves user information from Google using access token
func (s *ServiceImpl) getUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}
