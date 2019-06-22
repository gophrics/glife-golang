package social

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common/mongodb"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
)

/*
	The way it is, anyone can create an element right now, which is not very recommended, since database can be populated with all kinds of trip ids by random people querying stuff
	Also, since there's a way to update a sub element of mongodb, we should not probably flatten out the structure, and rather keep the structure as in travel service
*/
func Like(w http.ResponseWriter, r *http.Request) {
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

	var req LikeRequest
	json.Unmarshal(b, &req)

	var profileIdToLike = req.ProfileId
	if req.TripOrStep == 1 {
		var step Step
		filter := bson.D{{"tripid", req.TripId}, {"profileid", profileIdToLike}, {"stepId", req.StepId}}

		err = mongodb.Social.FindOne(context.TODO(), filter).Decode(&step)
		step.LikedBy, err2 = AppendIfNotExists(step.LikedBy, fmt.Sprintf("%s", profileId))

		if err2 != nil {
			render.JSON(w, r, err2.Error())
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			// step.ProfileId = profileIdToLike
			// step.TripId = req.TripId
			// step.StepId = req.StepId
			// _, err2 := mongodb.Social.InsertOne(context.TODO(), step)
			// if err2 != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }
			// render.JSON(w, r, "Success")
			// return
		}

		update := bson.D{{"$set", bson.D{{"likedby", step.LikedBy}}}}
		_, err = mongodb.Social.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		var trip Trip
		filter := bson.D{{"tripid", req.TripId}, {"profileid", profileIdToLike}}

		err = mongodb.Social.FindOne(context.TODO(), filter).Decode(&trip)
		trip.LikedBy, err2 = AppendIfNotExists(trip.LikedBy, fmt.Sprintf("%s", profileId))

		if err2 != nil {
			render.JSON(w, r, err2.Error())
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			// trip.ProfileId = profileIdToLike
			// trip.TripId = req.TripId
			// _, err2 := mongodb.Social.InsertOne(context.TODO(), trip)
			// if err2 != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }
			// render.JSON(w, r, "Success")
			// return
		}

		update := bson.D{{"$set", bson.D{{"likedby", trip.LikedBy}}}}
		_, err = mongodb.Social.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	render.JSON(w, r, "Success")
}

/*
	The way it is, anyone can create an element right now, which is not very recommended, since database can be populated with all kinds of trip ids by random people querying stuff
	Also, since there's a way to update a sub element of mongodb, we should not probably flatten out the structure, and rather keep the structure as in travel service
*/
func Unlike(w http.ResponseWriter, r *http.Request) {
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

	var req LikeRequest
	json.Unmarshal(b, &req)

	var profileIdToLike = req.ProfileId
	if req.TripOrStep == 1 {
		var step Step
		filter := bson.D{{"tripid", req.TripId}, {"profileid", profileIdToLike}, {"stepId", req.StepId}}

		err = mongodb.Social.FindOne(context.TODO(), filter).Decode(&step)
		step.LikedBy, err2 = RemoveIfExists(step.LikedBy, fmt.Sprintf("%s", profileId))

		if err2 != nil {
			render.JSON(w, r, err2.Error())
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			// step.ProfileId = profileIdToLike
			// step.TripId = req.TripId
			// step.StepId = req.StepId
			// _, err2 := mongodb.Social.InsertOne(context.TODO(), step)
			// if err2 != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }
			// render.JSON(w, r, "Success")
			// return
		}

		update := bson.D{{"$set", bson.D{{"likedby", step.LikedBy}}}}
		_, err = mongodb.Social.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		var trip Trip
		filter := bson.D{{"tripid", req.TripId}, {"profileid", profileIdToLike}}

		err = mongodb.Social.FindOne(context.TODO(), filter).Decode(&trip)
		trip.LikedBy, err2 = RemoveIfExists(trip.LikedBy, fmt.Sprintf("%s", profileId))

		if err2 != nil {
			render.JSON(w, r, err2.Error())
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			// trip.ProfileId = profileIdToLike
			// trip.TripId = req.TripId
			// _, err2 := mongodb.Social.InsertOne(context.TODO(), trip)
			// if err2 != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }
			// render.JSON(w, r, "Success")
			// return
		}

		update := bson.D{{"$set", bson.D{{"likedby", trip.LikedBy}}}}
		_, err = mongodb.Social.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	render.JSON(w, r, "Success")
}
