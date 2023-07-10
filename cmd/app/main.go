package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

const (
	port         = "localhost:3306"
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
	r.HandleFunc("/home", homePageHandler)
	r.HandleFunc("/room", handleCreateRoom)
	r.HandleFunc("/room/{id}", handleRoom(dbx))
	r.HandleFunc("/roomWS/{id}", roomWSHandler)
	r.HandleFunc("/game_field/id", gameField)
	r.HandleFunc("/sign_up", signUp)

	r.HandleFunc("/api/join_to_room", getJoinedUserData(dbx)).Methods("POST")
	r.HandleFunc("/api/sign_up", getRegisteredUserData(dbx)).Methods("POST")

	go handleRoomWSMessages()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Start server")
	srv := &http.Server{
		Handler: r,
		Addr:    "localhost:3000",
	}

	log.Fatal(srv.ListenAndServe())
}

func openDB() (*sql.DB, error) {
	return sql.Open(dbDriverName, "root:P@ssw0rd@tcp(localhost:3306)/dance_fusion?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
}
