package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UploadFileFromForm(c *fiber.Ctx, fieldName, uploadDir, publicBasePath string) (*string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil || file == nil {
		return nil, nil
	}

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to prepare upload directory: %w", err)
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	fullPath := filepath.Join(uploadDir, filename)

	if err := c.SaveFile(file, fullPath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	url := path.Join(publicBasePath, filename)
	return &url, nil
}
