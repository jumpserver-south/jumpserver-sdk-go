package assets

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

// Nodes URL constants.
const (
	NodeListURL   = "/api/v1/assets/nodes/"
	NodeDetailURL = "/api/v1/assets/nodes/%s/"
)

// NodesService handles /api/v1/assets/nodes.
type NodesService struct {
	client core.HTTPClient
}

// NewNodesService creates a new NodesService.
func NewNodesService(c core.HTTPClient) *NodesService {
	return &NodesService{client: c}
}

// List returns a paginated list of nodes.
func (s *NodesService) List(ctx context.Context, opts *core.ListOptions) ([]model.Node, *core.Response, error) {
	return sdkutil.List[model.Node](ctx, s.client, NodeListURL, opts)
}

// Get fetches a node by ID.
func (s *NodesService) Get(ctx context.Context, id string) (*model.Node, *core.Response, error) {
	return sdkutil.Get[model.Node](ctx, s.client, NodeDetailURL, id)
}

// Create creates a node.
func (s *NodesService) Create(ctx context.Context, req *model.NodeRequest) (*model.Node, *core.Response, error) {
	return sdkutil.Create[model.Node, model.NodeRequest](ctx, s.client, NodeListURL, req)
}

// Update patches a node.
func (s *NodesService) Update(ctx context.Context, id string, req *model.NodeRequest) (*model.Node, *core.Response, error) {
	return sdkutil.Update[model.Node, model.NodeRequest](ctx, s.client, NodeDetailURL, id, req)
}

// Delete deletes a node.
func (s *NodesService) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, NodeDetailURL, id)
}
