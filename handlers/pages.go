package handlers

import (
	"ioignition/internal/database"
	"ioignition/view"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request, u database.User) {
	view.Layout(view.Home(u.Email)).Render(r.Context(), w)
}

func (h *Handler) LandingPage(w http.ResponseWriter, r *http.Request) {
	view.Layout(view.Nav("")).Render(r.Context(), w)
}

func (h *Handler) DomainStatsPage(w http.ResponseWriter, r *http.Request, u database.User) {
	domain := chi.URLParam(r, "domain")

	if isHtmxRequest(r) {
		view.DomainStats(u.Email, domain).Render(r.Context(), w)
		return
	}

	view.Layout(view.DomainStats(u.Email, domain)).Render(r.Context(), w)
}

func (h *Handler) SignupForm(w http.ResponseWriter, r *http.Request) {
	view.Layout(view.Signup()).Render(r.Context(), w)
}

func (h *Handler) LoginForm(w http.ResponseWriter, r *http.Request) {
	view.Layout(view.Login()).Render(r.Context(), w)
}

func (h *Handler) AddDomainPage(w http.ResponseWriter, r *http.Request, u database.User) {
	view.Layout(view.AddDomain(u.Email)).Render(r.Context(), w)
}
