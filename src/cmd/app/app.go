package main

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/xtsank/mypills-super-service/docs/swagger"
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

func (app *App) provideDB() {
	do.Provide(app.i, config.NewConfig)
	do.Provide(app.i, db.NewDB)
}

func (app *App) provideRepo() {
	do.Provide(app.i, repository.NewPostgresUserRepository)
	do.Provide(app.i, repository.NewPostgresMedicineRepository)
	do.Provide(app.i, repository.NewPostgresCabinetItemRepository)
}

func (app *App) provideService() {
	do.Provide(app.i, service.NewAuthService)
	do.Provide(app.i, service.NewAdminService)
	do.Provide(app.i, service.NewCabinetService)
	do.Provide(app.i, service.NewMedicineService)
	do.Provide(app.i, service.NewProfileService)
	do.Provide(app.i, service.NewBcryptHasher)
	do.Provide(app.i, service.NewJWTManager)
}

func (app *App) provideHandler() {
	do.Provide(app.i, handler.NewAuthHandler)
	do.Provide(app.i, handler.NewCabinetHandler)
	do.Provide(app.i, handler.NewProfileHandler)
	do.Provide(app.i, handler.NewMedicineHandler)
	do.Provide(app.i, handler.NewAdminHandler)
}

func (app *App) provideAll() {
	app.provideDB()
	app.provideRepo()
	app.provideService()
	app.provideHandler()
}

func (app *App) initSwagger() {
	app.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (app *App) initMiddlewares() {
	app.router.Use(middleware.Logger())
	app.router.Use(middleware.ResponseHandler())
	app.router.Use(middleware.ErrorHandler())
}

func (app *App) initRoutes() {
	api := app.router.Group("/")

	authHandler := do.MustInvoke[*handler.AuthHandler](app.i)
	authHandler.RegisterRoutes(api)

	cabinetHandler := do.MustInvoke[*handler.CabinetHandler](app.i)
	cabinetHandler.RegisterRoutes(api)

	profileHandler := do.MustInvoke[*handler.ProfileHandler](app.i)
	profileHandler.RegisterRoutes(api)

	medicineHandler := do.MustInvoke[*handler.MedicineHandler](app.i)
	medicineHandler.RegisterRoutes(api)

	adminHandler := do.MustInvoke[*handler.AdminHandler](app.i)
	adminHandler.RegisterRoutes(api)
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

	addr := cfg.ServerAddress
	if addr != "" && addr[0] != ':' {
		addr = ":" + addr
	}

	return app.router.Run(addr)
}
