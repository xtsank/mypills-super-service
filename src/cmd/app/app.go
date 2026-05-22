package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/xtsank/mypills-super-service/docs/swagger"
	"github.com/xtsank/mypills-super-service/src/internal/infra/email"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/config"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/db"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/repository"
	"github.com/xtsank/mypills-super-service/src/internal/service"
	"github.com/xtsank/mypills-super-service/src/internal/transport/handler"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"
)

type App struct {
	i      do.Injector
	router *gin.Engine
}

func (app *App) provideMiddleware() {
	do.Provide(app.i, middleware.NewLogger)
	do.Provide(app.i, config.NewConfig)
	do.Provide(app.i, db.NewDB)
}

func (app *App) provideRepo() {
	do.Provide(app.i, repository.NewPostgresUserRepository)
	do.Provide(app.i, repository.NewPostgresMedicineRepository)
	do.Provide(app.i, repository.NewPostgresCabinetItemRepository)
	do.Provide(app.i, repository.NewPostgresDictionaryRepository)
}

func (app *App) provideService() {
	do.Provide(app.i, service.NewAuthService)
	do.Provide(app.i, service.NewAdminService)
	do.Provide(app.i, service.NewCabinetService)
	do.Provide(app.i, service.NewMedicineService)
	do.Provide(app.i, service.NewProfileService)
<<<<<<< HEAD
	do.Provide(app.i, service.NewDictionaryService)
=======
>>>>>>> 1f83dea7bd71d6b52bdd54933e14f6e23c6bc04a
	do.Provide(app.i, email.NewSMTPSender)
	do.Provide(app.i, service.NewNotificationService)
	do.Provide(app.i, service.NewBcryptHasher)
	do.Provide(app.i, service.NewJWTManager)
}

func (app *App) provideHandler() {
	do.Provide(app.i, handler.NewAuthHandler)
	do.Provide(app.i, handler.NewCabinetHandler)
	do.Provide(app.i, handler.NewProfileHandler)
	do.Provide(app.i, handler.NewMedicineHandler)
	do.Provide(app.i, handler.NewAdminHandler)
	do.Provide(app.i, handler.NewDictionaryHandler)
}

func (app *App) provideAll() {
	app.provideMiddleware()
	app.provideRepo()
	app.provideService()
	app.provideHandler()
}

func (app *App) initSwagger() {
	swagger := app.router.Group("/swagger")
	swagger.Use(middleware.TokenVerifier(app.i))
	swagger.Use(middleware.AdminOnly())
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (app *App) initMiddlewares() {
	app.router.Use(middleware.Cors())
	app.router.Use(middleware.Logger(app.i))
	app.router.Use(middleware.ResponseHandler())
	app.router.Use(middleware.ErrorHandler())
}

func (app *App) initRoutes() {
	api := app.router.Group("/")

	protected := api.Group("/")
	protected.Use(middleware.TokenVerifier(app.i))
	protected.Use(middleware.AdminOnly())

	authHandler := do.MustInvoke[*handler.AuthHandler](app.i)
	authHandler.RegisterRoutes(api)

	dictionaryHandler := do.MustInvoke[*handler.DictionaryHandler](app.i)
	dictionaryHandler.RegisterRoutes(api)

	cabinetHandler := do.MustInvoke[*handler.CabinetHandler](app.i)
	cabinetHandler.RegisterRoutes(protected)

	profileHandler := do.MustInvoke[*handler.ProfileHandler](app.i)
	profileHandler.RegisterRoutes(protected)

	medicineHandler := do.MustInvoke[*handler.MedicineHandler](app.i)
	medicineHandler.RegisterRoutes(protected)

	adminHandler := do.MustInvoke[*handler.AdminHandler](app.i)
	adminHandler.RegisterRoutes(protected)
}

func (app *App) initAll() {
	app.initSwagger()
	app.initMiddlewares()
	app.initRoutes()
}

func NewApp() *App {
	gin.SetMode(gin.ReleaseMode)

	i := do.New()
	router := gin.New()

	app := &App{i, router}

	app.provideAll()
	app.initAll()

	return app
}

func (app *App) Run() error {
	cfg := do.MustInvoke[*config.Config](app.i)
	logger := do.MustInvoke[*slog.Logger](app.i)
	notificationService := do.MustInvoke[service.INotificationService](app.i)

<<<<<<< HEAD
	startNotificationTicker(notificationService, cfg, logger)
=======
	startNotificationTicker(notificationService, cfg)
>>>>>>> 1f83dea7bd71d6b52bdd54933e14f6e23c6bc04a

	addr := cfg.ServerAddress
	if addr != "" && addr[0] != ':' {
		addr = ":" + addr
	}

	logger.Info("server_start", slog.String("address", addr))

	return app.router.Run(addr)
}

<<<<<<< HEAD
func startNotificationTicker(service service.INotificationService, cfg *config.Config, logger *slog.Logger) {
	interval := cfg.NotificationCheckInterval
	logger.Info("notification_ticker_start", slog.String("interval", interval.String()))
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			logger.Debug("notification_ticker_tick")
=======
func startNotificationTicker(service service.INotificationService, cfg *config.Config) {
	ticker := time.NewTicker(cfg.NotificationCheckInterval)
	go func() {
		for range ticker.C {
>>>>>>> 1f83dea7bd71d6b52bdd54933e14f6e23c6bc04a
			service.SendDueNotifications(context.Background())
		}
	}()
}

func (app *App) Logger() *slog.Logger {
	return do.MustInvoke[*slog.Logger](app.i)
}
