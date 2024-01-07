package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID     primitive.ObjectID `bson:"_id"`
	Title  string             `bson:"_title"`
	Author string             `bson:"_author"`
	ISBN   string             `bson:"_isbn"`
}
