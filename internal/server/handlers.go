package server

import (
	"github.com/shahin-bayat/go-ssh-client/internal/models"
	"github.com/shahin-bayat/go-ssh-client/internal/utils"
	"log"
	"net/http"
)

func (s *Server) ServeHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}
func (s *Server) ServerAdminPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/admin.html")
}
func (s *Server) ServeUserPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/user.html")
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Registering user")
	// Parse form values
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	confirmPassword := r.PostFormValue("confirmPassword")

	// 1. validate the form values
	err := utils.ValidateUserForm(username, password, confirmPassword)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	// 2. check if the user already exists in db
	err = s.db.UserExists(username)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return

	}
	// 3. hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	// 4. insert the user into the db
	err = s.db.CreateUser(username, hashedPassword, "user")
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	// 5. create ssh user on the server
	err = utils.CreateSSHUser(username, password)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	// 6. return a json success message
	response := models.SuccessResponse{
		Message: "User created successfully",
	}
	utils.WriteJSON(w, http.StatusCreated, response, nil)

}
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {}
