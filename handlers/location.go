package handlers

import (
	"context"
	"fmt"
	"ioignition/internal/database"
	"ioignition/utils"
	"ioignition/view"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetLocation(ip string, sessionId uuid.UUID, domainId uuid.UUID) {
	ctx := context.Background()

	cachedInfo := h.client.HGetAll(ctx, ip).Val()

	if len(cachedInfo) == 0 {
		// get ip info
		info, err := h.ipClient.GetIPInfo(net.ParseIP(ip))
		if err != nil {
			log.Println("Error getting ip info: ", err)
			return
		}

		intermediate := map[string]string{
			"Emoji":       info.CountryFlag.Emoji,
			"Region":      info.Region,
			"CountryCode": info.Country,
			"Name":        info.CountryName,
		}

		for k, v := range intermediate {
			err := h.client.HSet(ctx, ip, k, v).Err()
			if err != nil {
				log.Print("Error during client.HSet: ", err)
				return
			}
		}

		cachedInfo = h.client.HGetAll(ctx, ip).Val()
		if len(cachedInfo) == 0 {
			log.Println("Cache is still empty")
			return
		}
	}

	param := database.CreateLocationParams{
		ID:              uuid.New(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		DomainSessionID: sessionId,
		DomainID:        domainId,
		Emoji:           cachedInfo["Emoji"],
		CountryCode:     cachedInfo["CountryCode"],
		Region:          cachedInfo["Region"],
		Name:            cachedInfo["Name"],
	}

	_, err := h.dbQueries.CreateLocation(ctx, param)
	if err != nil {
		log.Print("Error creating location info: ", err)
	}
}

func (h *Handler) LocationCount(w http.ResponseWriter, r *http.Request, user database.User) {
	id := chi.URLParam(r, "domainId")

	domainId, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid domain id %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	i := chi.URLParam(r, "interval")
	u := chi.URLParam(r, "unit")

	interval, _ := h.validator.TimeRange(i, u)

	param := database.GetLocationCountParams{
		DomainID:  domainId,
		CreatedAt: time.Now().AddDate(0, 0, -interval),
	}

	locations, err := h.dbQueries.GetLocationCount(r.Context(), param)
	if err != nil {
		log.Printf("Error getting platform stats: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	var stats []view.Stat
	// ordered desc
	var max int

	if len(locations) > 0 {
		max = int(locations[0].LocationCount)
	}

	for _, l := range locations {
		s := view.Stat{
			Value:   fmt.Sprintf("%s %s", l.Emoji, l.Name),
			Count:   int(l.LocationCount),
			Percent: h.getPercent(float64(l.LocationCount), float64(max)),
		}

		stats = append(stats, s)
	}

	err = view.Stats("Countries", stats).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering view.Stats:", err)
	}
}
