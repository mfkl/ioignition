package handlers

import (
	"ioignition/internal/database"
	"ioignition/view"
	"net/http"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request, u database.User) {
	view.Layout(view.Nav(u.Email)).Render(r.Context(), w)
}

func (h *Handler) SignupForm(w http.ResponseWriter, r *http.Request) {
	view.Layout(view.Signup()).Render(r.Context(), w)
}

func (h *Handler) LoginForm(w http.ResponseWriter, r *http.Request) {
	view.Layout(view.Login()).Render(r.Context(), w)
}
