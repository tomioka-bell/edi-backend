package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/internal/core/services"
	"backend/internal/handlers"
	"backend/internal/repositories"
)

func RoutesEmployee(db *gorm.DB) *fiber.App {
	if db == nil {
		panic("Database connection is nil")
	}

	app := fiber.New()

	EmployeeRepository := repositories.NewEmployeeRepositoryDB(db)
	EmployeeService := services.NewEmployeeService(EmployeeRepository)
	EmployeeHandler := handlers.NewEmployeeHandler(EmployeeService)

	app.Get("/get-employee-by-ad-logon", EmployeeHandler.GetEmployeeByADLogonHandler)
	app.Post("/login/ldap", EmployeeHandler.LoginWithLdapHandler)
	app.Post("/login/verify-code", EmployeeHandler.VerifyLoginCodeEmployeeHandler)

	return app
}
