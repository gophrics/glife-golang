package profile

import (
	"context"

	"../../common/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const GOOGLE_APP_ID = "249369235819-11cfia1ht584n1kmk6gh6kbba8ab429u.apps.googleusercontent.com"

func _RegisterUser(req RegisterUserRequest) (RegisterUserResponse, *mongo.InsertOneResult, error) {

	var result RegisterUserResponse
	result.Country = req.Country
	result.Email = req.Email
	result.Name = req.Name
	result.Phone = req.Phone
	result.ProfileId = primitive.NewObjectID()
	// BIG TODO: Hash Password
	// TODO: Assuming single email, that need not be the case, user can have multiple emails linked to same account
	// For example, registration with a non google email and trying to register later with a google email
	insertResult, err := mongodb.Profile.InsertOne(context.TODO(), req)

	return result, insertResult, err
}
