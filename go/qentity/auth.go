package qentity

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var str *string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != nil {
		ns.Valid = true
		ns.String = *str
	} else {
		ns.Valid = false
	}

	return nil
}

type NullTime struct {
	sql.NullTime
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time)
	}
	return json.Marshal(nil)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var time *time.Time
	if err := json.Unmarshal(data, &time); err != nil {
		return err
	}

	if time != nil {
		nt.Valid = true
		nt.Time = *time
	} else {
		nt.Valid = false
	}

	return nil
}

func NewNullString(value string) NullString {
	return NullString{
		NullString: sql.NullString{
			String: value,
			Valid:  true,
		},
	}
}

func NewNullTime(time time.Time, valid bool) NullTime {
	return NullTime{
		NullTime: sql.NullTime{
			Time:  time,
			Valid: valid,
		},
	}
}

type Recaptcha struct {
	Secret      string  `json:"secret"`
	Threshold   float64 `json:"threshold"`
	ValidateURL string  `json:"validate_url"`
	Action      string  `json:"action"`
}

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type SessionData struct {
	UserCompanyID   string         `json:"id"`
	UserID          string         `json:"user_id"`
	CompanyID       string         `json:"company_id"`
	ParentID        NullString     `json:"parent_id"`
	TreePath        pq.StringArray `json:"tree_path"`
	IsAdmin         bool           `json:"is_admin"`
	Email           string         `json:"email"`
	Picture         NullString     `json:"picture"`
	UserName        string         `json:"username"`
	Name            string         `json:"name"`
	CompanyName     string         `json:"company_name"`
	CompanyCodeName string         `json:"company_codename"`
	IsCurrent       bool           `json:"is_current"`
	Password        []byte         `json:"password,omitempty"`
	SessionID       string         `json:"session_id"`
}

type ResponseValidate struct {
	SessionData SessionData   `json:"session_data"`
	UserProject []UserProject `json:"user_project"`
}

type UserProject struct {
	UserRoleID    string    `json:"id"`
	UserCompanyID string    `json:"user_company_id"`
	Roles         []Role    `json:"roles"`
	CompanyID     string    `json:"company_id"`
	ProjectID     string    `json:"project_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	IsAdmin       bool      `json:"is_admin"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
}

type Role struct {
	RoleID   string `json:"id"`
	RoleName string `json:"name"`
}
