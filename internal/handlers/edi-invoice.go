package handlers

import (
	"backend/internal/core/models"
	services "backend/internal/core/ports/services"
	"backend/internal/pkgs/mailer"
	"backend/internal/pkgs/utils"
	"encoding/json"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
)

type EDIInvoiceHandler struct {
	EDIInvoiceSrv services.EDIInvoiceService
}

func NewEDIInvoiceHandler(srv services.EDIInvoiceService) *EDIInvoiceHandler {
	return &EDIInvoiceHandler{EDIInvoiceSrv: srv}
}

func (h *EDIInvoiceHandler) CreateNewInvoiceWithVersionHandler(c *fiber.Ctx) error {
	var req models.CreateInvoiceReq

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
		"uploads/document/invoice",
		"/uploads/document/invoice",
	)

	if err != nil {
		log.Println("Upload file error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to upload file",
		})
	}

	trim := func(s string) string { return strings.TrimSpace(s) }
	toUpper := func(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }

	var numberInvoice string
	if strings.TrimSpace(req.NumberInvoice) == "" {
		gen, err := h.EDIInvoiceSrv.GenerateRunningNumberService()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate running number",
			})
		}
		numberInvoice = gen
	} else {
		numberInvoice = req.NumberInvoice
	}

	header := models.EDIInvoiceResp{
		NumberInvoice: numberInvoice,
		NumberOrder:   req.NumberOrder,
		InvoiceType:   req.InvoiceType,
		VendorCode:    req.VendorCode,

		FileURL:               fileURL,
		CreatedByExternalID:   trim(req.CreatedByExternalID),
		CreatedBySourceSystem: toUpper(req.CreatedBySourceSystem),
	}

	if len(req.Versions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "at least one version is required",
		})
	}

	var fileURLStr string
	if fileURL != nil {
		fileURLStr = strings.TrimSpace(*fileURL)
	}

	v0 := req.Versions[0]
	firstVer := models.EDIInvoiceVersionResp{
		PeriodFrom:            v0.PeriodFrom,
		PeriodTo:              v0.PeriodTo,
		StatusInvoice:         v0.StatusInvoice,
		ReadInvoice:           v0.ReadInvoice,
		Note:                  v0.Note,
		Quantity:              v0.Quantity,
		SourceFileURL:         v0.SourceFileURL,
		CreatedByExternalID:   trim(v0.CreatedByExternalID),
		CreatedBySourceSystem: toUpper(v0.CreatedBySourceSystem),
	}

	dataVendorMetrics, err := h.EDIInvoiceSrv.GetEDIVendorNotificationRecipientByCompanyService(header.VendorCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get VendorMetrics",
		})
	}

	if fileURL != nil && firstVer.SourceFileURL == nil {
		firstVer.SourceFileURL = fileURL
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
		if err := mailer.SendStatusInvoiceEmployeeEmail(
			company,
			v0.StatusInvoice,
			header.NumberInvoice,
			fileURLStr,
			v0.Note,
			"INVOICE",
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

	created, err := h.EDIInvoiceSrv.CreateNewInvoiceWithVersion(c.Context(), &header, &firstVer)
	if err != nil {
		log.Println("CreateNewInvoiceWithVersion error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create Invoice",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *EDIInvoiceHandler) GetEDIInvoiceWithActiveTopHandler(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	vendorCode := c.Query("vendor_code")
	rows, err := h.EDIInvoiceSrv.GetEDIInvoiceWithActiveTopService(limit, vendorCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve Invoices"})
	}
	return c.JSON(rows)
}

func (h *EDIInvoiceHandler) GetStatusInvoiceSummaryByVendorCodeHandler(c *fiber.Ctx) error {
	number := c.Query("vendor_code")
	row, err := h.EDIInvoiceSrv.GetStatusInvoiceSummaryByVendorCodeService(number)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve forecast"})
	}
	return c.JSON(row)
}

func (h *EDIInvoiceHandler) MarkInvoiceAsRead(c *fiber.Ctx) error {
	idStr := c.Params("edi_invoice_id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}

	var mssqlID mssql.UniqueIdentifier
	copy(mssqlID[:], uid[:])

	if err := h.EDIInvoiceSrv.MarkInvoiceAsReadService(mssqlID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to update read_Invoice",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Invoice marked as read",
	})
}

func (h *EDIInvoiceHandler) GetEDIInvoiceDetailByNumberHandler(c *fiber.Ctx) error {
	number := c.Query("number_invoice")
	row, err := h.EDIInvoiceSrv.GetEDIInvoiceDetailByNumber(number)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve Invoice"})
	}
	return c.JSON(row)
}

func (h *EDIInvoiceHandler) GetInvoiceDetailByNumberOrderHandler(c *fiber.Ctx) error {
	number := c.Query("number_order")
	row, err := h.EDIInvoiceSrv.GetInvoiceDetailByNumberOrderService(number)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve Invoice"})
	}
	return c.JSON(row)
}

type UpdateInvoiceStatusReq struct {
	StatusInvoice string `json:"status_invoice"`
}

func (h *EDIInvoiceHandler) UpdateStatusInvoice(c *fiber.Ctx) error {
	req := UpdateInvoiceStatusReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	idStr := c.Params("edi_invoice_id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}

	var mssqlID mssql.UniqueIdentifier
	copy(mssqlID[:], uid[:])

	if err := h.EDIInvoiceSrv.UpdateStatusInvoiceService(mssqlID, req.StatusInvoice); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to update status_Invoice",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Invoice status updated",
	})
}

func (h *EDIInvoiceHandler) CreateEDIInvoiceVersionHandler(c *fiber.Ctx) error {

	// 1) Read JSON from FormData
	payloadStr := c.FormValue("payload")
	if payloadStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payload is required",
		})
	}

	var req models.EDIInvoiceVersionResp
	if err := json.Unmarshal([]byte(payloadStr), &req); err != nil {
		log.Println("Invalid JSON:", err)
		log.Println("Original payload:", payloadStr)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload JSON",
		})
	}

	log.Println("===== Received Invoice Version Payload =====")
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
			"uploads/document/invoice",
			"/uploads/document/invoice",
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

	basicInfo, err := h.EDIInvoiceSrv.GetInvoiceBasicByIDService(req.EDIInvoiceID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get basic info",
		})
	}

	var fileURLStr string
	if req.SourceFileURL != nil {
		fileURLStr = strings.TrimSpace(*req.SourceFileURL)
	}

	if err := mailer.SendModifyInvoiceVendorEmail(
		basicInfo.VendorCode,
		basicInfo.NumberInvoice,
		fileURLStr,
		req.Note,
		"INVOICE",
	); err != nil {
		log.Println("failed to send invoice modify email to vendor:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to send invoice modify email to vendor",
		})
	}

	if h.EDIInvoiceSrv == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Service is not available",
		})
	}

	ediInvoiceVersionID, err := h.EDIInvoiceSrv.CreateEDIInvoiceVersionService(req)
	if err != nil {
		log.Println("Error creating EDIInvoiceVersion:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create EDIInvoiceVersion",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":                "EDIInvoiceVersion created successfully",
		"edi_invoice_version_id": ediInvoiceVersionID,
	})
}
