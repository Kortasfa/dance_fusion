package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	port         = "localhost:3000"
	dbDriverName = "mysql"
)

func main() {
	/*db, err := openDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	dbx := sqlx.NewDb(db, dbDriverName)*/

	mux := mux.NewRouter()
	mux.HandleFunc("/home", gameField)
	mux.HandleFunc("/start", gamePhone)
	mux.HandleFunc("/api/start", checkButton).Methods(http.MethodPost)
	//mux.HandleFunc("/api/play", start)
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Start server")
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err.Error())
	}
}

/*func openDB() (*sql.DB, error) {
	return sql.Open(dbDriverName, "root:P@ssw0rd@tcp(localhost:3306)/danceFusion?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
}*/
