package utils

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os/exec"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Response struct {
	Data interface{} `json:"data"`
}

func ValidateUserForm(username, password, confirmPassword string) error {
	if username == "" || password == "" || confirmPassword == "" {
		return fmt.Errorf("username, password, and confirmPassword are required fields")
	}

	// TODO: sanitize the input to make sure it does maliscious code to run on the server

	if password != confirmPassword {
		return fmt.Errorf("password and confirmPassword do not match")
	}

	return nil
}

func Hash(val string) [32]byte {
	return sha256.Sum256([]byte(val))
}
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func PasswordMatch(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func Match(val1, val2 [32]byte) bool {
	return subtle.ConstantTimeCompare(val1[:], val2[:]) == 1
}

func CreateSSHUser(username, password string) error {
	cmd := exec.Command("useradd", "-s", "/usr/sbin/nologin", username)

	// Run the command and capture its output
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to create SSH user: %w", err)
	}

	// Set the password for the user
	cmd = exec.Command("sh", "-c", fmt.Sprintf("echo '%s:%s' | chpasswd", username, password))
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to set password for SSH user: %w", err)
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}, headers map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	for key, value := range headers {
		_, ok := headers[key]
		if !ok {
			continue
		}
		w.Header().Set(key, value)
	}
	w.WriteHeader(status)
	if status == http.StatusNoContent {
		return
	}
	response := Response{Data: v}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal(err)
	}
}

func WriteErrorJSON(w http.ResponseWriter, status int, err error, headers map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	for key, value := range headers {
		_, ok := headers[key]
		if !ok {
			continue
		}
		w.Header().Set(key, value)
	}
	w.WriteHeader(status)
	errorResponse := ErrorResponse{Error: err.Error()}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Fatal(err)
	}
}
