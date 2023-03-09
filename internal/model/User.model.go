package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserModel struct {
	UserID   primitive.ObjectID `json:"userId,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
}
