package model

// Organization is a JumpServer organization.
type Organization struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	IsRoot      bool     `json:"is_root"`
	IsDefault   bool     `json:"is_default"`
	Members     []string `json:"members,omitempty"`
	Comment     string   `json:"comment"`
	DateCreated string   `json:"date_created"`
	DateUpdated string   `json:"date_updated"`
	CreatedBy   string   `json:"created_by"`
}

// OrganizationRequest is the create/update payload.
type OrganizationRequest struct {
	ID      string   `json:"id,omitempty"`
	Name    string   `json:"name"`
	Members []string `json:"members,omitempty"`
	Comment string   `json:"comment,omitempty"`
}

// OrganizationPage is the paginated list envelope for Organizations.
type OrganizationPage = Page[Organization]
