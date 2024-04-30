package server

import (
	"github.com/shahin-bayat/go-ssh-client/internal/models"
	"github.com/shahin-bayat/go-ssh-client/internal/utils"
	"net/http"
	"time"
)

func (s *Server) ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	loginTmpl := s.tmpl.Lookup("login.html")
	err := loginTmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func (s *Server) ServerAdminPage(w http.ResponseWriter, r *http.Request) {
	adminTmpl := s.tmpl.Lookup("admin.html")
	err := adminTmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *Server) ServeUserPage(w http.ResponseWriter, r *http.Request) {
	userTmpl := s.tmpl.Lookup("user.html")
	err := userTmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	//	TODO: Implement registration
}
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	loginTmpl := s.tmpl.Lookup("login.html")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if username == "" || password == "" {
		utils.RenderError(w, loginTmpl, "validation-error", "username and password are required fields")
		return
	}
	existingUser, err := s.db.GetUser(username)
	if err != nil {
		utils.RenderError(w, loginTmpl, "validation-error", "invalid username or password")
		return
	}

	if !utils.PasswordMatch(existingUser.Password, password) {
		utils.RenderError(w, loginTmpl, "validation-error", "invalid password")
		return
	}

	session, err := s.createOrGetSession(username, existingUser.ID)
	if err != nil {
		utils.RenderError(w, loginTmpl, "validation-error", "failed to create session")
		return
	}

	http.SetCookie(
		w, &http.Cookie{
			Name:     "session",
			Value:    session.Token,
			Expires:  session.ExpiresAt,
			HttpOnly: true,
		},
	)
	redirectTo := "/user"
	if existingUser.Role == "admin" {
		redirectTo = "/admin"
	}
	w.Header().Set("HX-Redirect", redirectTo)
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		w.Header().Set("HX-Redirect", "/")
		return
	}
	http.SetCookie(
		w, &http.Cookie{
			Name:    "session",
			Value:   "",
			Expires: time.Now(),
		},
	)

	sessionToken := c.Value
	_, err = s.db.GetSession(sessionToken)
	if err != nil {
		w.Header().Set("HX-Redirect", "/")
		return
	}
	err = s.db.DeleteSession(sessionToken)
	if err != nil {
		w.Header().Set("HX-Redirect", "/")
		return
	}

	w.Header().Set("HX-Redirect", "/")
}

func (s *Server) createOrGetSession(username string, userID uint) (*models.Session, error) {
	session, err := s.db.GetSessionByUserId(userID)
	if err != nil || session.ExpiresAt.Before(time.Now()) {
		// If there's no session or the session has expired, create a new one
		return s.db.CreateUserSession(username, userID)
	}
	return session, nil
}
