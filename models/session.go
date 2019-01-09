package models

type Session struct {
	Username string `json:"username,omitempty"`
	Session  string `json:"session,omitempty"`
}
