package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type Person struct {
	Name string
	Age  int
}

func main() {

	indexFile, err := os.Open("Tests/websocket.html")
	if err != nil {
		fmt.Println(err)
	}
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		fmt.Println(err)
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Client subscribed")
		var myPerson = Person{"Bill", 21}
		for {
			time.Sleep(2 * time.Second)
			if myPerson.Age < 100 {
				myJson, err := json.Marshal(myPerson)
				if err != nil {
					fmt.Println(err)
					return
				}
				err = conn.WriteMessage(websocket.TextMessage, myJson)
				if err != nil {
					fmt.Println(err)
					break
				}
				myPerson.Age += 2
			} else {
				conn.Close()
				break
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})
	http.ListenAndServe(":3000", nil)
	fmt.Println("Client unsubscribed")
}
