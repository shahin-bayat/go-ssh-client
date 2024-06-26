package database

import (
	"errors"
	"github.com/google/uuid"
	"github.com/shahin-bayat/go-ssh-client/internal/models"
	"time"
)

var (
	ErrorGetUser = errors.New("error getting user")
)

func (s *service) CreateUser(username, password, role string) error {
	_, err := s.db.Exec(`INSERT INTO users (username, password, role) VALUES ($1, $2, $3)`, username, password, role)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreateUserSession(userId uint) (*models.Session, error) {
	expiresAt := time.Now().Add(time.Hour * 24 * 30)
	sessionToken := uuid.NewString()

	session := models.Session{}

	result, err := s.db.Exec(
		`INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)`, userId, sessionToken,
		expiresAt,
	)
	if err != nil {
		return &models.Session{}, err
	}
	sessionId, err := result.LastInsertId()
	if err != nil {
		return &models.Session{}, err
	}
	err = s.db.QueryRow(`SELECT * FROM sessions WHERE id = $1`, sessionId).Scan(
		&session.ID, &session.UserID, &session.Token, &session.ExpiresAt,
	)
	if err != nil {
		return &models.Session{}, err
	}

	return &session, nil
}

func (s *service) GetSession(token string) (*models.Session, error) {
	var session models.Session
	err := s.db.QueryRow(`SELECT * FROM sessions WHERE token = $1`, token).Scan(
		&session.ID, &session.UserID, &session.Token, &session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *service) DeleteSession(token string) error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetSessionByUserId(userId uint) (*models.Session, error) {
	var session models.Session
	err := s.db.QueryRow(`SELECT * FROM sessions WHERE user_id = $1`, userId).Scan(
		&session.ID, &session.UserID, &session.Token, &session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
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

func (s *service) UpdateUserPassword(userId uint, password string) error {
	_, err := s.db.Exec(`UPDATE users SET password = $1 WHERE id = $2`, password, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetUsers() ([]models.User, error) {
	var users []models.User
	rows, err := s.db.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *service) GetUserById(id uint) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(`SELECT * FROM users WHERE id = $1`, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, ErrorGetUser
	}
	return &user, nil
}
