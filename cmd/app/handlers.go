package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"html/template"
	"log"
	"math/rand"
	"net/http"
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

type userAvatarData struct {
	HatSrc  string
	FaceSrc string
	BodySrc string
}

type hatData struct {
	HatID    int    `db:"id"`
	HatLevel int    `db:"recommended_level"`
	HatSrc   string `db:"hat_src"`
}

type facesData struct {
	FaceID    int    `db:"id"`
	FaceLevel int    `db:"recommended_level"`
	FaceSrc   string `db:"face_src"`
}

type bodyData struct {
	BodyID    int    `db:"id"`
	BodyLevel int    `db:"recommended_level"`
	BodySrc   string `db:"body_src"`
}

type customPageData struct {
	Faces  []facesData
	Bodies []bodyData
	Hats   []hatData
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

func customPageHandler(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("userInfoCookie")
		if err != nil {
			http.Redirect(w, r, "/logIn", http.StatusFound)
			return
		}
		tmpl, err := template.ParseFiles("pages/userAccount.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
		faces, err := getFaceData(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		bodies, err := getBodyData(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		hats, err := getHatData(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		data := customPageData{
			Faces:  faces,
			Bodies: bodies,
			Hats:   hats,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
	}
}
