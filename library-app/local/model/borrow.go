package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Borrow struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId primitive.ObjectID `bson:"_userId"`
	BookId primitive.ObjectID `bson:"_bookId"`
	From   time.Time          `bson:"_from"`
	To     time.Time          `bson:"_to"`
}
