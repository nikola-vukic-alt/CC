package repository

import (
	"context"
	"fmt"
	"library-app/central/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MemberRepository struct {
	collection *mongo.Collection
}

func NewMemberRepository(database *mongo.Database) *MemberRepository {
	return &MemberRepository{
		collection: database.Collection("members"),
	}
}

func (r *MemberRepository) SaveMember(ctx context.Context, member model.Member) error {
	member.Id = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, member)
	if err != nil {
		log.Printf("Error saving member: %v\n", err)
		return err
	}
	return nil
}

func (r *MemberRepository) GetMemberByID(ctx context.Context, id primitive.ObjectID) (model.Member, error) {
	var member model.Member
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&member)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Member{}, fmt.Errorf("member not found")
		}
		log.Printf("Error getting member by ID: %v\n", err)
		return model.Member{}, err
	}
	return member, nil
}

func (r *MemberRepository) GetMemberBySSN(ctx context.Context, ssn string) (model.Member, error) {
	var member model.Member
	filter := bson.M{"_ssn": ssn}
	err := r.collection.FindOne(ctx, filter).Decode(&member)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Member{}, fmt.Errorf("member not found")
		}
		log.Printf("Error getting member by SSN: %v\n", err)
		return model.Member{}, err
	}
	return member, nil
}

func (r *MemberRepository) GetAllMembers(ctx context.Context) ([]model.Member, error) {
	var members []model.Member
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error getting all members: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &members)
	if err != nil {
		log.Printf("Error decoding all members: %v\n", err)
		return nil, err
	}
	return members, nil
}

func (r *MemberRepository) UpdateMember(ctx context.Context, id primitive.ObjectID, updatedMember model.Member) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updatedMember}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating member: %v\n", err)
		return err
	}
	return nil
}

func (r *MemberRepository) DeleteMember(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting member: %v\n", err)
		return err
	}
	return nil
}
