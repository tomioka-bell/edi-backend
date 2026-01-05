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

type EDIForecastHandler struct {
	EDIForecastSrv services.EDIForecastService
	EDIReadStatSrv services.EDIReadStatService
}

func NewEDIForecastHandler(srv services.EDIForecastService, ediReadStatSrv services.EDIReadStatService) *EDIForecastHandler {
	return &EDIForecastHandler{EDIForecastSrv: srv, EDIReadStatSrv: ediReadStatSrv}
}

func (h *EDIForecastHandler) MarkForecastAsRead(c *fiber.Ctx) error {
	idStr := c.Params("edi_forecast_id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}

	var mssqlID mssql.UniqueIdentifier
	copy(mssqlID[:], uid[:])

	f, err := h.EDIForecastSrv.MarkForecastAsReadService(mssqlID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update read_forecast",
		})
	}

	// now := time.Now().UTC()

	// readReq := models.EDIReadStatReq{
	// 	Number:     f.NumberForecast,
	// 	Type:       "FORECAST",
	// 	VendorCode: f.VendorCode,
	// 	Read:       true,
	// 	ReadAt:     now,
	// 	CreatedAt:  now,
	// }

	// if err := h.EDIReadStatSrv.TrackReadService(readReq); err != nil {
	// 	log.Println("failed to track forecast read stat:", err)
	// }

	return c.JSON(fiber.Map{
		"message":         "forecast marked as read",
		"vendor_code":     f.VendorCode,
		"number_forecast": f.NumberForecast,
	})
}

func (h *EDIForecastHandler) CreateNewForecastWithVersionHandler(c *fiber.Ctx) error {
	var req models.CreateForecastReq

	// -----------------------
	// อ่าน JSON จาก form field: "payload"
	// -----------------------
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

	// -----------------------
	// อัปโหลดไฟล์ (ถ้ามี)
	// -----------------------
	fileURL, err := utils.UploadFileFromForm(
		c,
		"file",
		"uploads/document/forecast",
		"/uploads/document/forecast",
	)

	if err != nil {
		log.Println("Upload file error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to upload file",
		})
	}

	trim := func(s string) string { return strings.TrimSpace(s) }
	toUpper := func(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }

	// -----------------------
	// Header
	// -----------------------
	var numberForecast string
	if strings.TrimSpace(req.NumberForecast) == "" {
		gen, err := h.EDIForecastSrv.GenerateRunningNumberService()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate running number",
			})
		}
		numberForecast = gen
	} else {
		numberForecast = req.NumberForecast
	}

	header := models.EDI_ForecastResp{
		NumberForecast:        numberForecast,
		VendorCode:            req.VendorCode,
		FileURL:               fileURL,
		CreatedByExternalID:   trim(req.CreatedByExternalID),
		CreatedBySourceSystem: toUpper(req.CreatedBySourceSystem),
	}

	// -----------------------
	// Version แรก (req.Versions[0])
	// -----------------------
	if len(req.Versions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "at least one version is required",
		})
	}

	v0 := req.Versions[0]
	firstVer := models.EDI_ForecastVersionResp{
		PeriodFrom:            v0.PeriodFrom,
		PeriodTo:              v0.PeriodTo,
		StatusForecast:        v0.StatusForecast,
		ReadForecast:          v0.ReadForecast,
		Note:                  v0.Note,
		SourceFileURL:         v0.SourceFileURL,
		CreatedByExternalID:   trim(v0.CreatedByExternalID),
		CreatedBySourceSystem: toUpper(v0.CreatedBySourceSystem),
	}

	if fileURL != nil && firstVer.SourceFileURL == nil {
		firstVer.SourceFileURL = fileURL
	}

	log.Printf("[DEBUG] header FK: ext=%q src=%q", header.CreatedByExternalID, header.CreatedBySourceSystem)

	created, err := h.EDIForecastSrv.CreateNewForecastWithVersion(c.Context(), &header, &firstVer)
	if err != nil {
		log.Println("CreateNewForecastWithVersion error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create forecast",
		})
	}

	var fileURLStr string
	if fileURL != nil {
		fileURLStr = strings.TrimSpace(*fileURL)
	}

	changedBy := strings.ToUpper(strings.TrimSpace(firstVer.CreatedBySourceSystem))

	fmt.Printf("ข้อมูลที่ส่งมา : %q\n", changedBy)

	dataVendorMetrics, err := h.EDIForecastSrv.GetEDIVendorNotificationRecipientByCompanyService(header.VendorCode)
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
		if err := mailer.SendStatusForecastVendorEmail(
			toEmails,
			v0.StatusForecast,
			company,
			created.NumberForecast,
			fileURLStr,
			v0.Note,
		); err != nil {
			log.Println("failed to send forecast status email to vendor:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to send forecast status email to vendor",
			})
		}
	} else {
		log.Println("No valid vendor email recipient found, skip sending vendor email")
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *EDIForecastHandler) GetEDIForecastWithActiveTopHandler(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	vendorCode := c.Query("vendor_code")
	rows, err := h.EDIForecastSrv.GetEDIForecastWithActiveTopService(limit, vendorCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve forecasts"})
	}
	return c.JSON(rows)
}

