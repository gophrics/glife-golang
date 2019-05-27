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
	router.Get("/api/v1/travel/searchlocation", GetLocationFromCoordinates)
	router.Get("/api/v1/travel/searchcoordinates", GetCoordinatesFromLocation)
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

	fmt.Println(latlon)

	res, err := redis.Instance.Get(latlon).Result()

	if err == nil {
		fmt.Printf("Serving from location cache\n")
		render.JSON(w, r, res)
		return
	}

	res2, err2 := http.Get(fmt.Sprintf("https://us1.locationiq.com/v1/reverse.php?key=%s&lat=%f&lon=%f&format=json", "daecd8873d0c8e", req.Latitude, req.Longitude))
	if err2 != nil {
		panic(err2)
	}

	defer res2.Body.Close()

	body, err := ioutil.ReadAll(res2.Body)

	var resultJSON struct {
		DisplayName string                 `json:"display_name"`
		Address     map[string]interface{} `json:"address"`
	}

	json.Unmarshal([]byte(body), &resultJSON)

	fmt.Printf("Serving from locationiq server\n")

	redis.Instance.Set(latlon, resultJSON.DisplayName, 99999999999)

	render.JSON(w, r, resultJSON.DisplayName)
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

	json.Unmarshal(b, &req)

	res, err := redis.Instance.Get(req.Location).Result()

	if err == nil {
		fmt.Printf("Serving from location cache\n")
		render.JSON(w, r, res)
		return
	}

	res2, err2 := http.Get(fmt.Sprintf("https://us1.locationiq.com/v1/search.php?key=%s&q=%s&format=json", "daecd8873d0c8e", req.Location))
	if err2 != nil {
		panic(err2)
	}

	defer res2.Body.Close()

	body, err := ioutil.ReadAll(res2.Body)

	//fmt.Printf("%s", body)

	var resultJSON []struct {
		Latitude    string `json:"lat"`
		Longitude   string `json:"lon"`
		DisplayName string `json:"display_name"`
	}

	json.Unmarshal([]byte(body), &resultJSON)

	fmt.Printf("Serving from locationiq server\n")

	var slice [][]byte
	for _, element := range resultJSON {
		res, _ := json.Marshal(element)
		slice = append(slice, res)
	}

	redis.Instance.Set(req.Location, fmt.Sprintf("%s", slice), 99999999999)
	render.JSON(w, r, fmt.Sprintf("%s", slice))
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
