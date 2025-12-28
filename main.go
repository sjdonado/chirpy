package main

import (
	"chirpy/handler"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	mux := http.NewServeMux()
	filepathRoot := http.Dir(".")
	port := "8080"

	cfg := handler.NewApiConfig()

	mux.Handle("GET /admin/metrics", cfg.GetMetrics())
	mux.Handle("POST /admin/reset", cfg.ResetMetrics())

	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(filepathRoot))))

	mux.HandleFunc("GET /api/healthz", handler.GetHealthz)
	mux.HandleFunc("POST /api/validate_chirp", handler.PostValidateChirp)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}
