package models

type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"username,omitempty"`
	First_name  string `json:"first_name,omitempty"`
	Last_name   string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Date_joined string `json:"date_joined,omitempty"`
	Password    string `json:"-"`
}
