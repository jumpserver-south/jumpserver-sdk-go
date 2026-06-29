package model

// Session is an audit record of a user's connection to an asset.
type Session struct {
	ID            string `json:"id"`
	User          string `json:"user"`
	Asset         string `json:"asset"`
	Account       string `json:"account"`
	Protocol      string `json:"protocol"`
	Type          LabelValue `json:"type"`
	LoginFrom     LabelValue `json:"login_from"`
	RemoteAddr    string `json:"remote_addr"`
	IsFinished    bool   `json:"is_finished"`
	IsSuccess     bool   `json:"is_success"`
	OrgID         string `json:"org_id"`
	DateStart     string `json:"date_start"`
	DateEnd       string `json:"date_end"`
}

// SessionPage is the paginated list envelope for Sessions.
type SessionPage = Page[Session]

// CommandPage is the paginated list envelope for Commands.
type CommandPage = Page[Command]

// FTPLogPage is the paginated list envelope for FTPLogs.
type FTPLogPage = Page[FTPLog]

// LoginLogPage is the paginated list envelope for LoginLogs.
type LoginLogPage = Page[LoginLog]

// OperateLogPage is the paginated list envelope for OperateLogs.
type OperateLogPage = Page[OperateLog]

// Command is a recorded command in a session.
type Command struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Asset     string `json:"asset"`
	Account   string `json:"account"`
	Session   string `json:"session"`
	Input     string `json:"input"`
	Output    string `json:"output"`
	RiskLevel int    `json:"risk_level"`
	OrgID     string `json:"org_id"`
	Timestamp int64  `json:"timestamp"`
}

// FTPLog is a file-transfer audit record.
type FTPLog struct {
	ID          string `json:"id"`
	User        string `json:"user"`
	Asset       string `json:"asset"`
	Account     string `json:"account"`
	Session     string `json:"session"`
	RemoteAddr  string `json:"remote_addr"`
	Operate     LabelValue `json:"operate"`
	Path        string `json:"path"`
	IsSuccess   bool   `json:"is_success"`
	OrgID       string `json:"org_id"`
	DateStart   string `json:"date_start"`
}

// LoginLog is a user login audit record.
type LoginLog struct {
	ID        any    `json:"id"`
	Username  string `json:"username"`
	Type      any    `json:"type"`
	IP        string `json:"ip"`
	City      string `json:"city"`
	UserAgent string `json:"user_agent"`
	MFA       any    `json:"mfa"`
	Status    any    `json:"status"`
	Backend   string `json:"backend"`
	Reason    string `json:"reason"`
	Datetime  string `json:"datetime"`
}

// OperateLog is an admin operation audit record.
type OperateLog struct {
	ID           any        `json:"id"`
	User         string     `json:"user"`
	Action       LabelValue `json:"action"`
	ResourceType string     `json:"resource_type"`
	Resource     string     `json:"resource"`
	RemoteAddr   string     `json:"remote_addr"`
	Datetime     string     `json:"datetime"`
}

