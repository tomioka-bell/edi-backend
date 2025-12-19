package handlers

import (
	"backend/internal/core/models"
	services "backend/internal/core/ports/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
)

type EDIVendorMetricsHandler struct {
	VendorMetricsSrv services.EDIVendorMetricsService
}

func NewEDIVendorMetricsHandler(insSrv services.EDIVendorMetricsService) *EDIVendorMetricsHandler {
	return &EDIVendorMetricsHandler{VendorMetricsSrv: insSrv}
}

func (h *EDIVendorMetricsHandler) CreateVendorMetricsHandler(c *fiber.Ctx) error {
	var req models.VendorMetricsResp
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := h.VendorMetricsSrv.CreateVendorMetricsService(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Vendor metrics created successfully",
	})
}

func (h *EDIVendorMetricsHandler) GetVendorMetricsByCompanyHandler(c *fiber.Ctx) error {
	company := c.Query("company")

	metrics, err := h.VendorMetricsSrv.GetVendorMetricsByCompanyService(company)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(metrics)
}

func (h *EDIVendorMetricsHandler) GetAllVendorMetricsHandler(c *fiber.Ctx) error {
	metrics, err := h.VendorMetricsSrv.GetAllVendorMetricservice()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(metrics)
}

func (h *EDIVendorMetricsHandler) GetAllEDIVendorMetricsTopHandler(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 1000)
	metrics, err := h.VendorMetricsSrv.GetAllEDIVendorMetricsTopService(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(metrics)
}

func (h *EDIVendorMetricsHandler) UpdateVendorMetricsHandler(c *fiber.Ctx) error {
	idStr := c.Params("vendor_metrics_id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}

	var mssqlID mssql.UniqueIdentifier
	copy(mssqlID[:], uid[:])

	var req models.VendorMetricsUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	updates := map[string]interface{}{}
	if req.Initials != nil {
		updates["initials"] = *req.Initials
	}
	if req.CompanyName != nil {
		updates["company_name"] = *req.CompanyName
	}
	if req.ReminderDays != 0 {
		updates["reminder_days"] = req.ReminderDays
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}

	if err := h.VendorMetricsSrv.UpdateVendorMetricsService(mssqlID, updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Vendor metrics updated successfully",
	})
}

func (h *EDIVendorMetricsHandler) DeleteVendorMetricsHandler(c *fiber.Ctx) error {
	idStr := c.Params("vendor_metrics_id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}

	var mssqlID mssql.UniqueIdentifier
	copy(mssqlID[:], uid[:])

	if err := h.VendorMetricsSrv.DeleteVendorMetricsService(mssqlID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Vendor metrics deleted successfully",
	})
}
