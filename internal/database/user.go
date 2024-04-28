package database

import (
	"database/sql"
	"fmt"
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
