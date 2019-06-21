package profile

import (
	"context"
	"fmt"

	"../../common/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateProfileId() string {
	return primitive.NewObjectID().Hex()
}

func _GenerateUsername() string {
	return fmt.Sprintf("randomdog%s", 2)
}

func GetEmailFromProfileId() string {
	return "stub"
}

func _RegisterUser(req User) (User, *mongo.InsertOneResult, error) {

	var result User
	result.Country = req.Country
	result.Email = req.Email
	result.Name = req.Name
	result.Phone = req.Phone
	result.ProfileId = GenerateProfileId()
	// BIG TODO: Hash Password
	// TODO: Assuming single email, that need not be the case, user can have multiple emails linked to same account
	// For example, registration with a non google email and trying to register later with a google email
	insertResult, err := mongodb.Profile.InsertOne(context.TODO(), req)

	return result, insertResult, err
}
