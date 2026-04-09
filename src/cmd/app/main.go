package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xtsank/mypills-super-service/internal/transport/handler"

	"github.com/xtsank/mypills-super-service/internal/infra/postgres"
	"github.com/xtsank/mypills-super-service/internal/service"
)

func main() {
	log.Println("Application is starting...")

	userRepo := postgres.NewPostgresUserRepository("mama")
	log.Println("Repository layer initialized")

	authService := service.NewAuthService(userRepo)
	log.Println("Service layer initialized")

	authHandler := handler.NewAuthHandler(authService)
	log.Println("Transport layer initialized")

	router := gin.Default()
	log.Println("Gin router created")

	apiV1Group := router.Group("/api/v1")
	{
		authGroup := apiV1Group.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
		}
	}
	log.Println("Routes registered")

	serverAddress := ":8080"
	log.Printf("Starting server on %s", serverAddress)

	if err := router.Run(serverAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
