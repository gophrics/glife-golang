package profile

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewProfileId() primitive.ObjectID {
	return primitive.NewObjectID()
}

func _GenerateUsername() string {
	return fmt.Sprintf("randomdog%s", 2)
}

func GetEmailFromProfileId() string {
	return "stub"
}
