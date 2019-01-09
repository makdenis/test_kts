package models

import "time"

type TopicLike struct {
	Topic_id int64     `json:"topic_id"`
	User_id  int64     `json:"user_id"`
	Created  time.Time `json:"-"`
}
