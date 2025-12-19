package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/app/middlewares"
	"backend/internal/core/services"
	"backend/internal/handlers"
	"backend/internal/repositories"
)

func RoutesUser(db *gorm.DB) *fiber.App {
	if db == nil {
		panic("Database connection is nil")
	}

	app := fiber.New()

	UserRepository := repositories.NewUserRepositoryDB(db)
	UserService := services.NewUserService(UserRepository)
	UserHandler := handlers.NewUserHandler(UserService)

	app.Post("/create-user", middlewares.JWTProtected(), UserHandler.CreateUserHandler)
	app.Post("/sign-out", UserHandler.LogoutDBHandler)
	app.Get("/get-by-id/:user_id", middlewares.JWTProtected(), UserHandler.GetUserByIDHandler)
	app.Get("/get-all-user", middlewares.JWTProtected(), UserHandler.GetAllUserHandler)
	app.Patch("/update-user/:user_id", middlewares.JWTProtected(), UserHandler.UpdateUserHandler)

	app.Post("/login", UserHandler.LoginHandler)
	// app.Get("/get-by-profile", UserHandler.GetProfileHandler)
	app.Get("/get-by-profile", middlewares.JWTProtected(), UserHandler.GetProfileHandler)
	app.Get("/get-user-count", middlewares.JWTProtected(), UserHandler.GetUserCountHandler)

	app.Post("/request-password-reset", UserHandler.RequestPasswordReset)
	app.Post("/reset-password", UserHandler.ResetPassword)

	app.Post("/login/start", UserHandler.LoginUserHandler)
	app.Post("/login/verify", UserHandler.VerifyLoginCode)

	app.Get("/get-by-group", middlewares.JWTProtected(), UserHandler.GetEDIPrincipalUserByGroupHandler)
	app.Get("/get-by-company", middlewares.JWTProtected(), UserHandler.GetEDIUserByGroupHandler)
	app.Get("/get-active-emails-by-group", middlewares.JWTProtected(), UserHandler.GetActiveEmailsByGroupHandler)
	app.Patch("/update-principal/:edi_principal_id", middlewares.JWTProtected(), UserHandler.UpdatePrincipalWithMapHandler)

	return app
}
