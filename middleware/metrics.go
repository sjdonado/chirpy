package middleware

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) GetMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf("Hits: %d\n", cfg.fileserverHits.Load())
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
