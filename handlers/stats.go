package handlers

import (
	"ioignition/internal/database"
	"ioignition/utils"
	"ioignition/view"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) DomainStats(w http.ResponseWriter, r *http.Request, u database.User) {
	domain := chi.URLParam(r, "domain")

	d, err := h.dbQueries.GetDomain(r.Context(), domain)
	if err != nil {
		log.Printf("Error decoding json: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	intervalString := chi.URLParam(r, "interval")
	unit := chi.URLParam(r, "unit")

	// default to 30 days
	if intervalString == "" || unit == "" {
		intervalString = "30"
		unit = "D"
	}

	interval, err := strconv.Atoi(intervalString)
	if err != nil {
		interval = 30
		unit = "D"
	}

	param := database.GetSessionStatsParams{
		DomainID:        d.ID,
		Interval:        time.Now().AddDate(0, 0, -interval),
		CompareInterval: time.Now().AddDate(0, 0, -(2 * interval)),
	}

	stats, err := h.dbQueries.GetSessionStats(r.Context(), param)
	if err != nil {
		log.Printf("Error getting stats: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	urlStats, err := h.dbQueries.GetPageViewCount(r.Context(), database.GetPageViewCountParams(param))
	if err != nil {
		log.Printf("Error getting url stats: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	view.Metrics(stats, urlStats).Render(r.Context(), w)
}

func (h *Handler) GraphStats(w http.ResponseWriter, r *http.Request, u database.User) {
	domain := chi.URLParam(r, "domain")

	d, err := h.dbQueries.GetDomain(r.Context(), domain)
	if err != nil {
		log.Printf("Error decoding json: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	intervalString := chi.URLParam(r, "interval")
	unit := chi.URLParam(r, "unit")

	// default to 30 days
	if intervalString == "" || unit == "" {
		intervalString = "30"
		unit = "D"
	}

	interval, err := strconv.Atoi(intervalString)
	if err != nil {
		interval = 30
		unit = "D"
	}

	// FIXME: right now I'm hard coding the step but depending on the unit, choose
	// a strategy that will return months, hours, years etc
	param := database.GetGraphStatsParams{
		DomainID:      d.ID,
		Step:          5,
		IntervalStart: time.Now().AddDate(0, 0, -interval),
	}

	stats, err := h.dbQueries.GetGraphStats(r.Context(), param)
	if err != nil {
		log.Printf("Error getting stats: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	var labels []string
	var dataPoints []int

	for _, s := range stats {
		labels = append(labels, formattedTime(s.DrStartDate))
		dataPoints = append(dataPoints, int(s.SessionCount))
	}

	view.Graph(labels, dataPoints, "Visitors").Render(r.Context(), w)
}

func formattedTime(t time.Time) string {
	layout := "02 Jan"

	return t.Format(layout)
}
