package handlers

import (
	"invento/oauth/auth_server/models"
	"invento/oauth/auth_server/services"
	"invento/oauth/auth_server/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// UserValidationInput - user validation input
type UserValidationInput struct {
	Name     string `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required,passlength,notcontainsname,complexity"`
}

// HandleUserGeneration - Handles User Generation
func HandleUserGeneration(dbSVC *services.DBService, cs *services.ConfigService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			name := r.FormValue("name")
			username := r.FormValue("username")
			password := r.FormValue("password")

			var passMinLength int
			var err error
			if obj, ok := cs.Get("password-min-length"); ok == true {
				passMinLength, err = strconv.Atoi(obj.Value)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

			pc := cs.GetPasswordComplexityConfig()

			// Define validation rules
			validate := validator.New()
			validate.RegisterValidation("notcontainsname", utils.ValidateContainsNotHaveField)
			validate.RegisterValidation("complexity", utils.ValidateComplexity(&pc))
			validate.RegisterValidation("passlength", utils.ValidateMinLength(passMinLength))

			// Create an instance of the input struct
			input := &UserValidationInput{
				Password: password,
				Name:     name,
				Username: username,
			}

			// Validate the UserValidationInput struct
			if err := validate.Struct(input); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			userModel := models.NewUserModel(dbSVC.DB)
			user, err := userModel.GenerateUser(name, username, password)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			utils.WriteJSONResponse(w, user)
		}
	}
}
