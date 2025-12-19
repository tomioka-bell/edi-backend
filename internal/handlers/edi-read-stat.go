package handlers

import (
	"backend/internal/core/models"
	"backend/internal/core/services"

	"github.com/gofiber/fiber/v2"
)

func NewEDIReadStatHandler(s *services.EDIReadStatService) *EDIReadStatHandler {
	return &EDIReadStatHandler{service: s}
}

type EDIReadStatHandler struct {
	service *services.EDIReadStatService
}

func (h *EDIReadStatHandler) TrackRead(c *fiber.Ctx) error {
	var req models.EDIReadStatReq

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	err := h.service.TrackReadService(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to track read",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Read tracked successfully",
	})
}

func (h *EDIReadStatHandler) GetReadStatByVendorCodeHandler(c *fiber.Ctx) error {
	vendorCode := c.Query("vendor_code")
	if vendorCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "กรุณาระบุ vendor_code"})
	}

	readStatView, err := h.service.GetReadStatByVendorCodeService(vendorCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(readStatView)
}
