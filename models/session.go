package models

type Session struct {
	Username string `json:"username,omitempty"`
	Session  string `json:"session,omitempty"`
	User_id  int64 `json:"User_id,omitempty"`
}
