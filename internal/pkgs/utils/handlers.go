package utils

import (
	"backend/internal/pkgs/errs"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func HandleError(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case errs.AppError:
		fmt.Fprintln(c, e)
		return c.SendStatus(e.Code)
	case error:
		fmt.Fprintln(c, e)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return nil
}
