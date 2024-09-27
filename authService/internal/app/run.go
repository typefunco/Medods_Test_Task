package app

import (
	"log"
	"medods_auth/authService/internal/config"
	"medods_auth/authService/internal/handlers"
	"medods_auth/authService/internal/utils"

	"github.com/gin-gonic/gin"
)

func RunServer(config *config.Config) {
	jwtService := utils.NewJWTService(config.JWTSecret)
	postgresRepo, err := utils.NewPostgresRepository(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL repository: %v", err)
	}

	authHandler := &handlers.AuthHandler{
		UserRepo:  postgresRepo,
		TokenRepo: postgresRepo,
		JWT:       *jwtService,
	}

	router := gin.Default()
	handlers.RegisterRoutes(router, authHandler)

	log.Printf("Starting server on %s", config.ServerAddress)
	if err := router.Run(config.ServerAddress); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
