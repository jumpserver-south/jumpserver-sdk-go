package auth

import (
	"context"
	"fmt"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

// authentication URL constants.
const (
	TokenURL              = "/api/v1/authentication/tokens/"
	TokenConfirmURL       = "/api/v1/authentication/login-confirm-ticket/status/"
	MFASelectURL          = "/api/v1/authentication/mfa/select/"
	ConfirmURL            = "/api/v1/authentication/confirm/"
	SuperConnectionToken  = "/api/v1/authentication/super-connection-token/"
	SuperConnectionSecret = "/api/v1/authentication/super-connection-token/secret/"
	ConnectionTokenURL    = "/api/v1/authentication/connection-token/"
	SSOLoginURLPath       = "/api/v1/authentication/sso/login-url/"
)

// Service handles authentication-related endpoints (login, MFA,
// super-connection-token).
type Service struct {
	client core.HTTPClient
}

// NewService creates a new auth Service.
func NewService(c core.HTTPClient) *Service {
	return &Service{client: c}
}

// CreateToken performs username/password login and returns a Bearer
// token suitable for subsequent API calls.
func (s *Service) CreateToken(ctx context.Context, req *model.TokenRequest) (*model.Token, *core.Response, error) {
	return sdkutil.Create[model.Token, model.TokenRequest](ctx, s.client, TokenURL, req)
}

// ConfirmLoginStatus polls the login-confirm ticket status.
func (s *Service) ConfirmLoginStatus(ctx context.Context, ticketID string) (map[string]any, *core.Response, error) {
	path := fmt.Sprintf("%s?ticket_id=%s", TokenConfirmURL, ticketID)
	httpReq, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	out := map[string]any{}
	resp, err := s.client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// SelectMFA switches the MFA type mid-login flow.
func (s *Service) SelectMFA(ctx context.Context, ticketID, mfaType string) (map[string]any, *core.Response, error) {
	body := map[string]string{"mfa_type": mfaType}
	path := MFASelectURL + "?ticket_id=" + ticketID
	httpReq, err := s.client.NewRequest(ctx, "POST", path, body)
	if err != nil {
		return nil, nil, err
	}
	out := map[string]any{}
	resp, err := s.client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// CreateConnectionToken creates a connection token for accessing an
// asset. Requires user, asset, account, protocol, and connect_method.
func (s *Service) CreateConnectionToken(ctx context.Context, req *model.ConnectionTokenRequest) (*model.ConnectionToken, *core.Response, error) {
	return sdkutil.Create[model.ConnectionToken, model.ConnectionTokenRequest](ctx, s.client, ConnectionTokenURL, req)
}

// SSOLoginURL returns an SSO login URL for the given user. This is an
// enterprise-only feature.
func (s *Service) SSOLoginURL(ctx context.Context, req *model.SSOLoginRequest) (map[string]any, *core.Response, error) {
	return sdkutil.MapAction(ctx, s.client, SSOLoginURLPath, req)
}

// CreateSuperConnectionToken creates a super connection token (requires
// elevated privileges / API key). Used for SSO-based asset access.
func (s *Service) CreateSuperConnectionToken(ctx context.Context, req *model.ConnectionTokenRequest) (*model.ConnectionToken, *core.Response, error) {
	return sdkutil.Create[model.ConnectionToken, model.ConnectionTokenRequest](ctx, s.client, SuperConnectionToken, req)
}

// GetSuperConnectionTokenSecret retrieves the secret/auth info for a
// super connection token.
func (s *Service) GetSuperConnectionTokenSecret(ctx context.Context, tokenID string) (map[string]any, *core.Response, error) {
	body := map[string]any{"id": tokenID, "expire_now": false}
	return sdkutil.MapAction(ctx, s.client, SuperConnectionSecret, body)
}

// GetClientURL returns the client connection URL for a connection token
// that can be used to launch the local client directly.
func (s *Service) GetClientURL(ctx context.Context, tokenID string) (string, *core.Response, error) {
	path := fmt.Sprintf("/api/v1/authentication/connection-token/%s/client-url/", tokenID)
	httpReq, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return "", nil, err
	}
	out := map[string]any{}
	resp, err := s.client.Do(ctx, httpReq, &out)
	if err != nil {
		return "", resp, err
	}
	url, _ := out["url"].(string)
	return url, resp, nil
}
