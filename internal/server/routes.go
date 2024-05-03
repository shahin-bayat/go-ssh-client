package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/static", "web/static")

	// Template engine
	e.Renderer = NewTemplate()

	// Health check
	e.GET("/health", s.healthHandler)

	// Template routes
	e.GET("/", s.ServeLoginPage)
	e.GET("/admin/dashboard", s.ServerAdminPage, s.AdminOnly)
	e.GET("/admin/users", s.ServeAdminUsersPage, s.AdminOnly)
	e.GET("/user", s.ServeUserPage, s.Auth)

	// API routes
	e.POST("/register", s.Register, s.AdminOnly)
	e.GET("/users", s.GetUsers, s.AdminOnly)
	e.POST("/change-password", s.ChangePassword, s.Auth)

	e.POST("/login", s.Login)
	e.POST("/logout", s.Logout)

	return e
}

func (s *Server) healthHandler(c echo.Context) error {
	response := s.db.Health()
	return c.JSON(http.StatusOK, response)
}
