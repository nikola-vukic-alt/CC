package data

import (
	"library-app/local/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Borrows = []*model.Borrow{
	{
		Id:     getObjectId("62706d1b624b3da748f63fe3"),
		Title:  "War and Peace",
		Author: "F. M. Dostoyevsky",
		ISBN:   "978-3-16-148410-0",
	},
	{
		Id:     getObjectId("62706d1b623b3da748f63fa1"),
		Title:  "Pere Goriot",
		Author: "Honore de Balzac",
		ISBN:   "978-3-16-142411-0",
	},
	{
		Id:     getObjectId("55306d1b623b3da748f63fa1"),
		Title:  "Madame Bovary",
		Author: "Gustav Flaubert",
		ISBN:   "696-9-16-148410-0",
	},
	{
		Id:     getObjectId("62706d1b623b4da748f63bc3"),
		Title:  "Anna Karenina",
		Author: "L. N. Tolstoy",
		ISBN:   "420-3-69-148410-0",
	},
	{
		Id:     getObjectId("55306d1b615b3da748f63fa1"),
		Title:  "Waiting for Godot",
		Author: "Samuel Beckett",
		ISBN:   "978-3-16-148111-7",
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
