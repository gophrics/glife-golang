package profile

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetUserRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Country string `json:"country"`
}

type GetUserResponse struct {
	ProfileId primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Country   string             `json:"country"`
}

type RegisterUserWithGoogleRequest struct {
	Token string `json:"token"`
}

type LoginUserWithGoogleRequest struct {
	RegisterUserWithGoogleRequest
	Email string `json:"email"`
}

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Country  string `json:"country"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	ProfileId primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Country   string             `json:"country"`
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
