package handlers

import (
	"ioignition/internal/database"
	"ioignition/token"
	"net/http"
)

const UniqueViolationCode = "23505"

// used in the middleware auth to validate requests
type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

type Handler struct {
	db    *database.Queries
	token *token.Token
}

func NewHandler(db *database.Queries, jwtSecret string) *Handler {
	t := token.NewToken(jwtSecret)

	return &Handler{db, t}
}
