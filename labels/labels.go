package labels

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

const (
	ListURL          = "/api/v1/labels/labels/"
	DetailURL        = "/api/v1/labels/labels/%s/"
	ResourcesURL     = "/api/v1/labels/labels/%s/resource-types/%d/resources/"
	AssetListURL     = "/api/v1/assets/labels/"
	AssetDetailURL   = "/api/v1/assets/labels/%s/"
)

// Service handles /api/v1/labels/labels (v3.10+).
type Service struct {
	client core.HTTPClient
}

// NewService creates a new labels Service.
func NewService(c core.HTTPClient) *Service {
	return &Service{client: c}
}

// List returns a paginated list of labels.
func (s *Service) List(ctx context.Context, opts *core.ListOptions) ([]model.Label, *core.Response, error) {
	return sdkutil.List[model.Label](ctx, s.client, ListURL, opts)
}

// Get fetches a label by ID.
func (s *Service) Get(ctx context.Context, id string) (*model.Label, *core.Response, error) {
	return sdkutil.Get[model.Label](ctx, s.client, DetailURL, id)
}

// Create creates a label.
func (s *Service) Create(ctx context.Context, req *model.LabelRequest) (*model.Label, *core.Response, error) {
	return sdkutil.Create[model.Label, model.LabelRequest](ctx, s.client, ListURL, req)
}

// Update patches a label.
func (s *Service) Update(ctx context.Context, id string, req *model.LabelRequest) (*model.Label, *core.Response, error) {
	return sdkutil.Update[model.Label, model.LabelRequest](ctx, s.client, DetailURL, id, req)
}

// Delete deletes a label.
func (s *Service) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, DetailURL, id)
}
