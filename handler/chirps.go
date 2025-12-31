package handler

import (
	"chirpy/internal/database"
	"chirpy/lib"
	"chirpy/serializer"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type ChirpsHandler struct {
	queries *database.Queries
}

func NewChirpsHandler(queries *database.Queries) *ChirpsHandler {
	return &ChirpsHandler{queries: queries}
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

func (h *ChirpsHandler) CreateChirp() http.Handler {
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

		chirp, err := h.queries.CreateChirp(r.Context(), database.CreateChirpParams{Body: replaceNotAllowedWords(payload.Body), UserID: payload.UserID})
		if err != nil {
			lib.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		lib.RespondWithJSON(w, http.StatusCreated, serializer.SerializeChirp(chirp))
	})
}

func (h *ChirpsHandler) GetAllChirps() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirps, err := h.queries.GetAllChirps(r.Context())
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

func (h *ChirpsHandler) GetOneChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			lib.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		chirp, err := h.queries.GetOneChirp(r.Context(), id)
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
