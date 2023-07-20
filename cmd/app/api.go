package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getJoinedUserData(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Parsing error", 500)
			log.Println(err.Error())
			return
		}
		var data struct {
			UserID int
			RoomID string
		}
		err = json.Unmarshal(reqData, &data)
		if err != nil {
			http.Error(w, "JSON parsing error", 500)
			log.Println(err.Error())
			return
		}
		_, exists := roomIDDict[data.RoomID]
		if !exists {
			w.WriteHeader(404)
			return
		}
		currRoomID, found := retrieveUserRoom(fmt.Sprintf("%d", data.UserID))
		if len(roomIDDict[data.RoomID]) >= maxUsers || (found && currRoomID == data.RoomID && (len(roomIDDict[data.RoomID])-1) >= maxUsers) { // Если пользователь подключается к той же комнате, к которой уже подключен, то позволить
			w.WriteHeader(409)
			return
		}

		if found {
			if currRoomID == data.RoomID {
				w.WriteHeader(200)
				return
			} else { // Удалить пользователя из комнаты, прописать функцию удаления пользователя через ws
				broadcastJoiningUserID <- []string{"remove", currRoomID, strconv.Itoa(data.UserID)}
				roomIDDict[currRoomID] = removeValueFromSlice(roomIDDict[currRoomID], fmt.Sprintf("%d", data.UserID))
			}
		}
		roomIDDict[data.RoomID] = append(roomIDDict[data.RoomID], fmt.Sprintf("%d", data.UserID)) //////////// Мы добавляем туда ID пользователя
		userName, imgSrc, err := getUserInfo(db, fmt.Sprintf("%d", data.UserID))
		if err != nil {
			http.Error(w, "Internal server error", 500)
			log.Println(err.Error())
			return
		}
		broadcastJoiningUserID <- []string{"add", data.RoomID, fmt.Sprintf("%d", data.UserID), userName, imgSrc}
		w.WriteHeader(200)
	}
}

func getRegisteredUserData(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Parsing error", 500)
			log.Println(err.Error())
			return
		}
		var data struct {
			UserName string
			Password string
		}
		err = json.Unmarshal(reqData, &data)
		if err != nil {
			http.Error(w, "JSON parsing error", 500)
			log.Println(err.Error())
			return
		}
		userName := data.UserName
		exists, err := userExists(db, data.UserName)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
		if exists {
			w.WriteHeader(409)
			return
		}
		userID, err := insertNewUser(db, userName, data.Password)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		_, imgSrc, err := getUserInfo(db, fmt.Sprintf("%d", userID))
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
		user := userInfo{
			UserID:   userID,
			UserName: userName,
			ImgSrc:   imgSrc,
		}
		err = setJsonCookie(w, "userInfoCookie", user, 24*time.Hour)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
		w.WriteHeader(200)
	}
}

func getLoginUserData(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Parsing error", 500)
			log.Println(err.Error())
			return
		}
		var data struct {
			UserName string
			Password string
		}
		err = json.Unmarshal(reqData, &data)
		if err != nil {
			http.Error(w, "JSON parsing error", 500)
			log.Println(err.Error())
			return
		}
		userID, exists, err := credentialExists(db, data.UserName, data.Password)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
		if exists {
			_, imgSrc, err := getUserInfo(db, fmt.Sprintf("%d", userID))
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
				log.Println(err.Error())
				return
			}
			user := userInfo{
				UserID:   userID,
				UserName: data.UserName,
				ImgSrc:   imgSrc,
			}
			err = setJsonCookie(w, "userInfoCookie", user, 24*time.Hour)
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
				log.Println(err.Error())
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(409)
		}
	}
}

func getMotion(w http.ResponseWriter, r *http.Request) {
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Parsing error", 500)
		log.Println(err.Error())
		return
	}
	var data struct {
		Name         string
		MotionString string
		UserID       int
	}
	err = json.Unmarshal(reqData, &data)
	if err != nil {
		http.Error(w, "JSON parsing error", 500)
		log.Println(err.Error())
		return
	}
	selectedRoomID, found := retrieveUserRoom(strconv.Itoa(data.UserID))
	if !found {
		http.Error(w, "User not found", 500)
		return
	}
	for conn, gameFieldID := range gameFieldWSDict {
		if gameFieldID == selectedRoomID {
			//fmt.Println("отправляем", gameFieldID, data.SelectedRoomID)
			err := conn.WriteMessage(websocket.TextMessage, reqData)
			if err != nil {
				err := conn.Close()
				if err != nil {
					w.WriteHeader(409)
					return
				}
				delete(roomWSDict, conn)
			}
			w.WriteHeader(200)
			return
		}
	}
	w.WriteHeader(409)
}

func exitFromGame(w http.ResponseWriter, r *http.Request) {
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Parsing error", 500)
		log.Println(err.Error())
		return
	}
	var data struct {
		UserID int
	}
	err = json.Unmarshal(reqData, &data)
	if err != nil {
		http.Error(w, "JSON parsing error", 500)
		log.Println(err.Error())
		return
	}
	selectedRoomID, found := retrieveUserRoom(strconv.Itoa(data.UserID))
	if !found {
		http.Error(w, "User not found", 500)
		return
	}
	message := struct {
		ExitingUser string
	}{
		ExitingUser: fmt.Sprintf("%d", data.UserID),
	}
	messageData, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "JSON marshal error", 500)
		log.Println("JSON marshal error:", err)
		return
	}
	for conn, gameFieldID := range gameFieldWSDict {
		if gameFieldID == selectedRoomID {
			err := conn.WriteMessage(websocket.TextMessage, messageData)
			if err != nil {
				err := conn.Close()
				if err != nil {
					w.WriteHeader(409)
					return
				}
				delete(roomWSDict, conn)
			}
			w.WriteHeader(200)
			return
		}
	}
}

func exitFromRoom(r *http.Request) (int, string, error) {
	var user userInfo
	err := getJsonCookie(r, "userInfoCookie", &user)
	if err != nil {
		return 0, "", err
	}
	roomID, found := retrieveUserRoom(strconv.Itoa(user.UserID))
	if found {
		roomIDDict[roomID] = removeValueFromSlice(roomIDDict[roomID], strconv.Itoa(user.UserID))
		broadcastJoiningUserID <- []string{"remove", roomID, strconv.Itoa(user.UserID)}
	}
	return user.UserID, roomID, nil
}

func exitFromRoomAPI(w http.ResponseWriter, r *http.Request) {
	exitFromGame(w, r)
	userID, roomID, err := exitFromRoom(r)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	fmt.Println(userID, "left room", roomID)
	w.WriteHeader(200)
}

func exitFromAccount(w http.ResponseWriter, r *http.Request) {
	userID, _, err := exitFromRoom(r)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "userInfoCookie",
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, -1),
	})
	fmt.Println(userID, "logged out")
}
