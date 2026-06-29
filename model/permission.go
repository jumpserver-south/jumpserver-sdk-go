package model

// AssetPermission is an asset-authorization rule.
type AssetPermission struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Users          IDNameList   `json:"users"`
	UserGroups     IDNameList   `json:"user_groups"`
	Assets         IDNameList   `json:"assets"`
	Nodes          IDNameList   `json:"nodes"`
	Accounts       []string     `json:"accounts"`
	Protocols      []string     `json:"protocols"`
	Actions        []LabelValue `json:"actions"`
	IsActive       bool         `json:"is_active"`
	DateStart      string       `json:"date_start"`
	DateExpired    string       `json:"date_expired"`
	Comment        string       `json:"comment"`
	OrgID          string       `json:"org_id"`
	OrgName        string       `json:"org_name"`
	CreatedBy      string       `json:"created_by"`
	DateCreated    string       `json:"date_created"`
	DateUpdated    string       `json:"date_updated"`
}

// AssetPermissionRequest is the create/update payload.
type AssetPermissionRequest struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Users       []string `json:"users,omitempty"`
	UserGroups  []string `json:"user_groups,omitempty"`
	Assets      []string `json:"assets,omitempty"`
	Nodes       []string `json:"nodes,omitempty"`
	Accounts    []string `json:"accounts,omitempty"`
	Protocols   []string `json:"protocols,omitempty"`
	Actions     []string `json:"actions,omitempty"`
	IsActive    bool     `json:"is_active,omitempty"`
	DateStart   string   `json:"date_start,omitempty"`
	DateExpired string   `json:"date_expired,omitempty"`
	Comment     string   `json:"comment,omitempty"`
}

// AssetPermissionPage is the paginated list envelope for AssetPermissions.
type AssetPermissionPage = Page[AssetPermission]
