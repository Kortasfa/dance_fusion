package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func roomWSHandler(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		websocketID := vars["id"]
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade websocket connection:", err)
			return
		}
		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {
			}
		}(conn)

		roomWSDict[conn] = websocketID

		_, message, err := conn.ReadMessage() // Чтение названия песни
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			delete(roomWSDict, conn)
			return
		}

		motionListPaths, err := getMotionListPath(db, string(message))
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println(err.Error())
			delete(roomWSDict, conn)
			return
		}
		for roomID, userSlice := range roomIDDict {
			if roomID == websocketID {
				i := 0
				for {
					if i >= len(userSlice) {
						break
					}
					userID := userSlice[i]
					userIDInt, err := strconv.Atoi(userID)
					if err != nil {
						http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
						log.Println(err.Error())
						delete(roomWSDict, conn)
						return
					}
					if userIDInt < 0 {
						i++
						continue
					}
					motionListPath := motionListPaths[i]
					fileContent, err := ioutil.ReadFile(motionListPath)
					if err != nil {
						http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
						log.Println("Ошибка при открытии файла для игрока", userID, ":", err)
						delete(roomWSDict, conn)
						return
					}
					broadcastJoinPageWSMessage <- []string{userID, string(fileContent)}
					i++
				}
				break
			}
		}

		for { // Чтение действия (pause / resume) + end game
			_, message, err = conn.ReadMessage()
			if err != nil {
				delete(roomWSDict, conn)
				break
			}
			for roomID, userSlice := range roomIDDict {
				if roomID == websocketID {
					for _, userID := range userSlice {
						broadcastJoinPageWSMessage <- []string{userID, string(message)}
					}
					break
				}
			}
		}

	}
}

func joinPageWSHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	websocketID := vars["id"]
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade websocket connection:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	joinPageWSDict[conn] = websocketID

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(joinPageWSDict, conn)
			break
		}
	}
}

func neuralWSHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameFieldID := vars["id"]
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 2,
		WriteBufferSize: 1024 * 2,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade websocket connection:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	gameFieldWSDict[conn] = gameFieldID
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(gameFieldWSDict, conn)
			break
		}
	}
}

func handleRoomWSMessages() {
	for mesArr := range broadcastJoiningUserID {
		for wsConnect := range roomWSDict {
			action := mesArr[0]
			roomID := mesArr[1]
			userID := mesArr[2]
			message := ""
			if action == "add" {
				userName := mesArr[3]
				hatSrc := mesArr[4]
				faceSrc := mesArr[5]
				bodySrc := mesArr[6]
				message = action + "|" + userID + "|" + userName + "|" + hatSrc + "|" + faceSrc + "|" + bodySrc
			} else {
				message = action + "|" + userID
			}
			if roomWSDict[wsConnect] == roomID {
				err := wsConnect.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					err := wsConnect.Close()
					delete(roomWSDict, wsConnect)
					if err != nil {
						return
					}
				}
			}
		}
	}
}

func handleJoinPageWSMessages() { // broadcastJoinPageWSMessage <- []string{UserID, Data}
	for mesArr := range broadcastJoinPageWSMessage {
		for wsConnect := range joinPageWSDict {
			userID := mesArr[0]
			data := mesArr[1]
			if joinPageWSDict[wsConnect] == userID {
				err := wsConnect.WriteMessage(websocket.TextMessage, []byte(data))
				if err != nil {
					err := wsConnect.Close()
					delete(joinPageWSDict, wsConnect)
					if err != nil {
						return
					}
				}
			}
		}
	}
}

func danceInfoHandleMessages() {
	for mesArr := range broadcastGameFieldWSMessage {
		found := false
		roomID := mesArr[0]
		data := mesArr[1]
		for {
			for conn, gameFieldID := range gameFieldWSDict {
				if gameFieldID == roomID {
					err := conn.WriteMessage(websocket.TextMessage, []byte(data))
					if err != nil {
						err := conn.Close()
						delete(gameFieldWSDict, conn)
						if err != nil {
							return
						}
					}
					found = true
				}
			}
			if found {
				break
			}
		}

	}
}
