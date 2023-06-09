package handlers

import (
	"invento/oauth/auth_server/models"
	"invento/oauth/auth_server/services"
	"invento/oauth/auth_server/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// HandleUserGeneration - Handles User Generation
func HandleUserGeneration(dbSVC *services.DBService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			username := r.FormValue("username")
			password := r.FormValue("password")

			// Define validation rules
			validate := validator.New()
			err := validate.Var(username, "required")
			if err != nil {
				http.Error(w, "'username' field is required", http.StatusBadRequest)
				return
			}

			err = validate.Var(password, "required")
			if err != nil {
				http.Error(w, "'password' field is required", http.StatusBadRequest)
				return
			}

			userModel := models.NewUserModel(dbSVC.DB)
			user, err := userModel.GenerateUser(username, password)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			utils.WriteJSONResponse(w, user)
		}
	}
}
