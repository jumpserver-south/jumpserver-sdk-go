package model

// Zone is a network zone for grouping assets.
type Zone struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Assets       IDNameList `json:"assets"`
	Gateways     IDNameList `json:"gateways"`
	AssetsAmount int        `json:"assets_amount,omitempty"`
	Labels       IDNameList `json:"labels,omitempty"`
	Comment      string     `json:"comment"`
	OrgID        string     `json:"org_id"`
	OrgName      string     `json:"org_name"`
	CreatedBy    string     `json:"created_by,omitempty"`
	UpdatedBy    string     `json:"updated_by,omitempty"`
	DateCreated   string     `json:"date_created"`
	DateUpdated   string     `json:"date_updated"`
}

// ZoneRequest is the create/update payload.
type ZoneRequest struct {
	ID       string   `json:"id,omitempty"`
	Name     string   `json:"name"`
	Assets   []string `json:"assets,omitempty"`
	Gateways []string `json:"gateways,omitempty"`
	Labels   []string `json:"labels,omitempty"`
	Comment  string   `json:"comment,omitempty"`
}

// ZonePage is the paginated list envelope for Zones.
type ZonePage = Page[Zone]
