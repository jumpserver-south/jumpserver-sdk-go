package model

// Account is an asset account (credential).
type Account struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Username     string     `json:"username"`
	Asset        any        `json:"asset"`
	SecretType   LabelValue `json:"secret_type"`
	Secret       string     `json:"secret,omitempty"`
	Privileged   bool       `json:"privileged"`
	IsActive     bool       `json:"is_active"`
	Connectivity any        `json:"connectivity"`
	SuFrom       any        `json:"su_from"`
	Version      int        `json:"version"`
	Comment      string     `json:"comment"`
	CreatedBy    string     `json:"created_by"`
	UpdatedBy    string     `json:"updated_by"`
	DateCreated  string     `json:"date_created"`
	DateUpdated  string     `json:"date_updated"`
}

// AccountRequest is the create/update payload.
type AccountRequest struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Username   string `json:"username"`
	Asset      string `json:"asset"`
	SecretType string `json:"secret_type"`
	Secret     string `json:"secret,omitempty"`
	PushNow    bool   `json:"push_now,omitempty"`
	Privileged bool   `json:"privileged,omitempty"`
	IsActive   bool   `json:"is_active,omitempty"`
	SuFrom     string `json:"su_from,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// AccountPage is the paginated list envelope for Accounts.
type AccountPage = Page[Account]

// AccountTemplate is a reusable account credential template.
type AccountTemplate struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Username      string     `json:"username"`
	SecretType    LabelValue `json:"secret_type"`
	Secret        string     `json:"secret,omitempty"`
	Privileged    bool       `json:"privileged"`
	IsActive      bool       `json:"is_active"`
	SuFrom        IDName     `json:"su_from"`
	AutoPush      bool       `json:"auto_push"`
	PushParams    any        `json:"push_params"`
	OrgID         string     `json:"org_id"`
	OrgName       string     `json:"org_name"`
	Comment       string     `json:"comment"`
	DateCreated   string     `json:"date_created"`
	DateUpdated   string     `json:"date_updated"`
}

// AccountTemplateRequest is the create/update payload for account templates.
type AccountTemplateRequest struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	SecretType string `json:"secret_type"`
	Secret     string `json:"secret,omitempty"`
	Privileged bool   `json:"privileged,omitempty"`
	IsActive   bool   `json:"is_active,omitempty"`
	SuFrom     string `json:"su_from,omitempty"`
	AutoPush   bool   `json:"auto_push,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// AccountTemplatePage is the paginated list envelope for AccountTemplates.
type AccountTemplatePage = Page[AccountTemplate]

// AccountBulkByTemplateRequest adds accounts to assets using a template.
type AccountBulkByTemplateRequest struct {
	Assets    []string `json:"assets"`
	Template  string   `json:"template"`
	OnInvalid string   `json:"on_invalid,omitempty"`
	IsActive  bool     `json:"is_active,omitempty"`
}

// ChangeSecretAutomation represents an automated secret rotation policy.
type ChangeSecretAutomation struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Accounts       []any  `json:"accounts"`
	Assets         []any  `json:"assets"`
	Nodes          []any  `json:"nodes"`
	SecretType     any    `json:"secret_type"`
	SecretStrategy any    `json:"secret_strategy"`
	IsActive       bool   `json:"is_active"`
	Periodic       bool   `json:"is_periodic"`
	Crontab        string `json:"crontab"`
	Interval       int    `json:"interval"`
	Recipients     []any  `json:"recipients"`
	OrgID          string `json:"org_id"`
	OrgName        string `json:"org_name"`
	DateCreated    string `json:"date_created"`
	DateUpdated    string `json:"date_updated"`
}

// ChangeSecretAutomationRequest is the create/update payload.
type ChangeSecretAutomationRequest struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name"`
	Accounts       []string `json:"accounts"`
	Assets         []string `json:"assets,omitempty"`
	Nodes          []string `json:"nodes,omitempty"`
	SecretType     string   `json:"secret_type"`
	SecretStrategy string   `json:"secret_strategy,omitempty"`
	IsActive       bool     `json:"is_active,omitempty"`
	IsPeriodic     bool     `json:"is_periodic,omitempty"`
	Crontab        string   `json:"crontab,omitempty"`
	Interval       int      `json:"interval,omitempty"`
	Recipients     []string `json:"recipients,omitempty"`
	Comment        string   `json:"comment,omitempty"`
}

// ChangeSecretAutomationPage is the paginated list envelope for ChangeSecretAutomations.
type ChangeSecretAutomationPage = Page[ChangeSecretAutomation]

// AccountBackupPlan represents an account backup schedule.
type AccountBackupPlan struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Accounts       []any  `json:"accounts"`
	Assets         []any  `json:"assets"`
	Nodes          []any  `json:"nodes"`
	SecretType     any    `json:"secret_type"`
	SecretStrategy any    `json:"secret_strategy"`
	IsActive       bool   `json:"is_active"`
	IsPeriodic     bool   `json:"is_periodic"`
	Crontab        string `json:"crontab"`
	Interval       int    `json:"interval"`
	Recipients     []any  `json:"recipients"`
	BackupType     any    `json:"backup_type"`
	OrgID          string `json:"org_id"`
	OrgName        string `json:"org_name"`
	DateCreated    string `json:"date_created"`
	DateUpdated    string `json:"date_updated"`
}

// AccountBackupPlanRequest is the create/update payload.
type AccountBackupPlanRequest struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name"`
	Accounts       []string `json:"accounts"`
	Assets         []string `json:"assets,omitempty"`
	Nodes          []string `json:"nodes,omitempty"`
	SecretType     string   `json:"secret_type"`
	SecretStrategy string   `json:"secret_strategy,omitempty"`
	IsActive       bool     `json:"is_active,omitempty"`
	IsPeriodic     bool     `json:"is_periodic,omitempty"`
	Crontab        string   `json:"crontab,omitempty"`
	Interval       int      `json:"interval,omitempty"`
	Recipients     []string `json:"recipients,omitempty"`
	Comment        string   `json:"comment,omitempty"`
}

// AccountBackupPlanPage is the paginated list envelope for AccountBackupPlans.
type AccountBackupPlanPage = Page[AccountBackupPlan]

// AccountVerifyResult is the result of an account connectivity check (v4).
type AccountVerifyResult struct {
	Account      string `json:"account"`
	Asset        string `json:"asset"`
	IsValid      bool   `json:"is_valid"`
	Connectivity string `json:"connectivity"`
	Error        string `json:"error,omitempty"`
	DateVerified string `json:"date_verified"`
}

// AccountVerifyTaskRequest is the payload for creating a verification task.
type AccountVerifyTaskRequest struct {
	Accounts []string `json:"accounts"`
}

// AccountVerifyTask is the result of a verification task creation.
type AccountVerifyTask struct {
	Task string `json:"task"`
}
