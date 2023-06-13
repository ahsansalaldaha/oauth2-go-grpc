package handlers

import (
	"fmt"
	"invento/oauth/auth_server/services"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

// TokenIntrospectHandler - Token validation handler
func TokenIntrospectHandler(rs *services.RedisService) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("received Introspect Handler")
		if r.Method == http.MethodPost {
			// Retrieve a single parameter using FormValue
			tokenString := r.FormValue("token")
			logrus.Println("token:", tokenString)
			// Check if the token is cached in Redis
			_, err := rs.Get(tokenString).Result()

			if err == redis.Nil {
				// Parse and validate the token
				token, jwtErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return services.GetPublicKey(), nil
				})
				// Check if the token is valid
				if jwtErr == nil && token.Valid {
					// Cache the valid token in Redis with an expiration time matching the token's expiration
					claims := token.Claims.(jwt.MapClaims)
					exp := claims["exp"].(float64)
					expTime := time.Unix(int64(exp), 0)
					duration := expTime.Sub(time.Now())

					err := rs.Set(tokenString, "valid", duration).Err()
					if err != nil {
						log.Printf("Failed to cache token in Redis: %v", err)

						return
					}

					w.WriteHeader(http.StatusOK) // 200 OK
					return
				}

				http.Error(w, jwtErr.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				// Redis error, log it and return false
				log.Printf("Redis error: %v", err)
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			// Token found in cache, return true
			w.WriteHeader(http.StatusOK) // 200 OK
			return
		}
	}

}
