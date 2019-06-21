package travel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common/redis"
	"github.com/go-chi/render"
)

func GetWeatherByLocation(w http.ResponseWriter, r *http.Request) {
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
	// w for weather
	var latlon = fmt.Sprintf("%.2f%.2fw", req.Latitude, req.Longitude)

	res, err := redis.Instance.Get(latlon).Result()

	var resultJSON map[string]interface{}

	if err == nil {
		fmt.Printf("Serving from location cache\n")
		json.Unmarshal([]byte(res), &resultJSON)
		render.JSON(w, r, resultJSON)
		return
	}

	fmt.Printf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&APPID=%s", req.Latitude, req.Longitude, "b4d2d441850a65f8ae5bbb81492d6125")
	res2, err2 := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&APPID=%s", req.Latitude, req.Longitude, "b4d2d441850a65f8ae5bbb81492d6125"))
	if err2 != nil {
		panic(err2)
	}

	defer res2.Body.Close()

	body, err := ioutil.ReadAll(res2.Body)

	if err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(body), &resultJSON)

	redis.Instance.Set(latlon, fmt.Sprintf("%s", body), 99999999999)
	render.JSON(w, r, resultJSON)
}
