package audits

import (
	"context"
	"io"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

// Session URL constants.
const (
	SessionListURL   = "/api/v1/terminal/sessions/"
	SessionDetailURL = "/api/v1/terminal/sessions/%s/"
	SessionReplayURL = "/api/v1/terminal/sessions/%s/replay/"
)

// Service handles session, command, FTP, login, and operate log endpoints.
type Service struct {
	client core.HTTPClient
}

// NewService creates a new audits Service.
func NewService(c core.HTTPClient) *Service {
	return &Service{client: c}
}

// ListSessions returns a paginated list of sessions.
func (s *Service) ListSessions(ctx context.Context, opts *core.ListOptions) ([]model.Session, *core.Response, error) {
	return sdkutil.List[model.Session](ctx, s.client, SessionListURL, opts)
}

// GetSession fetches a session by ID.
func (s *Service) GetSession(ctx context.Context, id string) (*model.Session, *core.Response, error) {
	return sdkutil.Get[model.Session](ctx, s.client, SessionDetailURL, id)
}

// DownloadReplay streams the session replay archive into w.
func (s *Service) DownloadReplay(ctx context.Context, sessionID string, w io.Writer) (*core.Response, error) {
	httpReq, err := s.client.NewRequest(ctx, "GET", sdkutil.Spath(SessionReplayURL, sessionID), nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "*/*")
	resp, err := s.client.DoRaw(ctx, httpReq, w)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
