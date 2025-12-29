package handler

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/lib"
	"chirpy/serializer"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func NewApiConfig() (*ApiConfig, func(), error) {
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, err
	}

	cfg := &ApiConfig{db: database.New(db)}
	cleanup := func() {
		_ = db.Close()
	}

	return cfg, cleanup, nil
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
		if os.Getenv("PLATFORM") == "dev" {
			err := cfg.db.DeleteAllUsers(r.Context())
			if err != nil {
				lib.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}
		}
		w.WriteHeader(http.StatusOK)
		cfg.fileserverHits.Store(0)
	})
}

func (cfg *ApiConfig) CreateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := json.NewDecoder(r.Body)
		defer r.Body.Close()

		payload := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		if err := body.Decode(&payload); err != nil {
			lib.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		hashedPassword, err := auth.HashPassword(payload.Password)
		if err != nil {
			lib.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
			Email:          payload.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			lib.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		lib.RespondWithJSON(w, http.StatusCreated, serializer.SerializeUser(user))
	})
}

var blacklist = []string{"kerfuffle", "sharbert", "fornax"}

func replaceNotAllowedWords(body string) string {
	for word := range strings.SplitSeq(body, " ") {
		for _, badWord := range blacklist {
			if strings.ToLower(word) == badWord {
				body = strings.ReplaceAll(body, word, "****")
			}
		}
	}
	return body
}

func (cfg *ApiConfig) CreateChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := json.NewDecoder(r.Body)
		defer r.Body.Close()

		payload := struct {
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}{}

		if err := body.Decode(&payload); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			lib.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
			return
		}

		if len(payload.Body) > 140 {
			lib.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
			return
		}

		chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: replaceNotAllowedWords(payload.Body), UserID: payload.UserID})
		if err != nil {
			lib.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		lib.RespondWithJSON(w, http.StatusCreated, serializer.SerializeChirp(chirp))
	})
}

func (cfg *ApiConfig) GetAllChirps() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirps, err := cfg.db.GetAllChirps(r.Context())
		if err != nil {
			lib.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		serializedChirps := []serializer.Chirp{}
		for _, chirp := range chirps {
			serializedChirps = append(serializedChirps, serializer.SerializeChirp(chirp))
		}

		lib.RespondWithJSON(w, http.StatusOK, serializedChirps)
	})
}

func (cfg *ApiConfig) GetOneChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			lib.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		chirp, err := cfg.db.GetOneChirp(r.Context(), id)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				lib.RespondWithError(w, http.StatusNotFound, "Chirp not found")
			default:
				lib.RespondWithError(w, http.StatusNotFound, err.Error())
			}
			return
		}

		lib.RespondWithJSON(w, http.StatusOK, serializer.SerializeChirp(chirp))
	})
}
