package handler

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/lib"
	"chirpy/serializer"
	"encoding/json"
	"net/http"
)

type UsersHandler struct {
	queries *database.Queries
}

func NewUsersHandler(queries *database.Queries) *UsersHandler {
	return &UsersHandler{queries: queries}
}

func (h *UsersHandler) CreateUser() http.Handler {
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

		user, err := h.queries.CreateUser(r.Context(), database.CreateUserParams{
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

func (h *UsersHandler) Login() http.Handler {
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
		user, err := h.queries.GetUserByEmail(r.Context(), payload.Email)
		if err != nil {
			lib.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		if success, err := auth.CheckPasswordHash(payload.Password, user.HashedPassword); err != nil || !success {
			lib.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		lib.RespondWithJSON(w, http.StatusOK, serializer.SerializeUser(user))
	})
}