func (h *EDIForecastHandler) GetEDIForecastWithActiveByNumberHandler(c *fiber.Ctx) error {
	number := c.Query("number_forecast")
	row, err := h.EDIForecastSrv.GetEDIForecastDetailByNumber(number)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve forecast"})
	}
	return c.JSON(row)
}

func (h *EDIForecastHandler) GetStatusSummaryByVendorCodeHandler(c *fiber.Ctx) error {
	number := c.Query("vendor_code")
	row, err := h.EDIForecastSrv.GetStatusSummaryByVendorCodeService(number)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve forecast"})
	}
	return c.JSON(row)
}

func (h *EDIForecastHandler) GetEDIForecastVersionByIDHandler(c *fiber.Ctx) error {
	versionID := c.Params("version_id")
	version, err := h.EDIForecastSrv.GetEDIForecastVersionByIDService(versionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve forecast version"})
	}
	return c.JSON(version)
}

// func (h *EDIForecastHandler) CreateEDIForecastVersionHandler(c *fiber.Ctx) error {
// 	var req models.EDI_ForecastVersionResp

// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request payload",
// 		})
// 	}

// 	if h.EDIForecastSrv == nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Service is not available",
// 		})
// 	}

// 	err := h.EDIForecastSrv.CreateEDIForecastVersionService(req)
// 	if err != nil {
// 		log.Println("Error creating  EDIForecastVersion:", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to create  EDIForecastVersion",
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "EDIForecastVersion created successfully",
// 	})
// }

func (h *EDIForecastHandler) CreateEDIForecastVersionHandler(c *fiber.Ctx) error {

	// 1) Read JSON from FormData
	payloadStr := c.FormValue("payload")
	if payloadStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payload is required",
		})
	}

	var req models.EDI_ForecastVersionResp
	if err := json.Unmarshal([]byte(payloadStr), &req); err != nil {
		log.Println("Invalid JSON:", err)
		log.Println("Original payload:", payloadStr)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload JSON",
		})
	}

	log.Println("===== Received Forecast Version Payload =====")
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
			"uploads/document/forecast",
			"/uploads/document/forecast",
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

	basicInfo, err := h.EDIForecastSrv.GetForecastBasicByIDService(req.EDIForecastID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get basic info",
		})
	}

	dataVendorMetrics, err := h.EDIForecastSrv.GetEDIVendorNotificationRecipientByCompanyService(basicInfo.VendorCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get VendorMetrics",
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
		if err := mailer.SendModifyForecastVendorEmail(
			toEmails,
			company,
			basicInfo.NumberForecast,
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

	if h.EDIForecastSrv == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Service is not available",
		})
	}

	ediForecastVersionID, err := h.EDIForecastSrv.CreateEDIForecastVersionService(req)
	if err != nil {
		log.Println("Error creating EDIForecastVersion:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create EDIForecastVersion",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":                 "EDIForecastVersion created successfully",
		"edi_forecast_version_id": ediForecastVersionID,
	})
}

type UpdateForecastStatusReq struct {
	StatusForecast string `json:"status_forecast"`
}

func (h *EDIForecastHandler) UpdateStatusForecast(c *fiber.Ctx) error {
	req := UpdateForecastStatusReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	idStr := c.Params("edi_forecast_id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}

	var mssqlID mssql.UniqueIdentifier
	copy(mssqlID[:], uid[:])

	if err := h.EDIForecastSrv.UpdateStatusForecastService(mssqlID, req.StatusForecast); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to update status_forecast",
		})
	}

	return c.JSON(fiber.Map{
		"message": "forecast status updated",
	})
}
