package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/config"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/db"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/repository"
	"github.com/xtsank/mypills-super-service/src/internal/transport/handler"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"

	_ "github.com/xtsank/mypills-super-service/docs/swagger"
	"github.com/xtsank/mypills-super-service/src/internal/service"
)

func main() {
	i := do.New()

	// DB
	do.Provide(i, config.NewConfig)
	do.Provide(i, db.NewDB)

	// Repo
	do.Provide(i, func(i do.Injector) (user.IUserRepository, error) {
		return repository.NewPostgresUserRepository(i)
	})
	do.Provide(i, func(i do.Injector) (medicine.IMedicineRepository, error) {
		return repository.NewPostgresMedicineRepository(i)
	})
	do.Provide(i, func(i do.Injector) (cabinet_item.ICabinetItemRepository, error) {
		return repository.NewPostgresCabinetItemRepository(i)
	})

	// Services
	do.Provide(i, func(i do.Injector) (service.IAuthService, error) {
		return service.NewAuthService(i)
	})
	do.Provide(i, func(i do.Injector) (service.IAdminService, error) {
		return service.NewAdminService(i)
	})
	do.Provide(i, func(i do.Injector) (service.ICabinetService, error) {
		return service.NewCabinetService(i)
	})
	do.Provide(i, func(i do.Injector) (service.IMedicineService, error) {
		return service.NewMedicineService(i)
	})
	do.Provide(i, func(i do.Injector) (service.IProfileService, error) {
		return service.NewProfileService(i)
	})

	do.Provide(i, handler.NewAuthHandler)
	do.Provide(i, func(i do.Injector) (service.IPasswordHasher, error) {
		return service.NewBcryptHasher(i)
	})
	do.Provide(i, service.NewJWTManager)

	router := gin.New()
	log.Println("Gin router created")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("Swagger registered")

	router.Use(middleware.Logger())
	router.Use(middleware.ResponseHandler())
	router.Use(middleware.ErrorHandler())
	log.Println("Middlewares registered")

	authHandler := do.MustInvoke[*handler.AuthHandler](i)
	authHandler.RegisterRoutes(router.Group("/"))
	log.Println("Routes registered")

	serverAddress := ":8080"
	log.Printf("Starting server on %s", serverAddress)

	err := router.Run(serverAddress)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
