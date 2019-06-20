package profile

import "go.mongodb.org/mongo-driver/bson/primitive"

func NewProfileId() primitive.ObjectID {
	return primitive.NewObjectID()
}
