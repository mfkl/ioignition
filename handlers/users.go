package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"ioignition/internal/database"
	"ioignition/token"
	"ioignition/utils"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email                string `json:"email,omitempty"`
		Password             string `json:"password,omitempty"`
		PasswordConfirmation string `json:"password_confirmation,omitempty"`
	}

	d := json.NewDecoder(r.Body)
	body := &reqBody{}

	err := d.Decode(body)
	if err != nil {
		log.Printf("Error decoding req (%s): %+v", r.Body, err)
		utils.RespondWithInternalServerError(w)
		return
	}

	email, err := mail.ParseAddress(body.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("email is invalid"))
		return
	}

	// simple password validator
	if strings.Compare(body.Password, body.PasswordConfirmation) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("passwords do not match"))
		return
	}

	if len(body.Password) < 12 {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("password should be 12 characters or more"))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password (%s): %+v", body.Password, err)
		utils.RespondWithInternalServerError(w)
		return
	}

	createUserParam := database.CreateUserParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        email.Address,
		PasswordHash: string(hash),
	}

	u, err := h.dbQueries.CreateUser(r.Context(), createUserParam)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == UniqueViolationCode {
				utils.RespondWithError(w, http.StatusBadRequest, errors.New("email already exists"))
				return
			}
		}

		log.Printf("Error creating user: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	accessToken, err := h.token.Create(u.ID, token.Access)
	if err != nil {
		log.Printf("Error getting access token: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	h.SetSesionCookie(w, r, accessToken)
	w.Header().Set("HX-Replace-Url", "/")
	// render home page
	h.HomePage(w, r, u)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	}

	d := json.NewDecoder(r.Body)
	body := &reqBody{}

	err := d.Decode(body)
	if err != nil {
		log.Printf("Error decoding req (%s): %+v", r.Body, err)
		utils.RespondWithInternalServerError(w)
		return
	}

	email, err := mail.ParseAddress(body.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("email is invalid"))
		return
	}

	u, err := h.dbQueries.GetUser(r.Context(), email.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusBadRequest, errors.New("email does not exists"))
			return
		}

		log.Printf("Error getting user: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(body.Password))
	if err != nil {
		utils.RespondWithError(w, http.StatusForbidden, errors.New("unable to login"))
		return
	}

	accessToken, err := h.token.Create(u.ID, token.Access)
	if err != nil {
		log.Printf("Error getting access token: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	// Possible Todos:
	// * handle refresh token, needed?
	// * save refresh token
	// * save token by device type?
	h.SetSesionCookie(w, r, accessToken)

	w.Header().Set("HX-Replace-Url", "/")

	// render home page
	h.HomePage(w, r, u)
}
