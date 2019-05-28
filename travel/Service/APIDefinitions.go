package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common/mongodb"
	"../../common/redis"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/api/v1/travel/searchcoordinates", GetLocationFromCoordinates)
	router.Post("/api/v1/travel/searchlocation", GetCoordinatesFromLocation)
	router.Post("/api/v1/travel/getinfo", GetTravelInfo)
	return router
}

func GetLocationFromCoordinates(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	json.Unmarshal(b, &req)

	// Truncing down to two digits
	var latlon = fmt.Sprintf("%.2f%.2f", req.Latitude, req.Longitude)

	res, err := redis.Instance.Get(latlon).Result()

	var resultJSON struct {
		DisplayName string                 `json:"display_name"`
		Address     map[string]interface{} `json:"address"`
		Error       string                 `json:"error"`
	}

	if err == nil {
		fmt.Printf("Serving from location cache\n")
		json.Unmarshal([]byte(res), &resultJSON)
		render.JSON(w, r, resultJSON)
		return
	}

	res2, err2 := http.Get(fmt.Sprintf("https://us1.locationiq.com/v1/reverse.php?key=%s&lat=%f&lon=%f&format=json", "daecd8873d0c8e", req.Latitude, req.Longitude))
	if err2 != nil {
		panic(err2)
	}

	defer res2.Body.Close()

	body, err := ioutil.ReadAll(res2.Body)
	json.Unmarshal([]byte(body), &resultJSON)
	fmt.Printf("%s", resultJSON)

	fmt.Printf("Serving from locationiq server\n")

	if resultJSON.Error == "" {
		redis.Instance.Set(latlon, fmt.Sprintf("%s", body), 99999999999)
	}
	render.JSON(w, r, resultJSON)
}

func GetCoordinatesFromLocation(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req struct {
		Location string `json:"location"`
	}

	var resultJSON []struct {
		Latitude    string `json:"lat"`
		Longitude   string `json:"lon"`
		DisplayName string `json:"display_name"`
	}

	json.Unmarshal(b, &req)

	res, err := redis.Instance.Get(req.Location).Result()

	if err == nil {
		fmt.Printf("Serving from location cache\n")
		json.Unmarshal([]byte(res), &resultJSON)
		render.JSON(w, r, resultJSON)
		return
	}

	res2, err2 := http.Get(fmt.Sprintf("https://us1.locationiq.com/v1/search.php?key=%s&q=%s&format=json", "daecd8873d0c8e", req.Location))
	if err2 != nil {
		panic(err2)
	}

	defer res2.Body.Close()

	body, err := ioutil.ReadAll(res2.Body)

	json.Unmarshal([]byte(body), &resultJSON)

	fmt.Printf("Serving from locationiq server\n")

	var slice [][]byte
	for _, element := range resultJSON {
		res, _ := json.Marshal(element)
		slice = append(slice, res)
	}

	redis.Instance.Set(req.Location, fmt.Sprintf("%s", body), 99999999999)
	render.JSON(w, r, resultJSON)
}

func GetTravelInfo(w http.ResponseWriter, r *http.Request) {
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

	var profileIdString = fmt.Sprintf("%s", req.ProfileId)
	var collection = mongodb.Instance.Database("travel").Collection(profileIdString)

	collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	render.JSON(w, r, result)
}

func SaveTravelInfo(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req SaveTravelInfoType
	json.Unmarshal(b, &req)

	filter := bson.D{{"id", req.TravelId}}
	var profileIdString = fmt.Sprintf("%s", req.ProfileId)

	var collection = mongodb.Instance.Database("travel").Collection(profileIdString)

	collection.UpdateOne(context.TODO(), filter, req.TravelInfo)

	response := make(map[string]string)
	response["OperationStatus"] = "Success"
	render.JSON(w, r, response)
}
