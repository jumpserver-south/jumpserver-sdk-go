package tickets

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

const (
	ListURL       = "/api/v1/tickets/tickets/"
	DetailURL     = "/api/v1/tickets/tickets/%s/"
	ApproveURL    = "/api/v1/tickets/tickets/%s/approve/"
	FlowListURL   = "/api/v1/tickets/flows/"
	FlowDetailURL = "/api/v1/tickets/flows/%s/"
)

// Service handles /api/v1/tickets.
type Service struct {
	client core.HTTPClient
}

// NewService creates a new tickets Service.
func NewService(c core.HTTPClient) *Service {
	return &Service{client: c}
}

// List returns a paginated list of tickets.
func (s *Service) List(ctx context.Context, opts *core.ListOptions) ([]model.Ticket, *core.Response, error) {
	return sdkutil.List[model.Ticket](ctx, s.client, ListURL, opts)
}

// Get fetches a ticket by ID.
func (s *Service) Get(ctx context.Context, id string) (*model.Ticket, *core.Response, error) {
	return sdkutil.Get[model.Ticket](ctx, s.client, DetailURL, id)
}

// Create opens an asset-application ticket.
func (s *Service) Create(ctx context.Context, req *model.TicketRequest) (*model.Ticket, *core.Response, error) {
	return sdkutil.Create[model.Ticket, model.TicketRequest](ctx, s.client, ListURL, req)
}

// Approve approves a ticket with action "approve" or "reject".
func (s *Service) Approve(ctx context.Context, id, action string) (*core.Response, error) {
	body := map[string]string{"action": action}
	httpReq, err := s.client.NewRequest(ctx, "POST", sdkutil.Spath(ApproveURL, id), body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, httpReq, nil)
}

// ListFlows returns a paginated list of ticket flows (workflow definitions).
func (s *Service) ListFlows(ctx context.Context, opts *core.ListOptions) ([]model.TicketFlow, *core.Response, error) {
	return sdkutil.List[model.TicketFlow](ctx, s.client, FlowListURL, opts)
}

// UpdateFlow patches a ticket flow definition.
func (s *Service) UpdateFlow(ctx context.Context, id string, req *model.TicketFlowRequest) (*model.TicketFlow, *core.Response, error) {
	return sdkutil.Update[model.TicketFlow, model.TicketFlowRequest](ctx, s.client, FlowDetailURL, id, req)
}
