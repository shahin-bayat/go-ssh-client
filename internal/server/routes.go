package server

import (
	"fmt"
	"github.com/shahin-bayat/go-ssh-client/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", s.healthHandler)

	r.Get(
		"/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/index.html")
		},
	)

	// Admin dashboard
	r.Get(
		"/admin", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/admin.html")
		},
	)

	// User dashboard
	r.Get(
		"/user", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/user.html")
		},
	)

	r.Post("/register", s.RegisterUser)
	r.Post(
		"/login", func(w http.ResponseWriter, r *http.Request) {
			// Implementation of user login
		},
	)

	// Add a new route for password change
	r.Post(
		"/change-password", func(w http.ResponseWriter, r *http.Request) {
			// Parse form values
			username := r.PostFormValue("username")
			currentPassword := r.PostFormValue("currentPassword")
			newPassword := r.PostFormValue("newPassword")
			confirmNewPassword := r.PostFormValue("confirmNewPassword")

			// Authenticate user (you can use session management or token-based authentication here)

			// Implement logic to change password
			// Check if new password matches the confirmation
			if newPassword != confirmNewPassword {
				http.Error(w, "Passwords do not match", http.StatusBadRequest)
				return
			}

			fmt.Println(currentPassword)

			// Update the database with the new password (replace with your actual database update logic)
			_, err := s.db.UpdateUserPassword(username, newPassword)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Password changed successfully
			fmt.Fprintln(w, "Password changed successfully")
		},
	)

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(
		w, http.StatusCreated, s.db.Health(), nil,
	)
}
