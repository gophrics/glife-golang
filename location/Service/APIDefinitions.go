package location

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common/redis"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/api/v1/location/updatelocation", UpdateLocation)
	router.Get("/api/v1/location/nearme/{latitude}/{longitude}", NearMe)
	return router
}

func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req Location

	json.Unmarshal(b, &req)

	res, err := redis.Instance.GeoAdd("LastKnown", &redis.GeoLocation{
		Latitude:  req.latitude,
		Longitude: req.longitude,
		Name:      req.profileId,
	}).Result()

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	var response map[string]string
	response["OperationStatus"] = "Success"
	response["Result"] = fmt.Sprintf("%s", res)

	render.JSON(w, r, response)
}

func NearMe(w http.ResponseWriter, r *http.Request) {

}
