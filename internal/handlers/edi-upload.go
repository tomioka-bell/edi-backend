package handlers

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func ServeUploadFile(c *fiber.Ctx) error {
	folder := c.Query("folder")
	filename := c.Query("filename")

	if folder == "" || filename == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid path")
	}

	filePath := filepath.Join("uploads", folder, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.SendFile(filePath)
}

// func TestSendEmailHandler(c *fiber.Ctx) error {
// 	email := c.Query("email")
// 	if email == "" {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "email is required",
// 		})
// 	}

// 	err := mailer.SendTestEmail(email)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "Test email sent successfully",
// 	})
// }
