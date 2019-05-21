package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"../../common/mysql"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/api/v1/chat/getandpost", ChatWebsocketAPI)
	return router
}

func ChatWebsocketAPI(w http.ResponseWriter, r *http.Request) {

	var req ConnectionMessage

	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024) // error ignored for sake of simplicity

	if err != nil {
		fmt.Printf("Cannot upgrade to websocket " + err.Error())
		return
	}

	msgType, msg, err := conn.ReadMessage()

	if err != nil { // Close connection
		return
	}
	json.Unmarshal(msg, &req)
	lastUpdateTimestamp := req.Timestamp

	sqlQuery := fmt.Sprintf("SELECT * from %s where Timestamp > '%s'", req.ChatroomID, lastUpdateTimestamp)
	rows, err := mysql.Instance.Query(sqlQuery)

	if err != nil {
		panic(err)
		return
	}

	var message ChatMessage
	for rows.Next() {
		rows.Scan(&message.ProfileId, &message.Timestamp, &message.Message)

		i, err := json.Marshal(message)
		fmt.Printf("%s", message)
		if err != nil {
			panic(err)
			return
		}

		if err = conn.WriteMessage(msgType, i); err != nil {
			panic(err)
			return
		}
	}
	lastUpdateTimestamp = message.Timestamp

	var messageReq ChatMessage

	go func() {
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil { // Close connection
				return
			}
			json.Unmarshal(msg, &messageReq)

			sqlQuery := fmt.Sprintf("SELECT * from %s where Timestamp > '%s'", messageReq.ChatroomID, lastUpdateTimestamp)
			rows, err := mysql.Instance.Query(sqlQuery)

			if err != nil {
				panic(err)
				return
			}

			var message ChatMessage
			for rows.Next() {
				rows.Scan(&message.ProfileId, &message.Timestamp, &message.Message)

				i, err := json.Marshal(message)

				if err != nil {
					panic(err)
					return
				}

				if err = conn.WriteMessage(msgType, i); err != nil {
					panic(err)
					return
				}
			}
			lastUpdateTimestamp = message.Timestamp
			time.Sleep(500)
		}
	}()
	for {

		msgType, msg, err := conn.ReadMessage()
		if err != nil { // Close connection
			return
		}
		json.Unmarshal(msg, &messageReq)

		prep, _ := mysql.Instance.Prepare("INSERT INTO ? VALUES (?,?,?)")
		res, err := prep.Exec(messageReq.ChatroomID, messageReq.ProfileId, messageReq.Timestamp, messageReq.Message)

		if err != nil {
			return
		}

		i, err := json.Marshal(res)

		if err = conn.WriteMessage(msgType, i); err != nil {
			return
		}

		lastUpdateTimestamp = req.Timestamp
	}
}
