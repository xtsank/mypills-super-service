package main

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/swaggo/files"
	swaggerFilesV2 "github.com/swaggo/files/v2"
	"github.com/swaggo/gin-swagger"
	_ "github.com/xtsank/mypills-super-service/docs/swagger"
	"github.com/xtsank/mypills-super-service/src/internal/config"
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
}

func (app *App) provideRepo() {
	do.Provide(app.i, pgRepo.NewPostgresUserRepository)
	do.Provide(app.i, pgRepo.NewPostgresMedicineRepository)
	do.Provide(app.i, pgRepo.NewPostgresMedicineCatalogRepository)
	do.Provide(app.i, pgRepo.NewPostgresLegacyCabinetItemRepository)
	do.Provide(app.i, pgRepo.NewPostgresCabinetRepository)
	do.Provide(app.i, pgRepo.NewPostgresCabinetItemRepository)
	do.Provide(app.i, pgRepo.NewPostgresRestockRepository)
	do.Provide(app.i, pgRepo.NewPostgresDictionaryRepository)
}

func (app *App) provideService() {
	do.Provide(app.i, service.NewAuthService)
	do.Provide(app.i, service.NewAdminService)
	do.Provide(app.i, service.NewLegacyCabinetService)
	do.Provide(app.i, service.NewCabinetService)
	do.Provide(app.i, service.NewCabinetItemsService)
	do.Provide(app.i, service.NewRestockService)
	do.Provide(app.i, service.NewLegacyMedicineService)
	do.Provide(app.i, service.NewMedicineCatalogService)
	do.Provide(app.i, service.NewRecommendationService)
	do.Provide(app.i, service.NewProfileService)
	do.Provide(app.i, service.NewDictionaryService)
	do.Provide(app.i, service.NewBcryptHasher)
	do.Provide(app.i, service.NewJWTManager)
}

func (app *App) provideHandler() {
	do.Provide(app.i, handler.NewLegacyAuthHandler)
	do.Provide(app.i, handler.NewLegacyCabinetHandler)
	do.Provide(app.i, handler.NewCabinetHandler)
	do.Provide(app.i, handler.NewCabinetItemsHandler)
	do.Provide(app.i, handler.NewAuthHandler)
	do.Provide(app.i, handler.NewProfileHandler)
	do.Provide(app.i, handler.NewReferenceHandler)
	do.Provide(app.i, handler.NewRestockHandler)
	do.Provide(app.i, handler.NewMedicineCatalogHandler)
	do.Provide(app.i, handler.NewRecommendationHandler)
	do.Provide(app.i, handler.NewLegacyProfileHandler)
	do.Provide(app.i, handler.NewLegacyMedicineHandler)
	do.Provide(app.i, handler.NewLegacyAdminHandler)
	do.Provide(app.i, handler.NewLegacyDictionaryHandler)
}

func (app *App) provideAll() {
	app.provideMiddleware()
	app.provideRepo()
	app.provideService()
	app.provideHandler()
}

func (app *App) initSwagger() {
	app.router.GET("/api/spec/v2.yaml", func(c *gin.Context) {
		c.File("docs/openapi/v2.yaml")
	})
	app.router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/")
	})
	swagger := app.router.Group("/swagger")
	v2Files := http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerFilesV2.FS)))
	swagger.GET("/*any", func(c *gin.Context) {
		if c.Param("any") == "/swagger-initializer.js" {
			c.Data(http.StatusOK, "application/javascript", []byte(`window.onload = function () {
  const ui = SwaggerUIBundle({
    url: "/api/spec/v2.yaml",
    dom_id: "#swagger-ui",
    deepLinking: true,
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    layout: "StandaloneLayout",
    responseInterceptor: function (response) {
      if (response.status === 200 && response.url &&
          new URL(response.url, window.location.origin).pathname === "/api/v2/sessions") {
        try {
          const body = response.body || JSON.parse(response.text || "{}");
          if (typeof body.token === "string") {
            ui.preauthorizeApiKey("bearerAuth", body.token);
          }
        } catch (error) {
          console.warn("Could not authorize Swagger UI with the issued token", error);
        }
      }
      return response;
    }
  });
};`))
			return
		}
		v2Files.ServeHTTP(c.Writer, c.Request)
	})
	app.router.GET("/swagger-v1", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger-v1/index.html")
	})
	legacySwagger := app.router.Group("/swagger-v1")
	v1Files := *swaggerFiles.Handler
	v1Swagger := ginSwagger.WrapHandler(&v1Files)
	legacySwagger.GET("/*any", func(c *gin.Context) {
		if c.Param("any") == "/" {
			c.Redirect(http.StatusFound, "/swagger-v1/index.html")
			return
		}
		v1Swagger(c)
	})
}

