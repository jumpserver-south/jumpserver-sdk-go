package assets

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

// Platforms URL constants.
const (
	PlatformListURL   = "/api/v1/assets/platforms/"
	PlatformDetailURL = "/api/v1/assets/platforms/%d/"
	ProtocolListURL   = "/api/v1/assets/protocols/"
)

// PlatformsService handles /api/v1/assets/platforms.
type PlatformsService struct {
	client core.HTTPClient
}

// NewPlatformsService creates a new PlatformsService.
func NewPlatformsService(c core.HTTPClient) *PlatformsService {
	return &PlatformsService{client: c}
}

// List returns a paginated list of platforms.
func (s *PlatformsService) List(ctx context.Context, opts *core.ListOptions) ([]model.Platform, *core.Response, error) {
	return sdkutil.List[model.Platform](ctx, s.client, PlatformListURL, opts)
}

// Get fetches a platform by ID.
func (s *PlatformsService) Get(ctx context.Context, id int) (*model.Platform, *core.Response, error) {
	return sdkutil.Get[model.Platform](ctx, s.client, PlatformDetailURL, id)
}
