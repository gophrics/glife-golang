package social

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common"
	"../../common/dgraph"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/dgraph-io/dgo/y"
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

	var socialData struct {
		Uid       string   `json:"profileId"`
		Following []string `json:"following"`
	}
	json.Unmarshal(b, &socialData)

	socialData.Uid = fmt.Sprintf("%s", profileId)

	encodedSocialData, err := json.Marshal(socialData)

	mu := &api.Mutation{
		SetJson: encodedSocialData,
	}

	txn := dgraph.Instance.NewTxn()
	_, err = txn.Mutate(context.TODO(), mu)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = txn.Commit(context.TODO())
	if err == y.ErrAborted {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	defer txn.Discard(context.TODO())
}

func Like(w http.ResponseWriter, r *http.Request) {

}

func Comment(w http.ResponseWriter, r *http.Request) {

}

func GetFollowersList(w http.ResponseWriter, r *http.Request) {

}

func GetFollowingList(w http.ResponseWriter, r *http.Request) {

}
