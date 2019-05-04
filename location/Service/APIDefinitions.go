package location

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common/redis"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/api/v1/location/updatelocation", UpdateLocation)
	router.Get("/api/v1/location/nearme", NearMe)
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

	/*
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	*/

	var req Location

	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024) // error ignored for sake of simplicity

	if err != nil {
		fmt.Printf("Cannot upgrade to websocket " + err.Error())
		return
	}

	// json.Unmarshal(b, &req)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil { // Close connection
			return
		}
		json.Unmarshal(msg, &req)

		res, err := redis.Instance.GeoRadius("LastKnown", req.latitude, req.longitude, &redis.GeoRadiusQuery{
			// Update Radius! WTF
			Radius:    10000,
			WithCoord: true,
		}).Result()

		for _, element := range res {
			profile := &Location{
				profileId: element.Name,
				latitude:  element.Latitude,
				longitude: element.Longitude,
			}
			i, err := json.Marshal(profile)
			if err != nil {
				return
			}
			if err = conn.WriteMessage(msgType, i); err != nil {
				return
			}
		}
	}
}
