package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shahin-bayat/go-ssh-client/internal/models"
	"github.com/shahin-bayat/go-ssh-client/internal/utils"
	"github.com/shahin-bayat/go-ssh-client/views/components"
	"github.com/shahin-bayat/go-ssh-client/views/pages"
)

func (s *Server) ServeLoginPage(c echo.Context) error {
	component := pages.Login()
	return component.Render(c.Request().Context(), c.Response().Writer)

}
func (s *Server) ServerAdminPage(c echo.Context) error {
	component := pages.Admin()
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) ServeAdminUsersPage(c echo.Context) error {
	users, err := s.db.GetUsers()
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", utils.ErrorResponse{Error: "failed to get users"})
	}
	component := pages.Users(users)
	return component.Render(c.Request().Context(), c.Response().Writer)

}
func (s *Server) ServeUserPage(c echo.Context) error {
	return c.Render(http.StatusOK, "user-dashboard", nil)
}

func (s *Server) ChangePassword(c echo.Context) error {
	username := c.FormValue("username")
	currentPassword := c.FormValue("current-password")
	newPassword := c.FormValue("password")
	confirmPassword := c.FormValue("confirm-password")

	component := components.ChangePassword()

	if username == "" || currentPassword == "" || newPassword == "" || confirmPassword == "" {
		return c.Render(
			http.StatusBadRequest, "change-password-error", utils.ErrorResponse{Error: "all fields are required"},
		)
	}
	if newPassword != confirmPassword {
		return c.Render(
			http.StatusBadRequest, "change-password-error", utils.ErrorResponse{Error: "passwords do not match"},
		)
	}
	if len(newPassword) < 6 {
		return c.Render(
			http.StatusBadRequest, "change-password-error",
			utils.ErrorResponse{Error: "password must be at least 6 characters"},
		)
	}
	existingUser, err := s.db.GetUser(username)
	if err != nil {
		return c.Render(
			http.StatusUnauthorized, "change-password-error", utils.ErrorResponse{Error: "invalid username"},
		)
	}
	if !utils.PasswordMatch(existingUser.Password, currentPassword) {
		return c.Render(
			http.StatusUnauthorized, "change-password-error", utils.ErrorResponse{Error: "invalid password"},
		)
	}

	err = s.db.UpdateUserPassword(existingUser.ID, newPassword)
	if err != nil {
		return c.Render(
			http.StatusInternalServerError, "change-password-error",
			utils.ErrorResponse{Error: "failed to update password"},
		)
	}

	// if user is not an admin, update the SSH user password
	if existingUser.Role != "admin" {
		// TODO: Implement password change for SSH user
	}
	return err
}

func (s *Server) Register(c echo.Context) error {
	//	TODO: Implement registration
	return nil
}
func (s *Server) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Render(
			http.StatusBadRequest, "error", utils.ErrorResponse{Error: "username and password are required"},
		)
	}
	existingUser, err := s.db.GetUser(username)
	if err != nil {
		return c.Render(http.StatusUnauthorized, "error", utils.ErrorResponse{Error: "invalid username"})
	}

	if !utils.PasswordMatch(existingUser.Password, password) {
		return c.Render(http.StatusUnauthorized, "error", utils.ErrorResponse{Error: "invalid password"})
	}

	session, err := s.createOrGetSession(existingUser.ID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", utils.ErrorResponse{Error: "failed to create session"})
	}

	c.SetCookie(
		&http.Cookie{
			Name:     "session",
			Value:    session.Token,
			Expires:  session.ExpiresAt,
			HttpOnly: true,
		},
	)
	redirectTo := "/user/dashboard"
	if existingUser.Role == "admin" {
		redirectTo = "/admin/dashboard"
	}
	c.Response().Header().Set("HX-Redirect", redirectTo)
	return nil
}
func (s *Server) Logout(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return err
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
		c.Response().Header().Set("HX-Redirect", "/")
		return err
	}
	err = s.db.DeleteSession(sessionToken)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return nil
}
func (s *Server) GetUsers(c echo.Context) error {
	users, err := s.db.GetUsers()
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", utils.ErrorResponse{Error: "failed to get users"})
	}
	return c.Render(http.StatusOK, "users", users)
}

func (s *Server) createOrGetSession(userID uint) (*models.Session, error) {
	session, err := s.db.GetSessionByUserId(userID)
	if err != nil || session.ExpiresAt.Before(time.Now()) {
		// If there's no session or the session has expired, create a new one
		return s.db.CreateUserSession(userID)
	}
	return session, nil
}
