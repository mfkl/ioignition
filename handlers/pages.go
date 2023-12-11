package handlers

import (
	"ioignition/internal/database"
	"ioignition/view"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request, u database.User) {
	if isHtmxRequest(r) {
		err := view.Home(u.Email).Render(r.Context(), w)
		if err != nil {
			log.Print("Error rendering view.Home:", err)
		}
		return
	}

	err := view.Layout(view.Home(u.Email)).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering HomePage:", err)
	}
}

func (h *Handler) LandingPage(w http.ResponseWriter, r *http.Request) {
	err := view.Layout(view.Nav("")).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering LandingPage:", err)
	}
}

func (h *Handler) DomainStatsPage(w http.ResponseWriter, r *http.Request, u database.User) {
	domain := chi.URLParam(r, "domain")

	if isHtmxRequest(r) {
		err := view.DomainStats(u.Email, domain).Render(r.Context(), w)
		if err != nil {
			log.Print("Error rendering view.DomainStats:", err)
		}

		return
	}

	err := view.Layout(view.DomainStats(u.Email, domain)).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering DomainStatsPage:", err)
	}
}

func (h *Handler) SignupForm(w http.ResponseWriter, r *http.Request) {
	err := view.Layout(view.Signup()).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering SignupForm:", err)
	}
}

func (h *Handler) LoginForm(w http.ResponseWriter, r *http.Request) {
	err := view.Layout(view.Login()).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering LoginForm:", err)
	}
}

func (h *Handler) AddDomainPage(w http.ResponseWriter, r *http.Request, u database.User) {
	if isHtmxRequest(r) {
		err := view.AddDomain(u.Email).Render(r.Context(), w)
		if err != nil {
			log.Print("Error rendering view.AddDomain:", err)
		}

		return
	}

	err := view.Layout(view.AddDomain(u.Email)).Render(r.Context(), w)
	if err != nil {
		log.Print("Error rendering AddDomainPage:", err)
	}
}
