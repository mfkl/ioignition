package handlers

import (
	"database/sql"
	"ioignition/internal/database"
	"ioignition/token"
	"net/http"
)

const UniqueViolationCode = "23505"

// used in the middleware auth to validate requests
type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

type Handler struct {
	db        *sql.DB
	dbQueries *database.Queries
	token     *token.Token
}

func NewHandler(db *sql.DB, dbQueries *database.Queries, jwtSecret string) *Handler {
	t := token.NewToken(jwtSecret)

	return &Handler{db, dbQueries, t}
}

// useful utils
func (h *Handler) NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// Render whole page with layout if not
func isHtmxRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
