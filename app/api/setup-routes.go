package api

import (
	routes "backend/app/api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, redisClient *redis.Client) {
	if db == nil {
		panic("Database connection is nil")
	}

	api := app.Group("/api", logger.New())
	api.Mount("/auth", routes.RoutesAuth(db))
	api.Mount("/user", routes.RoutesUser(db))
	api.Mount("/employee", routes.RoutesEmployee(db))
	api.Mount("/questionnaire", routes.RoutesQuestionnaire(db))
	api.Mount("/uploads", routes.RoutesUpload())
	api.Mount("/forecasts", routes.RoutesForecast(db))
	api.Mount("/company", routes.RoutesCompany(db))
	api.Mount("/vendor-metrics", routes.RoutesVendorMetrics(db))
	api.Mount("/orders", routes.RoutesEDIOrder(db))
	api.Mount("/invoice", routes.RoutesEDIInvoice(db))
	api.Mount("/summary-data", routes.RoutesSummaryData(db, redisClient))
}
