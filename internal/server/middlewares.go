package server

import (
	"github.com/labstack/echo/v4"
	"github.com/shahin-bayat/go-ssh-client/internal/utils"
	"net/http"
	"time"
)

func (s *Server) requireRole(role []string, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session")
		redirectTo := "/"
		if err != nil {
			return c.Redirect(http.StatusSeeOther, redirectTo)
		}

		sessionToken := cookie.Value
		session, err := s.db.GetSession(sessionToken)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, redirectTo)
		}

		if session.ExpiresAt.Before(time.Now()) {
			return c.Redirect(http.StatusSeeOther, redirectTo)
		}

		existingUser, err := s.db.GetUserById(session.UserID)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, redirectTo)
		}

		// Check if user has the required role
		if !utils.SliceHas(existingUser.Role, role) {
			c.Redirect(http.StatusSeeOther, redirectTo)
		}

		return next(c)
	}
}

func (s *Server) AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return s.requireRole([]string{"admin"}, next)
}

func (s *Server) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return s.requireRole([]string{"admin", "user"}, next)
}
