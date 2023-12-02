package handlers

import (
	"fmt"
	"net/http"
	"time"
)

const (
	sessionCookieName = "session"
	cookieExpiry      = 30 * 24 * 60 * time.Minute
)

func (h *Handler) SetSesionCookie(w http.ResponseWriter, r *http.Request, t string) {
	fmt.Println(cookieExpiry)
	cookie := &http.Cookie{
		Name:  sessionCookieName,
		Value: t,
		Path:  "/",

		MaxAge:   int(cookieExpiry.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}
