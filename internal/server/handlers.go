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

func (s *Server) ChangePassword(c echo.Context) error {
	username := c.FormValue("username")
	currentPassword := c.FormValue("current-password")
	newPassword := c.FormValue("password")
	confirmPassword := c.FormValue("confirm-password")

	if username == "" || currentPassword == "" || newPassword == "" || confirmPassword == "" {
		component := components.Error("all fields are required")
		c.Response().WriteHeader(http.StatusBadRequest)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	if newPassword != confirmPassword {
		component := components.Error("passwords do not match")
		c.Response().WriteHeader(http.StatusBadRequest)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	if len(newPassword) < 6 {
		component := components.Error("password must be at least 6 characters long")
		c.Response().WriteHeader(http.StatusBadRequest)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	existingUser, err := s.db.GetUser(username)
	if err != nil {
		component := components.Error("invalid username")
		c.Response().WriteHeader(http.StatusUnauthorized)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	if !utils.PasswordMatch(existingUser.Password, currentPassword) {
		component := components.Error("invalid password")
		c.Response().WriteHeader(http.StatusUnauthorized)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		component := components.Error("internal server error")
		c.Response().WriteHeader(http.StatusInternalServerError)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	err = s.db.UpdateUserPassword(existingUser.ID, hashedPassword)
	if err != nil {
		component := components.Error("internal server error")
		c.Response().WriteHeader(http.StatusInternalServerError)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	// if user is not an admin, update the SSH user password
	if existingUser.Role != "admin" {
		// TODO: Implement password change for SSH user
	}
	component := components.Success("password updated successfully")
	c.Response().WriteHeader(http.StatusOK)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		component := components.Error("all fields are required")
		c.Response().WriteHeader(http.StatusBadRequest)
		return component.Render(c.Request().Context(), c.Response().Writer)

	}
	existingUser, err := s.db.GetUser(username)
	if err != nil {
		component := components.Error("invalid credentials")
		c.Response().WriteHeader(http.StatusUnauthorized)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	if !utils.PasswordMatch(existingUser.Password, password) {
		component := components.Error("invalid credentials")
		c.Response().WriteHeader(http.StatusUnauthorized)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	session, err := s.createOrGetSession(existingUser.ID)
	if err != nil {
		component := components.Error("internal server error")
		c.Response().WriteHeader(http.StatusInternalServerError)
		return component.Render(c.Request().Context(), c.Response().Writer)
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

func (s *Server) createOrGetSession(userID uint) (*models.Session, error) {
	session, err := s.db.GetSessionByUserId(userID)
	if err != nil || session.ExpiresAt.Before(time.Now()) {
		// If there's no session or the session has expired, create a new one
		return s.db.CreateUserSession(userID)
	}
	return session, nil
}
