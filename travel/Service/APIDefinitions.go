package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common/mongodb"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/api/v1/travel/getinfo", GetTravelInfo)
	return router
}

func GetTravelInfo(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req TravelInfo
	var result []TravelResponse

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
