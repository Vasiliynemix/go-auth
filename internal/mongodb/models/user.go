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
	ID          primitive.ObjectID `bson:"_id" json:"-"`
	GUID        string             `bson:"guid" json:"guid"`
	Login       string             `bson:"login" json:"login"`
	LoginType   LoginType          `bson:"login_type" json:"login_type"`
	Name        string             `bson:"name" json:"name"`
	LastName    string             `bson:"last_name" json:"last_name"`
	LastLoginAt time.Time          `bson:"last_login_at" json:"last_login_at"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}
