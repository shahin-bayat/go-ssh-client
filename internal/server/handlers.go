package server

import (
	"github.com/labstack/echo/v4"
	"github.com/shahin-bayat/go-ssh-client/internal/models"
	"github.com/shahin-bayat/go-ssh-client/internal/utils"
	"net/http"
	"time"
)

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	//	TODO: Implement registration
}
func (s *Server) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.JSON(
			http.StatusBadRequest, models.ErrorResponse{
				Message: "username and password are required",
			},
		)
	}
	existingUser, err := s.db.GetUser(username)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "invalid username or password"})
	}

	if !utils.PasswordMatch(existingUser.Password, password) {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "invalid username or password"})
	}

	session, err := s.createOrGetSession(username, existingUser.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "error creating session"})
	}

	c.SetCookie(
		&http.Cookie{
			Name:     "session",
			Value:    session.Token,
			Expires:  session.ExpiresAt,
			HttpOnly: true,
		},
	)
	//redirectTo := "/user"
	//if existingUser.Role == "admin" {
	//	redirectTo = "/admin"
	//}
	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "login successful"})
	//return c.Redirect(http.StatusFound, "/admin")
}

func (s *Server) Logout(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err != nil {
		return c.Redirect(http.StatusFound, "/")
	}
	c.SetCookie(
		&http.Cookie{
			Name:    "session",
			Value:   "",
			Expires: time.Now(),
		},
	)

	sessionToken := cookie.Value
	_, err = s.db.GetSession(sessionToken)
	if err != nil {
		return c.Redirect(http.StatusFound, "/")
	}
	err = s.db.DeleteSession(sessionToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "error deleting session"})
	}

	return c.Redirect(http.StatusFound, "/")
}

func (s *Server) createOrGetSession(username string, userID uint) (*models.Session, error) {
	session, err := s.db.GetSessionByUserId(userID)
	if err != nil || session.ExpiresAt.Before(time.Now()) {
		// If there's no session or the session has expired, create a new one
		return s.db.CreateUserSession(username, userID)
	}
	return session, nil
}
