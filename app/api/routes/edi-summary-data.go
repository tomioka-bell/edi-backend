// api/summary_routes.go
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"backend/app/middlewares"
	"backend/internal/core/services"
	"backend/internal/handlers"
	"backend/internal/jobs"
	"backend/internal/repositories"
)

func RoutesSummaryData(db *gorm.DB, redisClient *redis.Client) *fiber.App {
	if db == nil {
		panic("Database connection is nil")
	}

	app := fiber.New()

	summaryRepo := repositories.NewEDISummaryDataRepositoryDB(db, redisClient)
	summarySvc := services.NewEDISummaryDataService(summaryRepo)
	summaryHandler := handlers.NewEDISummaryDataHandler(summarySvc)

	// ---------- HTTP routes ----------
	app.Get("/get-status-summary", middlewares.JWTProtected(), summaryHandler.GetAllStatusSummaryDataHandler)
	app.Get("/get-number-count-summary", middlewares.JWTProtected(), summaryHandler.GetAllTotalCountSummaryHandler)
	app.Get("/get-status-total-summary", middlewares.JWTProtected(), summaryHandler.GetAllStatusTotalSummaryHandler)
	app.Get("/get-monthly-status-summary", middlewares.JWTProtected(), summaryHandler.GetAllMonthlyStatusSummaryHandler)
	app.Get("/count-vendor", middlewares.JWTProtected(), summaryHandler.CountUserHandler)

	app.Get("/get-vendor-flat-summary", middlewares.JWTProtected(), summaryHandler.GetVendorFlatSummaryHandler)
	app.Get("/period-alert-forecast", middlewares.JWTProtected(), summaryHandler.GetForecastPeriodAlertsHandler)
	app.Get("/period-alert-order", middlewares.JWTProtected(), summaryHandler.GetOrderPeriodAlertsHandler)
	app.Get("/period-alert-invoice", middlewares.JWTProtected(), summaryHandler.GetInvoicePeriodAlertsHandler)

	// ---------- cron job ----------
	_ = jobs.StartForecastPeriodAlertCron(summarySvc)
	_ = jobs.StartOrderPeriodAlertCron(summarySvc)

	return app
}
