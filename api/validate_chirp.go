package api

import (
	"chirpy/lib"
	"encoding/json"
	"log"
	"net/http"
)

func PostValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body := json.NewDecoder(r.Body)
	payload := struct {
		Body string `json:"body"`
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

	lib.RespondWithJSON(w, http.StatusOK, map[string]bool{"valid": true})
}
