package server

import (
	"github.com/shahin-bayat/go-ssh-client/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", s.healthHandler)

	// Template routes
	r.Get("/", s.ServeHomePage)
	r.With(s.AdminOnly).Get("/admin", s.ServerAdminPage)
	r.With(s.Auth).Get("/user", s.ServeUserPage)

	r.With(s.AdminOnly).Post("/register", s.Register)
	r.Post("/login", s.Login)

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(
		w, http.StatusCreated, s.db.Health(), nil,
	)
}
