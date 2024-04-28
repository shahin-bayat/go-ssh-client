package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/shahin-bayat/go-ssh-client/internal/models"
)

var (
	ErrorGetUser = errors.New("error getting user")
)

func (s *service) UserExists(username string) error {
	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, username).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("user already exists")
	}

	return nil
}

func (s *service) CreateUser(username, password, role string) error {
	_, err := s.db.Exec(`INSERT INTO users (username, password, role) VALUES ($1, $2, $3)`, username, password, role)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateUserPassword(username, newPassword string) (sql.Result, error) {
	result, err := s.db.Exec(`UPDATE users SET password = $1 WHERE username = $2`, newPassword, username)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) GetUser(username string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(`SELECT * FROM users WHERE username = $1`, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, ErrorGetUser
	}
	return &user, nil
}
