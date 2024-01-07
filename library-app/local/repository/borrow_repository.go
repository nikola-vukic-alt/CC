package repository

import (
	"context"
	"library-app/local/model"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BorrowRepository struct {
	collection *mongo.Collection
}

func NewBorrowRepository(database *mongo.Database) *BorrowRepository {
	return &BorrowRepository{
		collection: database.Collection("borrows"),
	}
}

func (r *BorrowRepository) SaveBorrow(ctx context.Context, borrow model.Borrow) error {
	borrow.Id = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, borrow)
	if err != nil {
		log.Printf("Error saving borrow: %v\n", err)
		return err
	}
	return nil
}
