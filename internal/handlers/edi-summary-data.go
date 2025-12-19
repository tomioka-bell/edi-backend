package handlers

import (
	services "backend/internal/core/ports/services"

	"github.com/gofiber/fiber/v2"
)

type EDISummaryDataHandler struct {
	EDISummaryDataService services.EDISummaryDataService
}

func NewEDISummaryDataHandler(srv services.EDISummaryDataService) *EDISummaryDataHandler {
	return &EDISummaryDataHandler{EDISummaryDataService: srv}
}

func (h *EDISummaryDataHandler) GetAllStatusSummaryDataHandler(c *fiber.Ctx) error {
	result, err := h.EDISummaryDataService.GetAllStatusSummaryData()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *EDISummaryDataHandler) GetAllTotalCountSummaryHandler(c *fiber.Ctx) error {
	result, err := h.EDISummaryDataService.GetAllTotalCountSummary()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *EDISummaryDataHandler) GetAllStatusTotalSummaryHandler(c *fiber.Ctx) error {
	result, err := h.EDISummaryDataService.GetAllStatusTotalSummary()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *EDISummaryDataHandler) GetAllMonthlyStatusSummaryHandler(c *fiber.Ctx) error {
	result, err := h.EDISummaryDataService.GetAllMonthlyStatusSummary()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *EDISummaryDataHandler) CountUserHandler(c *fiber.Ctx) error {
	count, err := h.EDISummaryDataService.CountUserService()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"vendor": count,
	})
}

func (h *EDISummaryDataHandler) GetVendorFlatSummaryHandler(c *fiber.Ctx) error {
	vendorCode := c.Query("vendorCode")

	data, err := h.EDISummaryDataService.GetVendorFlatSummary(c.Context(), vendorCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

func (h *EDISummaryDataHandler) GetForecastPeriodAlertsHandler(c *fiber.Ctx) error {
	data, err := h.EDISummaryDataService.GetForecastPeriodAlerts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

func (h *EDISummaryDataHandler) GetOrderPeriodAlertsHandler(c *fiber.Ctx) error {
	data, err := h.EDISummaryDataService.GetOrderPeriodAlerts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

func (h *EDISummaryDataHandler) GetInvoicePeriodAlertsHandler(c *fiber.Ctx) error {
	data, err := h.EDISummaryDataService.GetInvoicePeriodAlerts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}
