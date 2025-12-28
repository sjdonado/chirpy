package handler

import (
	"io"
	"log"
	"net/http"
)

func GetHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, "OK\n"); err != nil {
		log.Fatal("Response could not be written")
	}
}
