package server

import (
	"errors"
	"github.com/shahin-bayat/go-ssh-client/internal/utils"
	"net/http"
)

var (
	ErrorUnauthorized = errors.New("unauthorized")
)

func (s *Server) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Extract the username and password from the request
			username, password, ok := r.BasicAuth()
			if ok {
				existingUser, err := s.db.GetUser(username)
				if err != nil {
					utils.WriteErrorJSON(w, http.StatusUnauthorized, err, nil)
					return
				}

				usernameHash := utils.Hash(username)
				existingUsernameHash := utils.Hash(existingUser.Username)
				existingPasswordHash := existingUser.Password

				usernameMatch := utils.Match(usernameHash, existingUsernameHash)
				passwordMatch := utils.PasswordMatch(existingPasswordHash, password)

				if usernameMatch && passwordMatch && existingUser.Role == "admin" {
					next.ServeHTTP(w, r)
					return
				}
			}

			headers := map[string]string{
				"WWW-Authenticate": `Basic realm="restricted", charset="UTF-8"`,
			}
			utils.WriteErrorJSON(w, http.StatusUnauthorized, ErrorUnauthorized, headers)
		},
	)
}
