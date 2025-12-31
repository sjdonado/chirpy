package handler

import (
	"chirpy/internal/api"
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/middleware"
	"chirpy/serializer"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type UsersHandler struct {
	cfg     *api.Config
	queries *database.Queries
}

func NewUsersHandler(config *api.Config, queries *database.Queries) *UsersHandler {
	return &UsersHandler{cfg: config, queries: queries}
}

func (h *UsersHandler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := json.NewDecoder(r.Body)
		defer r.Body.Close()

		payload := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		if err := body.Decode(&payload); err != nil {
			api.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		hashedPassword, err := auth.HashPassword(payload.Password)
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		user, err := h.queries.CreateUser(r.Context(), database.CreateUserParams{
			Email:          payload.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		api.RespondWithJSON(w, http.StatusCreated, serializer.SerializeUser(user))
	}
}

func (h *UsersHandler) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := json.NewDecoder(r.Body)
		defer r.Body.Close()

		payload := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		if err := body.Decode(&payload); err != nil {
			api.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		userID, err := middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "User ID not found in context")
			return
		}

		hashedPassword, err := auth.HashPassword(payload.Password)
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		user, err := h.queries.UpdateUser(r.Context(), database.UpdateUserParams{
			ID:             userID,
			Email:          payload.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		api.RespondWithJSON(w, http.StatusOK, serializer.SerializeUser(user))
	}
}

func (h *UsersHandler) UpgradeUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := json.NewDecoder(r.Body)
		defer r.Body.Close()

		payload := struct {
			Event string `json:"event"`
			Data  struct {
				UserID uuid.UUID `json:"user_id"`
			} `json:"data"`
		}{}

		if err := body.Decode(&payload); err != nil {
			api.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if payload.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		user, err := h.queries.GetUser(r.Context(), payload.Data.UserID)
		if err != nil {
			api.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		if user.IsChirpyRed {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		_, err = h.queries.UpgradeUser(r.Context(), database.UpgradeUserParams{
			ID:          payload.Data.UserID,
			IsChirpyRed: true,
		})
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
