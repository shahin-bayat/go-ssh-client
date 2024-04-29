package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/shahin-bayat/go-ssh-client/internal/database"
)

type Server struct {
	port      int
	db        database.Service
	loginTmpl *template.Template
	adminTmpl *template.Template
	userTmpl  *template.Template
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:      port,
		db:        database.New(),
		loginTmpl: template.Must(template.ParseFiles("web/pages/login.html")),
		adminTmpl: template.Must(template.ParseFiles("web/pages/admin.html")),
		userTmpl:  template.Must(template.ParseFiles("web/pages/user.html")),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
