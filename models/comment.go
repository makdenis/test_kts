package models

import "time"

type Comment struct {
	ID         int64     `json:"id"`
	Body       string    `json:"body,omitempty"`
	Creator_id int64     `json:"creator_id,omitempty"`
	Created    time.Time `json:"created,omitempty"`
	Topic_id   int64     `json:"-"`
}
