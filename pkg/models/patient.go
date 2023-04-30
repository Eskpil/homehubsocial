package models

type Patient struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Age  int64  `json:"age" bson:"age"`
}
