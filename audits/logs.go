package audits

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

// Login / Operate log URL constants.
const (
	LoginLogListURL          = "/api/v1/audits/login-logs/"
	OperateLogListURL        = "/api/v1/audits/operate-logs/"
	PasswordChangeLogListURL = "/api/v1/audits/password-change-logs/"
	JobLogListURL            = "/api/v1/audits/job-logs/"
)

// ListLoginLogs returns a paginated list of login logs.
func (s *Service) ListLoginLogs(ctx context.Context, opts *core.ListOptions) ([]model.LoginLog, *core.Response, error) {
	return sdkutil.List[model.LoginLog](ctx, s.client, LoginLogListURL, opts)
}

// ListOperateLogs returns a paginated list of operate logs.
func (s *Service) ListOperateLogs(ctx context.Context, opts *core.ListOptions) ([]model.OperateLog, *core.Response, error) {
	return sdkutil.List[model.OperateLog](ctx, s.client, OperateLogListURL, opts)
}
