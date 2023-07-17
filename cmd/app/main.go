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
	port         = ":3000"
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
	r.HandleFunc("/test", test)
	r.HandleFunc("/room", handleCreateRoom)
	r.HandleFunc("/room/{id}", handleRoom(dbx))
	r.HandleFunc("/roomWS/{id}", roomWSHandler)
	r.HandleFunc("/gameField/id", gameField)
	r.HandleFunc("/signUp", signUp)
	r.HandleFunc("/login", logIn)

	r.HandleFunc("/api/joinToRoom", getJoinedUserData(dbx)).Methods("POST")
	r.HandleFunc("/api/signUp", getRegisteredUserData(dbx)).Methods("POST")
	r.HandleFunc("/api/logIn", getLoginUserData(dbx)).Methods("POST")
	r.HandleFunc("/protected", protectedHandler)
	r.HandleFunc("/clear", clearCookie(dbx))
	r.HandleFunc("/api/post", create)

	go handleRoomWSMessages()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// Start file server in a separate Goroutine
	go func() {
		port := ":8082"
		directory := "/home/korta/Downloads/dance_fusion/static/test" // Update this line with the desired directory path
		http.Handle("/", http.FileServer(http.Dir(directory)))

		log.Printf("File server running at http://localhost%s\n", port)
		log.Fatal(http.ListenAndServe(port, nil))
	}()

	fmt.Println("Start server")
	srv := &http.Server{
		Handler: r,
		Addr:    port,
	}

	log.Printf("Main server running at http://localhost%s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func openDB() (*sql.DB, error) {
	return sql.Open(dbDriverName, "root:root123321@tcp(localhost:3306)/dance_fusion?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
}
