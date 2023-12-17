package main

import (
	"ioignition/handlers"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router, h *handlers.Handler) {
	r.Get("/", h.Authed(h.HomePage))

	// Account routes
	r.Get("/sign-up", h.SignupForm)
	r.Get("/login", h.LoginForm)

	// User routes
	r.Post("/sign-up", h.CreateUser)
	r.Post("/login", h.Login)

	// Domain routes
	r.Get("/domain", h.Authed(h.AddDomainPage))
	r.Post("/domain", h.Authed(h.AddDomain))
	r.Get("/domains", h.Authed(h.ListDomains))

	// Domain Stats Routes
	r.Get("/{domain}", h.Authed(h.DomainStatsPage))
	r.Get("/{domainId}/stats/{interval}-{unit}", h.Authed(h.DomainStats))
	r.Get("/{domainId}/graph/{interval}-{unit}", h.Authed(h.GraphStats))
	r.Get("/{domainId}/online", h.Authed(h.GetOnlineCount))

	r.Get("/{domainId}/urlvisits/{interval}-{unit}", h.Authed(h.UrlVisits))
	r.Get("/{domainId}/referers/{interval}-{unit}", h.Authed(h.RefererCount))
	r.Get("/{domainId}/platforms/{interval}-{unit}", h.Authed(h.PlatformCount))
	r.Get("/{domainId}/browsers/{interval}-{unit}", h.Authed(h.BrowserCount))
	r.Get("/{domainId}/locations/{interval}-{unit}", h.Authed(h.LocationCount))
}

// the following requests are made my the script running on users client
func registerApiRoutes(r chi.Router, h *handlers.Handler) {
	// start recording a new session
	r.Post("/event/{clientId}", h.StatEvent)
	// update the url change in the session created above
	r.Post("/event/{clientId}/url", h.StatUrlUpdate)
	// end the session created in /event/{clientId} route
	r.Post("/event/{clientId}/end", h.StatEndSession)
}
