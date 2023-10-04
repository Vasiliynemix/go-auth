package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type LoginType struct {
	ID   int    `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	GUID        string             `bson:"guid,omitempty" json:"guid,omitempty"`
	Login       string             `bson:"login,omitempty" json:"login,omitempty"`
	LoginType   LoginType          `bson:"login_type,omitempty" json:"login_type,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	LastName    string             `bson:"last_name,omitempty" json:"last_name,omitempty"`
	LastLoginAt time.Time          `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}
