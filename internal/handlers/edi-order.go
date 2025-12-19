package handlers

import (
	"backend/internal/core/models"
	services "backend/internal/core/ports/services"
	"backend/internal/pkgs/mailer"
	"backend/internal/pkgs/utils"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
)

type EDIOrderHandler struct {
	EDIOrderSrv services.EDIOrderService
}

func NewEDIOrderHandler(srv services.EDIOrderService) *EDIOrderHandler {
	return &EDIOrderHandler{EDIOrderSrv: srv}
}

func (h *EDIOrderHandler) CreateNewOrderWithVersionHandler(c *fiber.Ctx) error {
	var req models.CreateOrderReq

	payload := c.FormValue("payload")
	if strings.TrimSpace(payload) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payload is required",
		})
	}

	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		log.Println("Unmarshal payload error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid payload json",
		})
	}

	fileURL, err := utils.UploadFileFromForm(
		c,
		"file",
		"uploads/document/order",
		"/uploads/document/order",
	)
	if err != nil {
		log.Println("Upload file error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to upload file",
		})
	}

	trim := func(s string) string { return strings.TrimSpace(s) }
	toUpper := func(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }

	var numberOrder string
	if strings.TrimSpace(req.NumberOrder) == "" {
		gen, err := h.EDIOrderSrv.GenerateRunningNumberService()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate running number",
			})
		}
		numberOrder = gen
	} else {
		numberOrder = req.NumberOrder
	}

	header := models.EDIOrderResp{
		NumberOrder:           numberOrder,
		NumberForecast:        req.NumberForecast,
		VendorCode:            req.VendorCode,
		FileURL:               fileURL,
		CreatedByExternalID:   trim(req.CreatedByExternalID),
		CreatedBySourceSystem: toUpper(req.CreatedBySourceSystem),
	}

	if len(req.Versions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "at least one version is required",
		})
	}

	v0 := req.Versions[0]
	firstVer := models.EDIOrderVersionResp{
		PeriodFrom:            v0.PeriodFrom,
		PeriodTo:              v0.PeriodTo,
		StatusOrder:           v0.StatusOrder,
		ReadOrder:             v0.ReadOrder,
		Note:                  v0.Note,
		SourceFileURL:         v0.SourceFileURL,
		CreatedByExternalID:   trim(v0.CreatedByExternalID),
		CreatedBySourceSystem: toUpper(v0.CreatedBySourceSystem),
	}

	if fileURL != nil && firstVer.SourceFileURL == nil {
		firstVer.SourceFileURL = fileURL
	}

	var fileURLStr string
	if fileURL != nil {
		fileURLStr = strings.TrimSpace(*fileURL)
	}

	changedBy := strings.ToUpper(strings.TrimSpace(firstVer.CreatedBySourceSystem))

	fmt.Printf("ข้อมูลที่ส่งมา : %q\n", changedBy)

	dataVendorMetrics, err := h.EDIOrderSrv.GetEDIVendorNotificationRecipientByCompanyService(header.VendorCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get VendorMetrics",
		})
	}

	var toEmails []string
	var company string

	for _, v := range dataVendorMetrics {
		if v.Principal == nil {
			continue
		}
		email := strings.TrimSpace(v.Principal.Email)
		if email == "" {
			continue
		}
		toEmails = append(toEmails, email)

		if company == "" {
			company = v.Company
		}
	}

	if len(toEmails) > 0 && company != "" {
		if err := mailer.SendStatusOrderVendorEmail(
			toEmails,
			v0.StatusOrder,
			company,
			header.NumberOrder,
			fileURLStr,
			v0.Note,
		); err != nil {
			log.Println("failed to send order status email to vendor:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to send order status email to vendor",
			})
		}
	} else {
		log.Println("No valid vendor email recipient found, skip sending vendor email")
	}

	log.Printf("[DEBUG] header FK: ext=%q src=%q", header.CreatedByExternalID, header.CreatedBySourceSystem)

	created, err := h.EDIOrderSrv.CreateNewOrderWithVersion(c.Context(), &header, &firstVer)
	if err != nil {
		log.Println("CreateNewOrderWithVersion error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create Order",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *EDIOrderHandler) GetEDIOrderWithActiveTopHandler(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	vendorCode := c.Query("vendor_code")
	rows, err := h.EDIOrderSrv.GetEDIOrderWithActiveTopService(limit, vendorCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve Orders"})
	}
	return c.JSON(rows)
}

func (h *EDIOrderHandler) GetStatusOrderSummaryByVendorCodeHandler(c *fiber.Ctx) error {
	number := c.Query("vendor_code")
	row, err := h.EDIOrderSrv.GetStatusOrderSummaryByVendorCodeService(number)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve forecast"})
	}
	return c.JSON(row)
}

