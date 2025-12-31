package handler

import (
	"chirpy/internal/api"
	"chirpy/internal/database"
	"chirpy/middleware"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type MetricsHandler struct {
	queries *database.Queries
}

func NewMetricsHandler(queries *database.Queries) *MetricsHandler {
	return &MetricsHandler{queries: queries}
}

func (h *MetricsHandler) GetMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileserverHits, err := middleware.GetFileServerHitsFromContext(r.Context())
		if err != nil {
			log.Fatal("Failed to get file server hits")
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf(`
		<html>
		  <body>
		    <h1>Welcome, Chirpy Admin</h1>
		    <p>Chirpy has been visited %d times!</p>
		  </body>
		</html>`, fileserverHits.Load())
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
				api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}
		}
		fileserverHits, err := middleware.GetFileServerHitsFromContext(r.Context())
		if err != nil {
			log.Fatal("Failed to get file server hits")
		}

		fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
	})
}
