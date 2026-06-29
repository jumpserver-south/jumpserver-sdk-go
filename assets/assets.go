package assets

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

// Assets URL constants.
const (
	AssetListURL          = "/api/v1/assets/assets/"
	AssetDetailURL        = "/api/v1/assets/assets/%s/"
	AssetPermUsersURL     = "/api/v1/assets/assets/%s/perm-users/"
	AssetPermUserPermsURL = "/api/v1/assets/assets/%s/perm-users/%s/permissions/"
)

// AssetsService handles the generic /api/v1/assets/assets endpoints.
type AssetsService struct {
	client core.HTTPClient
}

// NewAssetsService creates a new AssetsService.
func NewAssetsService(c core.HTTPClient) *AssetsService {
	return &AssetsService{client: c}
}

// List returns a paginated list of assets. Pass nil filters for no
// resource-specific filtering; common pagination goes in opts.
func (s *AssetsService) List(ctx context.Context, filters map[string]string, opts *core.ListOptions) ([]model.Asset, *core.Response, error) {
	params := map[string]string{}
	if opts != nil {
		opts.Apply(params)
	}
	for k, v := range filters {
		if v != "" {
			params[k] = v
		}
	}
	path := sdkutil.AppendQuery(AssetListURL, params)
	httpReq, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	var page model.AssetPage
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

// Get fetches an asset by ID.
func (s *AssetsService) Get(ctx context.Context, id string) (*model.Asset, *core.Response, error) {
	return sdkutil.Get[model.Asset](ctx, s.client, AssetDetailURL, id)
}

// Delete deletes an asset by ID.
func (s *AssetsService) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, AssetDetailURL, id)
}

// PermUsers returns the users permitted to access an asset.
func (s *AssetsService) PermUsers(ctx context.Context, assetID string, opts *core.ListOptions) ([]model.User, *core.Response, error) {
	params := map[string]string{}
	if opts != nil {
		opts.Apply(params)
	}
	path := sdkutil.AppendQuery(sdkutil.Spath(AssetPermUsersURL, assetID), params)
	httpReq, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	var out []model.User
	resp, err := s.client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}
