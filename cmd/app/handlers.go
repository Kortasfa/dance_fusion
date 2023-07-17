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

type userInfo struct {
	UserID       int
	UserName     string
	ImgSrc       string
	SelectedRoom string
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
		MotionString string
		MotionName   string
	}
	err = json.Unmarshal(reqData, &data)
	if err != nil {
		http.Error(w, "JSON parsing error", 500)
		log.Println(err.Error())
		return
	}
	fmt.Println(data.MotionString, data.MotionName)
	fmt.Println(data)
	w.WriteHeader(200)
}

func neuralWSHandler(w http.ResponseWriter, r *http.Request) {
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
	message := `{
		"userID": "1",
		"motion": "0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 6.0000, 1.7000, 7.8000, -10.4000, 9.6000, 24.0000, 6.2000, 1.5000, 8.0000, -13.3000, 11.2000, 28.6000, 6.2000, 1.5000, 8.2000, -16.9000, 12.9000, 35.8000, 6.2000, 1.7000, 8.3000, -24.8000, 16.3000, 45.3000, 6.7000, 2.0000, 7.8000, -27.2000, 18.7000, 53.7000, 7.2000, 2.0000, 7.1000, -27.3000, 19.0000, 54.4000, 7.2000, 2.0000, 7.1000, -27.3000, 19.0000, 54.4000, 7.4000, 1.9000, 6.9000, -27.6000, 19.8000, 51.6000, 8.0000, 1.8000, 6.9000, -39.8000, 20.3000, 48.6000, 8.6000, 1.9000, 7.0000, -54.4000, 16.5000, 47.7000, 8.6000, 1.9000, 7.0000, -54.4000, 16.5000, 47.7000, 9.3000, 2.2000, 7.7000, -78.6000, -0.8000, 60.7000, 9.3000, 2.2000, 7.7000, -78.6000, -0.8000, 60.7000, 9.4000, 2.1000, 6.8000, -79.0000, -10.5000, 67.0000, 8.8000, 2.1000, 4.1000, -56.5000, -28.7000, 51.8000, 8.2000, 2.0000, 3.8000, -48.6000, -32.4000, 33.6000, 7.8000, 2.1000, 4.2000, -50.3000, -36.3000, 21.7000, 7.1000, 2.0000, 3.9000, -54.7000, -40.7000, 14.9000, 6.1000, 1.8000, 3.6000, -57.1000, -41.3000, 8.2000, 6.1000, 1.8000, 3.6000, -57.1000, -41.3000, 8.2000, 5.0000, 1.5000, 4.0000, -75.2000, -27.8000, 0.8000, 5.2000, 1.3000, 4.5000, -87.9000, -18.5000, 2.0000, 5.2000, 1.3000, 4.5000, -87.9000, -18.5000, 2.0000, 4.7000, 0.6000, 5.0000, -104.7000, 1.6000, 12.7000, 4.7000, 0.6000, 5.0000, -104.7000, 1.6000, 12.7000, 4.9000, 0.6000, 5.0000, -102.2000, 14.6000, 22.6000, 2.1000, -0.6000, 2.9000, -77.7000, 36.9000, 50.0000, 1.3000, -1.8000, 1.5000, -74.3000, 61.5000, 51.7000, 1.1000, -2.9000, 0.3000, -63.2000, 96.3000, 44.8000, 2.3000, -4.1000, 0.0000, -49.1000, 137.0000, 33.0000, 5.9000, -5.1000, 0.2000, -33.1000, 181.8000, 15.8000, 5.9000, -5.1000, 0.2000, -33.1000, 181.8000, 15.8000, 5.9000, -5.1000, 0.2000, -33.1000, 181.8000, 15.8000, 11.8000, -5.0000, 0.3000, -43.4000, 207.7000, -5.7000, 28.4000, -3.5000, 6.3000, -107.2000, 217.5000, -89.3000, 28.4000, -3.5000, 6.3000, -107.2000, 217.5000, -89.3000, 37.0000, 0.5000, 13.6000, -114.4000, 215.2000, -85.7000, 41.5000, 3.9000, 15.8000, -77.1000, 155.4000, -65.6000, 36.6000, 4.0000, 17.0000, 43.1000, 56.2000, -39.4000, 24.4000, 3.8000, 15.0000, 142.3000, -13.4000, 11.8000, 5.7000, 4.8000, 5.6000, 47.5000, -27.9000, 104.7000, 5.7000, 4.8000, 5.6000, 47.5000, -27.9000, 104.7000, 1.3000, 5.9000, 2.4000, -46.3000, -19.7000, 138.0000, -3.0000, 7.3000, -3.7000, -107.7000, -11.9000, 148.6000, -3.5000, 7.4000, -6.8000, -91.4000, -22.1000, 116.7000, -3.5000, 7.4000, -6.8000, -91.4000, -22.1000, 116.7000, -5.8000, 4.9000, -9.4000, -51.1000, -59.0000, -3.8000, -8.7000, 3.1000, -6.5000, -62.7000, -55.1000, -51.2000, -8.7000, 3.1000, -6.5000, -62.7000, -55.1000, -51.2000, -10.2000, 1.4000, -3.2000, -90.7000, -31.3000, -69.1000, -10.6000, -0.7000, 0.4000, -126.1000, 43.0000, -47.2000, -10.4000, -1.5000, 1.2000, -119.9000, 99.1000, -29.7000, -6.5000, -2.4000, 1.0000, -94.3000, 158.1000, -16.2000, 0.0000, -3.8000, 0.5000, -62.6000, 198.3000, -9.5000, 8.5000, -7.1000, 0.1000, -42.1000, 217.8000, -8.1000, 18.1000, -14.2000, 0.7000, -32.8000, 230.1000, -7.9000, 18.1000, -14.2000, 0.7000, -32.8000, 230.1000, -7.9000, 27.5000, -21.3000, 0.7000, -16.2000, 285.4000, -4.1000, 27.5000, -21.3000, 0.7000, -16.2000, 285.4000, -4.1000, 63.0000, -10.5000, -3.7000, 87.2000, 303.1000, -90.0000, 63.0000, -10.5000, -3.7000, 87.2000, 303.1000, -90.0000, 64.9000, 0.0000, 2.1000, 105.5000, 85.3000, -123.2000, 47.7000, -3.7000, 11.8000, 76.4000, -120.8000, -90.0000, 19.7000, -14.1000, 14.2000, 133.1000, -175.1000, 212.0000, 14.8000, -14.1000, 2.8000, 102.1000, -166.2000, 310.3000, 7.5000, -10.8000, -6.1000, 34.8000, -156.4000, 306.4000, 1.3000, -6.0000, -9.2000, -31.5000, -141.2000, 243.8000, -4.0000, -2.1000, -7.9000, -95.5000, -121.7000, 174.4000, -6.5000, 0.3000, -5.0000, -132.3000, -103.8000, 130.4000, -7.9000, 1.2000, -4.8000, -121.8000, -107.9000, 101.4000, -7.9000, 1.2000, -4.8000, -121.8000, -107.9000, 101.4000, -12.5000, -0.1000, -8.4000, -78.0000, -128.0000, -2.9000, -14.0000, -0.9000, -7.1000, -62.7000, -97.7000, -70.7000, -14.0000, -0.9000, -7.1000, -62.7000, -97.7000, -70.7000, -13.3000, -1.2000, 0.3000, -33.2000, 11.8000, -122.2000, -13.3000, -1.2000, 0.3000, -33.2000, 11.8000, -122.2000, -10.6000, -1.8000, 2.6000, -4.7000, 146.9000, -88.7000, -6.2000, -3.2000, 3.5000, 8.0000, 217.7000, -61.5000, 0.6000, -5.3000, 2.6000, 19.9000, 265.8000, -34.0000, 9.5000, -9.0000, 0.8000, 27.7000, 283.2000, -17.1000, 20.1000, -14.5000, -1.3000, 15.7000, 271.7000, -17.1000, 20.1000, -14.5000, -1.3000, 15.7000, 271.7000, -17.1000, 29.6000, -22.0000, -1.9000, -30.4000, 279.9000, -31.2000, 29.6000, -22.0000, -1.9000, -30.4000, 279.9000, -31.2000, 42.9000, -14.6000, -2.4000, -87.1000, 346.8000, -102.6000, 42.9000, -14.6000, -2.4000, -87.1000, 346.8000, -102.6000, 52.8000, -5.5000, 0.1000, -69.9000, 225.1000, -148.2000, 36.4000, -9.9000, 13.9000, 181.1000, -214.5000, -96.4000, 36.4000, -9.9000, 13.9000, 181.1000, -214.5000, -96.4000, 22.4000, -15.5000, 14.6000, 300.3000, -261.7000, 46.7000, 12.7000, -15.8000, -0.9000, 230.3000, -198.1000, 249.4000, 6.7000, -11.2000, -5.3000, 106.7000, -175.1000, 252.7000, 6.7000, -11.2000, -5.3000, -20.1000, -146.7000, 226.9000, -3.5000, -2.9000, -5.0000, -126.6000, -112.0000, 196.7000, -3.9000, -0.8000, -3.8000, -173.5000, -96.7000, 176.9000, -3.9000, -0.8000, -3.8000, -173.5000, -96.7000, 176.9000, -7.5000, 0.2000, -9.0000, -125.1000, -134.2000, 103.0000, -10.4000, -0.5000, -9.9000, -75.7000, -148.7000, 23.9000, -10.4000, -0.5000, -9.9000, -75.7000, -148.7000, 23.9000, -14.3000, -1.8000, -5.4000, 19.7000, -66.2000, -118.4000, -14.3000, -1.8000, -5.4000, 19.7000, -66.2000, -118.4000, -13.4000, -2.0000, -2.5000, 44.5000, 4.9000, -144.4000, -10.6000, -3.3000, -0.5000, 54.0000, 147.8000, -151.8000, -10.6000, -3.3000, -0.5000, 54.0000, 147.8000, -151.8000, -1.7000, -7.3000, 0.8000, 23.1000, 284.1000, -146.3000, 9.5000, -11.5000, 1.4000, 12.4000, 316.8000, -141.6000, 9.5000, -11.5000, 1.4000, 12.4000, 316.8000, -141.6000, 20.6000, -17.0000, 2.3000, -35.4000, 292.1000, -127.4000, 20.6000, -17.0000, 2.3000, -35.4000, 292.1000, -127.4000, 38.7000, -26.7000, 8.2000, -145.8000, 326.0000, -77.9000, 48.2000, -13.4000, 6.6000, -149.0000, 301.0000, -59.6000, 48.2000, -13.4000, 6.6000, -149.0000, 301.0000, -59.6000, 49.9000, -3.3000, 13.3000, 39.4000, -132.8000, -76.0000, 30.1000, -14.0000, 17.1000, 195.2000, -270.2000, 30.4000, 17.0000, -19.2000, 11.4000, 243.1000, -253.8000, 172.9000, 11.4000, -19.1000, 0.5000, 170.9000, -189.9000, 259.2000, 7.4000, -14.6000, -6.1000, 56.5000, -146.8000, 257.5000, 4.3000, -9.0000, -6.6000, -27.6000, -132.2000, 223.2000"
	}`

	//for {
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Failed to write message to websocket:", err)
		//break
	}
	/*	_, answer, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message from websocket:", err)
			break
		}
		fmt.Println(string(answer))
		break
	}*/
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
			WssURL:  "wss://" + r.Host + "/roomWS/" + roomID,
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
		if len(roomIDDict[data.RoomID]) >= maxUsers {
			w.WriteHeader(409)
			return
		}
		slice := roomIDDict[data.RoomID]
		roomIDDict[data.RoomID] = append(slice, fmt.Sprintf("%d", data.UserID))
		userName, imgSrc, err := getUserInfo(db, fmt.Sprintf("%d", data.UserID))
		if err != nil {
			http.Error(w, "Internal server error", 500)
			log.Println(err.Error())
			return
		}
		broadcastJoiningUserID <- []string{data.RoomID, fmt.Sprintf("%d", data.UserID), userName, imgSrc}

		var user userInfo
		err = getJsonCookie(r, "userInfoCookie", &user)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		user.SelectedRoom = data.RoomID
		err = setJsonCookie(w, "userInfoCookie", user, 24*time.Hour)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		w.WriteHeader(200)
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
			UserID:       userID,
			UserName:     userName,
			ImgSrc:       imgSrc,
			SelectedRoom: "",
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
				UserID:       userID,
				UserName:     data.UserName,
				ImgSrc:       imgSrc,
				SelectedRoom: "",
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

func clearCookie(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "userInfoCookie",
			Path:    "/",
			Expires: time.Now().AddDate(0, 0, -1),
		})
		fmt.Println("Cookie is deleted")
		w.WriteHeader(200)
	}
}
