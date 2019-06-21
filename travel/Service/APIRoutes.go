package travel

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common"
	"../../common/mongodb"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(common.JWT_SECRET_KEY), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		router.Post("/api/v1/travel/getinfo", GetTravelInfo)
	})

	router.Group(func(r chi.Router) {
		router.Post("/api/v1/travel/searchcoordinates", GetLocationFromCoordinates)
		router.Post("/api/v1/travel/searchlocation", GetCoordinatesFromLocation)
		router.Post("/api/v1/travel/searchweatherbylocation", GetWeatherByLocation)
	})

	return router
}

func GetTravelInfo(w http.ResponseWriter, r *http.Request) {
	_, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		fmt.Printf("%s", err2.Error())
		return
	}

	var username string = fmt.Sprintf("%s", claims["username"])

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req TravelInfo
	var result []TravelData

	json.Unmarshal(b, &req)
	filter := bson.D{{"id", req.TravelId}}

	var collection = mongodb.Instance.Database("travel").Collection(username)

	collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	render.JSON(w, r, result)
}

func SaveTravelInfo(w http.ResponseWriter, r *http.Request) {
	_, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		fmt.Printf("%s", err2.Error())
		return
	}

	var username string = fmt.Sprintf("%s", claims["username"])

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req SaveTravelInfoType
	json.Unmarshal(b, &req)

	filter := bson.D{{"id", req.TravelId}}

	var collection = mongodb.Instance.Database("travel").Collection(username)

	collection.UpdateOne(context.TODO(), filter, req.TravelInfo)

	response := make(map[string]string)
	response["OperationStatus"] = "Success"
	render.JSON(w, r, response)
}
