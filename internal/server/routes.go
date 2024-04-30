package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shahin-bayat/go-ssh-client/web"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	// Register Web Handlers
	web.RegisterHandlers(e)
	// Register Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register API Handlers
	e.GET("/health", s.healthHandler)
	//e.Post("/register", s.Register)
	e.POST("/login", s.Login)
	e.POST("/logout", s.Logout)

	return e
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
