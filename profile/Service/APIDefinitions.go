package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../../common/mongodb"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/api/v1/profile/register", RegisterUser)
	router.Post("/api/v1/profile/registerWithGoogle", RegisterUserWithGoogle)
	router.Post("/api/v1/profile/loginWithGoogle", LoginUserWithGoogle)
	router.Get("/api/v1/profile/search/{searchstring}", FindUser)
	router.Get("/api/v1/profile/getuser/{profileId}", GetUser)
	return router
}

func NewProfileId() primitive.ObjectID {
	return primitive.NewObjectID()
}

func LoginUserWithGoogle(w http.ResponseWriter, r *http.Request) {

}

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

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("1")
	var resp2 GoogleAuthVerification
	json.Unmarshal(body, &resp2)

	fmt.Printf("%s", resp2.Aud)
	if resp2.Aud != "249369235819-11cfia1ht584n1kmk6gh6kbba8ab429u.apps.googleusercontent.com" {
		http.Error(w, err.Error(), 500)
	}

	fmt.Printf("3\n", resp2)
	// Check if user exist in database
	filter := bson.D{
		{"email", resp2.Email},
	}

	var profileInDB GetUserResponse
	err = mongodb.Profile.FindOne(context.TODO(), filter).Decode(&profileInDB)

	fmt.Printf("4\n", resp2)
	if err == nil {
		http.Error(w, "Profile already exist", 500)
	}

	var result RegisterUserResponse
	result.Country = ""
	result.Email = resp2.Email
	result.Name = resp2.Name
	result.Phone = ""
	result.ProfileId = primitive.NewObjectID()

	insertResult, err := mongodb.Profile.InsertOne(context.TODO(), result)

	fmt.Printf("5\n", resp2)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	fmt.Println("Inserted document from google auth: ", insertResult)

	response := make(map[string]string)
	response["OperationStatus"] = "Success"
	response["Result"] = fmt.Sprintf("%s", result)
	render.JSON(w, r, response)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req RegisterUserRequest

	json.Unmarshal(b, &req)

	var result RegisterUserResponse
	result.Country = req.Country
	result.Email = req.Email
	result.Name = req.Name
	result.Phone = req.Phone
	result.ProfileId = primitive.NewObjectID()

	insertResult, err := mongodb.Profile.InsertOne(context.TODO(), result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	response := make(map[string]string)
	response["OperationStatus"] = "Success"
	response["Result"] = fmt.Sprintf("%s", result)
	render.JSON(w, r, response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// If struct not initialzed, inner variables don't exist

	profileId := fmt.Sprintf("%s", chi.URLParam(r, "profileId"))

	filter := bson.D{{
		"_id",
		bson.D{{
			"$in",
			bson.A{profileId},
		}},
	}}

	var result GetUserResponse

	err := mongodb.Profile.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	render.JSON(w, r, result)
}

func FindUser(w http.ResponseWriter, r *http.Request) {

	ss := fmt.Sprintf("%s", chi.URLParam(r, "searchstring"))
	ss = fmt.Sprintf(".*%s.*", ss)

	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"name", primitive.Regex{ss, ""}}},
			bson.D{{"country", ss}},
			bson.D{{"phone", ss}},
		}}}
	findOptions := options.Find()

	var result []GetUserResponse

	cur, err := mongodb.Profile.Find(context.TODO(), filter, findOptions)

	if err != nil {
		log.Fatal(err)
	}

	var x GetUserResponse
	for cur.Next(context.TODO()) {
		err := cur.Decode(&x)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", x)
		result = append(result, x)
	}

	render.JSON(w, r, result)

}
