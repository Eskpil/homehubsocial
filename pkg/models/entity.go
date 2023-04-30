package models

import "time"

type EntityState struct {
	Id        string `json:"-" bson:"_id"`
	HistoryId string `json:"-" bson:"history_id"`

	State      string            `json:"state" bson:"state"`
	Attributes map[string]string `json:"attributes" bson:"attributes"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Entity struct {
	Id     string `json:"id" bson:"_id"`
	UserId string `json:"user_id" bson:"user_id"`
	Name   string `json:"name" bson:"name"`

	HistoryId string        `json:"-" bson:"history_id"`
	History   []EntityState `json:"history" bson:"-"`
}
