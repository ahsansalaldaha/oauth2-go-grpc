package handlers

import (
	"fmt"
	"invento/oauth/auth_server/services"
	"net/http"
)

// AuthorizeHandler - authorizes the request
func AuthorizeHandler(rs *services.RedisService, credSVC *services.CredentialService) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		// Validate the client_id and redirect_uri
		clientID := r.URL.Query().Get("client_id")
		redirectURI := r.URL.Query().Get("redirect_uri")

		if !credSVC.ValidateClientID(clientID) || !credSVC.ValidateRedirectURI(clientID, redirectURI) {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// You'll need to implement your own user authentication and consent logic here.
		// In this example, we'll skip the authentication and consent process and
		// assume the user has granted permission.
		// Check if the user has submitted the login form
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			// Validate the username and password
			if credSVC.ValidateUserCredentials(username, password) {
				// Generate an authorization code
				cs := services.NewCodeService(rs)
				code := cs.Get()
				// Redirect the user back to the client application with the authorization code
				http.Redirect(w, r, fmt.Sprintf("%s?code=%s", redirectURI, code), http.StatusFound)
			}
		}

		// Display the login form
		loginForm := `
		<html>
			<head>
				<title>Experia Login</title>
			</head>
			<body>
				<h1>Login</h1>
				<form method="post">
					<label>Username:</label>
					<input type="text" name="username" required>
					<br>
					<label>Password:</label>
					<input type="password" name="password" required>
					<br>
					<button type="submit">Authorize</button>
				</form>
			</body>
		</html>
	`
		w.Write([]byte(loginForm))
	}
}
