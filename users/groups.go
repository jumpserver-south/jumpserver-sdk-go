package users

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

const (
	GroupListURL     = "/api/v1/users/groups/"
	GroupDetailURL   = "/api/v1/users/groups/%s/"
	GroupRelationURL = "/api/v1/users/users-groups-relations/"
)

// GroupsService handles /api/v1/users/groups.
type GroupsService struct {
	client core.HTTPClient
}

// NewGroupsService creates a new user groups Service.
func NewGroupsService(c core.HTTPClient) *GroupsService {
	return &GroupsService{client: c}
}

// List returns a paginated list of groups.
func (s *GroupsService) List(ctx context.Context, opts *core.ListOptions) ([]model.Group, *core.Response, error) {
	return sdkutil.List[model.Group](ctx, s.client, GroupListURL, opts)
}

// Get fetches a group by ID.
func (s *GroupsService) Get(ctx context.Context, id string) (*model.Group, *core.Response, error) {
	return sdkutil.Get[model.Group](ctx, s.client, GroupDetailURL, id)
}

// Create creates a group.
func (s *GroupsService) Create(ctx context.Context, req *model.GroupRequest) (*model.Group, *core.Response, error) {
	return sdkutil.Create[model.Group, model.GroupRequest](ctx, s.client, GroupListURL, req)
}

// Update patches a group.
func (s *GroupsService) Update(ctx context.Context, id string, req *model.GroupRequest) (*model.Group, *core.Response, error) {
	return sdkutil.Update[model.Group, model.GroupRequest](ctx, s.client, GroupDetailURL, id, req)
}

// Delete deletes a group.
func (s *GroupsService) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, GroupDetailURL, id)
}

// BindUsers assigns a set of users to a group via the relation endpoint.
func (s *GroupsService) BindUsers(ctx context.Context, relations []model.UserGroupRelation) (*core.Response, error) {
	httpReq, err := s.client.NewRequest(ctx, "POST", GroupRelationURL, relations)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, httpReq, nil)
}

// ListUsers lists users belonging to a group via the users endpoint
// with a group_id filter: /api/v1/users/users/?group_id=<groupID>.
func (s *GroupsService) ListUsers(ctx context.Context, groupID string, opts *core.ListOptions) ([]model.User, *core.Response, error) {
	params := map[string]string{"group_id": groupID}
	if opts != nil {
		opts.Apply(params)
	}
	path := sdkutil.AppendQuery(ListURL, params)
	httpReq, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	var page model.UserPage
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
