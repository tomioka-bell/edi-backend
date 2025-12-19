package handlers

import (
	"backend/internal/core/models"
	services "backend/internal/core/ports/services"

	"github.com/gofiber/fiber/v2"
)

type EDICompanyHandler struct {
	CompanySrv services.EDICompanyService
}

func NewEDICompanyHandler(insSrv services.EDICompanyService) *EDICompanyHandler {
	return &EDICompanyHandler{CompanySrv: insSrv}
}

func (h *EDICompanyHandler) CreateCompanyHandler(c *fiber.Ctx) error {
	var req models.EDICompanyResp
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload data"})
	}

	if err := h.CompanySrv.CreateCompanyService(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "company created successfully"})
}

func (h *EDICompanyHandler) CreateNotificationRecipientHandler(c *fiber.Ctx) error {
	var req models.EDICompanyNotificationRecipientResp
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	if err := h.CompanySrv.CreateNotificationRecipientService(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "notification recipient created successfully"})
}

func (h *EDICompanyHandler) GetCompanyByCompanyIDHandler(c *fiber.Ctx) error {
	companyID := c.Query("company_id")
	if companyID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "company_id is required"})
	}

	company, err := h.CompanySrv.GetCompanyByCompanyIDService(companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(company)
}

//=================================================================================================

func (h *EDICompanyHandler) GetEDIVendorNotificationRecipientByCompanyHandler(c *fiber.Ctx) error {
	company := c.Query("company")
	if company == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "company is required"})
	}

	recipients, err := h.CompanySrv.GetEDIVendorNotificationRecipientByCompanyService(company)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(recipients)
}

func (h *EDICompanyHandler) CreateVendorNotificationRecipientHandler(c *fiber.Ctx) error {
	var req models.EDIVendorNotificationRecipientResp
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	if err := h.CompanySrv.CreateVendorNotificationRecipientService(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "vendor notification recipient created successfully"})
}

func (h *EDICompanyHandler) DeleteNotificationRecipientVendorHandler(c *fiber.Ctx) error {
	vendorNotificationRecipientID := c.Query("vendor_notification_recipient_id")
	if vendorNotificationRecipientID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "vendor_notification_recipient_id is required"})
	}

	if err := h.CompanySrv.DeleteNotificationRecipientVendorService(vendorNotificationRecipientID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "vendor notification recipient deleted successfully"})
}
