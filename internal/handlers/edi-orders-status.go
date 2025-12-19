package handlers

import (
	"backend/internal/core/models"
	"backend/internal/pkgs/mailer"
	"backend/internal/pkgs/utils"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *EDIOrderHandler) CreateEDIOrderVersionStatusLogHandler(c *fiber.Ctx) error {
	var form models.EDIOrderVersionStatusLogForm
	var StatusOrder string

	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form payload",
		})
	}

	if h.EDIOrderSrv == nil {
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
		"uploads/document/Order",
		"/uploads/document/Order",
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
	case "PENDING APPROVED":
		StatusOrder = "Cont, Pending Approved"
	case "FULLY CONFIRMED":
		StatusOrder = "Fully Confirmed"
	case "APPROVED":
		StatusOrder = "Confirmed"
	default:
		StatusOrder = form.NewStatus
	}

	// -----------------------
	// แปลง string -> mssql.UniqueIdentifier
	// -----------------------
	ediOrderVersionID, err := ToMSSQLUUID(form.EDIOrderVersionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid edi_Order_version_id",
		})
	}

	ediOrderID, err := ToMSSQLUUID(form.EDIOrderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid edi_Order_id",
		})
	}

	// -----------------------
	// ประกอบ struct สำหรับ service
	// -----------------------
	var req models.EDIOrderVersionStatusLogResp
	req.EDIOrderVersionID = ediOrderVersionID
	req.EDIOrderID = ediOrderID
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
	if err := h.EDIOrderSrv.CreateEDIOrderVersionStatusLogService(req); err != nil {
		log.Println("Error creating EDIOrderVersionStatusLog:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create EDIOrderVersionStatusLog",
		})
	}

	// ดึงข้อมูล version ล่าสุด
	dataVersion, err := h.EDIOrderSrv.GetEDIOrderVersionByIDService(req.EDIOrderVersionID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get EDIOrderVersion",
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
	// -----------------------
	var company string
	var companyVendor string
	changedBy := strings.ToUpper(strings.TrimSpace(form.ChangedBySourceSystem))

	// ====== เคส APP_EMPLOYEE (พนักงานเราเปลี่ยน) -> ส่งหา Vendor (+ optional ส่งหา Employee) ======
	if changedBy == "APP_EMPLOYEE" {
		company = dataVersion.VendorCode
		dataVendorMetrics, err := h.EDIOrderSrv.GetEDIVendorNotificationRecipientByCompanyService(dataVersion.VendorCode)
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
			if err := mailer.SendStatusOrderVendorEmail(
				toEmails,
				StatusOrder,
				company,
				dataVersion.NumberOrder,
				fileURLStr,
				req.Note,
			); err != nil {
				log.Println("failed to send forecast status email to vendor:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to send forecast status email to vendor",
				})
			}
		} else {
			log.Println("No valid vendor email recipient found, skip sending vendor email")
		}

	} else if changedBy == "APP_USER" {
		companyVendor = dataVersion.VendorCode
		if err := mailer.SendStatusOrderEmployeeEmail(
			companyVendor,
			StatusOrder,
			dataVersion.NumberOrder,
			fileURLStr,
			req.Note,
			"order",
		); err != nil {
			log.Println("failed to send order status email to employee (from vendor change):", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to send order status email to employee",
			})
		}
	} else {
		log.Printf("Unknown ChangedBySourceSystem=%q, default to APP_EMPLOYEE behavior\n", changedBy)
		companyVendor = dataVersion.VendorCode
		if err := mailer.SendStatusOrderEmployeeEmail(
			companyVendor,
			StatusOrder,
			dataVersion.NumberOrder,
			fileURLStr,
			req.Note,
			"order",
		); err != nil {
			log.Println("failed to send order status email to employee (fallback):", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to send order status email to employee",
			})
		}
	}

	if err := h.EDIOrderSrv.UpdateStatusOrderService(req.EDIOrderID, StatusOrder); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update status_order",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "EDIOrderVersionStatusLog created successfully",
		"file_url": fileURL,
	})
}

func (h *EDIOrderHandler) GetOrderVersionStatusLogByOrderVersionIDHandler(c *fiber.Ctx) error {
	OrderVersionID := c.Params("order_version_id")

	row, err := h.EDIOrderSrv.GetOrderVersionStatusLogByOrderVersionIDService(OrderVersionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve Order version status log",
		})
	}

	if row == nil {
		return c.JSON(fiber.Map{
			"data": nil,
		})
	}

	return c.JSON(row)
}

func (h *EDIOrderHandler) GetOrderVersionStatusLogByOrderVersionIDAndApprovedHandler(c *fiber.Ctx) error {
	OrderVersionID := c.Params("order_version_id")

	row, err := h.EDIOrderSrv.GetOrderVersionStatusLogByOrderVersionIDAndApprovedService(OrderVersionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve Order version status log",
		})
	}

	if row == nil {
		return c.JSON(fiber.Map{
			"data": nil,
		})
	}

	return c.JSON(row)
}
