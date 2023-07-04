package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", "192.168.138.39:3000", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	http.HandleFunc("/gyroscope", gyrPage)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
