package handlers

import (
	"context"
	"errors"
	"fmt"
	"invento/oauth/server/common/proto"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// TokenURL - Token url
const (
	TokenURL         = "http://auth-server:8099/token"
	IntrospectionURL = "http://auth-server:8099/introspect"
)

// AuthServer - auth server
type AuthServer struct {
	proto.UnimplementedAuthServiceServer
}

// Authenticate - utilize to authenticate
func (s *AuthServer) Authenticate(ctx context.Context, req *proto.AuthRequest) (*proto.AuthResponse, error) {

	logrus.Info("Authenticate - Server called")
	// Initialize your OAuth2 provider configuration
	conf := &oauth2.Config{
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: TokenURL,
		},
		RedirectURL: req.RedirectUri,
	}
	logrus.Info(conf)
	logrus.Info("Token exchange called")

	logrus.Info("Req Code:", req.Code)

	// Exchange the authorization code for an access token
	token, err := conf.Exchange(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	logrus.Info("Token exchange responded")

	// Return the token data in the response
	return &proto.AuthResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    int32(token.Expiry.Sub(time.Now()).Seconds()),
		RefreshToken: token.RefreshToken,
		Scope:        "all",
		// Scope:        strings.Join(token.Extra("scope").([]string), " "),
	}, nil
}

// AuthenticateImplicit - authenticate implicit
func (s *AuthServer) AuthenticateImplicit(ctx context.Context, req *proto.AuthImplicitRequest) (*proto.AuthResponse, error) {

	// logrus.Info("reached authenticateImplicit")
	// Initialize your OAuth2 provider configuration
	conf := &oauth2.Config{
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: TokenURL,
		},
	}

	// logrus.Info("Username Req:", req.GetUsername())
	// logrus.Info("Password Req:", req.GetPassword())
	token, err := conf.PasswordCredentialsToken(oauth2.NoContext, req.GetUsername(), req.GetPassword())
	if err != nil {
		// logrus.Info("Issue while generating token")
		return nil, err
	}
	// Return the token data in the response
	return &proto.AuthResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    int32(token.Expiry.Sub(time.Now()).Seconds()),
		RefreshToken: token.RefreshToken,
		Scope:        "all",
		// Scope:        strings.Join(token.Extra("scope").([]string), " "),
	}, nil
}

// ValidateToken - validate token
func (s *AuthServer) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	logrus.Info("Into Validate Token")
	// Replace with the URL of your OAuth2 provider's introspection endpoint

	// Replace with your OAuth2 client ID and secret
	clientID := req.GetClientId()
	clientSecret := req.GetClientSecret()

	// Prepare the request to the introspection endpoint
	data := url.Values{}
	data.Set("token", req.GetToken())

	newReq, err := http.NewRequest("POST", IntrospectionURL, strings.NewReader(data.Encode()))
	if err != nil {
		return &proto.ValidateTokenResponse{
			Active: false,
		}, err
	}

	newReq.SetBasicAuth(clientID, clientSecret)
	newReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(newReq)
	if err != nil {
		return &proto.ValidateTokenResponse{
			Active: false,
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		logrus.Info("Found OK the token is valid")
		return &proto.ValidateTokenResponse{
			Active: true,
		}, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return &proto.ValidateTokenResponse{
			Active: false,
		}, errors.New("Error reading response body: " + err.Error())
	}

	return &proto.ValidateTokenResponse{
		Active: false,
	}, errors.New(string(body))

}
