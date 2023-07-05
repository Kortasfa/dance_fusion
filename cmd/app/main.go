package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/echo", echo)
	r.HandleFunc("/", home)
	r.HandleFunc("/gyroscope", gyrPage)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Start server")
	srv := &http.Server{
		Handler: r,
		Addr:    "192.168.138.39:3000",
	}

	log.Fatal(srv.ListenAndServe())
}
