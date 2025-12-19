package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Options struct {
	Dir          string
	AllowedMIMEs []string
	MaxSize      int64
	BaseURL      string
	Required     bool
}

func UploadFromForm(c *fiber.Ctx, formField string, opt Options) (string, string, error) {
	fileHeader, err := c.FormFile(formField)
	if err != nil || fileHeader == nil {
		if opt.Required {
			return "", "", errors.New("file is required")
		}
		return "", "", nil
	}

	if opt.MaxSize > 0 && fileHeader.Size > opt.MaxSize {
		return "", "", errors.New("file too large")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	buf := make([]byte, 512)
	n, _ := io.ReadFull(src, buf)
	detected := http.DetectContentType(buf[:n])

	if _, err := src.Seek(0, 0); err != nil {
		return "", "", err
	}

	if len(opt.AllowedMIMEs) > 0 {
		ok := false
		for _, m := range opt.AllowedMIMEs {
			if strings.EqualFold(m, detected) {
				ok = true
				break
			}
		}
		if !ok {
			return "", "", errors.New("unsupported file type: " + detected)
		}
	}

	if opt.Dir == "" {
		opt.Dir = "./uploads"
	}
	if err := os.MkdirAll(opt.Dir, os.ModePerm); err != nil {
		return "", "", err
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		if exts, _ := mime.ExtensionsByType(detected); len(exts) > 0 {
			ext = exts[0]
		}
	}
	random := make([]byte, 8)
	if _, err := rand.Read(random); err != nil {
		return "", "", err
	}
	newName := strings.Join([]string{
		time.Now().Format("20060102_150405"),
		hex.EncodeToString(random),
	}, "_") + ext

	relPath := filepath.Join(opt.Dir, newName)

	if err := c.SaveFile(fileHeader, relPath); err != nil {
		return "", "", err
	}

	publicURL := ""
	if opt.BaseURL != "" {
		webPath := strings.ReplaceAll(relPath, "\\", "/")
		publicURL = strings.TrimRight(opt.BaseURL, "/") + "/" + strings.TrimLeft(webPath, "./")
	}

	return relPath, publicURL, nil
}
