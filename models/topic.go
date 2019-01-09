package models

type Topic struct {
	ID                 int64  `json:"id"`
	Title              string `json:"title,omitempty"`
	Body               string `json:"body,omitempty"`
	Number_of_comments int64  `json:"number_of_comments"`
	Number_of_likes    int64  `json:"number_of_likes"`
	Creator_id         int64  `json:"creator_id,omitempty"`
	Created            string `json:"created,omitempty"`
}
