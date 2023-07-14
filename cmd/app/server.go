package main

import (
	"log"
	"net/http"
)

func main() {
	port := ":8082"
	directory := "/home/korta/Downloads/dance_fusion-back-Valerian/browser" // Update this line with the desired directory path
	http.Handle("/", http.FileServer(http.Dir(directory)))

	log.Printf("Running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
