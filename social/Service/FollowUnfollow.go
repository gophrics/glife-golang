package social

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	neo4jd "../../common/neo4j"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

/*
MERGE (n {name: '3'}) //Create if a node with name='3' does not exist else match it
MERGE (test2 {name:'2'}) //Create if a node with name='2' does not exist else match it
MERGE (n)-[:know {r:'123'}]->(test2) //Create the relation between these nodes if it does not already exist
*/

/*
	METHOD: POST
	REQ: FollowRequest {
		Following []string
	}
	RESPONSE ""
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

func Unfollow(w http.ResponseWriter, r *http.Request) {

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
		_, err = neo4jd.Session.Run(`	(a: Person {profileId: $myid})
										(m: Person {profileId: $id})
										MATCH(a)-[r:relationshipName]->(m)
										DELETE r`, map[string]interface{}{
			"id":   followProfile,
			"myid": profileId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // handle error
		}
	}

	render.JSON(w, r, "")
}

func _GetFollowingList(profileId string) ([]string, error) {
	dbresult, err := neo4jd.Session.Run(`	(a: Person {profileId: $id})
										MATCH(a)-[r:relationshipName]->(m)
										RETURN m`, map[string]interface{}{
		"id": profileId,
	})

	if err != nil {
		return nil, err
	}

	var i = 0
	var result []string
	for dbresult.Next() {
		friendprofileid, ok := dbresult.Record().GetByIndex(i).(string)
		if ok != true {
			return nil, errors.New("There was a problem trying to get the followers")
		}
		result = append(result, friendprofileid)
		i += 1
	}

	return result, nil
}
