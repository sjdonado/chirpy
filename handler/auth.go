package handler

import (
	"chirpy/internal/api"
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/serializer"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct {
	queries *database.Queries
}

func NewAuthHandler(queries *database.Queries) *AuthHandler {
	return &AuthHandler{queries: queries}
}

func (h *AuthHandler) Login() http.HandlerFunc {
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
		user, err := h.queries.GetUserByEmail(r.Context(), payload.Email)
		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		if success, err := auth.CheckPasswordHash(payload.Password, user.HashedPassword); err != nil || !success {
			api.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		accessToken, err := auth.MakeJWT(user.ID, os.Getenv("JWT_SECRET"), time.Hour)
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = h.queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			UserID:    user.ID,
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		})
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		api.RespondWithJSON(w, http.StatusOK, serializer.SerializeLoginResponse(user, accessToken, refreshToken))
	}
}

func (h *AuthHandler) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawRefreshToken, err := api.GetBearerToken(r.Header)
		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		refreshToken, err := h.queries.GetRefreshtoken(r.Context(), rawRefreshToken)
		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if refreshToken.RevokedAt.Valid {
			api.RespondWithError(w, http.StatusUnauthorized, "Token revoked")
			return
		}

		if refreshToken.ExpiresAt.Before(time.Now()) {
			api.RespondWithError(w, http.StatusUnauthorized, "Token expired")
			return
		}

		accessToken, err := auth.MakeJWT(refreshToken.UserID, os.Getenv("JWT_SECRET"), time.Hour)
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		api.RespondWithJSON(w, http.StatusOK, map[string]string{
			"token": accessToken,
		})
	}
}

func (h *AuthHandler) RevokeToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawRefreshToken, err := api.GetBearerToken(r.Header)
		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		refreshToken, err := h.queries.GetRefreshtoken(r.Context(), rawRefreshToken)
		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// already revoked
		if refreshToken.RevokedAt.Valid {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		err = h.queries.RevokeRefreshToken(r.Context(), rawRefreshToken)
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
