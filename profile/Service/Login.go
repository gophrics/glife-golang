package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../../common"
	"../../common/mongodb"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req LoginUserRequest
	json.Unmarshal(b, &req)

	filter := bson.D{
		{"email", req.Email},
	}

	var profileInDB common.User

	err = mongodb.Profile.FindOne(context.TODO(), filter).Decode(&profileInDB)

	if err != nil {
		http.Error(w, "Login failed, user not found", http.StatusBadRequest)
		return
	}

	if profileInDB.Password != req.Password {
		http.Error(w, "Login failed due to password mismatch", http.StatusBadRequest)
		return
	}

	claim := jwt.MapClaims{"profileid": profileInDB.ProfileId}

	_, token, err := tokenAuth.Encode(claim)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", token)
	var response common.Token
	response.Token = token
	render.JSON(w, r, response)
}

func LoginUserWithGoogle(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req LoginUserWithGoogleRequest

	json.Unmarshal(b, &req)

	// Verify google auth token is valid for our client
	resp, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", req.Token))

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var resp2 GoogleAuthVerification
	json.Unmarshal(body, &resp2)

	fmt.Printf("%s", resp2.Aud)
	if resp2.Aud != common.GOOGLE_APP_ID || resp2.Email != req.Email {
		http.Error(w, "Authentication failure because of, you know, reasons", 500)
	}

	filter := bson.D{
		{"email", resp2.Email},
	}

	var profileInDB common.User
	err = mongodb.Profile.FindOne(context.TODO(), filter).Decode(&profileInDB)

	if err != nil {
		http.Error(w, "Profile donot exist", http.StatusBadRequest)
		return
	}

	// BIG TODO: Use JWT Token - this is hackable just by tampering response, and unsecure

	claims := jwt.MapClaims{
		"profileid": profileInDB.ProfileId,
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
