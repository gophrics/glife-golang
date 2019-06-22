package social

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common"
	"../../common/mongodb"
	neo4jd "../../common/neo4j"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/api/v1/social/follow", Follow)
		r.Post("/api/v1/social/like", Like)
	})
	return router
}

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(common.JWT_SECRET_KEY), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

/*
MERGE (n {name: '3'}) //Create if a node with name='3' does not exist else match it
MERGE (test2 {name:'2'}) //Create if a node with name='2' does not exist else match it
MERGE (n)-[:know {r:'123'}]->(test2) //Create the relation between these nodes if it does not already exist
*/

func Follow(w http.ResponseWriter, r *http.Request) {

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

	var followProfiles FollowRequest

	json.Unmarshal(b, &followProfiles)

	for _, followProfile := range followProfiles.Following {
		_, err = neo4jd.Session.Run(`	MERGE(a: Person {profileId: $myid})
										MATCH(m: Person {profileId: $id})
										MERGE(a)-[:relationshipName]->(m)`, map[string]interface{}{
			"id":   followProfile,
			"myid": profileId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // handle error
		}
	}

	render.JSON(w, r, "")
}

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

func Comment(w http.ResponseWriter, r *http.Request) {

}

func GetFollowersList(w http.ResponseWriter, r *http.Request) {
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

	var followProfiles FollowRequest

	json.Unmarshal(b, &followProfiles)

	for _, followProfile := range followProfiles.Following {
		_, err = neo4jd.Session.Run(`	MERGE(a: Person {profileId: $myid})
										MATCH(m: Person {profileId: $id})
										MERGE(a)-[:relationshipName]->(m)`, map[string]interface{}{
			"id":   followProfile,
			"myid": profileId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // handle error
		}
	}

	render.JSON(w, r, "")
}

func GetFollowingList(w http.ResponseWriter, r *http.Request) {

}
