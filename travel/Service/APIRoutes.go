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
	"github.com/go-chi/chi/middleware"
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

	router.Use(middleware.Recoverer)

	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/api/v1/travel/getalltrips/{profileid}", GetAllTrips)
		r.Post("/api/v1/travel/gettrip", GetTrip)
		r.Post("/api/v1/travel/savetrip", SaveTrip)
		r.Post("/api/v1/travel/gettriphash", CheckHashPerTrip)
	})

	router.Group(func(r chi.Router) {
		r.Post("/api/v1/travel/searchcoordinates", GetLocationFromCoordinates)
		r.Post("/api/v1/travel/searchlocation", GetCoordinatesFromLocation)
		r.Post("/api/v1/travel/searchweatherbylocation", GetWeatherByLocation)
	})

	return router
}

func GetTrip(w http.ResponseWriter, r *http.Request) {
	_, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		fmt.Printf("%s", err2.Error())
		return
	}
	var profileId = claims["profileid"]
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tripData TripInfo
	json.Unmarshal(b, &tripData)

	var result Trip

	fmt.Printf("%s %s", tripData.TripId, profileId)
	filter := bson.D{{"tripid", tripData.TripId}, {"profileid", profileId}}

	err = mongodb.Travel.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, result)
}

func GetAllTrips(w http.ResponseWriter, r *http.Request) {
	_, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	var profileId = claims["profileid"]
	var myOwn = true
	// Check if Profile is there in the parameter
	ss := fmt.Sprintf("%s", chi.URLParam(r, "profileid"))

	if ss != "" {
		profileId = ss
		myOwn = false
	}

	var result []Trip

	filter := bson.D{{"profileid", profileId}}

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
		if myOwn || elem.Public {
			result = append(result, elem)
		}
	}

	render.JSON(w, r, result)
}

func SaveTrip(w http.ResponseWriter, r *http.Request) {
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
		{"profileid", req.ProfileId},
	}

	_, err = mongodb.Travel.UpdateOne(context.TODO(), filter, req)

	if err != nil {
		_, err2 := mongodb.Travel.InsertOne(context.TODO(), req)
		if err2 != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	filter := bson.D{{"tripId", req.TripId}, {"profileid", profileId}}

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
