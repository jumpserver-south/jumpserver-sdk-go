package model

// Group is a user group.
type Group struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Comment     string   `json:"comment"`
	Users       []string `json:"users,omitempty"`
	OrgID       string   `json:"org_id"`
	OrgName     string   `json:"org_name"`
	DateCreated string   `json:"date_created"`
	DateUpdated string   `json:"date_updated"`
}

// GroupRequest is the create/update payload.
type GroupRequest struct {
	ID      string   `json:"id,omitempty"`
	Name    string   `json:"name"`
	Comment string   `json:"comment,omitempty"`
	Users   []string `json:"users,omitempty"`
}

// GroupPage is the paginated list envelope for Groups.
type GroupPage = Page[Group]

// UserGroupRelation binds a user to a group.
type UserGroupRelation struct {
	User      string `json:"user"`
	UserGroup string `json:"usergroup"`
	OrgID     string `json:"org_id,omitempty"`
}
