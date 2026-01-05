package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/app/middlewares"
	"backend/internal/core/services"
	"backend/internal/handlers"
	"backend/internal/repositories"
)

func RoutesEDIOrder(db *gorm.DB) *fiber.App {
	if db == nil {
		panic("Database connection is nil")
	}

	app := fiber.New()

	EDIOrderRepository := repositories.NewEDIOrderRepositoryDB(db)
	EDIOrderService := services.NewEDIOrderService(EDIOrderRepository)
	EDIOrderHandler := handlers.NewEDIOrderHandler(EDIOrderService)

	app.Post("/create-edi-order", middlewares.JWTProtected(), EDIOrderHandler.CreateNewOrderWithVersionHandler)
	app.Get("/get-orders-active-top", middlewares.JWTProtected(), EDIOrderHandler.GetEDIOrderWithActiveTopHandler)
	app.Patch("/mark-order-as-read/:edi_order_id", middlewares.JWTProtected(), EDIOrderHandler.MarkOrderAsRead)
	app.Get("/get-order-detail-by-number", middlewares.JWTProtected(), EDIOrderHandler.GetEDIOrderDetailByNumberHandler)
	app.Post("/create-order-version", middlewares.JWTProtected(), EDIOrderHandler.CreateEDIOrderVersionHandler)
	app.Get("/get-status-order-summary-by-vendor-code", middlewares.JWTProtected(), EDIOrderHandler.GetStatusOrderSummaryByVendorCodeHandler)
	app.Get("/get-order-by-vendor-code", middlewares.JWTProtected(), EDIOrderHandler.GetOrderHeaderByVendorCodeHandler)
	app.Get("/get-order-by-forecaset", middlewares.JWTProtected(), EDIOrderHandler.GetOrderHeaderByNumberForecastHandler)
	app.Get("/get-order-version-by-order-number/:order_number", EDIOrderHandler.GetOrderVersionStatusLogByOrderNumberAndApprovedHandler)
	// =============================================== Status Log =========================================================
	app.Post("/create-order-version-status-log", middlewares.JWTProtected(), EDIOrderHandler.CreateEDIOrderVersionStatusLogHandler)
	app.Get("/get-status-log-by-version-id/:order_version_id", middlewares.JWTProtected(), EDIOrderHandler.GetOrderVersionStatusLogByOrderVersionIDHandler)
	app.Get("/get-status-log-by-version-id-and-approved/:order_version_id", middlewares.JWTProtected(), EDIOrderHandler.GetOrderVersionStatusLogByOrderVersionIDAndApprovedHandler)

	return app
}
