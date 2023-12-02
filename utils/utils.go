package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

var InternalServerError = errors.New("internal server error")

type errRes struct {
	Error string `json:"error,omitempty"`
}

func RespondWithInternalServerError(w http.ResponseWriter) {
	RespondWithError(w, http.StatusInternalServerError, InternalServerError)
}

func RespondWithError(w http.ResponseWriter, status int, err error) {
	d := errRes{Error: err.Error()}

	RespondWithJson(w, status, d)
}

func RespondWithJson(w http.ResponseWriter, status int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Marshal error %s", err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}

func GetAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	headerComponents := strings.Split(authHeader, " ")

	if len(headerComponents) != 2 || headerComponents[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if headerComponents[1] == "" {
		return "", errors.New("missing access token")
	}

	return headerComponents[1], nil
}
