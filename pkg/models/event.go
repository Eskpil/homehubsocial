package models

import "time"

type Event struct {
	Id      string `json:"id" bson:"_id"`
	Content string `json:"content" bson:"content"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
