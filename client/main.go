package main

import (
	"context"
	"fmt"
	"invento/oauth/client/common/proto"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// Remove this to un implement implicit grant
	ImplementImplicitGrant()

	// CodeBasedAuthorization()
}

// ImplementImplicitGrant - implementation of oauth package
func ImplementImplicitGrant() {
	conn, err := grpc.Dial("server:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	ClientID := "7cdcd81a-229f-4f68-8fcc-6f39b75a096a"
	ClientSecret := "747306e9-2ee3-455b-9e37-41b59439dee6"

	client := proto.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &proto.AuthImplicitRequest{
		ClientId:     ClientID,
		ClientSecret: ClientSecret,
		Username:     "ali",
		Password:     "Aa#12345678",
	}

	resp, err := client.AuthenticateImplicit(ctx, req)
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}
	logrus.Printf("Access Token: %s\n", resp.AccessToken)
	logrus.Printf("Token Type: %s\n", resp.TokenType)
	logrus.Printf("Expires In: %d\n", resp.ExpiresIn)
	logrus.Printf("Refresh Token: %s\n", resp.RefreshToken)
	logrus.Printf("Scope: %s\n", resp.Scope)

	logrus.Printf("Validating token now\n")
	valResp, err := client.ValidateToken(ctx, &proto.ValidateTokenRequest{
		ClientId:     ClientID,
		ClientSecret: ClientSecret,
		Token:        resp.AccessToken,
	})

	if err != nil {
		log.Fatalf("Failed to parse token: %v", err)
	}

	logrus.Info("Token is Active?: ", valResp.Active)
}

var authCodeChan = make(chan string)

// CodeBasedAuthorization - implementation
func CodeBasedAuthorization() {
	// Start the local HTTP server to handle the RedirectUri
	http.HandleFunc("/oauth2/callback", callbackHandler)
	go http.ListenAndServe(":3000", nil)

	const ClientID = "7cdcd81a-229f-4f68-8fcc-6f39b75a096a"
	const ClientSecret = "747306e9-2ee3-455b-9e37-41b59439dee6"

	// Open the authorization URL in the user's browser
	authorizationURL := "http://0.0.0.0:8099/authorize" +
		"?response_type=code&client_id=" + ClientID + "&redirect_uri=http://0.0.0.0:3000/oauth2/callback&scope=basic"
	logrus.Printf("Open this URL in your browser to authorize the application:\n%s\n", authorizationURL)

	// Wait for the authorization code to be received by the HTTP server
	authCode := <-authCodeChan
	logrus.Info("Auth code found", authCode)
	conn, err := grpc.Dial("server:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logrus.Info("Reached Server with code: ", authCode)
	req := &proto.AuthRequest{
		ClientId:     ClientID,
		ClientSecret: ClientSecret,
		Code:         authCode,
		RedirectUri:  "http://0.0.0.0:3000/oauth2/callback",
		GrantType:    "authorization_code",
	}
	resp, err := client.Authenticate(ctx, req)
	if err != nil {
		logrus.Fatalf("Failed to authenticate: %v", err)
	}

	logrus.Printf("Access Token: %s\n", resp.AccessToken)
	logrus.Printf("Token Type: %s\n", resp.TokenType)
	logrus.Printf("Expires In: %d\n", resp.ExpiresIn)
	logrus.Printf("Refresh Token: %s\n", resp.RefreshToken)
	logrus.Printf("Scope: %s\n", resp.Scope)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Callback handler called")
	// Get the authorization code from the query parameters
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	// Send the authorization code to the main function
	authCodeChan <- code

	// Display a success message to the user and close the HTTP server
	fmt.Fprintf(w, "Authorization successful. You can close this window.")
}
