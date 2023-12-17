package handlers

import (
	"ioignition/internal/database"
	"ioignition/utils"
	"ioignition/view"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func getFaviconUrl(host string) string {
	return "https://www.google.com/s2/favicons?domain=" + host
}

func formattedTime(t time.Time) string {
	layout := "02 Jan"

	return t.Format(layout)
}

func (h *Handler) DomainStats(w http.ResponseWriter, r *http.Request, u database.User) {
	id := chi.URLParam(r, "domainId")

	domainId, err := uuid.Parse(id)
	if err != nil {
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
		DomainID:        domainId,
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

	err = view.Metrics(stats, urlStats).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering view.Metrics:", err)
	}
}

func (h *Handler) GraphStats(w http.ResponseWriter, r *http.Request, u database.User) {
	id := chi.URLParam(r, "domainId")

	domainId, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid domain id %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	intervalString := chi.URLParam(r, "interval")
	unit := chi.URLParam(r, "unit")

	interval, _ := h.validator.TimeRange(intervalString, unit)

	// FIXME: right now I'm hard coding the step but depending on the unit, choose
	// a strategy that will return months, hours, years etc
	param := database.GetGraphStatsParams{
		DomainID:      domainId,
		Step:          5,
		IntervalStart: int32(interval),
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

	err = view.Graph(labels, dataPoints, "Visitors").Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering view.Graph:", err)
	}
}

func (h *Handler) GetOnlineCount(w http.ResponseWriter, r *http.Request, u database.User) {
	id := chi.URLParam(r, "domainId")

	domainId, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid domain id %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	param := database.GetCurrentlyActiveUsersParams{
		DomainID:         domainId,
		SessionStartTime: time.Now().AddDate(0, 0, -1),
	}

	active, err := h.dbQueries.GetCurrentlyActiveUsers(r.Context(), param)
	if err != nil {
		log.Printf("Error getting active users: %+v", err)

		err = view.CurrentlyOnline(0).Render(r.Context(), w)
		log.Printf("Error rendering active users: %+v", err)
		return
	}

	err = view.CurrentlyOnline(int(active.Count)).Render(r.Context(), w)
	if err != nil {
		log.Printf("Error rendering active users: %+v", err)
	}
}

func (h *Handler) RefererCount(w http.ResponseWriter, r *http.Request, user database.User) {
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

	param := database.GetRefererStatsParams{
		DomainID:  domainId,
		CreatedAt: time.Now().AddDate(0, 0, -interval),
	}

	referers, err := h.dbQueries.GetRefererStats(r.Context(), param)
	if err != nil {
		log.Printf("Error getting stats: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	var stats []view.Stat
	var max int

	// ordered desc
	if len(referers) > 0 {
		max = int(referers[0].RefererCount)
	}

	for _, r := range referers {
		u, err := url.Parse(r.Referer.String)
		if err != nil {
			log.Print("Invalid url saved: ", err)
		}

		s := view.Stat{
			Value:    u.Host,
			Count:    int(r.RefererCount),
			Percent:  h.getPercent(float64(r.RefererCount), float64(max)),
			HasImage: true,
			Img:      getFaviconUrl(u.Host),
		}

		stats = append(stats, s)
	}

	err = view.Stats("Top Sources", stats).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering view.Stats:", err)
	}
}

func (h *Handler) PlatformCount(w http.ResponseWriter, r *http.Request, user database.User) {
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

	param := database.GetPlatformStatsParams{
		DomainID:  domainId,
		CreatedAt: time.Now().AddDate(0, 0, -interval),
	}

	platforms, err := h.dbQueries.GetPlatformStats(r.Context(), param)
	if err != nil {
		log.Printf("Error getting platform stats: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	var stats []view.Stat
	// ordered desc
	var max int

	if len(platforms) > 0 {
		max = int(platforms[0].PlatformCount)
	}

	for _, p := range platforms {
		s := view.Stat{
			Value:   p.Platform.String,
			Count:   int(p.PlatformCount),
			Percent: h.getPercent(float64(p.PlatformCount), float64(max)),
		}

		stats = append(stats, s)
	}

	err = view.Stats("Operating Systems", stats).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering view.Stats:", err)
	}
}

func (h *Handler) BrowserCount(w http.ResponseWriter, r *http.Request, user database.User) {
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

	param := database.GetBrowserStatsParams{
		DomainID:  domainId,
		CreatedAt: time.Now().AddDate(0, 0, -interval),
	}

	browsers, err := h.dbQueries.GetBrowserStats(r.Context(), param)
	if err != nil {
		log.Printf("Error getting platform stats: %+v", err)
		utils.RespondWithInternalServerError(w)
		return
	}

	var stats []view.Stat
	// ordered desc
	var max int

	if len(browsers) > 0 {
		max = int(browsers[0].BrowserCount)
	}

	for _, b := range browsers {
		s := view.Stat{
			Value:   b.Browser.String,
			Count:   int(b.BrowserCount),
			Percent: h.getPercent(float64(b.BrowserCount), float64(max)),
		}

		stats = append(stats, s)
	}

	err = view.Stats("Browsers", stats).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering view.Stats:", err)
	}
}
