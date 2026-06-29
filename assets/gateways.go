package assets

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

// Gateways URL constants.
const (
	GatewayListURL   = "/api/v1/assets/gateways/"
	GatewayDetailURL = "/api/v1/assets/gateways/%s/"
)

// GatewaysService handles /api/v1/assets/gateways.
type GatewaysService struct {
	client core.HTTPClient
}

// NewGatewaysService creates a new GatewaysService.
func NewGatewaysService(c core.HTTPClient) *GatewaysService {
	return &GatewaysService{client: c}
}

// List returns a paginated list of gateways.
func (s *GatewaysService) List(ctx context.Context, opts *core.ListOptions) ([]model.Gateway, *core.Response, error) {
	return sdkutil.List[model.Gateway](ctx, s.client, GatewayListURL, opts)
}

// Get fetches a gateway by ID.
func (s *GatewaysService) Get(ctx context.Context, id string) (*model.Gateway, *core.Response, error) {
	return sdkutil.Get[model.Gateway](ctx, s.client, GatewayDetailURL, id)
}

// Create creates a gateway.
func (s *GatewaysService) Create(ctx context.Context, req *model.GatewayRequest) (*model.Gateway, *core.Response, error) {
	return sdkutil.Create[model.Gateway, model.GatewayRequest](ctx, s.client, GatewayListURL, req)
}

// Update patches a gateway.
func (s *GatewaysService) Update(ctx context.Context, id string, req *model.GatewayRequest) (*model.Gateway, *core.Response, error) {
	return sdkutil.Update[model.Gateway, model.GatewayRequest](ctx, s.client, GatewayDetailURL, id, req)
}

// Delete deletes a gateway.
func (s *GatewaysService) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, GatewayDetailURL, id)
}
