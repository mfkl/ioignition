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
}

func registerApiRoutes(r chi.Router, h *handlers.Handler) {
	// start recording a new session
	r.Post("/event/{clientId}", h.StatEvent)
	// update the url change in the session created above
	r.Post("/event/{clientId}/url", h.StatUrlUpdate)
	// end the session created in /event/{clientId} route
	r.Post("/event/{clientId}/end", h.StatEndSession)
}
