package perms

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

const (
	AssetPermissionListURL   = "/api/v1/perms/asset-permissions/"
	AssetPermissionDetailURL = "/api/v1/perms/asset-permissions/%s/"
	SelfAssetAccountsURL     = "/api/v1/perms/users/self/assets/%s/"
)

// Service handles /api/v1/perms/asset-permissions.
type Service struct {
	client core.HTTPClient
}

// NewService creates a new permissions Service.
func NewService(c core.HTTPClient) *Service {
	return &Service{client: c}
}

// List returns a paginated list of asset permissions.
func (s *Service) List(ctx context.Context, opts *core.ListOptions) ([]model.AssetPermission, *core.Response, error) {
	return sdkutil.List[model.AssetPermission](ctx, s.client, AssetPermissionListURL, opts)
}

// Get fetches an asset permission by ID.
func (s *Service) Get(ctx context.Context, id string) (*model.AssetPermission, *core.Response, error) {
	return sdkutil.Get[model.AssetPermission](ctx, s.client, AssetPermissionDetailURL, id)
}

// Create creates an asset permission.
func (s *Service) Create(ctx context.Context, req *model.AssetPermissionRequest) (*model.AssetPermission, *core.Response, error) {
	return sdkutil.Create[model.AssetPermission, model.AssetPermissionRequest](ctx, s.client, AssetPermissionListURL, req)
}

// Update patches an asset permission.
func (s *Service) Update(ctx context.Context, id string, req *model.AssetPermissionRequest) (*model.AssetPermission, *core.Response, error) {
	return sdkutil.Update[model.AssetPermission, model.AssetPermissionRequest](ctx, s.client, AssetPermissionDetailURL, id, req)
}

// Delete deletes an asset permission.
func (s *Service) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, AssetPermissionDetailURL, id)
}

// GetSelfAssetAccounts returns the accounts available to the current
// user for a specific asset. Corresponds to the Python
// get_self_asset_accounts(asset_id) call.
func (s *Service) GetSelfAssetAccounts(ctx context.Context, assetID string) (map[string]any, *core.Response, error) {
	httpReq, err := s.client.NewRequest(ctx, "GET", sdkutil.Spath(SelfAssetAccountsURL, assetID), nil)
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
