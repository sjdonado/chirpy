package handler

import (
	"chirpy/internal/database"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func NewApiConfig() *ApiConfig {
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return &ApiConfig{db: database.New(db)}
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) GetMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf(`
		<html>
		  <body>
		    <h1>Welcome, Chirpy Admin</h1>
		    <p>Chirpy has been visited %d times!</p>
		  </body>
		</html>`, cfg.fileserverHits.Load())
		if _, err := io.WriteString(w, response); err != nil {
			log.Fatal("Response could not be written")
		}
	})
}

func (cfg *ApiConfig) ResetMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		cfg.fileserverHits.Store(0)
	})
}
