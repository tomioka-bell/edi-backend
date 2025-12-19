package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/app/middlewares"
	"backend/internal/core/services"
	"backend/internal/handlers"
	"backend/internal/repositories"
)

func RoutesEDIInvoice(db *gorm.DB) *fiber.App {
	if db == nil {
		panic("Database connection is nil")
	}

	app := fiber.New()

	EDIInvoiceRepository := repositories.NewEDIInvoiceRepositoryDB(db)
	EDIInvoiceService := services.NewEDIInvoiceService(EDIInvoiceRepository)
	EDIInvoiceHandler := handlers.NewEDIInvoiceHandler(EDIInvoiceService)

	app.Post("/create-edi-invoice", middlewares.JWTProtected(), EDIInvoiceHandler.CreateNewInvoiceWithVersionHandler)
	app.Get("/get-invoices-active-top", middlewares.JWTProtected(), EDIInvoiceHandler.GetEDIInvoiceWithActiveTopHandler)
	app.Patch("/mark-invoice-as-read/:edi_invoice_id", middlewares.JWTProtected(), EDIInvoiceHandler.MarkInvoiceAsRead)
	app.Get("/get-invoice-detail-by-number", middlewares.JWTProtected(), EDIInvoiceHandler.GetEDIInvoiceDetailByNumberHandler)
	app.Post("/create-invoice-version", middlewares.JWTProtected(), EDIInvoiceHandler.CreateEDIInvoiceVersionHandler)
	app.Get("/get-order-detail-by-number-order", middlewares.JWTProtected(), EDIInvoiceHandler.GetInvoiceDetailByNumberOrderHandler)
	app.Get("/get-status-invoice-summary-by-vendor-code", middlewares.JWTProtected(), EDIInvoiceHandler.GetStatusInvoiceSummaryByVendorCodeHandler)

	// =============================================== Status Log =========================================================
	app.Post("/create-invoice-version-status-log", middlewares.JWTProtected(), EDIInvoiceHandler.CreateEDIInvoiceVersionStatusLogHandler)
	app.Get("/get-status-log-by-version-id/:invoice_version_id", middlewares.JWTProtected(), EDIInvoiceHandler.GetInvoiceVersionStatusLogByInvoiceVersionIDHandler)
	app.Get("/get-status-log-by-version-id-and-approved/:invoice_version_id", middlewares.JWTProtected(), EDIInvoiceHandler.GetInvoiceVersionStatusLogByInvoiceVersionIDAndApprovedHandler)

	return app
}
