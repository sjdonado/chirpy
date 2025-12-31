package handler

import (
	"chirpy/internal/database"
	"chirpy/lib"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type MetricsHandler struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
}

func NewMetricsHandler(queries *database.Queries) *MetricsHandler {
	return &MetricsHandler{queries: queries}
}

func (h *MetricsHandler) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (h *MetricsHandler) GetMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf(`
		<html>
		  <body>
		    <h1>Welcome, Chirpy Admin</h1>
		    <p>Chirpy has been visited %d times!</p>
		  </body>
		</html>`, h.fileserverHits.Load())
		if _, err := io.WriteString(w, response); err != nil {
			log.Fatal("Response could not be written")
		}
	})
}

func (h *MetricsHandler) ResetMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("PLATFORM") == "dev" {
			err := h.queries.DeleteAllUsers(r.Context())
			if err != nil {
				lib.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}
		}
		w.WriteHeader(http.StatusOK)
		h.fileserverHits.Store(0)
	})
}
