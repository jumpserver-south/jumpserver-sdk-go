package model

// Platform describes a JumpServer platform template (Linux, Windows, etc.).
type Platform struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Category    LabelValue   `json:"category"`
	Type        LabelValue   `json:"type"`
	Charset     LabelValue   `json:"charset"`
	Internal    bool         `json:"internal"`
	Domain      bool         `json:"domain_enabled"`
	SuEnabled   bool         `json:"su_enabled"`
	Protocols   []NamePort   `json:"protocols"`
	Comment     string       `json:"comment"`
	CreatedBy   string       `json:"created_by"`
	UpdatedBy   string       `json:"updated_by"`
	DateCreated string       `json:"date_created"`
	DateUpdated string       `json:"date_updated"`
}

// PlatformPage is the paginated list envelope for Platforms.
type PlatformPage = Page[Platform]
