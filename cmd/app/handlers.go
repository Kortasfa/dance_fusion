package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var roomWSDict = make(map[*websocket.Conn]string)
var broadcastJoiningUserID = make(chan []string)
var roomIDDict = make(map[string]int)
var maxUsers int = 4

type stylesData struct {
	StyleID   int    `db:"id"`
	StyleName string `db:"name"`
}

type songsData struct {
	SongID          int    `db:"id"`
	SongName        string `db:"song_name"`
	SongAuthor      string `db:"author_name"`
	PreviewVideoSrc string `db:"preview_video_src"`
	ImageSrc        string `db:"image_src"`
	StyleID         int    `db:"style_id"`
}

type userData struct {
	UserID   int    `db:"id"`
	UserName string `db:"name"`
	ImgSrc   string `db:"img_src"`
}

type menuPageData struct {
	Styles  []stylesData
	Songs   []songsData
	RoomKey string
	WssURL  string
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
			WssURL:  "ws://" + r.Host + "/roomWS/" + roomID,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
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
	roomIDDict[fmt.Sprintf("%d", roomID)] = 0
	http.Redirect(w, r, fmt.Sprintf("/room/%d", roomID), http.StatusFound)
}

func joinPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("pages/gamePhone.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	data := struct {
		UserID string
	}{
		UserID: "1",
	}
	err = tmpl.Execute(w, data)
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
	data := new(userData)
	err := row.Scan(&data.UserName, &data.ImgSrc)
	if err != nil {
		return "", "", err
	}
	return data.UserName, data.ImgSrc, nil
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
			UserID string
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
		if roomIDDict[data.RoomID] >= maxUsers {
			w.WriteHeader(409)
			return
		}
		roomIDDict[data.RoomID] = roomIDDict[data.RoomID] + 1
		w.WriteHeader(200)
		userName, imgSrc, _ := getUserInfo(db, data.UserID)
		broadcastJoiningUserID <- []string{data.RoomID, data.UserID, userName, imgSrc}
		_, err = fmt.Fprintf(w, "Message sent: %s", data.UserID)
		if err != nil {
			return
		}
	}
}

func handleRoomWSMessages() {
	for mesArr := range broadcastJoiningUserID {
		for wsConnect := range roomWSDict {
			roomID := mesArr[0]
			userID := mesArr[1]
			userName := mesArr[2]
			imgSrc := mesArr[3]
			if roomWSDict[wsConnect] == roomID {
				message := userID + "|" + userName + "|" + imgSrc
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

func roomWSHandler(w http.ResponseWriter, r *http.Request) {
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

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(roomWSDict, conn)
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
		WssURL: "ws://" + r.Host + "/start/ws",
	}
	err = ts.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func signUp(w http.ResponseWriter, r *http.Request) {
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

func insertNewUser(db *sqlx.DB, userName string, password string) error {
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
	_, err := db.Exec(query, user.UserName, user.Password, user.UserImageSrc)
	if err != nil {
		return err
	}
	return nil
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
		} else {
			err = insertNewUser(db, userName, data.Password)
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
				log.Println(err.Error())
				return
			}
		}
		w.WriteHeader(200)
	}
}

func logIn(w http.ResponseWriter, r *http.Request) {
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

func credentialExists(db *sqlx.DB, userName string, password string) (bool, error) {
	const query = `
			SELECT COUNT(*)
			FROM users
			WHERE name = ? and password = ?`
	var count int
	err := db.QueryRow(query, userName, password).Scan(&count)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
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
		exists, err := credentialExists(db, data.UserName, data.Password)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
		if !exists {
			w.WriteHeader(409)
			return
		}
		w.WriteHeader(200)
	}
}
