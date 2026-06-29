package model

import "fmt"

// Phone represents a JumpServer phone number with country code. The
// phone field can be a string or number in the API, so it is kept as
// any.
type Phone struct {
	Code  string `json:"code"`
	Phone any    `json:"phone"`
}

// String returns "+<code><number>" or "".
func (p *Phone) String() string {
	if p == nil {
		return ""
	}
	var n string
	switch v := p.Phone.(type) {
	case string:
		n = v
	case float64:
		n = fmt.Sprintf("%.0f", v)
	case int, int64:
		n = fmt.Sprintf("%d", v)
	default:
		return ""
	}
	if n == "" {
		return ""
	}
	code := p.Code
	if code == "" {
		code = "+86"
	}
	return code + n
}

// MfaLevel is a labelled MFA level.
type MfaLevel struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

// User is a JumpServer user.
type User struct {
	ID                      string              `json:"id"`
	Name                    string              `json:"name"`
	Username                string              `json:"username"`
	Email                   string              `json:"email"`
	Wechat                  string              `json:"wechat"`
	Phone                   Phone               `json:"phone"`
	MfaLevel                MfaLevel            `json:"mfa_level"`
	Source                  LabelValue          `json:"source"`
	WecomID                 string              `json:"wecom_id"`
	DingtalkID              string              `json:"dingtalk_id"`
	FeishuID                string              `json:"feishu_id"`
	LarkID                  string              `json:"lark_id"`
	SlackID                 string              `json:"slack_id"`
	CreatedBy               string              `json:"created_by"`
	UpdatedBy               string              `json:"updated_by"`
	Comment                 string              `json:"comment"`
	AvatarURL               string              `json:"avatar_url"`
	Groups                  IDNameList          `json:"groups"`
	SystemRoles             IDisplayNames       `json:"system_roles"`
	OrgRoles                IDisplayNames       `json:"org_roles"`
	Labels                  []string            `json:"labels"`
	PasswordStrategy        LabelValue          `json:"password_strategy"`
	IsSuperuser             bool                `json:"is_superuser"`
	IsOrgAdmin              bool                `json:"is_org_admin"`
	IsServiceAccount        bool                `json:"is_service_account"`
	IsValid                 bool                `json:"is_valid"`
	IsExpired               bool                `json:"is_expired"`
	IsActive                bool                `json:"is_active"`
	IsOtpSecretKeyBound     bool                `json:"is_otp_secret_key_bound"`
	CanPublicKeyAuth        bool                `json:"can_public_key_auth"`
	MfaEnabled              bool                `json:"mfa_enabled"`
	NeedUpdatePassword      bool                `json:"need_update_password"`
	MfaForceEnabled         bool                `json:"mfa_force_enabled"`
	IsFirstLogin            bool                `json:"is_first_login"`
	LoginBlocked            bool                `json:"login_blocked"`
	DateExpired             string              `json:"date_expired"`
	DateJoined              string              `json:"date_joined"`
	LastLogin               string              `json:"last_login"`
	DateUpdated             string              `json:"date_updated"`
	DateAPIKeyLastUsed      string              `json:"date_api_key_last_used"`
	DatePasswordLastUpdated string              `json:"date_password_last_updated"`
	OrgsRoles               map[string][]string `json:"orgs_roles,omitempty"`
}

// String renders "name(username)".
func (u User) String() string { return fmt.Sprintf("%s(%s)", u.Name, u.Username) }

// UserRequest is the payload for user create/update.
type UserRequest struct {
	ID                 string   `json:"id,omitempty"`
	Name               string   `json:"name"`
	Username           string   `json:"username"`
	Email              string   `json:"email"`
	Groups             []string `json:"groups"`
	PasswordStrategy   string   `json:"password_strategy,omitempty"`
	NeedUpdatePassword bool     `json:"need_update_password,omitempty"`
	MfaLevel           int      `json:"mfa_level,omitempty"`
	Source             string   `json:"source"`
	SystemRoles        []string `json:"system_roles"`
	OrgRoles           []string `json:"org_roles"`
	IsActive           bool     `json:"is_active"`
	DateExpired        string   `json:"date_expired,omitempty"`
	Phone              string   `json:"phone,omitempty"`
	Password           string   `json:"password,omitempty"`
	Comment            string   `json:"comment,omitempty"`
}

// UserPage is the paginated list envelope for Users.
type UserPage = Page[User]
