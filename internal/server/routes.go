package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", s.healthHandler)

	// Template routes
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	r.Get("/", s.ServeLoginPage)
	r.With(s.AdminOnly).Get("/admin", s.ServerAdminPage)
	r.With(s.Auth).Get("/user", s.ServeUserPage)

	r.With(s.AdminOnly).Post("/register", s.Register)
	r.Post("/login", s.Login)

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := s.db.Health()
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
