package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Member struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"_name"`
	Surname   string             `bson:"_surname"`
	Address   string             `bson:"_address"`
	SSN       string             `bson:"_ssn"`
	BorrowCnt int                `bson:"_borrowCnt"`
}
