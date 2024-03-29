package main

import (
	"context"
	"invento/oauth/auth_server/handlers"
	"invento/oauth/auth_server/models"
	"invento/oauth/auth_server/services"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {

	if err := services.GenerateRSAKeyPairIfNotExists("common/keys"); err != nil {
		panic(err)
	}

	redisService := services.NewRedisService(context.Background())
	dbSVC := services.NewDBService()
	queueSVC := services.NewQueueService()

	cs := services.NewConfigService(dbSVC.DB, redisService)
	cs.StoreDefaultConfigs()

	credSVC := services.NewCredentialService(dbSVC, redisService, queueSVC, cs)

	dbSVC.DB.AutoMigrate(&models.User{}, &models.Client{}, &models.Redirect{}, &models.Config{}, &models.UserLock{})

	http.HandleFunc("/client", handlers.HandleClientGeneration(dbSVC))
	http.HandleFunc("/user", handlers.HandleUserGeneration(dbSVC, cs))
	http.HandleFunc("/token", handlers.TokenHandler(redisService, credSVC))
	http.HandleFunc("/authorize", handlers.AuthorizeHandler(redisService, credSVC))
	http.HandleFunc("/introspect", handlers.TokenIntrospectHandler(redisService))

	logrus.Info("Successfully started 🚀")
	log.Fatal(http.ListenAndServe(":8099", nil))

}
