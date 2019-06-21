package travel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common/redis"
	"github.com/go-chi/render"
)

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
