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
	port         = ":3306"
	dbDriverName = "mysql"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	dbx := sqlx.NewDb(db, dbDriverName)

	r := mux.NewRouter()
	r.HandleFunc("/echo", echo)
	r.HandleFunc("/", home)
	r.HandleFunc("/gyroscope", gyrPage)

	r.HandleFunc("/menu", menuPage(dbx))
	r.HandleFunc("/preview_video", previewVideoWS)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Start server")
	srv := &http.Server{
		Handler: r,
		Addr:    "192.168.138.39:3000",
	}

	log.Fatal(srv.ListenAndServe())
}

func openDB() (*sql.DB, error) {
	return sql.Open(dbDriverName, "root:P@ssw0rd@tcp(localhost:3306)/dance_fusion?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
}
