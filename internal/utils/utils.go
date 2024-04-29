package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
)

type ErrorResponse struct {
	Error string `json:"error"`
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

func RenderError(w http.ResponseWriter, template *template.Template, block, message string) {
	err := template.ExecuteTemplate(w, block, ErrorResponse{Error: message})
	if err != nil {
		log.Fatal(err)
	}
}

func SliceHas(val string, slice []string) bool {
	for _, r := range slice {
		if r == val {
			return true
		}
	}
	return false
}
