package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/static", "views/static")

	// Health check
	e.GET("/health", s.healthHandler)

	// Template routes
	e.GET("/", s.ServeLoginPage)
	e.GET("/admin/dashboard", s.ServerAdminPage, s.AdminOnly)
	e.GET("/admin/users", s.ServeAdminUsersPage, s.AdminOnly)

	// API routes
	e.POST("/change-password", s.ChangePassword, s.Auth)
	e.POST("/login", s.Login)
	e.POST("/logout", s.Logout)

	return e
}

func (s *Server) healthHandler(c echo.Context) error {
	response := s.db.Health()
	return c.JSON(http.StatusOK, response)
}
