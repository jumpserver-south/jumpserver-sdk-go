package orgs

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

const (
	ListURL   = "/api/v1/orgs/orgs/"
	DetailURL = "/api/v1/orgs/orgs/%s/"
)

// Service handles /api/v1/orgs/orgs.
type Service struct {
	client core.HTTPClient
}

// NewService creates a new orgs Service.
func NewService(c core.HTTPClient) *Service {
	return &Service{client: c}
}

// List returns a paginated list of organizations.
func (s *Service) List(ctx context.Context, opts *core.ListOptions) ([]model.Organization, *core.Response, error) {
	return sdkutil.List[model.Organization](ctx, s.client, ListURL, opts)
}

// Get fetches an organization by ID.
func (s *Service) Get(ctx context.Context, id string) (*model.Organization, *core.Response, error) {
	return sdkutil.Get[model.Organization](ctx, s.client, DetailURL, id)
}

// Create creates an organization.
func (s *Service) Create(ctx context.Context, req *model.OrganizationRequest) (*model.Organization, *core.Response, error) {
	return sdkutil.Create[model.Organization, model.OrganizationRequest](ctx, s.client, ListURL, req)
}

// Update patches an organization.
func (s *Service) Update(ctx context.Context, id string, req *model.OrganizationRequest) (*model.Organization, *core.Response, error) {
	return sdkutil.Update[model.Organization, model.OrganizationRequest](ctx, s.client, DetailURL, id, req)
}

// Delete deletes an organization.
func (s *Service) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, DetailURL, id)
}
