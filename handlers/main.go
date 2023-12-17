package handlers

import (
	"database/sql"
	"ioignition/internal/database"
	"ioignition/token"
	"ioignition/validator"
	"net/http"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/redis/go-redis/v9"
)

const (
	UniqueViolationCode = "23505"
)

// used in the middleware auth to validate requests
type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

type Handler struct {
	dbQueries *database.Queries
	client    *redis.Client
	ipClient  *ipinfo.Client
	token     *token.Token
	validator *validator.Validator
}

func NewHandler(
	dbQueries *database.Queries,
	client *redis.Client,
	ipClient *ipinfo.Client,
	jwtSecret string,
) *Handler {
	t := token.NewToken(jwtSecret)
	validator := validator.NewValidator()

	return &Handler{dbQueries, client, ipClient, t, validator}
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
