package model

// Ticket is a JumpServer approval ticket.
type Ticket struct {
	ID          string     `json:"id"`
	SerialNum   string     `json:"serial_num"`
	Type        LabelValue `json:"type"`
	State       LabelValue `json:"state"`
	Status      LabelValue `json:"status"`
	Title       string     `json:"title"`
	Comment     string     `json:"comment"`
	OrgID       string     `json:"org_id"`
	OrgName     string     `json:"org_name"`
	Applicant   any        `json:"applicant"`
	DateCreated string     `json:"date_created"`
	DateUpdated string     `json:"date_updated"`
}

// TicketRequest is the create payload (asset-application type).
type TicketRequest struct {
	Type           string   `json:"type"`
	Title          string   `json:"title"`
	ApplyAccounts  []string `json:"apply_accounts,omitempty"`
	ApplyAssets    []string `json:"apply_assets,omitempty"`
	ApplyNodes     []string `json:"apply_nodes,omitempty"`
	ApplyActions   []string `json:"apply_actions,omitempty"`
	ApplyDateStart string   `json:"apply_date_start,omitempty"`
	ApplyDateExpired string `json:"apply_date_expired,omitempty"`
	Comment        string   `json:"comment,omitempty"`
}

// TicketPage is the paginated list envelope for Tickets.
type TicketPage = Page[Ticket]

// TicketFlow is a workflow definition for ticket approval.
type TicketFlow struct {
	ID              string   `json:"id"`
	Type            any      `json:"type"`
	ApproveStrategy any      `json:"approve_strategy"`
	Applicants      []IDName `json:"applicants"`
	IsActive        bool     `json:"is_active"`
	DateCreated     string   `json:"date_created"`
	DateUpdated     string   `json:"date_updated"`
	CreatedBy       string   `json:"created_by"`
	OrgID           string   `json:"org_id"`
	OrgName         string   `json:"org_name"`
	Comment         string   `json:"comment"`
}

// TicketFlowRequest is the update payload for ticket flows.
type TicketFlowRequest struct {
	ApproveStrategy any  `json:"approve_strategy,omitempty"`
	Applicants      []string `json:"applicants,omitempty"`
	IsActive        bool `json:"is_active,omitempty"`
	Comment         string `json:"comment,omitempty"`
}

// TicketFlowPage is the paginated list envelope for TicketFlows.
type TicketFlowPage = Page[TicketFlow]
