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
	router.Post("/api/v1/profile/search", FindUser)
	router.Get("/api/v1/profile/getuser/{profileId}", GetUser)
	return router
}

func NewProfileId() string {
	return "1"
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
	result.ProfileId = NewProfileId()

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
	ss := fmt.Sprintf(".*%s.*", "Nitin")

	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"name", primitive.Regex{ss, ""}}},
			bson.D{{"country", ss}},
			bson.D{{"phone", ss}},
		}}}
	findOptions := options.Find()

	var result GetUserResponse

	cur, err := mongodb.Profile.Find(context.TODO(), filter, findOptions)

	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", result)
	}

	render.JSON(w, r, "")

}