func (app *App) initMiddlewares() {
	app.router.Use(middleware.Cors())
	app.router.Use(middleware.Logger(app.i))
	app.router.Use(middleware.ResponseHandler())
	app.router.Use(middleware.ErrorHandler())
}

func (app *App) initRoutes() {
	authHandler := do.MustInvoke[*handler.LegacyAuthHandler](app.i)
	dictionaryHandler := do.MustInvoke[*handler.LegacyDictionaryHandler](app.i)
	cabinetHandler := do.MustInvoke[*handler.LegacyCabinetHandler](app.i)
	profileHandler := do.MustInvoke[*handler.LegacyProfileHandler](app.i)
	medicineHandler := do.MustInvoke[*handler.LegacyMedicineHandler](app.i)
	adminHandler := do.MustInvoke[*handler.LegacyAdminHandler](app.i)
	for _, prefix := range []string{"", "/api/v1"} {
		legacy := app.router.Group(prefix)
		authHandler.RegisterRoutes(legacy)
		dictionaryHandler.RegisterRoutes(legacy)

		protected := legacy.Group("")
		protected.Use(middleware.TokenVerifier(app.i))
		cabinetHandler.RegisterRoutes(protected)
		profileHandler.RegisterRoutes(protected)
		medicineHandler.RegisterRoutes(protected)

		admin := protected.Group("")
		admin.Use(middleware.AdminOnly())
		adminHandler.RegisterRoutes(admin)
	}

	v2Public := app.router.Group("/api/v2")
	do.MustInvoke[*handler.AuthHandler](app.i).RegisterRoutes(v2Public)
	do.MustInvoke[*handler.ReferenceHandler](app.i).RegisterRoutes(v2Public)

	v2Protected := app.router.Group("/api/v2")
	v2Protected.Use(middleware.TokenVerifier(app.i))
	do.MustInvoke[*handler.ProfileHandler](app.i).RegisterRoutes(v2Protected)
	do.MustInvoke[*handler.CabinetHandler](app.i).RegisterRoutes(v2Protected)
	do.MustInvoke[*handler.CabinetItemsHandler](app.i).RegisterRoutes(v2Protected)
	do.MustInvoke[*handler.RestockHandler](app.i).RegisterRoutes(v2Protected)
	do.MustInvoke[*handler.MedicineCatalogHandler](app.i).RegisterRoutes(v2Protected)
	do.MustInvoke[*handler.RecommendationHandler](app.i).RegisterRoutes(v2Protected)
	v2Admin := v2Protected.Group("")
	v2Admin.Use(middleware.AdminOnly())
	do.MustInvoke[*handler.MedicineCatalogHandler](app.i).RegisterAdminRoutes(v2Admin)
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
	addr := cfg.ServerAddress
	if addr != "" && addr[0] != ':' {
		addr = ":" + addr
	}

	logger.Info("server_start", slog.String("address", addr))

	return app.router.Run(addr)
}

func (app *App) Logger() *slog.Logger {
	return do.MustInvoke[*slog.Logger](app.i)
}
