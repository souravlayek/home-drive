package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type MetaData struct {
	Id        primitive.ObjectID  `json:"_id" bson:"_id, omitempty"`
	Name      string              `json:"name" bson:"name"`
	Path      string              `json:"path" bson:"path"`
	UpdatedAt primitive.Timestamp `json:"updatedAt" bson:"updatedAt"`
}
