package api

import (
	"chirpy/lib"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var blacklist = []string{"kerfuffle", "sharbert", "fornax"}

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

	response := map[string]string{
		"cleaned_body": replaceNotAllowedWords(payload.Body),
	}

	lib.RespondWithJSON(w, http.StatusOK, response)
}

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
