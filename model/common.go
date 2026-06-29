package model

// Page is a generic paginated list envelope matching JumpServer's
// {count, next, previous, results} shape. All *Page types in this
// package are now aliases to this generic type.
type Page[T any] struct {
	Total       int    `json:"count"`
	NextURL     string `json:"next"`
	PreviousURL string `json:"previous"`
	Results     []T    `json:"results"`
}

// IDName is a compact (id, name) pair used throughout JumpServer for
// foreign-key references.
type IDName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// IDNameList is a list of IDName pairs.
type IDNameList []IDName

// IDs returns just the ids.
func (s IDNameList) IDs() []string {
	ids := make([]string, len(s))
	for i, x := range s {
		ids[i] = x.ID
	}
	return ids
}

// IDDisplayName adds a display_name field used by roles and similar.
type IDDisplayName struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// IDisplayNames is a slice of IDDisplayName.
type IDisplayNames []IDDisplayName

// IDs returns the ids.
func (s IDisplayNames) IDs() []string {
	ids := make([]string, len(s))
	for i, x := range s {
		ids[i] = x.ID
	}
	return ids
}

// LabelValue is the {"label","value"} pair that JumpServer returns for
// enumerated fields (e.g. asset category, source).
type LabelValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// NamePort represents a protocol binding (name + port).
type NamePort struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

// PlatformMini is a minimal platform reference (used by Asset).
type PlatformMini struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// JMSDefaultOrg is the ID of the default "ROOT" organization.
const JMSDefaultOrg = "ROOT"

// JMSGlobalOrg is the well-known ID of the global organization.
const JMSGlobalOrg = "00000000-0000-0000-0000-000000000002"
