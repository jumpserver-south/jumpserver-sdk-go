package acls

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

// Command filter URL constants.
const (
	CommandFilterListURL   = "/api/v1/acls/command-filter-acls/"
	CommandFilterDetailURL = "/api/v1/acls/command-filter-acls/%s/"
	CommandGroupListURL    = "/api/v1/acls/command-groups/"
	CommandGroupDetailURL  = "/api/v1/acls/command-groups/%s/"
	CommandReviewURL       = "/api/v1/acls/command-filter-acls/command-review/"
)

// Login ACL URL constants.
const (
	LoginACLListURL    = "/api/v1/acls/login-acls/"
	LoginACLDetailURL  = "/api/v1/acls/login-acls/%s/"
	LoginAssetCheckURL = "/api/v1/acls/login-asset/check/"
)

// CommandFiltersService handles /api/v1/acls/command-filter-acls and
// /api/v1/acls/command-groups.
type CommandFiltersService struct {
	client core.HTTPClient
}

// NewCommandFiltersService creates a new CommandFiltersService.
func NewCommandFiltersService(c core.HTTPClient) *CommandFiltersService {
	return &CommandFiltersService{client: c}
}

// List returns a paginated list of command filters.
func (s *CommandFiltersService) List(ctx context.Context, opts *core.ListOptions) ([]model.CommandFilter, *core.Response, error) {
	return sdkutil.List[model.CommandFilter](ctx, s.client, CommandFilterListURL, opts)
}

// Get fetches a command filter by ID.
func (s *CommandFiltersService) Get(ctx context.Context, id string) (*model.CommandFilter, *core.Response, error) {
	return sdkutil.Get[model.CommandFilter](ctx, s.client, CommandFilterDetailURL, id)
}

// Create creates a command filter.
func (s *CommandFiltersService) Create(ctx context.Context, req *model.CommandFilterRequest) (*model.CommandFilter, *core.Response, error) {
	return sdkutil.Create[model.CommandFilter, model.CommandFilterRequest](ctx, s.client, CommandFilterListURL, req)
}

// Update patches a command filter.
func (s *CommandFiltersService) Update(ctx context.Context, id string, req *model.CommandFilterRequest) (*model.CommandFilter, *core.Response, error) {
	return sdkutil.Update[model.CommandFilter, model.CommandFilterRequest](ctx, s.client, CommandFilterDetailURL, id, req)
}

// Delete deletes a command filter.
func (s *CommandFiltersService) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, CommandFilterDetailURL, id)
}

// ListGroups returns a paginated list of command groups.
func (s *CommandFiltersService) ListGroups(ctx context.Context, opts *core.ListOptions) ([]model.CommandGroup, *core.Response, error) {
	return sdkutil.List[model.CommandGroup](ctx, s.client, CommandGroupListURL, opts)
}

// GetGroup fetches a command group by ID.
func (s *CommandFiltersService) GetGroup(ctx context.Context, id string) (*model.CommandGroup, *core.Response, error) {
	return sdkutil.Get[model.CommandGroup](ctx, s.client, CommandGroupDetailURL, id)
}

// CreateGroup creates a command group.
func (s *CommandFiltersService) CreateGroup(ctx context.Context, req *model.CommandGroupRequest) (*model.CommandGroup, *core.Response, error) {
	return sdkutil.Create[model.CommandGroup, model.CommandGroupRequest](ctx, s.client, CommandGroupListURL, req)
}

// UpdateGroup patches a command group.
func (s *CommandFiltersService) UpdateGroup(ctx context.Context, id string, req *model.CommandGroupRequest) (*model.CommandGroup, *core.Response, error) {
	return sdkutil.Update[model.CommandGroup, model.CommandGroupRequest](ctx, s.client, CommandGroupDetailURL, id, req)
}

// DeleteGroup deletes a command group.
func (s *CommandFiltersService) DeleteGroup(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, CommandGroupDetailURL, id)
}

// LoginACLsService handles /api/v1/acls/login-acls.
type LoginACLsService struct {
	client core.HTTPClient
}

// NewLoginACLsService creates a new LoginACLsService.
func NewLoginACLsService(c core.HTTPClient) *LoginACLsService {
	return &LoginACLsService{client: c}
}

// List returns a paginated list of login ACLs.
func (s *LoginACLsService) List(ctx context.Context, opts *core.ListOptions) ([]model.LoginACL, *core.Response, error) {
	return sdkutil.List[model.LoginACL](ctx, s.client, LoginACLListURL, opts)
}

// Get fetches a login ACL by ID.
func (s *LoginACLsService) Get(ctx context.Context, id string) (*model.LoginACL, *core.Response, error) {
	return sdkutil.Get[model.LoginACL](ctx, s.client, LoginACLDetailURL, id)
}
