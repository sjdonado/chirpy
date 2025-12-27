package main

import (
	"chirpy/api"
	"chirpy/middleware"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	filepathRoot := http.Dir(".")
	port := "8080"

	cfg := &middleware.ApiConfig{}

	mux.Handle("GET /admin/metrics", cfg.GetMetrics())
	mux.Handle("POST /admin/reset", cfg.ResetMetrics())

	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(filepathRoot))))

	mux.HandleFunc("GET /api/healthz", api.GetHealthz)
	mux.HandleFunc("POST /api/validate_chirp", api.PostValidateChirp)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}
