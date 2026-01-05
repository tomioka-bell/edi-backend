package handlers

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"

	"backend/internal/core/models"
	"backend/internal/pkgs/mailer"
	"backend/internal/pkgs/utils"
)

func ToMSSQLUUID(idStr string) (mssql.UniqueIdentifier, error) {
	var uid mssql.UniqueIdentifier

	u, err := uuid.Parse(idStr)
	if err != nil {
		return uid, err
	}

	copy(uid[:], u[:])
	return uid, nil
}

func (h *EDIForecastHandler) CreateEDIForecastVersionStatusLogHandler(c *fiber.Ctx) error {
	var form models.EDI_ForecastVersionStatusLogForm
	var StatusForecast string

	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form payload",
		})
	}

	if h.EDIForecastSrv == nil {
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
		"uploads/document/forecast",
		"/uploads/document/forecast",
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
		StatusForecast = "Rejected"
	case "FULLY_CONFIRMED":
		StatusForecast = "Confirmed"
	case "APPROVED":
		StatusForecast = "Approved"
	default:
		StatusForecast = form.NewStatus
	}

	// -----------------------
	// แปลง string -> mssql.UniqueIdentifier
	// -----------------------
	ediForecastVersionID, err := ToMSSQLUUID(form.EDIForecastVersionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid edi_forecast_version_id",
		})
	}

	ediForecastID, err := ToMSSQLUUID(form.EDIForecastID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid edi_forecast_id",
		})
	}

	// -----------------------
	// ประกอบ struct สำหรับ service
	// -----------------------
	var req models.EDI_ForecastVersionStatusLogResp
	req.EDIForecastVersionID = ediForecastVersionID
	req.EDIForecastID = ediForecastID
	req.OldStatus = form.OldStatus
	req.NewStatus = StatusForecast
	req.Note = form.Note
	req.ChangedByExternalID = form.ChangedByExternalID
	req.ChangedBySourceSystem = form.ChangedBySourceSystem
	req.FileURL = fileURL
	req.CreatedAt = time.Now()

	// -----------------------
	// Service call: สร้าง status log
	// -----------------------
	if err := h.EDIForecastSrv.CreateEDIForecastVersionStatusLogService(req); err != nil {
		log.Println("Error creating EDIForecastVersionStatusLog:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create EDIForecastVersionStatusLog",
		})
	}

	// ดึงข้อมูล version ล่าสุด
	dataVersion, err := h.EDIForecastSrv.GetEDIForecastVersionByIDService(req.EDIForecastVersionID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get EDIForecastVersion",
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
	if StatusForecast != "Chang" {
		var company string
		var companyVendor string
		changedBy := strings.ToUpper(strings.TrimSpace(form.ChangedBySourceSystem))

		// ====== เคส APP_EMPLOYEE (พนักงานเราเปลี่ยน) -> ส่งหา Vendor (+ optional ส่งหา Employee) ======
		if changedBy == "APP_EMPLOYEE" {
			company = dataVersion.VendorCode
			dataVendorMetrics, err := h.EDIForecastSrv.GetEDIVendorNotificationRecipientByCompanyService(dataVersion.VendorCode)
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
				if err := mailer.SendStatusForecastVendorEmail(
					toEmails,
					StatusForecast,
					company,
					dataVersion.NumberForecast,
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
			if err := mailer.SendStatusForecastEmployeeEmail(
				companyVendor,
				StatusForecast,
				dataVersion.NumberForecast,
				fileURLStr,
				req.Note,
				"forecast",
			); err != nil {
				log.Println("failed to send forecast status email to employee (from vendor change):", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to send forecast status email to employee",
				})
			}
		} else {
			log.Printf("Unknown ChangedBySourceSystem=%q, default to APP_EMPLOYEE behavior\n", changedBy)
			companyVendor = dataVersion.VendorCode
			if err := mailer.SendStatusForecastEmployeeEmail(
				companyVendor,
				StatusForecast,
				dataVersion.NumberForecast,
				fileURLStr,
				req.Note,
				"forecast",
			); err != nil {
				log.Println("failed to send forecast status email to employee (fallback):", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to send forecast status email to employee",
				})
			}
		}
	} else {
		log.Println("StatusForecast is 'Chang', skip sending email notification")
	}

	if err := h.EDIForecastSrv.UpdateStatusForecastService(req.EDIForecastID, StatusForecast); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update status_forecast",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "EDIForecastVersionStatusLog created successfully",
		"file_url": fileURL,
	})
}

func (h *EDIForecastHandler) GetForecastVersionStatusLogByForecastVersionIDHandler(c *fiber.Ctx) error {
	forecastVersionID := c.Params("forecast_version_id")

	row, err := h.EDIForecastSrv.GetForecastVersionStatusLogByForecastVersionIDService(forecastVersionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve forecast version status log",
		})
	}

	if row == nil {
		return c.JSON(fiber.Map{
			"data": nil,
		})
	}

	return c.JSON(row)
}

func (h *EDIForecastHandler) GetForecastVersionStatusLogByForecastVersionIDAndApprovedHandler(c *fiber.Ctx) error {
	forecastVersionID := c.Params("forecast_version_id")

	row, err := h.EDIForecastSrv.GetForecastVersionStatusLogByForecastVersionIDAndApprovedService(forecastVersionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve forecast version status log",
		})
	}

	if row == nil {
		return c.JSON(fiber.Map{
			"data": nil,
		})
	}

	return c.JSON(row)
}
