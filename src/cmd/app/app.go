package main

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "github.com/xtsank/mypills-super-service/docs/swagger"
	"github.com/xtsank/mypills-super-service/src/internal/config"
	"github.com/xtsank/mypills-super-service/src/internal/infra/email"
	mongoDB "github.com/xtsank/mypills-super-service/src/internal/infra/mongo/db"
	mongoRepo "github.com/xtsank/mypills-super-service/src/internal/infra/mongo/repository"
	pgDB "github.com/xtsank/mypills-super-service/src/internal/infra/postgres/db"
	pgRepo "github.com/xtsank/mypills-super-service/src/internal/infra/postgres/repository"
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

	do.Provide(app.i, pgDB.NewPostgresDB)
	do.Provide(app.i, mongoDB.NewMongoDB)
}

func (app *App) provideRepo() {
	cfg := do.MustInvoke[*config.Config](app.i)
	switch strings.ToLower(cfg.DBProvider) {
	case "mongo":
		do.Provide(app.i, mongoRepo.NewMongoUserRepository)
		do.Provide(app.i, mongoRepo.NewMongoMedicineRepository)
		do.Provide(app.i, mongoRepo.NewMongoCabinetItemRepository)
		do.Provide(app.i, mongoRepo.NewMongoDictionaryRepository)
	case "postgres", "":
		do.Provide(app.i, pgRepo.NewPostgresUserRepository)
		do.Provide(app.i, pgRepo.NewPostgresMedicineRepository)
		do.Provide(app.i, pgRepo.NewPostgresCabinetItemRepository)
		do.Provide(app.i, pgRepo.NewPostgresDictionaryRepository)
	default:
		logger := do.MustInvoke[*slog.Logger](app.i)
		logger.Error("unknown_db_provider", slog.String("db_provider", cfg.DBProvider))
		panic("unknown DB provider")
	}
}

func (app *App) provideService() {
	do.Provide(app.i, service.NewAuthService)
	do.Provide(app.i, service.NewAdminService)
	do.Provide(app.i, service.NewCabinetService)
	do.Provide(app.i, service.NewMedicineService)
	do.Provide(app.i, service.NewProfileService)
	do.Provide(app.i, service.NewDictionaryService)
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

	startNotificationTicker(notificationService, cfg, logger)

	addr := cfg.ServerAddress
	if addr != "" && addr[0] != ':' {
		addr = ":" + addr
	}

	logger.Info("server_start", slog.String("address", addr))

	return app.router.Run(addr)
}

func startNotificationTicker(service service.INotificationService, cfg *config.Config, logger *slog.Logger) {
	interval := cfg.NotificationCheckInterval
	startDelay := time.Duration(0)
	if cfg.NotificationStartTime > 0 {
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		nextStart := midnight.Add(cfg.NotificationStartTime)
		if !nextStart.After(now) {
			nextStart = nextStart.Add(24 * time.Hour)
		}
		startDelay = time.Until(nextStart)
	}

	logger.Info(
		"notification_ticker_start",
		slog.String("interval", interval.String()),
		slog.String("start_delay", startDelay.String()),
	)

	go func() {
		if startDelay > 0 {
			timer := time.NewTimer(startDelay)
			<-timer.C
		}

		service.SendDueNotifications(context.Background())

		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			logger.Debug("notification_ticker_tick")
			service.SendDueNotifications(context.Background())
		}
	}()
}

func (app *App) Logger() *slog.Logger {
	return do.MustInvoke[*slog.Logger](app.i)
}
