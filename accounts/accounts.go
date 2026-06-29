package accounts

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

const (
	ListURL   = "/api/v1/accounts/accounts/"
	DetailURL = "/api/v1/accounts/accounts/%s/"
	BulkURL   = "/api/v1/accounts/accounts/bulk/"
	SecretURL = "/api/v1/accounts/account-secrets/%s/"
)

// Account connectivity testing URL constants (v4).
const (
	VerifyURL       = "/api/v1/accounts/accounts/verify/"
	VerifyDetailURL = "/api/v1/accounts/accounts/%s/verify/"
	VerifyTaskURL   = "/api/v1/accounts/accounts/verify/"
)

// Service handles /api/v1/accounts.
type Service struct {
	client core.HTTPClient
}

// NewService creates a new accounts Service.
func NewService(c core.HTTPClient) *Service {
	return &Service{client: c}
}

// List returns a paginated list of accounts.
func (s *Service) List(ctx context.Context, opts *core.ListOptions) ([]model.Account, *core.Response, error) {
	return sdkutil.List[model.Account](ctx, s.client, ListURL, opts)
}

// Get fetches an account by ID.
func (s *Service) Get(ctx context.Context, id string) (*model.Account, *core.Response, error) {
	return sdkutil.Get[model.Account](ctx, s.client, DetailURL, id)
}

// Create creates an account.
func (s *Service) Create(ctx context.Context, req *model.AccountRequest) (*model.Account, *core.Response, error) {
	return sdkutil.Create[model.Account, model.AccountRequest](ctx, s.client, ListURL, req)
}

// CreateBulk adds the same account to many assets in one call.
func (s *Service) CreateBulk(ctx context.Context, reqs []model.AccountRequest) (*core.Response, error) {
	httpReq, err := s.client.NewRequest(ctx, "POST", BulkURL, reqs)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, httpReq, nil)
}

// CreateBulkByTemplate adds accounts to assets using an account template.
func (s *Service) CreateBulkByTemplate(ctx context.Context, req *model.AccountBulkByTemplateRequest) (*core.Response, error) {
	httpReq, err := s.client.NewRequest(ctx, "POST", BulkURL, req)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, httpReq, nil)
}

// Update patches an account.
func (s *Service) Update(ctx context.Context, id string, req *model.AccountRequest) (*model.Account, *core.Response, error) {
	return sdkutil.Update[model.Account, model.AccountRequest](ctx, s.client, DetailURL, id, req)
}

// Delete deletes an account.
func (s *Service) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, DetailURL, id)
}

// GetSecret fetches the decrypted account secret.
func (s *Service) GetSecret(ctx context.Context, id string) (*model.Account, *core.Response, error) {
	return sdkutil.Get[model.Account](ctx, s.client, SecretURL, id)
}

// Verify returns the connectivity verification result for an account (v4).
func (s *Service) Verify(ctx context.Context, id string) (*model.AccountVerifyResult, *core.Response, error) {
	return sdkutil.Get[model.AccountVerifyResult](ctx, s.client, VerifyDetailURL, id)
}

// CreateVerifyTask creates a connectivity verification task (v4).
func (s *Service) CreateVerifyTask(ctx context.Context, req *model.AccountVerifyTaskRequest) (*model.AccountVerifyTask, *core.Response, error) {
	return sdkutil.Create[model.AccountVerifyTask, model.AccountVerifyTaskRequest](ctx, s.client, VerifyTaskURL, req)
}
