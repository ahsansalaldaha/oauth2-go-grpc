package handlers

import (
	"encoding/json"
	"invento/oauth/auth_server/services"
	"invento/oauth/auth_server/utils"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// TokenHandler - Token handler
func TokenHandler(rs *services.RedisService, credSVC *services.CredentialService) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		// Validate the client credentials and authorization code
		clientID := r.FormValue("client_id")
		clientSecret := r.FormValue("client_secret")

		grantType := r.FormValue("grant_type")

		// Define validation rules
		validate := validator.New()
		err := validate.Var(clientID, "required")
		if err != nil {
			http.Error(w, "'clientID' field is required", http.StatusBadRequest)
			return
		}

		err = validate.Var(clientSecret, "required")
		if err != nil {
			http.Error(w, "'clientSecret' field is required", http.StatusBadRequest)
			return
		}

		logrus.Info("Grant Type: ", grantType)
		if grantType == "authorization_code" {
			code := r.FormValue("code")
			logrus.Info("Code: ", code)
			err = validate.Var(code, "required")
			if err != nil {
				http.Error(w, "'code' field is required", http.StatusBadRequest)
				return
			}
			cs := services.NewCodeService(rs)
			logrus.Info("credSVC.ValidateClientAndSecret", credSVC.ValidateClientAndSecret(clientID, clientSecret))
			logrus.Info("cs.Verify(code)", cs.Verify(code))
			// You'll need to implement your own validation logic here
			if !credSVC.ValidateClientAndSecret(clientID, clientSecret) || !cs.Verify(code) {
				logrus.Info("Invalid request: authorization_code")
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
		} else if grantType == "password" {
			if !credSVC.ValidateClientAndSecret(clientID, clientSecret) {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
		}

		// Generate an access token and refresh token
		accessToken, refreshToken, err := services.GenerateJWTTokens(services.GetPrivateKey())
		if err != nil {
			http.Error(w, "Token generation failed", http.StatusInternalServerError)
			return
		}

		// Prepare the token response
		resp := &oauth2.Token{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			Expiry:       time.Now().Add(utils.JWTAuthTokenExpiryTime),
			RefreshToken: refreshToken,
		}

		// Set the response headers and body
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
