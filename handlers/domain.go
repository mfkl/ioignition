package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"ioignition/internal/database"
	"ioignition/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type domainResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Url       string    `json:"url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func mapDomainResponse(d *database.Domain) *domainResponse {
	return &domainResponse{
		ID:        d.ID,
		Url:       d.Url,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (h *Handler) AddDomain(w http.ResponseWriter, r *http.Request, user database.User) {
	type reqBody struct {
		Domain string `json:"domain,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	body := &reqBody{}

	err := decoder.Decode(body)
	if err != nil {
		log.Printf("Error decoding json: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	// url.Parse does not work as expected without scheme
	if !strings.Contains(body.Domain, "://") {
		body.Domain = "https://" + body.Domain
	}

	domain, err := url.Parse(body.Domain)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("domain is not valid"))
		return
	}

	param := database.CreateDomainParams{
		ID:        uuid.New(),
		Url:       fmt.Sprintf("%s://%s", domain.Scheme, domain.Host),
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	d, err := h.db.CreateDomain(r.Context(), param)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == UniqueViolationCode {
				utils.RespondWithError(w, http.StatusBadRequest, errors.New("domain already exists"))
				return
			}
		}
		log.Printf("Error saving domain: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, mapDomainResponse(&d))
}
