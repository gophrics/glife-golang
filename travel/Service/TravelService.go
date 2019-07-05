package travel

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo/options"

	"../../common"
	"../../common/mongodb"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	router.Use(middleware.Logger)

	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/api/v1/travel/getalltrips/{username}", GetAllTrips)
		r.Post("/api/v1/travel/gettrip", GetTrip)
		r.Post("/api/v1/travel/savetrip", SaveTrip)
		r.Post("/api/v1/travel/gettriphash", CheckHashPerTrip)
		r.Get("/api/v1/travel/search/{searchstring}", Search)
		r.Post("/api/v1/travel/toggletripprivacy", TogglePrivacy)
	})

	router.Group(func(r chi.Router) {
		r.Post("/api/v1/travel/searchcoordinates", GetLocationFromCoordinates)
		r.Post("/api/v1/travel/searchlocation", GetCoordinatesFromLocation)
		r.Post("/api/v1/travel/searchweatherbylocation", GetWeatherByLocation)
	})

	return router
}

func TogglePrivacy(w http.ResponseWriter, r *http.Request) {
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

	type RequestModal struct {
		TripInfo
		Value bool
	}

	var request RequestModal

	json.Unmarshal(b, &request)

	filter := bson.D{{"tripid", request.TripInfo.TripId}, {"profileid", profileId}}

	update := bson.D{{"$set", bson.D{{"likedby", request.Value}}}}

	_, err = mongodb.Travel.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, nil)
}

/*
METHOD: POST
Body: TripInfo
		TripId
*/
func GetTrip(w http.ResponseWriter, r *http.Request) {
	_, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		fmt.Printf("%s", err2.Error())
		return
	}
	var profileId = fmt.Sprintf("%s", claims["profileid"])
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tripData TripInfo
	json.Unmarshal(b, &tripData)

	if tripData.ProfileId != "" {
		profileId = tripData.ProfileId
	}
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
	token, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	var profileId = fmt.Sprintf("%s", claims["profileid"])
	var myOwn = true
	// Check if Profile is there in the parameter
	ss := fmt.Sprintf("%s", chi.URLParam(r, "username"))

	if ss != "" {

		req, err := http.NewRequest("GET", fmt.Sprintf("http://profile:8080/api/v1/profile/getuserwithusername/%s", ss), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%s", token)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Raw))
		resp, err := (&http.Client{}).Do(req)

		fmt.Printf("%s", resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var user common.User
		json.Unmarshal(b, &user)

		profileId = user.ProfileId
		myOwn = false
	}

	fmt.Printf("\n%s\n", profileId)

	var result []Trip

	filter := bson.D{{"profileid", profileId}}

	cur, err := mongodb.Travel.Find(context.TODO(), filter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for cur.Next(context.TODO()) {
		var elem Trip
		err := cur.Decode(&elem)
		fmt.Printf("\n%s\n", elem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if myOwn || elem.Public {
			result = append(result, elem)
		}
	}

	fmt.Printf("\n%s\n", result)

	render.JSON(w, r, result)
}

func SaveTrip(w http.ResponseWriter, r *http.Request) {
	_, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	var profileId string = fmt.Sprintf("%s", claims["profileid"])

	fmt.Printf("%s\n", profileId)
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
		{"profileid", profileId},
	}

	var TripUpdateFilter = bson.D{
		{"$set", bson.D{
			{"tripid", req.TripId},
			{"profileid", profileId},
			{"tripname", req.TripName},
			{"steps", req.Steps},
			{"public", req.Public},
			{"masterimage", req.MasterImage},
			{"startdate", req.StartDate},
			{"enddate", req.EndDate},
			{"temperature", req.Temperature},
			{"countrycode", req.CountryCode},
			{"daysoftravel", req.DaysOfTravel},
			{"activities", req.Activities},
			{"location", req.Location},
		}},
	}

	var upsert bool = true
	options := &options.UpdateOptions{Upsert: &upsert}
	_, err = mongodb.Travel.UpdateOne(context.TODO(), filter, TripUpdateFilter, options)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := make(map[string]string)
	response["OperationStatus"] = "Success"
	render.JSON(w, r, response)
}

func Search(w http.ResponseWriter, r *http.Request) {
	ss := fmt.Sprintf("%s", chi.URLParam(r, "searchstring"))

	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"tripname", bson.D{{"$regex", primitive.Regex{Pattern: ".*" + ss + ".*", Options: "i"}}}}},
				bson.D{{"stepname", bson.D{{"$regex", primitive.Regex{Pattern: ".*" + ss + ".*", Options: "i"}}}}},
			}},
	}

	var result []TripMeta

	cur, err := mongodb.Travel.Find(context.TODO(), filter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var x TripMeta
	for cur.Next(context.TODO()) {
		err := cur.Decode(&x)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result = append(result, x)
	}

	render.JSON(w, r, result)
}

// @title CheckHashPerTrip
// @version 1.0
// @description This take tripid, will compute hash value of a trip and send back response
// @BasePath /api/v1
// @Body: { "tripId": ""}
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

	// Fix hash mismatch
	result.Hash = fmt.Sprintf("%x", Hash(trip))

	render.JSON(w, r, result)
}
