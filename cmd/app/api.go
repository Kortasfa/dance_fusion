package main

import (
	"encoding/json"
	"fmt"
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

	broadcastGameFieldWSMessage <- []string{selectedRoomID, string(reqData)}
	w.WriteHeader(200)
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
	// создать список комнат, в котоорых в данный момент происходит игра. Если этой комнаты в игре нет, то ничего не посылать
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
	if containsInSlice(activeGameRooms, selectedRoomID) {
		if getConnectedUsersCount(selectedRoomID) < 1 {
			fmt.Println("Сворачиваем игру: ", selectedRoomID)
			endGame(selectedRoomID)
		}
		broadcastGameFieldWSMessage <- []string{selectedRoomID, string(messageData)}
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

	broadcastJoinPageWSMessage <- []string{strconv.Itoa(data.UserID), string(reqData)}
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
	fmt.Println("Записываем color и maxPoint в broadcastGameFieldWSMessage", data.RoomID, data.MaxPoint, data.ColorID)
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

func deletePlayerFromGame(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	selectedRoomID, found := retrieveUserRoom(userID)
	if found {
		data := struct {
			Exit bool
		}{
			Exit: true,
		}
		messageData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Error marshal struct", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		broadcastJoinPageWSMessage <- []string{userID, string(messageData)}

		//Выход из команты
		roomIDDict[selectedRoomID] = removeValueFromSlice(roomIDDict[selectedRoomID], userID)
		broadcastJoiningUserID <- []string{"remove", selectedRoomID, userID}

		//Выход из игры
		message := struct {
			ExitingUser string
		}{
			ExitingUser: userID,
		}
		messageData, err = json.Marshal(message)
		if err != nil {
			http.Error(w, "Error marshal struct", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		if containsInSlice(activeGameRooms, selectedRoomID) {
			if getConnectedUsersCount(selectedRoomID) < 1 {
				fmt.Println("Сворачиваем игру: ", selectedRoomID)
				endGame(selectedRoomID)
			}
			broadcastGameFieldWSMessage <- []string{selectedRoomID, string(messageData)}
		}
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusNotFound)
}

func addUserScore(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Parsing error", 500)
			log.Println(err.Error())
			return
		}
		var data struct {
			UserID    int `json:"user_id"`
			UserScore int `json:"score"`
		}
		err = json.Unmarshal(reqData, &data)
		if err != nil {
			http.Error(w, "JSON parsing error", 500)
			log.Println(err.Error())
			return
		}
		err = addUserScoreSQL(db, data.UserID, data.UserScore)
		if err != nil {
			http.Error(w, "JSON parsing error", 500)
			log.Println(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func addBot(w http.ResponseWriter, r *http.Request) {
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Parsing error", 500)
		log.Println(err.Error())
		return
	}
	var data struct {
		RoomID string `json:"room_id"`
		BotID  string `json:"bot_id"`
	}
	err = json.Unmarshal(reqData, &data)
	if err != nil {
		http.Error(w, "JSON parsing error", 500)
		log.Println(err.Error())
		return
	}
	_, exists := roomIDDict[data.RoomID]
	if !exists {
		http.Error(w, "Room dont exist", 404)
		log.Println(err.Error())
		return
	}
	roomIDDict[data.RoomID] = append(roomIDDict[data.RoomID], data.BotID)
	w.WriteHeader(http.StatusOK)
}

func removeBot(w http.ResponseWriter, r *http.Request) {
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Parsing error", 500)
		log.Println(err.Error())
		return
	}
	var data struct {
		RoomID string `json:"room_id"`
		BotID  string `json:"bot_id"`
	}
	err = json.Unmarshal(reqData, &data)
	if err != nil {
		http.Error(w, "JSON parsing error", 500)
		log.Println(err.Error())
		return
	}
	_, exists := roomIDDict[data.RoomID]
	if !exists {
		http.Error(w, "Room dont exist", 404)
		log.Println(err.Error())
		return
	}
	roomIDDict[data.RoomID] = removeValueFromSlice(roomIDDict[data.RoomID], data.BotID)
	w.WriteHeader(http.StatusOK)
}

func startGameAPI(w http.ResponseWriter, r *http.Request) {
	roomID := r.FormValue("room_id")
	if !containsInSlice(activeGameRooms, roomID) {
		activeGameRooms = append(activeGameRooms, roomID)
	}
	w.WriteHeader(http.StatusOK)
}

func endGame(roomID string) {
	if containsInSlice(activeGameRooms, roomID) {
		activeGameRooms = removeValueFromSlice(activeGameRooms, roomID)
	}
}

func endGameAPI(w http.ResponseWriter, r *http.Request) {
	roomID := r.FormValue("room_id")
	endGame(roomID)
	w.WriteHeader(http.StatusOK)
}

func checkForAchievements(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Parsing error", 500)
			log.Println(err.Error())
			return
		}

		var data struct {
			UserID int   `json:"user_id"`
			SongID int   `json:"song_id"`
			BotIDs []int `json:"bot_ids"`
			BossID int   `json:"boss_id"`
		}

		err = json.Unmarshal(reqData, &data)
		if err != nil {
			http.Error(w, "JSON parsing error", 500)
			log.Println(err.Error())
			return
		}

		data.BotIDs = append(data.BotIDs, 0)

		query := `
			SELECT
				user_achievement_id,
				user_id,
				progress,
				max_progress
			FROM
				user_achievements
			WHERE
				user_id = ? AND song_id = ? AND bot_id IN (?) AND boss_id = ?
		`
		var achievementProgressData []struct {
			UserAchievementID int `db:"user_achievement_id"`
			UserID            int `db:"user_id"`
			Progress          int `db:"progress"`
			MaxProgress       int `db:"max_progress"`
		}

		err = db.Select(&achievementProgressData, query, data.UserID, data.SongID, data.BotIDs, data.BossID)
		if err != nil {
			http.Error(w, "Database error", 500)
			log.Println(err.Error())
			return
		}

		for _, progressInfo := range achievementProgressData {
			userAchievementID := progressInfo.UserAchievementID
			progress := progressInfo.Progress
			maxProgress := progressInfo.MaxProgress
			err = addHasNewAchievement(db, progressInfo.UserID)
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
				log.Println(err)
				return
			}

			if progress < maxProgress {
				progress++
			}

			completed := 0
			if progress >= maxProgress {
				completed = 1
			}

			updateQuery := `
				UPDATE
				    user_achievements
				SET
				    progress = ?,
					completed = ?
				WHERE
				    user_achievement_id = ?
			`

			_, err = db.Exec(updateQuery, progress, completed, userAchievementID)
			if err != nil {
				http.Error(w, "user achievements table update error", 500)
				log.Println(err.Error())
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func earnPointsForAchievements(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		achievementID := r.FormValue("achievement_id")
		achievementIDInt, err := strconv.Atoi(achievementID)
		if err != nil {
			http.Error(w, "invalid achievement id", 500)
			log.Println(err.Error())
			return
		}
		fmt.Println("Получил id ", achievementIDInt)
		var user userInfo
		err = getJsonCookie(r, "userInfoCookie", &user)
		if err != nil {
			http.Redirect(w, r, "/logIn", http.StatusFound)
			log.Println(err.Error())
			return
		}
		query := `
			SELECT
				user_achievement_id,
				completed,
				collected,
				score
			FROM
				user_achievements
			WHERE
				user_id = ? AND achievement_id = ?
		`
		var achievementData struct {
			UserAchievementID int `db:"user_achievement_id"`
			Completed         int `db:"completed"`
			Collected         int `db:"collected"`
			Score             int `db:"score"`
		}
		err = db.QueryRow(query, user.UserID, achievementIDInt).Scan(&achievementData.UserAchievementID, &achievementData.Completed, &achievementData.Collected, &achievementData.Score)
		if err != nil {
			http.Error(w, "Database error", 500)
			log.Println(err.Error())
			return
		}
		fmt.Println("Получил информацию о ачивке ", achievementData)
		if achievementData.Completed != 1 {
			http.Error(w, "achievement not yet completed", 409)
			return
		}
		if achievementData.Collected != 0 {
			http.Error(w, "award already received", 409)
			return
		}
		fmt.Println("Прошёл бэк проверку")
		updateQuery := `
				UPDATE
				    user_achievements
				SET
					collected = 1
				WHERE
				    user_achievement_id = ?
			`

		_, err = db.Exec(updateQuery, achievementData.UserAchievementID)
		if err != nil {
			http.Error(w, "user achievements table update error", 500)
			log.Println(err.Error())
			return
		}
		fmt.Println("Поменял collected на 1", achievementData.UserAchievementID)

		err = addUserScoreSQL(db, user.UserID, achievementData.Score)
		if err != nil {
			http.Error(w, "achievement scoring error", 500)
			log.Println(err.Error())
			return
		}
		fmt.Println("Добавил score пользователю", user.UserID, achievementData.Score)
		w.WriteHeader(http.StatusOK)
	}
}
