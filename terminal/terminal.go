package terminal

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
)

const (
	RegisterURL   = "/api/v1/terminal/terminal-registrations/"
	ConfigURL     = "/api/v1/terminal/terminals/config/"
	HeartbeatURL  = "/api/v1/terminal/terminals/status/"
	TaskURL       = "/api/v1/terminal/tasks/%s/"
	ConnMethodURL = "/api/v1/terminal/components/connect-methods/"
)

// Service handles /api/v1/terminal endpoints used by terminal
// components (Koko, Lion, Magnus, etc.).
type Service struct {
	client core.HTTPClient
}

// NewService creates a new terminal Service.
func NewService(c core.HTTPClient) *Service {
	return &Service{client: c}
}

// Register registers a new terminal component with the server.
func (s *Service) Register(ctx context.Context, name, typeName, comment string) (map[string]any, *core.Response, error) {
	body := map[string]string{"name": name, "type": typeName, "comment": comment}
	return sdkutil.MapAction(ctx, s.client, RegisterURL, body)
}

// Config returns the terminal configuration blob.
func (s *Service) Config(ctx context.Context) (map[string]any, *core.Response, error) {
	httpReq, err := s.client.NewRequest(ctx, "GET", ConfigURL, nil)
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

// Heartbeat posts a terminal status heartbeat.
func (s *Service) Heartbeat(ctx context.Context, statuses any) (*core.Response, error) {
	httpReq, err := s.client.NewRequest(ctx, "POST", HeartbeatURL, statuses)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, httpReq, nil)
}

// ConnectMethods returns the per-component connect-method map.
func (s *Service) ConnectMethods(ctx context.Context) (map[string]any, *core.Response, error) {
	httpReq, err := s.client.NewRequest(ctx, "GET", ConnMethodURL, nil)
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
