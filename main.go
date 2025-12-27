package main

import (
	"chirpy/middleware"
	"io"
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

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := io.WriteString(w, "OK\n"); err != nil {
			log.Fatal("Response could not be written")
		}
	})

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}
