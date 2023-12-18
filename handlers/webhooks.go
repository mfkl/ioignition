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
)

// helper
func (h *Handler) PersistUrl(r *http.Request, url string, event string, sessionId uuid.UUID) (database.Url, error) {
	urlParam := database.CreateSessionUrlParams{
		ID:              uuid.New(),
		Url:             url,
		EventName:       event,
		DomainSessionID: sessionId,
		UpdatedAt:       time.Now(),
		CreatedAt:       time.Now(),
	}

	return h.dbQueries.CreateSessionUrl(r.Context(), urlParam)
}

// webhook
func (h *Handler) StatEvent(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		SessionId string `json:"sessionId,omitempty"`
		Event     string `json:"event,omitempty"`
		// Note: domain is the registered domain against which you would check if
		// it's registered to be using ioignition analytics
		Domain string `json:"domain,omitempty"`
		// URL is the url the analytics data was sent from
		Url      string `json:"url,omitempty"`
		Referrer string `json:"referrer,omitempty"`
		Width    int    `json:"width,omitempty"`
		Browser  string `json:"browser,omitempty"`
		Platform string `json:"platform,omitempty"`
	}

	clientId := chi.URLParam(r, "clientId")

	if clientId == "" {
		log.Print("Client id missing")
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("client id cannot be empty"))
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

	// TODO: refactor, duplicated from domains.go
	// strip 'www.'
	body.Domain = strings.TrimPrefix(body.Domain, "www.")

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
		ClientID:         clientId,
		SessionID:        body.SessionId,
		DomainID:         domain.ID,
		SessionStartTime: time.Now(),
		Referer:          h.NewNullString(body.Referrer),
		DeviceWidth:      sql.NullInt32{Int32: int32(body.Width), Valid: true},
		Browser:          h.NewNullString(body.Browser),
		Platform:         h.NewNullString(body.Platform),
		UpdatedAt:        time.Now(),
		CreatedAt:        time.Now(),
	}

	session, err := h.dbQueries.CreateDomainSession(r.Context(), sessionParam)
	if err != nil {
		log.Printf("Error creating session: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	ip := r.Header.Get("X-FORWARDED-FOR")
	// call get location, can be concurrent as info not needed immediately
	go h.GetLocation(ip, session.ID, domain.ID)

	_, err = h.PersistUrl(r, body.Url, body.Event, session.ID)
	if err != nil {
		log.Printf("Error creating stat entry: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, nil)
}

// webhook
func (h *Handler) StatUrlUpdate(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		SessionId string `json:"sessionId,omitempty"`
		Event     string `json:"event,omitempty"`
		Url       string `json:"url,omitempty"`
	}

	clientId := chi.URLParam(r, "clientId")

	if clientId == "" {
		log.Print("Client id missing")
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("client id cannot be empty"))
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

	getSessionParam := database.GetDomainSessionParams{
		SessionID: body.SessionId,
		ClientID:  clientId,
	}

	session, err := h.dbQueries.GetDomainSession(r.Context(), getSessionParam)
	if err != nil {
		log.Printf("Error getting session: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	_, err = h.PersistUrl(r, body.Url, body.Event, session.ID)
	if err != nil {
		log.Printf("Error creating stat entry: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, nil)
}

// webhook
func (h *Handler) StatEndSession(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		SessionId string `json:"sessionId,omitempty"`
		Event     string `json:"event,omitempty"`
	}

	clientId := chi.URLParam(r, "clientId")

	if clientId == "" {
		log.Print("Client id missing")
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("client id cannot be empty"))
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

	getSessionParam := database.GetDomainSessionParams{
		SessionID: body.SessionId,
		ClientID:  clientId,
	}

	session, err := h.dbQueries.GetDomainSession(r.Context(), getSessionParam)
	if err != nil {
		log.Printf("Error getting session: %+v", err)
		return
	}

	// there is no session end time as the StatEvent api only registers the
	// start of a session
	updateSessionParam := database.UpdateDomainSessionParams{
		ID:             session.ID,
		SessionEndTime: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:      time.Now(),
	}

	err = h.dbQueries.UpdateDomainSession(r.Context(), updateSessionParam)
	if err != nil {
		log.Printf("Error ending session: %+v", err)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, nil)
}
