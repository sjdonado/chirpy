package main

import (
	"log"
	"net/http"
	"io"
)

func main() {
	mux := http.NewServeMux()

	filepathRoot := http.Dir(".")
	port := "8080"

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(filepathRoot)))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := io.WriteString(w, "OK\n"); err != nil {
			log.Fatal("Response could not be written")
		}
	})

	s := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}
