package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var roomWSDict = make(map[*websocket.Conn]string) // {WSConnection: roomID, WSConnection: roomID, WSConnection: roomID, ...}
var broadcastJoiningUserID = make(chan []string)  // [RoomID, UserID, userName, imgSrc]
var roomIDDict = make(map[string][]string)        // {roomID: [userID, userID, userID, userID]}
var maxUsers int = 4

var joinPageWSDict = make(map[*websocket.Conn]string) // {WSConnection: UserID, WSConnection: UserID, WSConnection: UserID, ...}
var broadcastJoinPageWSMessage = make(chan []string)  // [UserID, Data]

var gameFieldWSDict = make(map[*websocket.Conn]string) // {WSConnection: gameFieldID, WSConnection: gameFieldID, WSConnection: gameFieldID, ...}
//var broadcastGameFieldWSMessage = make(chan []string)

type stylesData struct {
	StyleID   int    `db:"id"`
	StyleName string `db:"name"`
}

type songsData struct {
	SongID          int    `db:"id"`
	SongName        string `db:"song_name"`
	SongAuthor      string `db:"author_name"`
	PreviewVideoSrc string `db:"preview_video_src"`
	VideoSrc        string `db:"video_src"`
	ImageSrc        string `db:"image_src"`
	StyleID         int    `db:"style_id"`
}

type userInfo struct {
	UserID   int
	UserName string
	ImgSrc   string
}

type menuPageData struct {
	Styles  []stylesData
	Songs   []songsData
	RoomKey string
	WssURL  string
}

