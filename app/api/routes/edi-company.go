package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/app/middlewares"
	"backend/internal/core/services"
	"backend/internal/handlers"
	"backend/internal/repositories"
)

func RoutesCompany(db *gorm.DB) *fiber.App {
	if db == nil {
		panic("Database connection is nil")
	}

	app := fiber.New()

	CompanyRepository := repositories.NewEDICompanyRepositoryDB(db)
	CompanyService := services.NewEDICompanyService(CompanyRepository)
	CompanyHandler := handlers.NewEDICompanyHandler(CompanyService)

	app.Get("/get-company-by-company-id", middlewares.JWTProtected(), CompanyHandler.GetCompanyByCompanyIDHandler)
	app.Post("/create-company", middlewares.JWTProtected(), CompanyHandler.CreateCompanyHandler)
	// app.Post("/create-notification-recipient", CompanyHandler.CreateNotificationRecipientHandler)
	app.Get("/get-email-by-company", middlewares.JWTProtected(), CompanyHandler.GetEDIVendorNotificationRecipientByCompanyHandler)
	app.Post("/create-notification-recipient", middlewares.JWTProtected(), CompanyHandler.CreateVendorNotificationRecipientHandler)
	app.Delete("/delete-notification-recipient-vendor", middlewares.JWTProtected(), CompanyHandler.DeleteNotificationRecipientVendorHandler)

	return app
}
