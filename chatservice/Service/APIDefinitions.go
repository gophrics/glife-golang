package location

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../../common/mysql"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/api/v1/chat/nearme", ChatWebsocketAPI)
	return router
}

func ChatWebsocketAPI(w http.ResponseWriter, r *http.Request) {

	/*
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	*/

	var req ChatMessage

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

		prep, _ := mysql.Instance.Prepare("INSERT INTO ? VALUES (?,?,?)")
		res, err := prep.Exec(req.ChatroomID, req.ProfileId, req.Timestamp, req.Message)

		if res == nil || err == nil {
		}
	}
}
