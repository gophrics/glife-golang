package profile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../common/mysql"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/api/v1/profile/register", RegisterUser)
	router.Get("/api/v1/profile/getuser/{profileId}", GetUser)
	return router
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req User

	json.Unmarshal(b, &req)
	prep, _ := mysql.Instance.Prepare("INSERT INTO User VALUES (?,?,?,?)")
	res, err := prep.Exec(req.ProfileId, req.Name, req.Phone, req.Country)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	response := make(map[string]string)
	response["OperationStatus"] = "Success"
	response["Result"] = fmt.Sprintf("%s", res)
	render.JSON(w, r, response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// If struct not initialzed, inner variables don't exist

	profileId := chi.URLParam(r, "profileId")

	rows, err := mysql.Instance.Query("SELECT * from User where Id=?", profileId)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	var user User
	for rows.Next() {
		rows.Scan(&user.ProfileId, &user.Name, &user.Phone, &user.Country)
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	render.JSON(w, r, user)
}
