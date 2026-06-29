package model

// Role is an RBAC role. Scope is "org" or "system".
type Role struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Scope       LabelValue `json:"scope"`
	DisplayName string     `json:"display_name"`
	Comment     string     `json:"comment"`
	Internal    bool       `json:"builtin"`
	DateCreated string     `json:"date_created"`
	DateUpdated string     `json:"date_updated"`
}

// RolePage is the paginated list envelope for Roles.
type RolePage = Page[Role]
