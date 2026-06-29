package assets

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

// Zones URL constants.
const (
	ZoneListURL   = "/api/v1/assets/zones/"
	ZoneDetailURL = "/api/v1/assets/zones/%s/"
)

// ZonesService handles /api/v1/assets/zones network zones.
type ZonesService struct {
	client core.HTTPClient
}

// NewZonesService creates a new ZonesService.
func NewZonesService(c core.HTTPClient) *ZonesService {
	return &ZonesService{client: c}
}

// List returns a paginated list of zones.
func (s *ZonesService) List(ctx context.Context, opts *core.ListOptions) ([]model.Zone, *core.Response, error) {
	params := map[string]string{}
	if opts != nil {
		opts.Apply(params)
	}
	path := sdkutil.AppendQuery(ZoneListURL, params)
	httpReq, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	var page model.ZonePage
	resp, err := s.client.Do(ctx, httpReq, &page)
	if err != nil {
		return nil, resp, err
	}
	if resp != nil {
		resp.Count = page.Total
		resp.NextURL = page.NextURL
		resp.PreviousURL = page.PreviousURL
	}
	return page.Results, resp, nil
}

// Get fetches a zone by ID.
func (s *ZonesService) Get(ctx context.Context, id string) (*model.Zone, *core.Response, error) {
	return sdkutil.Get[model.Zone](ctx, s.client, ZoneDetailURL, id)
}

// Create creates a zone.
func (s *ZonesService) Create(ctx context.Context, req *model.ZoneRequest) (*model.Zone, *core.Response, error) {
	return sdkutil.Create[model.Zone, model.ZoneRequest](ctx, s.client, ZoneListURL, req)
}

// Update patches a zone.
func (s *ZonesService) Update(ctx context.Context, id string, req *model.ZoneRequest) (*model.Zone, *core.Response, error) {
	return sdkutil.Update[model.Zone, model.ZoneRequest](ctx, s.client, ZoneDetailURL, id, req)
}

// Delete deletes a zone by ID.
func (s *ZonesService) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, ZoneDetailURL, id)
}
