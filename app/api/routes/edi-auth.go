package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/internal/core/services"
	"backend/internal/handlers"
	"backend/internal/repositories"
)

func RoutesAuth(db *gorm.DB) *fiber.App {
	if db == nil {
		panic("Database connection is nil")
	}

	app := fiber.New()
	UserRepository := repositories.NewUserRepositoryDB(db)
	EmployeeRepository := repositories.NewEmployeeRepositoryDB(db)
	UserService := services.NewUserService(UserRepository)
	EmployeeService := services.NewEmployeeService(EmployeeRepository)
	UserHandler := handlers.NewAuthHandler(UserService, EmployeeService)

	// app.Post("/login/start", UserHandler.StartLoginHandler)
	app.Post("/login/verify", UserHandler.VerifyLoginHandler)

	return app
}
