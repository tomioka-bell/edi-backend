package api

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/handlers"
)

func RoutesUpload() *fiber.App {
	app := fiber.New()

	app.Get("/get-file", handlers.ServeUploadFile)
	// app.Get("/test-send-email", handlers.TestSendEmailHandler)

	return app
}
