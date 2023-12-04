package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"ioignition/internal/database"
	"ioignition/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// webhook
func (h *Handler) StatEvent(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Event         string `json:"event,omitempty"`
		Url           string `json:"url,omitempty"`
		Domain        string `json:"domain,omitempty"`
		Referrer      string `json:"referrer,omitempty"`
		Width         int    `json:"width,omitempty"`
		Agent         string `json:"agent,omitempty"`
		Sessionsstart string `json:"sessionsstart,omitempty"`
	}

	eventId := chi.URLParam(r, "eventId")

	if eventId == "" {
		log.Print("Event id missing")
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("event id cannot be empty"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	body := &reqBody{}

	err := decoder.Decode(body)
	if err != nil {
		log.Printf("Error decoding json: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	startTime, err := time.Parse(time.RFC1123, body.Sessionsstart)
	if err != nil {
		log.Printf("Error parsing session start time: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	// TODO: refactor, duplicated from domains.go
	// strip 'www.'
	body.Domain = strings.Trim(body.Domain, "www.")

	// url.Parse does not work as expected without scheme
	if !strings.Contains(body.Domain, "://") {
		body.Domain = "https://" + body.Domain
	}

	u, err := url.Parse(body.Domain)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("domain is not valid"))
		return
	}

	domain, err := h.dbQueries.GetDomain(r.Context(), u.Host)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusBadRequest, errors.New("domain is not registered"))
			return
		}

		log.Printf("Error getting domain: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	// there is no session end time as the StatEvent api only registers the
	// start of a session
	sessionParam := database.CreateDomainSessionParams{
		ID:               uuid.New(),
		EventID:          eventId,
		DomainID:         domain.ID,
		SessionStartTime: startTime,
		UpdatedAt:        time.Now(),
		CreatedAt:        time.Now(),
	}

	tx, err := h.db.Begin()
	if err != nil {
		log.Print("Error starting transaction: err")
		utils.RespondWithInternalServerError(w)
		return
	}

	defer tx.Rollback()

	qtx := h.dbQueries.WithTx(tx)
	session, err := qtx.CreateDomainSession(r.Context(), sessionParam)
	if err != nil {
		log.Printf("Error creating session: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	statParam := database.CreateDomainStatParams{
		ID:              uuid.New(),
		Url:             u.String(),
		Referer:         h.NewNullString(body.Referrer),
		DeviceWidth:     sql.NullInt32{Int32: int32(body.Width), Valid: true},
		DomainSessionID: session.ID,
		UserAgent:       h.NewNullString(body.Agent),
		UpdatedAt:       time.Now(),
		CreatedAt:       time.Now(),
	}

	_, err = qtx.CreateDomainStat(r.Context(), statParam)
	if err != nil {
		log.Printf("Error creating stat entry: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error commiting transaction: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}
}
