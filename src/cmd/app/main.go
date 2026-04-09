package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xtsank/mypills-super-service/internal/transport/handler"
	"github.com/xtsank/mypills-super-service/internal/transport/middleware"

	"github.com/xtsank/mypills-super-service/internal/infra/postgres"
	"github.com/xtsank/mypills-super-service/internal/service"
)

func main() {
	log.Println("Application is starting...")

	userRepo := postgres.NewPostgresUserRepository("mama")
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	log.Println("Layers initialized")

	router := gin.New()
	log.Println("Gin router created")

	router.Use(middleware.Logger())
	router.Use(middleware.ResponseHandler())
	router.Use(middleware.ErrorHandler())
	log.Println("Middlewares registered")

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
