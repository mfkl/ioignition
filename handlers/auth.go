package handlers

import (
	"ioignition/token"
	"ioignition/utils"
	"log"
	"net/http"
)

func (h *Handler) Authed(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := token.GetSessionCookie(r, sessionCookieName)
		if err != nil {
			utils.RespondWithJson(w, http.StatusUnauthorized, err)
			return
		}

		userId, err := h.token.VerifyToken(t, token.Access)
		if err != nil {
			utils.RespondWithJson(w, http.StatusUnauthorized, err)
			return
		}

		user, err := h.db.GetUserById(r.Context(), userId)
		if err != nil {
			log.Printf("Error getting user by id: %+v", err)
			utils.RespondWithInternalServerError(w)
			return
		}

		handler(w, r, user)
	}
}
