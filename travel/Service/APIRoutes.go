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

		router.Get("/api/v1/travel/getalltrips", GetAllTrips)
		router.Post("/api/v1/travel/gettrip", GetTrip)
		router.Post("/api/v1/travel/gettriphash", CheckHashPerTrip)
	})

	router.Group(func(r chi.Router) {
		router.Post("/api/v1/travel/searchcoordinates", GetLocationFromCoordinates)
		router.Post("/api/v1/travel/searchlocation", GetCoordinatesFromLocation)
		router.Post("/api/v1/travel/searchweatherbylocation", GetWeatherByLocation)
	})

	return router
}

func GetTrip(w http.ResponseWriter, r *http.Request) {
	_, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		fmt.Printf("%s", err2.Error())
		return
	}
	var profileId string = fmt.Sprintf("%s", claims["profileid"])
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var tripData TripInfo
	json.Unmarshal(b, &tripData)

	var result Trip

	filter := bson.D{{"tripId", tripData.TripId}, {"profileId", profileId}}

	err = mongodb.Travel.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	render.JSON(w, r, result)
}

func GetAllTrips(w http.ResponseWriter, r *http.Request) {
	_, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		fmt.Printf("%s", err2.Error())
		return
	}

	var profileId string = fmt.Sprintf("%s", claims["profileid"])

	var result []Trip

	filter := bson.D{{"profileId", profileId}}

	cur, err := mongodb.Travel.Find(context.TODO(), filter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	for cur.Next(context.TODO()) {
		var elem Trip
		err := cur.Decode(&elem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		result = append(result, elem)
	}

	render.JSON(w, r, result)
}

func SaveTravelInfo(w http.ResponseWriter, r *http.Request) {
	_, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	var profileId string = fmt.Sprintf("%s", claims["profileid"])

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req Trip
	json.Unmarshal(b, &req)
	req.ProfileId = profileId

	filter := bson.D{
		{"tripId", req.TripId},
		{"profileId", req.ProfileId},
	}

	_, err = mongodb.Travel.UpdateOne(context.TODO(), filter, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response := make(map[string]string)
	response["OperationStatus"] = "Success"
	render.JSON(w, r, response)
}

func CheckHashPerTrip(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var profileId = claims["profileid"]

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req TripInfo
	json.Unmarshal(b, &req)
	filter := bson.D{{"tripId", req.TripId}, {"profileId", profileId}}

	var trip Trip
	mongodb.Travel.FindOne(context.TODO(), filter).Decode(&trip)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	var result struct {
		Hash string
	}

	result.Hash = fmt.Sprintf("%s", Hash(trip))

	render.JSON(w, r, result)
}
