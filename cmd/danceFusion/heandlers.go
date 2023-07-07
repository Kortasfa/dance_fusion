package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type buttonRequest struct {
	IsBtnClicked string `json:"button"`
}

/*func gameFieldWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade websocket connection:", err)
		return
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte(""))
	if err != nil {
		log.Println("Failed to send welcome message:", err)
	}
}*/

/*func handleMessages() {
	for _ = range broadcast {
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte("True"))
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}*/

func gameFieldWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade websocket connection:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			break
		}
	}
}

func gameField(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/gameField.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	data := struct {
		WssURL string
	}{
		WssURL: "wss://" + r.Host + "/start/ws",
	}
	err = ts.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func gamePhone(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/gamePhone.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func checkButton(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte("True"))
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

/*func start(w http.ResponseWriter, r *http.Request) {
	emp := &buttonRequest{isBtnClicked: true}
	e, err := json.Marshal(emp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))

	w.WriteHeader(200)
}*/
