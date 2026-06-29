package assets

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

// Category assets URL constants.
const (
	AssetCategoryListURL   = "/api/v1/assets/%ss/"
	AssetCategoryDetailURL = "/api/v1/assets/%ss/%s/"
)

// CategoryService is a typed facade over a single asset category
// (hosts, devices, databases, webs, clouds, customs).
type CategoryService struct {
	client   core.HTTPClient
	category string
}

// NewCategoryService creates a new CategoryService for the given category.
func NewCategoryService(c core.HTTPClient, category string) *CategoryService {
	return &CategoryService{client: c, category: category}
}

// List returns a paginated list of assets in this category.
func (s *CategoryService) List(ctx context.Context, filters map[string]string, opts *core.ListOptions) ([]model.Asset, *core.Response, error) {
	params := map[string]string{}
	if opts != nil {
		opts.Apply(params)
	}
	for k, v := range filters {
		if v != "" {
			params[k] = v
		}
	}
	path := sdkutil.AppendQuery(sdkutil.Spath(AssetCategoryListURL, s.category), params)
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

// Get fetches a category-scoped asset.
func (s *CategoryService) Get(ctx context.Context, id string) (*model.Asset, *core.Response, error) {
	url := sdkutil.Spath(AssetCategoryDetailURL, s.category, id)
	httpReq, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	var out model.Asset
	resp, err := s.client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Create creates a category-scoped asset.
func (s *CategoryService) Create(ctx context.Context, req *model.AssetRequest) (*model.Asset, *core.Response, error) {
	url := sdkutil.Spath(AssetCategoryListURL, s.category)
	httpReq, err := s.client.NewRequest(ctx, "POST", url, req)
	if err != nil {
		return nil, nil, err
	}
	var out model.Asset
	resp, err := s.client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Update patches a category-scoped asset.
func (s *CategoryService) Update(ctx context.Context, id string, req *model.AssetRequest) (*model.Asset, *core.Response, error) {
	url := sdkutil.Spath(AssetCategoryDetailURL, s.category, id)
	httpReq, err := s.client.NewRequest(ctx, "PATCH", url, req)
	if err != nil {
		return nil, nil, err
	}
	var out model.Asset
	resp, err := s.client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Replace replaces a category-scoped asset.
func (s *CategoryService) Replace(ctx context.Context, id string, req *model.AssetRequest) (*model.Asset, *core.Response, error) {
	url := sdkutil.Spath(AssetCategoryDetailURL, s.category, id)
	httpReq, err := s.client.NewRequest(ctx, "PUT", url, req)
	if err != nil {
		return nil, nil, err
	}
	var out model.Asset
	resp, err := s.client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Delete deletes a category-scoped asset.
func (s *CategoryService) Delete(ctx context.Context, id string) (*core.Response, error) {
	url := sdkutil.Spath(AssetCategoryDetailURL, s.category, id)
	httpReq, err := s.client.NewRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, httpReq, nil)
}
