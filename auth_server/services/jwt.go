package services

import (
	"crypto/rsa"
	"fmt"
	"invento/oauth/auth_server/utils"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateJWTTokens - generates JWT token
func GenerateJWTTokens(privateKey *rsa.PrivateKey) (string, string, error) {
	// Set token expiration times
	accessTokenExp := time.Now().Add(utils.AccessTokenExpCache)
	refreshTokenExp := time.Now().Add(utils.RefreshTokenExpCache)

	// Create the access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		ExpiresAt: accessTokenExp.Unix(),
		Issuer:    "experia",
	})

	// Sign the access token
	accessTokenString, err := accessToken.SignedString(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %v", err)
	}

	// Create the refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		ExpiresAt: refreshTokenExp.Unix(),
		Issuer:    "experia",
	})

	// Sign the refresh token
	refreshTokenString, err := refreshToken.SignedString(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %v", err)
	}

	return accessTokenString, refreshTokenString, nil
}
