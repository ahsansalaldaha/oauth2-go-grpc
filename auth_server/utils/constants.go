package utils

import "time"

const (
	// ConfigRememberTime - represents config remember time in cache
	ConfigRememberTime = 1 * time.Hour
	// CodeExpiryTime - represents the time for how long newly generated code is stored in redis for verification
	CodeExpiryTime = 24 * time.Hour
	// JWTAuthTokenExpiryTime - represents the time for how long expiry of JWT token is
	JWTAuthTokenExpiryTime = 1 * time.Hour
	// CredsClientCache - represents client creds cache time
	CredsClientCache = 1 * time.Hour
	// AccessTokenExpCache - represents access token expiry cache time
	AccessTokenExpCache = 15 * time.Minute
	// RefreshTokenExpCache - represents refresh token expiry cache time
	RefreshTokenExpCache = 24 * time.Minute
)
