package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/app/middlewares"
	"backend/internal/core/services"
	"backend/internal/handlers"
	"backend/internal/repositories"
)

func RoutesForecast(db *gorm.DB) *fiber.App {
	if db == nil {
		panic("Database connection is nil")
	}

	app := fiber.New()

	EDIForecastRepository := repositories.NewEDIForecastRepositoryDB(db)
	EDIForecastService := services.NewEDIForecastService(EDIForecastRepository)
	EDIReadStatRepository := repositories.NewEDIReadStatRepositoryDB(db)
	EDIReadStatService := services.NewEDIReadStatService(EDIReadStatRepository)
	EDIForecastHandler := handlers.NewEDIForecastHandler(EDIForecastService, EDIReadStatService)

	app.Post("/create-forecast", middlewares.JWTProtected(), EDIForecastHandler.CreateNewForecastWithVersionHandler)
	app.Get("/get-forecast-active-top", middlewares.JWTProtected(), EDIForecastHandler.GetEDIForecastWithActiveTopHandler)
	app.Get("/get-forecast-active-by-number", middlewares.JWTProtected(), EDIForecastHandler.GetEDIForecastWithActiveByNumberHandler)
	app.Post("/create-forecast-version", middlewares.JWTProtected(), EDIForecastHandler.CreateEDIForecastVersionHandler)
	app.Patch("/mark-forecast-as-read/:edi_forecast_id", middlewares.JWTProtected(), EDIForecastHandler.MarkForecastAsRead)
	app.Get("/get-forecast-version-by-id/:version_id", EDIForecastHandler.GetEDIForecastVersionByIDHandler)
	app.Get("/get-status-summary-by-vendor-code", EDIForecastHandler.GetStatusSummaryByVendorCodeHandler)

	// =============================================== Status Log =========================================================
	app.Post("/create-forecast-version-status-log", middlewares.JWTProtected(), EDIForecastHandler.CreateEDIForecastVersionStatusLogHandler)
	app.Get("/get-status-log-by-version-id/:forecast_version_id", middlewares.JWTProtected(), EDIForecastHandler.GetForecastVersionStatusLogByForecastVersionIDHandler)
	app.Get("/get-status-log-by-version-id-and-approved/:forecast_version_id", middlewares.JWTProtected(), EDIForecastHandler.GetForecastVersionStatusLogByForecastVersionIDAndApprovedHandler)

	return app
}
