package main

import (
	"chirpy/handler"
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

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := database.New(db)

	metric_middleware := middleware.NewMetricMiddleware()

	metrics_handler := handler.NewMetricsHandler(queries, metric_middleware)
	users_handler := handler.NewUsersHandler(queries)
	chirps_handler := handler.NewChirpsHandler(queries)

	mux.Handle("GET /admin/metrics", metrics_handler.GetMetrics())
	mux.Handle("POST /admin/reset", metrics_handler.ResetMetrics())

	mux.Handle("/app/", metric_middleware.FileServerHits(http.StripPrefix("/app", http.FileServer(filepathRoot))))

	mux.Handle("POST /api/users", users_handler.CreateUser())
	mux.Handle("POST /api/login", users_handler.Login())

	mux.Handle("POST /api/chirps", chirps_handler.CreateChirp())
	mux.Handle("GET /api/chirps", chirps_handler.GetAllChirps())
	mux.Handle("GET /api/chirps/{id}", chirps_handler.GetOneChirp())

	mux.HandleFunc("GET /api/healthz", handler.GetHealthz)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}
