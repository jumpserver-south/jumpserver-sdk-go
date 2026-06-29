package accounts

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

const (
	ChangeSecretListURL    = "/api/v1/accounts/change-secret-automations/"
	ChangeSecretDetailURL  = "/api/v1/accounts/change-secret-automations/%s/"
	ChangeSecretExecuteURL = "/api/v1/accounts/change-secret-executions/"
)

// ChangeSecretService handles /api/v1/accounts/change-secret-automations.
type ChangeSecretService struct {
	client core.HTTPClient
}

// NewChangeSecretService creates a new ChangeSecretService.
func NewChangeSecretService(c core.HTTPClient) *ChangeSecretService {
	return &ChangeSecretService{client: c}
}

// List returns a paginated list of change secret automations.
func (s *ChangeSecretService) List(ctx context.Context, opts *core.ListOptions) ([]model.ChangeSecretAutomation, *core.Response, error) {
	return sdkutil.List[model.ChangeSecretAutomation](ctx, s.client, ChangeSecretListURL, opts)
}

// Get fetches a change secret automation by ID.
func (s *ChangeSecretService) Get(ctx context.Context, id string) (*model.ChangeSecretAutomation, *core.Response, error) {
	return sdkutil.Get[model.ChangeSecretAutomation](ctx, s.client, ChangeSecretDetailURL, id)
}

// Create creates a change secret automation.
func (s *ChangeSecretService) Create(ctx context.Context, req *model.ChangeSecretAutomationRequest) (*model.ChangeSecretAutomation, *core.Response, error) {
	return sdkutil.Create[model.ChangeSecretAutomation, model.ChangeSecretAutomationRequest](ctx, s.client, ChangeSecretListURL, req)
}

// Update patches a change secret automation.
func (s *ChangeSecretService) Update(ctx context.Context, id string, req *model.ChangeSecretAutomationRequest) (*model.ChangeSecretAutomation, *core.Response, error) {
	return sdkutil.Update[model.ChangeSecretAutomation, model.ChangeSecretAutomationRequest](ctx, s.client, ChangeSecretDetailURL, id, req)
}

// Delete deletes a change secret automation.
func (s *ChangeSecretService) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, ChangeSecretDetailURL, id)
}

// Execute triggers a change secret execution for the given automation.
func (s *ChangeSecretService) Execute(ctx context.Context, automationID string) (map[string]any, *core.Response, error) {
	body := map[string]string{"automation": automationID}
	return sdkutil.MapAction(ctx, s.client, ChangeSecretExecuteURL, body)
}
