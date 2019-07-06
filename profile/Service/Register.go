package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"../../common"
	"../../common/mongodb"
	"../../common/redis"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
)

func RegisterUserWithGoogle(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req RegisterUserWithGoogleRequest

	json.Unmarshal(b, &req)

	// Verify google auth token is valid for our client
	resp, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", req.Token))

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var resp2 GoogleAuthVerification
	json.Unmarshal(body, &resp2)

	if resp2.Aud != "249369235819-11cfia1ht584n1kmk6gh6kbba8ab429u.apps.googleusercontent.com" {
		http.Error(w, err.Error(), 500)
	}

	// Check if user exist in database
	filter := bson.D{
		{"email", resp2.Email},
	}

	var profileInDB common.User
	err = mongodb.Profile.FindOne(context.TODO(), filter).Decode(&profileInDB)

	if err == nil {
		claims := jwt.MapClaims{
			"profileid": profileInDB.ProfileId,
		}
		_, token, err := tokenAuth.Encode(claims)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var response common.Token
		response.Token = token
		render.JSON(w, r, response)

		return
	}

	var result common.User
	result.Country = ""
	result.Email = resp2.Email
	result.Name = resp2.Name
	result.Phone = ""
	result.ProfileId = GenerateProfileId()
	result.Username = _GenerateUsername()

	_, err = mongodb.Profile.InsertOne(context.TODO(), result)

	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	redis.Instance.Set(result.Username, true, 999999999999)

	claims := jwt.MapClaims{
		"profileid": result.ProfileId,
	}

	_, token, err := tokenAuth.Encode(claims)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response common.Token
	response.Token = token
	render.JSON(w, r, response)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req common.User

	json.Unmarshal(b, &req)
	if req.Email == "" || req.Phone == "" || req.Password == "" {
		http.Error(w, fmt.Sprintf("Need email, phone and password for registration. Recieved %s", req), http.StatusInternalServerError)
		return
	}

	req.ProfileId = GenerateProfileId()
	if req.Username == "" {
		req.Username = _GenerateUsername()
	}

	// BIG TODO: Hash Password
	// TODO: Assuming single email, that need not be the case, user can have multiple emails linked to same account
	// For example, registration with a non google email and trying to register later with a google email
	_, err = mongodb.Profile.InsertOne(context.TODO(), req)
	if err != nil {
		http.Error(w, "Registration failed, MongoDB unavailable at the moment", 500)
		return
	}

	redis.Instance.Set(req.Username, true, 999999999999)

	claims := jwt.MapClaims{
		"profileid": req.ProfileId,
	}

	_, token, err := tokenAuth.Encode(claims)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Printf("%s", token)

	var response common.Token
	response.Token = token
	render.JSON(w, r, response)
}
