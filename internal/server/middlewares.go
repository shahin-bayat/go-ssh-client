package server

import (
	"github.com/shahin-bayat/go-ssh-client/internal/utils"
	"net/http"
	"time"
)

func (s *Server) requireRole(role []string, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("session")
			redirectTo := "/"
			if err != nil {
				http.Redirect(w, r, redirectTo, http.StatusSeeOther)
				return
			}

			sessionToken := c.Value
			session, err := s.db.GetSession(sessionToken)
			if err != nil {
				http.Redirect(w, r, redirectTo, http.StatusSeeOther)
				return
			}

			if session.ExpiresAt.Before(time.Now()) {
				http.Redirect(w, r, redirectTo, http.StatusSeeOther)
				return
			}

			existingUser, err := s.db.GetUserById(session.UserID)
			if err != nil {
				http.Redirect(w, r, redirectTo, http.StatusSeeOther)
				return
			}

			// Check if user has the required role
			if !utils.SliceHas(existingUser.Role, role) {
				http.Redirect(w, r, redirectTo, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

func (s *Server) AdminOnly(next http.Handler) http.Handler {
	return s.requireRole([]string{"admin"}, next)
}

func (s *Server) Auth(next http.Handler) http.Handler {
	return s.requireRole([]string{"admin", "user"}, next)
}
