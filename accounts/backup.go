package accounts

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

const (
	BackupListURL    = "/api/v1/accounts/account-backup-plans/"
	BackupDetailURL  = "/api/v1/accounts/account-backup-plans/%s/"
	BackupExecuteURL = "/api/v1/accounts/account-backup-plan-executions/"
)

// BackupService handles /api/v1/accounts/account-backup-plans.
type BackupService struct {
	client core.HTTPClient
}

// NewBackupService creates a new BackupService.
func NewBackupService(c core.HTTPClient) *BackupService {
	return &BackupService{client: c}
}

// List returns a paginated list of account backup plans.
func (s *BackupService) List(ctx context.Context, opts *core.ListOptions) ([]model.AccountBackupPlan, *core.Response, error) {
	return sdkutil.List[model.AccountBackupPlan](ctx, s.client, BackupListURL, opts)
}

// Get fetches an account backup plan by ID.
func (s *BackupService) Get(ctx context.Context, id string) (*model.AccountBackupPlan, *core.Response, error) {
	return sdkutil.Get[model.AccountBackupPlan](ctx, s.client, BackupDetailURL, id)
}

// Create creates an account backup plan.
func (s *BackupService) Create(ctx context.Context, req *model.AccountBackupPlanRequest) (*model.AccountBackupPlan, *core.Response, error) {
	return sdkutil.Create[model.AccountBackupPlan, model.AccountBackupPlanRequest](ctx, s.client, BackupListURL, req)
}

// Update patches an account backup plan.
func (s *BackupService) Update(ctx context.Context, id string, req *model.AccountBackupPlanRequest) (*model.AccountBackupPlan, *core.Response, error) {
	return sdkutil.Update[model.AccountBackupPlan, model.AccountBackupPlanRequest](ctx, s.client, BackupDetailURL, id, req)
}

// Delete deletes an account backup plan.
func (s *BackupService) Delete(ctx context.Context, id string) (*core.Response, error) {
	return sdkutil.Delete(ctx, s.client, BackupDetailURL, id)
}

// Execute triggers a backup plan execution.
func (s *BackupService) Execute(ctx context.Context, planID string) (map[string]any, *core.Response, error) {
	body := map[string]string{"plan": planID}
	return sdkutil.MapAction(ctx, s.client, BackupExecuteURL, body)
}
