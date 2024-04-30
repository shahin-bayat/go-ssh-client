package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/shahin-bayat/go-ssh-client/internal/models"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type Service interface {
	Health() models.SuccessResponse
	CreateUser(username, password, role string) error
	GetUser(username string) (*models.User, error)
	GetUsers() ([]models.User, error)
	GetUserById(id uint) (*models.User, error)
	CreateAdminUser(username, password string) error
	CreateUserSession(userId uint) (*models.Session, error)
	GetSession(token string) (*models.Session, error)
	GetSessionByUserId(userId uint) (*models.Session, error)
	DeleteSession(token string) error
}

type service struct {
	db *sql.DB
}

var (
	dburl      = os.Getenv("DB_URL")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) Health() models.SuccessResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}
	response := models.SuccessResponse{
		Message: "It's healthy",
	}
	return response
}
