package handlers

import (
	"invento/oauth/auth_server/models"
	"invento/oauth/auth_server/services"
	"invento/oauth/auth_server/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// HandleClientGeneration - Handles Client Generation
func HandleClientGeneration(dbSVC *services.DBService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			logrus.Info("HandleClientGeneration Called")

			typeValue := r.FormValue("grant_type")
			redirects := r.Form["redirects"]

			// Define validation rules
			validate := validator.New()
			err := validate.Var(typeValue, "required")
			if err != nil {
				http.Error(w, "'type' field is required", http.StatusBadRequest)
				return
			}

			err = validate.Var(redirects, "required")
			if err != nil {
				http.Error(w, "'redirects' field is required", http.StatusBadRequest)
				return
			}

			var selectedGrantType models.GrantType
			switch typeValue {
			case "code":
				selectedGrantType = models.Code
			case "implicit":
				selectedGrantType = models.Implicit
			default:
				selectedGrantType = models.Code
			}

			logrus.Info("SelectedGrantType:", selectedGrantType)
			clientModel := models.NewClientModel(dbSVC.DB)
			client, err := clientModel.GenerateClient(selectedGrantType, redirects)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			utils.WriteJSONResponse(w, client)
		}
	}
}
