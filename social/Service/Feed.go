package social

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

/*
	Method: GET
	params: /{profileid}
	returns Array<TripModal>

	Currently returns all trips of all users following. Which is bad ?
*/
func GetFeeds(w http.ResponseWriter, r *http.Request) {
	token, claims, err2 := jwtauth.FromContext(r.Context())

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
	var result []Trip

	client := &http.Client{}
	for _, el := range followingList {

		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8082/api/v1/travel/getalltrip/%s", el), nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Set("Authorization", token.Signature)
		res, err := client.Do(req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var trips []Trip
		json.Unmarshal(b, &trips)

		result = append(result, trips...)
	}

	render.JSON(w, r, result)
}
