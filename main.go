package main

import (
	"chirpy/api"
	"chirpy/internal/database"
	"chirpy/middleware"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	mux := http.NewServeMux()
	filepathRoot := http.Dir(".")
	port := "8080"

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cfg := &middleware.ApiConfig{
		Queries: database.New(db),
	}

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
