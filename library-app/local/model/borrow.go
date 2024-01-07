package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Borrow struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId primitive.ObjectID `bson:"_userId"`
	Title  string             `bson:"_title"`
	Author string             `bson:"_author"`
	ISBN   string             `bson:"_isbn"`
	From   time.Time          `bson:"_from"`
	To     time.Time          `bson:"_to"`
}
