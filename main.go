package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	filepathRoot := http.Dir(".")
	port := "8080"

	mux.Handle("/", http.FileServer(filepathRoot))

	s := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}