func getMotion(w http.ResponseWriter, r *http.Request) {
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Parsing error", 500)
		log.Println(err.Error())
		return
	}
	var data struct {
		Name           string
		MotionString   string
		SelectedRoomID string
		UserID         int
	}
	err = json.Unmarshal(reqData, &data)
	if err != nil {
		http.Error(w, "JSON parsing error", 500)
		log.Println(err.Error())
		return
	}
	for conn, gameFieldID := range gameFieldWSDict {
		if gameFieldID == data.SelectedRoomID {
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

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("pages/homePage.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	var data = struct {
		BackgroundVideoSrc string
	}{
		BackgroundVideoSrc: "static/video/JDN_Landing_Video.mp4",
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func getStylesData(db *sqlx.DB) ([]stylesData, error) {
	const query = `
		SELECT
			id,
			name
		FROM
			styles
	`
	var data []stylesData

	err := db.Select(&data, query)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func getSongsData(db *sqlx.DB) ([]songsData, error) {
	const query = `
		SELECT
			id,
			song_name,
			author_name,
			video_src,
			preview_video_src,
			image_src,
			style_id
		FROM
			songs
	`
	var data []songsData

	err := db.Select(&data, query)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func sendConnectedUserInfo(db *sqlx.DB, roomID string) error {
	userIDs := roomIDDict[roomID]
	for _, userID := range userIDs {
		const query = `
			SELECT
				name,
				img_src
			FROM
				users
			WHERE id = ?`

		var data struct {
			UserName string `db:"name"`
			ImgSrc   string `db:"img_src"`
		}

		id, err := strconv.Atoi(userID)
		if err != nil {
			return err
		}
		err = db.Get(&data, query, id)
		if err != nil {
			return err
		}
		broadcastJoiningUserID <- []string{"add", roomID, userID, data.UserName, data.ImgSrc}
	}
	return nil
}

func handleRoom(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID := vars["id"]
		_, exists := roomIDDict[roomID]
		if !exists {
			w.WriteHeader(404)
			return
		}
		tmpl, err := template.ParseFiles("pages/mainRoom.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
		styles, err := getStylesData(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		songs, err := getSongsData(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		data := menuPageData{
			Styles:  styles,
			Songs:   songs,
			RoomKey: roomID,
			WssURL:  "wss://" + r.Host + "/roomWS/" + roomID,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		go func() {
			time.Sleep(1000 * time.Millisecond)
			err = sendConnectedUserInfo(db, roomID)
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
				log.Println(err.Error())
				return
			}
		}()
	}
}

func handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	var roomID int
	for {
		roomID = rand.Intn(100)
		_, exists := roomIDDict[fmt.Sprintf("%d", roomID)]
		if !exists {
			break
		}
	}
	roomIDDict[fmt.Sprintf("%d", roomID)] = []string{}
	http.Redirect(w, r, fmt.Sprintf("/room/%d", roomID), http.StatusFound)
}

func joinPageHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("userInfoCookie")
	if err != nil {
		http.Redirect(w, r, "/logIn", http.StatusFound)
		return
	}
	tmpl, err := template.ParseFiles("pages/gamePhone.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func getUserInfo(db *sqlx.DB, userID string) (string, string, error) {
	const query = `
		SELECT
			name,
			img_src
		FROM
			users
		WHERE
		   id=?
	`
	row := db.QueryRow(query, userID)
	//data := new(userData)
	data := new(struct {
		UserID   int    `db:"id"`
		UserName string `db:"name"`
		ImgSrc   string `db:"img_src"`
	})
	err := row.Scan(&data.UserName, &data.ImgSrc)
	if err != nil {
		return "", "", err
	}
	return data.UserName, data.ImgSrc, nil
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
			delete(roomWSDict, conn)
			break
		}
	}
}

func handleJoinPageWSMessages() { // broadcastJoinPageWSMessage <- []string{UserID, Data}
	for mesArr := range broadcastJoinPageWSMessage {
		for wsConnect := range joinPageWSDict {
			userID := mesArr[0]
			data := mesArr[1]

			if joinPageWSDict[wsConnect] == userID {
				//log.Println("ОТПРАВИЛ", userID, data)
				err := wsConnect.WriteMessage(websocket.TextMessage, []byte(data))
				if err != nil {
					err := wsConnect.Close()
					if err != nil {
						return
					}
					delete(roomWSDict, wsConnect)
				}
			}
		}
	}
}

func retrieveUserRoom(userID string) (string, bool) {
	for key, userIDSlice := range roomIDDict {
		for _, currUserID := range userIDSlice {
			if currUserID == userID {
				return key, true
			}
		}
	}
	return "", false
}

func removeValueFromSlice(words []string, valueToRemove string) []string {
	var result []string

	for _, word := range words {
		if word != valueToRemove {
			result = append(result, word)
		}
	}

	return result
}

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

func handleRoomWSMessages() {
	for mesArr := range broadcastJoiningUserID {
		for wsConnect := range roomWSDict {
			action := mesArr[0]
			roomID := mesArr[1]
			userID := mesArr[2]
			message := ""
			if action == "add" {
				userName := mesArr[3]
				imgSrc := mesArr[4]
				message = action + "|" + userID + "|" + userName + "|" + imgSrc
			} else {
				message = action + "|" + userID
			}

			if roomWSDict[wsConnect] == roomID {
				err := wsConnect.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					err := wsConnect.Close()
					if err != nil {
						return
					}
					delete(roomWSDict, wsConnect)
				}
			}
		}
	}
}

func getMotionListPath(db *sqlx.DB, songName string) (string, error) {
	const query = `
		SELECT
			motion_list_path
		FROM
			songs
		WHERE
		   song_name=?
	`

	var motionListPath string
	err := db.QueryRow(query, songName).Scan(&motionListPath)
	if err != nil {
		return "", err
	}

	return motionListPath, nil
}

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

		motionListPath, err := getMotionListPath(db, string(message))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			delete(roomWSDict, conn)
			return
		}

		fileContent, err := ioutil.ReadFile(motionListPath)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			fmt.Println("Ошибка при открытии файла:", err)
			delete(roomWSDict, conn)
			return
		}
		//log.Println(string(fileContent))
		for roomID, userSlice := range roomIDDict {
			if roomID == websocketID {
				for _, userID := range userSlice {
					broadcastJoinPageWSMessage <- []string{userID, string(fileContent)}
					//fmt.Println(userID, "Пишем")
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

func signUp(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("userInfoCookie")
	if err == nil {
		http.Redirect(w, r, "/join", http.StatusFound)
		return
	}
	ts, err := template.ParseFiles("pages/signUp.html")
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

func userExists(db *sqlx.DB, userName string) (bool, error) {
	const query = `
			SELECT COUNT(*)
			FROM users
			WHERE name = ?`
	var count int
	err := db.QueryRow(query, userName).Scan(&count)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func insertNewUser(db *sqlx.DB, userName string, password string) (int, error) {
	user := struct {
		UserName     string
		Password     string
		UserImageSrc string
	}{
		UserName:     userName,
		Password:     password,
		UserImageSrc: "static/img/user_1.png",
	}
	query := `
		INSERT INTO users(name, password, img_src)
		VALUES (?, ?, ?)`
	result, err := db.Exec(query, user.UserName, user.Password, user.UserImageSrc)
	if err != nil {
		return 0, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(userID), nil
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

func logIn(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("userInfoCookie")
	if err == nil {
		http.Redirect(w, r, "/join", http.StatusFound)
		return
	}
	ts, err := template.ParseFiles("pages/logIn.html")
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

func credentialExists(db *sqlx.DB, userName string, password string) (int, bool, error) {
	const query = `
		SELECT id
		FROM users
		WHERE name = ? and password = ?`
	var userIDs []int
	err := db.Select(&userIDs, query, userName, password)
	if len(userIDs) == 0 {
		return 0, false, nil
	} else if err != nil {
		log.Println(err.Error())
		return 0, false, err
	}
	return userIDs[0], true, nil
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

func setJsonCookie(w http.ResponseWriter, name string, value interface{}, expiration time.Duration) error {
	cookieValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	escapedValue := url.QueryEscape(string(cookieValue))
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   escapedValue,
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, 1),
	})
	return nil
}

func getJsonCookie(r *http.Request, name string, value interface{}) error {
	cookie, err := r.Cookie(name)
	if err != nil {
		return err
	}
	decodedValue, err := url.PathUnescape(cookie.Value)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(decodedValue), value)
	if err != nil {
		return err
	}
	return nil
}

func clearCookie(w http.ResponseWriter, r *http.Request) {
	var user userInfo
	err := getJsonCookie(r, "userInfoCookie", &user)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	roomID, found := retrieveUserRoom(strconv.Itoa(user.UserID))
	if found {
		roomIDDict[roomID] = removeValueFromSlice(roomIDDict[roomID], strconv.Itoa(user.UserID))
		broadcastJoiningUserID <- []string{"remove", roomID, strconv.Itoa(user.UserID)}
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "userInfoCookie",
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, -1),
	})
	fmt.Println("Cookie is deleted")
	w.WriteHeader(200)
}
