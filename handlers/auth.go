package handlers

import (
	"ioignition/token"
	"log"
	"net/http"
)

func (h *Handler) Authed(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := token.GetSessionCookie(r, sessionCookieName)
		if err != nil {
			h.LandingPage(w, r)
			return
		}

		userId, err := h.token.VerifyToken(t, token.Access)
		if err != nil {
			h.LandingPage(w, r)
			return
		}

		// this is an internal error, should be handled better
		user, err := h.dbQueries.GetUserById(r.Context(), userId)
		if err != nil {
			log.Printf("Error getting user by id: %+v", err)
			h.LandingPage(w, r)
			return
		}

		handler(w, r, user)
	}
}
