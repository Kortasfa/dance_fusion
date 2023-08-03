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

const hashCost = 8

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
var broadcastGameFieldWSMessage = make(chan []string)

var activeGameRooms []string

type activeRoomData struct {
	ActiveRoomID string
}

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
	Difficulty      int    `db:"difficulty"`
}

type userInfo struct {
	UserID   int
	UserName string
	HatSrc   string
	FaceSrc  string
	BodySrc  string
}

type bestPlayerInfo struct {
	UserID int `db:"best_player_id"`
	Score  int `db:"best_score"`
}

type botInfo struct {
	BotId         string `db:"bot_id"`
	BotScoresPath string `db:"bot_scores_path"`
	BotImgHat     string `db:"img_hat"`
	BotImgBody    string `db:"img_body"`
	BotImgFace    string `db:"img_face"`
	Difficulty    int    `db:"difficulty"`
}

type botNameData struct {
	BotName    string `db:"bot_name"`
	Difficulty int    `db:"difficulty"`
}

type bossInfo struct {
	BossId          string `db:"boss_id"`
	BossName        string `db:"boss_name"`
	BossHealthPoint string `db:"boss_health_point"`
	BossImgHat      string `db:"img_hat"`
	BossImgBody     string `db:"img_body"`
	BossImgFace     string `db:"img_face"`
}

type menuPageData struct {
	Styles         []stylesData
	Songs          []songsData
	RoomKey        string
	ConnectedUsers []userInfo
	Bots           []botNameData
	Bosses         []bossInfo
	WssURL         string
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
	Faces     []facesData
	Bodies    []bodyData
	Hats      []hatData
	UserScore int
}

type userAchievement struct {
	UserAchievementID int    `db:"user_achievement_id"`
	UserID            int    `db:"user_id"`
	AchievementID     int    `db:"achievement_id"`
	AchievementName   string `db:"achievement_name"`
	Progress          int    `db:"progress"`
	MaxProgress       int    `db:"max_progress"`
	Completed         int    `db:"completed"`
	Level             int    `db:"level"`
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

func handleRoom(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID := vars["id"]

		var room activeRoomData
		err := getJsonCookie(r, "activeRoomCookie", &room)
		if err != nil {
			http.Redirect(w, r, "/room", http.StatusFound)
		}
		if room.ActiveRoomID != roomID {
			http.Redirect(w, r, "/room/"+room.ActiveRoomID, http.StatusFound)
		}

		_, exists := roomIDDict[roomID]
		if !exists {
			w.WriteHeader(404)
			return
		}

		for _, userID := range roomIDDict[roomID] { // Удаляем ботов из комнаты
			userIDInt, err := strconv.Atoi(userID)
			if err == nil && userIDInt < 0 {
				roomIDDict[roomID] = removeValueFromSlice(roomIDDict[roomID], userID)
			}
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

		users, err := getConnectedUsers(roomID, db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		bots, err := getBotNames(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		bosses, err := getBossInfo(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		data := menuPageData{
			Styles:         styles,
			Songs:          songs,
			RoomKey:        roomID,
			ConnectedUsers: users,
			Bots:           bots,
			Bosses:         bosses,
			WssURL:         "wss://" + r.Host + "/roomWS/" + roomID,
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
	var room activeRoomData
	err := getJsonCookie(r, "activeRoomCookie", &room)
	if err != nil {
		rand.Seed(time.Now().UnixNano())
		var roomID int
		for {
			roomID = rand.Intn(100-10-1) + 10
			_, exists := roomIDDict[fmt.Sprintf("%d", roomID)]
			if !exists {
				break
			}
		}
		roomIDDict[fmt.Sprintf("%d", roomID)] = []string{}
		room.ActiveRoomID = fmt.Sprintf("%d", roomID)
		err := setJsonCookie(w, "activeRoomCookie", room, time.Hour*24)
		if err != nil {
			log.Println("Error setting active room cookie:", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/room/%d", roomID), http.StatusFound)
	}
	_, exists := roomIDDict[room.ActiveRoomID]
	if !exists {
		roomIDDict[room.ActiveRoomID] = []string{}
	}
	http.Redirect(w, r, "/room/"+room.ActiveRoomID, http.StatusFound)
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
		var user userInfo
		err := getJsonCookie(r, "userInfoCookie", &user)
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

		userScore, err := getScoreByUserID(db, user.UserID)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		data := customPageData{
			Faces:     faces,
			Bodies:    bodies,
			Hats:      hats,
			UserScore: userScore,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
	}
}

func achievementsPageHandler(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user userInfo
		err := getJsonCookie(r, "userInfoCookie", &user)
		if err != nil {
			http.Redirect(w, r, "/logIn", http.StatusFound)
			return
		}
		tmpl, err := template.ParseFiles("pages/achievements.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
		achievements, err := getUserAchievements(db, user.UserID)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		data := struct {
			Achievements []userAchievement
		}{
			Achievements: achievements,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
	}
}
