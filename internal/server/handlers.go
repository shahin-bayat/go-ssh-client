package server

import (
	"fmt"
	"github.com/shahin-bayat/go-ssh-client/views/layouts"
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
func (s *Server) ServeAdminPage(c echo.Context) error {
	component := pages.Admin()
	return layouts.Admin("Dashboard", component).Render(c.Request().Context(), c.Response().Writer)
}
func (s *Server) ServeUsersPage(c echo.Context) error {
	users, err := s.db.GetUsers()
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", utils.ErrorResponse{Error: "failed to get users"})
	}
	component := pages.Users(users)
	return layouts.Admin("Users", component).Render(c.Request().Context(), c.Response().Writer)

}
func (s *Server) ServeSettingsPage(c echo.Context) error {
	component := pages.Settings()
	return layouts.Admin("Settings", component).Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) ChangePassword(c echo.Context) error {
	username := c.FormValue("username")
	currentPassword := c.FormValue("current-password")
	newPassword := c.FormValue("password")
	confirmPassword := c.FormValue("confirm-password")

	if username == "" || currentPassword == "" || newPassword == "" || confirmPassword == "" {
		component := components.Alert("all fields are required", "")
		c.Response().WriteHeader(http.StatusBadRequest)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	if newPassword != confirmPassword {
		component := components.Alert("passwords do not match", "")
		c.Response().WriteHeader(http.StatusBadRequest)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	if len(newPassword) < 6 {
		component := components.Alert("password must be at least 6 characters long", "")
		c.Response().WriteHeader(http.StatusBadRequest)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	existingUser, err := s.db.GetUser(username)
	if err != nil {
		component := components.Alert("invalid credentials", "")
		c.Response().WriteHeader(http.StatusUnauthorized)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	if !utils.PasswordMatch(existingUser.Password, currentPassword) {
		component := components.Alert("invalid credentials", "")
		c.Response().WriteHeader(http.StatusUnauthorized)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		component := components.Alert("internal server error", "")
		c.Response().WriteHeader(http.StatusInternalServerError)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	err = s.db.UpdateUserPassword(existingUser.ID, hashedPassword)
	if err != nil {
		component := components.Alert("internal server error", "")
		c.Response().WriteHeader(http.StatusInternalServerError)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	// if user is not an admin, update the SSH user password
	if existingUser.Role != "admin" {
		// TODO: Implement password change for SSH user
	}
	component := components.Alert("", "password updated successfully")
	c.Response().WriteHeader(http.StatusOK)
	return component.Render(c.Request().Context(), c.Response().Writer)
}
func (s *Server) CreateUser(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")

	fmt.Println("Are we here?")
	fmt.Println(username, password, role)

	if username == "" || password == "" || role == "" {
		component := components.Alert("all fields are required", "")
		c.Response().WriteHeader(http.StatusBadRequest)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	if len(password) < 6 {
		component := components.Alert("password must be at least 6 characters long", "")
		c.Response().WriteHeader(http.StatusBadRequest)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	if role != "admin" && role != "user" {
		component := components.Alert("invalid role", "")
		c.Response().WriteHeader(http.StatusBadRequest)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		component := components.Alert("internal server error", "")
		c.Response().WriteHeader(http.StatusInternalServerError)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
	err = s.db.CreateUser(username, hashedPassword, role)
	if err != nil {
		component := components.Alert("internal server error", "")
		c.Response().WriteHeader(http.StatusInternalServerError)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	if role == "user" {
		err = utils.CreateSSHUser(username, password)
		if err != nil {
			component := components.Alert("internal server error", "")
			c.Response().WriteHeader(http.StatusInternalServerError)
			return component.Render(c.Request().Context(), c.Response().Writer)
		}
	}

	component := components.Alert("", "user created successfully")
	c.Response().Header().Set("HX-Refresh", "true")
	c.Response().WriteHeader(http.StatusOK)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		component := components.Alert("all fields are required", "")
		c.Response().WriteHeader(http.StatusBadRequest)
		return component.Render(c.Request().Context(), c.Response().Writer)

	}
	existingUser, err := s.db.GetUser(username)
	if err != nil {
		component := components.Alert("invalid credentials", "")
		c.Response().WriteHeader(http.StatusUnauthorized)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	if !utils.PasswordMatch(existingUser.Password, password) {
		component := components.Alert("invalid credentials", "")
		c.Response().WriteHeader(http.StatusUnauthorized)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	session, err := s.createOrGetSession(existingUser.ID)
	if err != nil {
		component := components.Alert("internal server error", "")
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
