package handlers

import (
	"ioignition/internal/database"
	"ioignition/utils"
	"ioignition/view"
	"log"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) getPercent(count, max float64) int {
	percent := math.Round((count / (math.Max(1, max))) * 100)

	return int(percent)
}

func (h *Handler) UrlVisits(w http.ResponseWriter, r *http.Request, user database.User) {
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

	param := database.GetPageViewsParams{
		DomainID:  domainId,
		CreatedAt: time.Now().AddDate(0, 0, -interval),
	}

	views, err := h.dbQueries.GetPageViews(r.Context(), param)
	if err != nil {
		log.Printf("Error getting stats: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	var stats []view.Stat
	// ordered desc
	max := views[0].UrlCount

	for _, v := range views {
		u, err := url.Parse(v.Url)
		if err != nil {
			log.Print("Invalid url saved: ", err)
		}

		s := view.Stat{
			Value:   u.Path,
			Count:   int(v.UrlCount),
			Percent: h.getPercent(float64(v.UrlCount), float64(max)),
		}

		stats = append(stats, s)
	}

	err = view.Stats("Page Views", stats).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering view.Stats:", err)
	}
}
