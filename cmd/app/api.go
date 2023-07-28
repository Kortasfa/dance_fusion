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
		user, err := getUserInfo(db, data.UserID)
		if err != nil {
			http.Error(w, "Internal server error", 500)
			log.Println(err.Error())
			return
		}
		broadcastJoiningUserID <- []string{"add", data.RoomID, fmt.Sprintf("%d", data.UserID), user.UserName, user.HatSrc, user.FaceSrc, user.BodySrc}
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

		user, err := getUserInfo(db, userID)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
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
			user, err := getUserInfo(db, userID)
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
				log.Println(err.Error())
				return
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
				delete(gameFieldWSDict, conn)
				if err != nil {
					w.WriteHeader(409)
					return
				}
			}
			w.WriteHeader(200)
			return
		}
	}
	w.WriteHeader(409)
}

func exitFromGameAPI(w http.ResponseWriter, r *http.Request) {
	userID, roomID, err := exitFromGame(r)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	fmt.Println("User", userID, "left the room", roomID)
	w.WriteHeader(200)
}

func exitFromGame(r *http.Request) (int, string, error) {
	var user userInfo
	err := getJsonCookie(r, "userInfoCookie", &user)
	if err != nil {
		log.Println(err.Error())
		return 0, "", err
	}

	selectedRoomID, found := retrieveUserRoom(strconv.Itoa(user.UserID))
	if !found {
		return user.UserID, "хз, он уже давно вышел", nil
	}
	//Выход из команты
	roomIDDict[selectedRoomID] = removeValueFromSlice(roomIDDict[selectedRoomID], strconv.Itoa(user.UserID))
	broadcastJoiningUserID <- []string{"remove", selectedRoomID, strconv.Itoa(user.UserID)}

	//Выход из игры
	message := struct {
		ExitingUser string
	}{
		ExitingUser: fmt.Sprintf("%d", user.UserID),
	}
	messageData, err := json.Marshal(message)
	if err != nil {
		return 0, "", err
	}
	for conn, gameFieldID := range gameFieldWSDict {
		if gameFieldID == selectedRoomID {
			err := conn.WriteMessage(websocket.TextMessage, messageData)
			if err != nil {
				err := conn.Close()
				delete(gameFieldWSDict, conn)
				if err != nil {
					return 0, "", err
				}
			}
			break
		}
	}
	return user.UserID, selectedRoomID, nil
}

func exitFromAccount(w http.ResponseWriter, r *http.Request) {
	userID, _, err := exitFromGame(r)
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

func getUserAvatar(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		var userAvatar userAvatarData
		err = json.Unmarshal(reqData, &userAvatar)
		if err != nil {
			http.Error(w, "Internal Server Error Unmarshall", 500)
			log.Println(err.Error())
			return
		}
		var user userInfo
		err = getJsonCookie(r, "userInfoCookie", &user)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		err = changeUserAvatar(db, userAvatar, user.UserID)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
		user.HatSrc = userAvatar.HatSrc
		user.FaceSrc = userAvatar.FaceSrc
		user.BodySrc = userAvatar.BodySrc
		err = setJsonCookie(w, "userInfoCookie", user, 24*time.Hour)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
	}
}

func sendPointToJoin(w http.ResponseWriter, r *http.Request) {
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err)
		return
	}

	var data struct {
		UserID int
		Point  int
	}

	err = json.Unmarshal(reqData, &data)
	if err != nil {
		http.Error(w, "Internal Server Error Unmarshall", 500)
		log.Println(err.Error())
		return
	}
	for conn, userID := range joinPageWSDict {
		if strconv.Itoa(data.UserID) == userID {
			err := conn.WriteMessage(websocket.TextMessage, reqData)
			if err != nil {
				err := conn.Close()
				delete(joinPageWSDict, conn)
				if err != nil {
					w.WriteHeader(409)
					log.Println(err.Error())
					return
				}
			}
			break
		}
	}
	w.WriteHeader(200)
}

func getDataSongJson(w http.ResponseWriter, r *http.Request) {
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Parsing error", 500)
		log.Println(err.Error())
		return
	}
	var data struct {
		RoomID   string
		MaxPoint int
		ColorID  string
	}
	err = json.Unmarshal(reqData, &data)
	if err != nil {
		http.Error(w, "JSON parsing error", 500)
		log.Println(err.Error())
		return
	}
	go func() {
		isSend := false
		for {
			for conn, gameFieldID := range gameFieldWSDict {
				if gameFieldID == data.RoomID {
					err := conn.WriteMessage(websocket.TextMessage, reqData)
					isSend = true
					if err != nil {
						delete(gameFieldWSDict, conn)
						err := conn.Close()
						if err != nil {
							w.WriteHeader(409)
							return
						}
					}
					break
				}
			}
			if isSend {
				break
			}
		}
	}()
	w.WriteHeader(200)
}
