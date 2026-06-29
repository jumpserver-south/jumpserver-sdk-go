package accounts

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

const (
	TemplateListURL   = "/api/v1/accounts/account-templates/"
	TemplateDetailURL = "/api/v1/accounts/account-templates/%s/"
)

// TemplatesService handles /api/v1/accounts/account-templates.
type TemplatesService struct {
	client core.HTTPClient
}

// NewTemplatesService creates a new TemplatesService.
func NewTemplatesService(c core.HTTPClient) *TemplatesService {
	return &TemplatesService{client: c}
}

// List returns a paginated list of account templates.
func (s *TemplatesService) List(ctx context.Context, opts *core.ListOptions) ([]model.AccountTemplate, *core.Response, error) {
	return sdkutil.List[model.AccountTemplate](ctx, s.client, TemplateListURL, opts)
}

// Get fetches an account template by ID.
func (s *TemplatesService) Get(ctx context.Context, id string) (*model.AccountTemplate, *core.Response, error) {
	return sdkutil.Get[model.AccountTemplate](ctx, s.client, TemplateDetailURL, id)
}

// Create creates an account template.
func (s *TemplatesService) Create(ctx context.Context, req *model.AccountTemplateRequest) (*model.AccountTemplate, *core.Response, error) {
	return sdkutil.Create[model.AccountTemplate, model.AccountTemplateRequest](ctx, s.client, TemplateListURL, req)
}

// Update patches an account template.
func (s *TemplatesService) Update(ctx context.Context, id string, req *model.AccountTemplateRequest) (*model.AccountTemplate, *core.Response, error) {
	return sdkutil.Update[model.AccountTemplate, model.AccountTemplateRequest](ctx, s.client, TemplateDetailURL, id, req)
}

// Delete deletes an account template.
func (s *TemplatesService) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, TemplateDetailURL, id)
}
