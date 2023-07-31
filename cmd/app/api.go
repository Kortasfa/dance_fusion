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
		_, exists, err := userExists(db, data.UserName)
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
	fmt.Println(data.RoomID, data.MaxPoint, data.ColorID)
	broadcastGameFieldWSMessage <- []string{data.RoomID, string(reqData)}
	w.WriteHeader(200)
}

func getBestPlayer(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		songID, err := strconv.Atoi(r.FormValue("song_id"))
		if err != nil {
			http.Error(w, "Invalid song ID", http.StatusBadRequest)
			log.Println(err.Error())
			return
		}

		bestPlayerData, err := getBestPlayerInfo(db, songID)
		if err != nil {
			http.Error(w, "Error getting best player information", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		if bestPlayerData.UserID == 0 {
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write([]byte("{}"))
			if err != nil {
				http.Error(w, "Json send error", http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
			return
		}
		fmt.Println(bestPlayerData.UserID)
		userData, err := getUserInfo(db, bestPlayerData.UserID)
		if err != nil {
			http.Error(w, "Error getting user information", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		response := struct {
			UserInfo  userInfo
			BestScore int
		}{
			UserInfo:  userData,
			BestScore: bestPlayerData.Score,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Json send error", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
	}
}

func updateBestPlayer(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Parsing error", 500)
			log.Println(err.Error())
			return
		}
		var data struct {
			SongID int `json:"song_id"`
			UserID int `json:"user_id"`
			Score  int `json:"score"`
		}
		err = json.Unmarshal(reqData, &data)
		if err != nil {
			http.Error(w, "JSON parsing error", 500)
			log.Println(err.Error())
			return
		}
		err = updateBestPlayerSQL(db, data.SongID, data.UserID, data.Score)
		if err != nil {
			http.Error(w, "Error updating best player info", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func changeUserName(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Parsing error", 500)
			log.Println(err.Error())
			return
		}
		var data struct {
			UserID      int    `json:"user_id"`
			NewUserName string `json:"new_user_name"`
		}
		err = json.Unmarshal(reqData, &data)
		if err != nil {
			http.Error(w, "JSON parsing error", 500)
			log.Println(err.Error())
			return
		}
		userID, exists, err := userExists(db, data.NewUserName)
		if err != nil {
			http.Error(w, "SQL request error", 500)
			log.Println(err.Error())
			return
		}
		if userID != data.UserID && exists {
			w.WriteHeader(http.StatusConflict)
			return
		}
		err = updateUserName(db, data.UserID, data.NewUserName)
		if err != nil {
			http.Error(w, "Error updating user name", 500)
			log.Println(err.Error())
			return
		}
		var user userInfo
		err = getJsonCookie(r, "userInfoCookie", &user)
		if err != nil {
			http.Error(w, "Error getting cookie", 500)
			log.Println(err.Error())
			return
		}
		user.UserName = data.NewUserName
		err = setJsonCookie(w, "userInfoCookie", user, 24*time.Hour)
		if err != nil {
			http.Error(w, "Error setting cookie", 500)
			log.Println(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func changeUserPassword(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Parsing error", 500)
			log.Println(err.Error())
			return
		}
		var data struct {
			UserID          int    `json:"user_id"`
			NewUserPassword string `json:"new_user_password"`
		}
		err = json.Unmarshal(reqData, &data)
		if err != nil {
			http.Error(w, "JSON parsing error", 500)
			log.Println(err.Error())
			return
		}
		err = updateUserPassword(db, data.UserID, data.NewUserPassword)
		if err != nil {
			http.Error(w, "Error updating user password", 500)
			log.Println(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getBotPath(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		botName := r.FormValue("bot_name")
		if botName == "" {
			// Обработка случая, когда поле "bot_name" не было отправлено в форме
			http.Error(w, "Field 'bot_name' is missing or empty", http.StatusBadRequest)
			return
		}
		botData, err := getBotInfo(db, botName)
		if err != nil {
			http.Error(w, "Error getting bot information", http.StatusConflict)
			log.Println(err.Error())
			return
		}
		fmt.Println(botData)
		response := struct {
			BotId         string
			BotScoresPath string
			BotImgHat     string
			BotImgBody    string
			BotImgFace    string
		}{
			BotId:         botData.BotId,
			BotScoresPath: botData.BotScoresPath,
			BotImgHat:     botData.BotImgHat,
			BotImgBody:    botData.BotImgBody,
			BotImgFace:    botData.BotImgFace,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
