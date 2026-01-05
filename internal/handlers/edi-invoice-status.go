package handlers

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend/internal/core/models"
	"backend/internal/pkgs/mailer"
	"backend/internal/pkgs/utils"
)

func (h *EDIInvoiceHandler) CreateEDIInvoiceVersionStatusLogHandler(c *fiber.Ctx) error {
	var form models.EDIInvoiceVersionStatusLogForm
	var StatusInvoice string

	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form payload",
		})
	}

	if h.EDIInvoiceSrv == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Service is not available",
		})
	}

	// -----------------------
	// Upload file (ใช้ helper)
	// -----------------------
	fileURL, err := utils.UploadFileFromForm(
		c,
		"file",
		"uploads/document/invoice",
		"/uploads/document/invoice",
	)
	if err != nil {
		log.Println("Error uploading file:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	// -----------------------
	// Mapping Status
	// -----------------------
	normalized := strings.ToUpper(strings.TrimSpace(form.NewStatus))

	switch normalized {
	case "REJECTED":
		StatusInvoice = "Rejected"
	case "FULLY_CONFIRMED":
		StatusInvoice = "Confirmed"
	case "APPROVED":
		StatusInvoice = "Approved"
	default:
		StatusInvoice = form.NewStatus
	}

	// -----------------------
	// แปลง string -> mssql.UniqueIdentifier
	// -----------------------
	ediInvoiceVersionID, err := ToMSSQLUUID(form.EDIInvoiceVersionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid edi_invoice_version_id",
		})
	}

	ediInvoiceID, err := ToMSSQLUUID(form.EDIInvoiceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid edi_invoice_id",
		})
	}

	// -----------------------
	// ประกอบ struct สำหรับ service
	// -----------------------
	var req models.EDIInvoiceVersionStatusLogResp
	req.EDIInvoiceVersionID = ediInvoiceVersionID
	req.EDIInvoiceID = ediInvoiceID
	req.OldStatus = form.OldStatus
	req.NewStatus = form.NewStatus
	req.Note = form.Note
	req.ChangedByExternalID = form.ChangedByExternalID
	req.ChangedBySourceSystem = form.ChangedBySourceSystem
	req.FileURL = fileURL
	req.CreatedAt = time.Now()

	// -----------------------
	// Service call
	// -----------------------
	if err := h.EDIInvoiceSrv.CreateEDIInvoiceVersionStatusLogService(req); err != nil {
		log.Println("Error creating EDIInvoiceVersionStatusLog:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create EDIInvoiceVersionStatusLog",
		})
	}

	// ดึงข้อมูล version ล่าสุด
	dataVersion, err := h.EDIInvoiceSrv.GetEDIInvoiceVersionByIDService(req.EDIInvoiceVersionID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get EDIInvoiceVersion",
		})
	}

	// -----------------------
	// เตรียมค่าไฟล์แนบ (string) แบบกัน panic
	// -----------------------
	var fileURLStr string
	if fileURL != nil {
		fileURLStr = strings.TrimSpace(*fileURL)
	}

	// -----------------------
	// ตรวจว่าใครเป็นคนเปลี่ยนสถานะ
	// ถ้าเป็น "Chang" จะไม่ส่งอีเมลแจ้งเตือน
	// -----------------------
	if StatusInvoice != "Chang" {
		var company string
		var companyVendor string
		changedBy := strings.ToUpper(strings.TrimSpace(form.ChangedBySourceSystem))

		// ====== เคส APP_EMPLOYEE (พนักงานเราเปลี่ยน) -> ส่งหา Vendor (+ optional ส่งหา Employee) ======
		if changedBy == "APP_EMPLOYEE" {
			company = dataVersion.VendorCode
			dataVendorMetrics, err := h.EDIInvoiceSrv.GetEDIVendorNotificationRecipientByCompanyService(dataVersion.VendorCode)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to get VendorMetrics",
				})
			}

			var toEmails []string

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
				if err := mailer.SendStatusInvoiceVendorEmail(
					toEmails,
					StatusInvoice,
					company,
					dataVersion.NumberInvoice,
					fileURLStr,
					req.Note,
				); err != nil {
					log.Println("failed to send invoice status email to vendor:", err)
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "failed to send invoice status email to vendor",
					})
				}
			} else {
				log.Println("No valid vendor email recipient found, skip sending vendor email")
			}

		} else if changedBy == "APP_USER" {
			companyVendor = dataVersion.VendorCode
			if err := mailer.SendStatusInvoiceEmployeeEmail(
				companyVendor,
				StatusInvoice,
				dataVersion.NumberInvoice,
				fileURLStr,
				req.Note,
				"invoice",
			); err != nil {
				log.Println("failed to send invoice status email to employee (from vendor change):", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to send invoice status email to employee",
				})
			}
		} else {
			log.Printf("Unknown ChangedBySourceSystem=%q, default to APP_EMPLOYEE behavior\n", changedBy)
			companyVendor = dataVersion.VendorCode
			if err := mailer.SendStatusInvoiceEmployeeEmail(
				companyVendor,
				StatusInvoice,
				dataVersion.NumberInvoice,
				fileURLStr,
				req.Note,
				"invoice",
			); err != nil {
				log.Println("failed to send invoice status email to employee (fallback):", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to send invoice status email to employee",
				})
			}
		}
	} else {
		log.Println("StatusInvoice is 'Chang', skip sending email notification")
	}

	if err := h.EDIInvoiceSrv.UpdateStatusInvoiceService(req.EDIInvoiceID, StatusInvoice); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update status_Invoice",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "EDIInvoiceVersionStatusLog created successfully",
		"file_url": fileURL,
	})
}

func (h *EDIInvoiceHandler) GetInvoiceVersionStatusLogByInvoiceVersionIDHandler(c *fiber.Ctx) error {
	InvoiceVersionID := c.Params("invoice_version_id")

	row, err := h.EDIInvoiceSrv.GetInvoiceVersionStatusLogByInvoiceVersionIDService(InvoiceVersionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve Invoice version status log",
		})
	}

	if row == nil {
		return c.JSON(fiber.Map{
			"data": nil,
		})
	}

	return c.JSON(row)
}

func (h *EDIInvoiceHandler) GetInvoiceVersionStatusLogByInvoiceVersionIDAndApprovedHandler(c *fiber.Ctx) error {
	InvoiceVersionID := c.Params("invoice_version_id")

	row, err := h.EDIInvoiceSrv.GetInvoiceVersionStatusLogByInvoiceVersionIDAndApprovedService(InvoiceVersionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve Invoice version status log",
		})
	}

	if row == nil {
		return c.JSON(fiber.Map{
			"data": nil,
		})
	}

	return c.JSON(row)
}
