package repository

import (
	"context"
	"fmt"
	"library-app/local/model"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
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

// return type is (borrow, err, statusCode)
func (r *BorrowRepository) GetMembersBorrow(ctx context.Context, memberId primitive.ObjectID, title string) (model.Borrow, error, int) {
	var borrow model.Borrow

	filter := bson.M{"_userId": memberId, "_title": title}
	err := r.collection.FindOne(ctx, filter).Decode(&borrow)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Borrow{}, fmt.Errorf("Borrow not found"), http.StatusNotFound
		}
		log.Printf("Error getting borrow: %v\n", err)
		return model.Borrow{}, err, http.StatusBadRequest
	}
	return borrow, nil, http.StatusOK
}

func (r *BorrowRepository) UpdateBorrow(ctx context.Context, id primitive.ObjectID, updatedBorrow model.Borrow) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updatedBorrow}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating member: %v\n", err)
		return err
	}
	return nil
}
