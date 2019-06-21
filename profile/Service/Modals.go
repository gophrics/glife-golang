package profile

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetUserRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Country string `json:"country"`
}

type User struct {
	ProfileId primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Country   string             `json:"country"`
	Username  string             `json:"username"`
	Password  string             `json:"password"`
}

type RegisterUserWithGoogleRequest struct {
	Token string `json:"token"`
}

type LoginUserWithGoogleRequest struct {
	RegisterUserWithGoogleRequest
	Email string `json:"email"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GoogleAuthVerification struct {
	Iss           string `json:"iss"`
	Sub           string `json:"sub"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Iat           string `json:"iap"`
	Exp           string `json:"exp"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}
