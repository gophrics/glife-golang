package social

import (
	"fmt"
	"net/http"

	neo4jd "../../common/neo4j"

	"../../common"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/api/v1/social/follow", Follow)
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

func Follow(w http.ResponseWriter, r *http.Request) {
	/*
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
	*/

	result, err := neo4jd.Session.Run("MERGE (n {name: $name})", map[string]interface{}{
		"name": "Item 2",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // handle error
		return
	}

	for result.Next() {
		fmt.Printf("Created Item with Name = '%s'\n", result.Record().GetByIndex(0).(string))
	}
	if err = result.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func Like(w http.ResponseWriter, r *http.Request) {

}

func Comment(w http.ResponseWriter, r *http.Request) {

}

func GetFollowersList(w http.ResponseWriter, r *http.Request) {

}

func GetFollowingList(w http.ResponseWriter, r *http.Request) {

}
