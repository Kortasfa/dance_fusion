package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const (
	port         = "localhost:3000"
	dbDriverName = "mysql"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	dbx := sqlx.NewDb(db, dbDriverName)

	r := mux.NewRouter()
	r.HandleFunc("/join", joinPageHandler).Methods("GET")
	r.HandleFunc("/", homePageHandler)
	r.HandleFunc("/home", homePageHandler)
	r.HandleFunc("/room", handleCreateRoom)
	r.HandleFunc("/room/", handleCreateRoom)
	r.HandleFunc("/room/{id}", handleRoom(dbx))
	r.HandleFunc("/gameField/id", gameField)
	r.HandleFunc("/signUp", signUp)
	r.HandleFunc("/logIn", logIn)
	r.HandleFunc("/custom", customPageHandler(dbx))

	r.HandleFunc("/roomWS/{id}", roomWSHandler(dbx))
	r.HandleFunc("/ws/joinToRoom/{id}", joinPageWSHandler)
	r.HandleFunc("/neuralWS/{id}", neuralWSHandler)

	r.HandleFunc("/api/joinToRoom", getJoinedUserData(dbx)).Methods("POST")
	r.HandleFunc("/api/signUp", getRegisteredUserData(dbx)).Methods("POST")
	r.HandleFunc("/api/logIn", getLoginUserData(dbx)).Methods("POST")
	r.HandleFunc("/api/motion", getMotion).Methods("POST")
	r.HandleFunc("/api/exitFromGame", exitFromGameAPI)
	r.HandleFunc("/api/exitFromAccount", exitFromAccount)
	r.HandleFunc("/api/custom", getUserAvatar(dbx)).Methods("POST")
	r.HandleFunc("/api/sendPoint", sendPointToJoin).Methods("POST")
	r.HandleFunc("/api/sendMaxPoint", getMaxScore).Methods("POST")
	r.HandleFunc("/api/getBestPlayer", getBestPlayer(dbx)).Methods("POST")
	r.HandleFunc("/api/updateBestPlayer", updateBestPlayer(dbx)).Methods("POST")

	go handleRoomWSMessages()
	go handleJoinPageWSMessages()
	go danceInfoHandleMessages()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Start server")
	srv := &http.Server{
		Handler: r,
		Addr:    port,
	}

	log.Fatal(srv.ListenAndServe())
}

func openDB() (*sql.DB, error) {
	return sql.Open(dbDriverName, "root:P@ssw0rd@tcp(localhost:3306)/dance_fusion?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
}
