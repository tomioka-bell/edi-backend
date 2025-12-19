package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/app/middlewares"
	"backend/internal/core/services"
	"backend/internal/handlers"
	"backend/internal/repositories"
)

func RoutesVendorMetrics(db *gorm.DB) *fiber.App {
	if db == nil {
		panic("Database connection is nil")
	}

	app := fiber.New()

	VendorMetricsRepository := repositories.NewEDIVendorMetricsRepositoryDB(db)
	VendorMetricsService := services.NewEDIVendorMetricsService(VendorMetricsRepository)
	VendorMetricsHandler := handlers.NewEDIVendorMetricsHandler(VendorMetricsService)

	app.Get("/get-company-by-company", middlewares.JWTProtected(), VendorMetricsHandler.GetVendorMetricsByCompanyHandler)
	app.Get("/get-vendor-metrics-top", middlewares.JWTProtected(), VendorMetricsHandler.GetAllEDIVendorMetricsTopHandler)
	app.Post("/create-vendor-metrics", middlewares.JWTProtected(), VendorMetricsHandler.CreateVendorMetricsHandler)
	app.Patch("/update-vendor-metrics/:vendor_metrics_id", middlewares.JWTProtected(), VendorMetricsHandler.UpdateVendorMetricsHandler)
	app.Delete("/delete-vendor-metrics/:vendor_metrics_id", middlewares.JWTProtected(), VendorMetricsHandler.DeleteVendorMetricsHandler)
	app.Get("/get-all-vendor-metrics", middlewares.JWTProtected(), VendorMetricsHandler.GetAllVendorMetricsHandler)

	return app
}
