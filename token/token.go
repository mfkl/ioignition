package token

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	issuer        = "IoIgnition"
	day           = 24 * time.Hour
	TokenValidFor = 30 * day
	Access        = "access"
	Refresh       = "refresh"
)

var (
	ErrNoSessionCookie = errors.New("no session cookie included in request")
	ErrInvalidCookie   = errors.New("cookie is not valid")
	ErrInvalidToken    = errors.New("token is not valid")
)

type Token struct {
	secret []byte
}

func NewToken(secret string) *Token {
	return &Token{[]byte(secret)}
}

// issuer type can be Access or Refresh
func (t *Token) Create(id uuid.UUID, issuerType string) (string, error) {
	now := time.Now()
	expires := now.Add(TokenValidFor)

	if issuerType == "refresh" {
		expires = now.Add(TokenValidFor)
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    issuer + "-" + issuerType,
		ExpiresAt: jwt.NewNumericDate(expires),
		IssuedAt:  jwt.NewNumericDate(now),
		Subject:   id.String(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString(t.secret)
}

func (t *Token) VerifyToken(token, issuerType string) (uuid.UUID, error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) { return t.secret, nil },
	)
	if err != nil {
		log.Print("ParseWithClaims: ", err)
		return uuid.Nil, err
	}

	if !jwtToken.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	i, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if i != issuer+"-"+issuerType {
		return uuid.Nil, ErrInvalidToken
	}

	idString, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	return id, nil
}

// Get Cookie
func GetSessionCookie(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", ErrNoSessionCookie
	}

	if e := c.Valid(); e != nil {
		return "", ErrInvalidCookie
	}

	return c.Value, nil
}
