package social

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common"
	neo4jd "../../common/neo4j"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/api/v1/social/follow", Follow)
		r.Post("/api/v1/social/unfollow", Unfollow)
		r.Post("/api/v1/social/like", Like)
		r.Post("/api/v1/social/unlike", Unlike)
		r.Get("/api/v1/social/feed", GetFeeds)
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

func Comment(w http.ResponseWriter, r *http.Request) {

}

// UnImplemented
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
	_, claims, err2 := jwtauth.FromContext(r.Context())

	if err2 != nil {
		fmt.Printf("%s", err2.Error())
		return
	}

	var profileId = fmt.Sprintf("%s", claims["profileid"])

	followingList, err := _GetFollowingList(profileId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, followingList)
}
