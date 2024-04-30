package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/shahin-bayat/go-ssh-client/internal/database"
)

type Server struct {
	port int
	db   database.Service
	tmpl *template.Template
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	tmpl, err := template.New("").ParseGlob("web/**/*.html")
	if err != nil {
		log.Fatal(err)
	}

	NewServer := &Server{
		port: port,
		db:   database.New(),
		tmpl: tmpl,
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
