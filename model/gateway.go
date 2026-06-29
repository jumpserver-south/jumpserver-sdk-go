package model

// Gateway is a bastion gateway attached to a domain/zone.
type Gateway struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Address     string       `json:"address"`
	Platform    PlatformMini `json:"platform"`
	Protocols   []NamePort   `json:"protocols"`
	Nodes       IDNameList   `json:"nodes"`
	IsActive    bool         `json:"is_active"`
	Comment     string       `json:"comment"`
	OrgID       string       `json:"org_id"`
	OrgName     string       `json:"org_name"`
	DateCreated string       `json:"date_created"`
	DateUpdated string       `json:"date_updated"`
}

// GatewayRequest is the create/update payload.
type GatewayRequest struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name"`
	Address   string     `json:"address"`
	Platform  int        `json:"platform"`
	Protocols []NamePort `json:"protocols,omitempty"`
	Nodes     []string   `json:"nodes,omitempty"`
	IsActive  bool       `json:"is_active,omitempty"`
	Comment   string     `json:"comment,omitempty"`
}

// GatewayPage is the paginated list envelope for Gateways.
type GatewayPage = Page[Gateway]