func (h *EDIOrderHandler) MarkOrderAsRead(c *fiber.Ctx) error {
	idStr := c.Params("edi_order_id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}

	var mssqlID mssql.UniqueIdentifier
	copy(mssqlID[:], uid[:])

	if err := h.EDIOrderSrv.MarkOrderAsReadService(mssqlID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to update read_order",
		})
	}

	return c.JSON(fiber.Map{
		"message": "order marked as read",
	})
}

func (h *EDIOrderHandler) GetEDIOrderDetailByNumberHandler(c *fiber.Ctx) error {
	number := c.Query("number_order")
	row, err := h.EDIOrderSrv.GetEDIOrderDetailByNumber(number)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve order"})
	}
	return c.JSON(row)
}

func (h *EDIOrderHandler) GetOrderHeaderByVendorCodeHandler(c *fiber.Ctx) error {
	VendorCode := c.Query("vendor_code")
	row, err := h.EDIOrderSrv.GetOrderHeaderByVendorCodeService(VendorCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve order"})
	}
	return c.JSON(row)
}

func (h *EDIOrderHandler) GetOrderHeaderByNumberForecastHandler(c *fiber.Ctx) error {
	VendorCode := c.Query("number_forecast")
	row, err := h.EDIOrderSrv.GetOrderHeaderByNumberForecastService(VendorCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve order"})
	}
	return c.JSON(row)
}

type UpdateOrderStatusReq struct {
	StatusOrder string `json:"status_order"`
}

func (h *EDIOrderHandler) UpdateStatusOrder(c *fiber.Ctx) error {
	req := UpdateOrderStatusReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	idStr := c.Params("edi_order_id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}

	var mssqlID mssql.UniqueIdentifier
	copy(mssqlID[:], uid[:])

	if err := h.EDIOrderSrv.UpdateStatusOrderService(mssqlID, req.StatusOrder); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to update status_order",
		})
	}

	return c.JSON(fiber.Map{
		"message": "order status updated",
	})
}

func (h *EDIOrderHandler) CreateEDIOrderVersionHandler(c *fiber.Ctx) error {

	// 1) Read JSON from FormData
	payloadStr := c.FormValue("payload")
	if payloadStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payload is required",
		})
	}

	var req models.EDIOrderVersionResp
	if err := json.Unmarshal([]byte(payloadStr), &req); err != nil {
		log.Println("Invalid JSON:", err)
		log.Println("Original payload:", payloadStr)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload JSON",
		})
	}

	log.Println("===== Received Order Version Payload =====")
	pretty, _ := json.MarshalIndent(req, "", "  ")
	log.Println(string(pretty))
	log.Println("=============================================")

	// 2) File handling
	fileHeader, err := c.FormFile("file")
	if err == nil && fileHeader != nil {
		log.Printf("File uploaded: %s (%d bytes)\n", fileHeader.Filename, fileHeader.Size)

		uploadedURL, upErr := utils.UploadFileFromForm(
			c,
			"file",
			"uploads/document/order",
			"/uploads/document/order",
		)

		if upErr != nil {
			log.Println("Error uploading file:", upErr)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to upload file",
			})
		}

		req.SourceFileURL = uploadedURL
		log.Println("File uploaded to:", *uploadedURL)

	} else {
		log.Println("No file uploaded:", err)
	}

	basicInfo, err := h.EDIOrderSrv.GetOrderBasicByIDService(req.EDIOrderID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get basic info",
		})
	}

	dataVendorMetrics, err := h.EDIOrderSrv.GetEDIVendorNotificationRecipientByCompanyService(basicInfo.VendorCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get VendorMetrics data for order",
		})
	}

	var fileURLStr string
	if req.SourceFileURL != nil {
		fileURLStr = strings.TrimSpace(*req.SourceFileURL)
	}

	var toEmails []string
	var company string

	for _, v := range dataVendorMetrics {
		if v.Principal == nil {
			continue
		}
		email := strings.TrimSpace(v.Principal.Email)
		if email == "" {
			continue
		}
		toEmails = append(toEmails, email)

		if company == "" {
			company = v.Company
		}
	}

	if len(toEmails) > 0 && company != "" {
		if err := mailer.SendModifyOrderVendorEmail(
			toEmails,
			company,
			basicInfo.NumberOrder,
			fileURLStr,
			req.Note,
		); err != nil {
			log.Println("failed to send order modify email to vendor:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to send order modify email to vendor",
			})
		}
	} else {
		log.Println("No valid vendor email recipient found, skip sending vendor email")
	}

	if h.EDIOrderSrv == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Service is not available",
		})
	}

	if err := h.EDIOrderSrv.CreateEDIOrderVersionService(req); err != nil {
		log.Println("Error creating EDIOrderVersion:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create EDIOrderVersion",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "EDIOrderVersion created successfully",
	})
}
