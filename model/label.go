package model

// Label is a tag label (v3.10+ API).
type Label struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	DisplayName   string `json:"display_name"`
	Comment       string `json:"comment"`
	OrgID         string `json:"org_id"`
	OrgName       string `json:"org_name"`
	DateCreated   string `json:"date_created"`
	DateUpdated   string `json:"date_updated"`
	ResAmount     int    `json:"res_amount"`
}

// LabelRequest is the create/update payload.
type LabelRequest struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Value   string `json:"value"`
	Comment string `json:"comment,omitempty"`
}

// LabelPage is the paginated list envelope for Labels.
type LabelPage = Page[Label]
