package services

import (
	"backend/internal/core/domains"
	"backend/internal/core/models"
	ports "backend/internal/core/ports/repositories"
	servicesports "backend/internal/core/ports/services"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
)

type EDIForecastService struct {
	ediForecastRepo ports.EDIForecastRepository
}

func NewEDIForecastService(ediForecastRepo ports.EDIForecastRepository) servicesports.EDIForecastService {
	return &EDIForecastService{ediForecastRepo: ediForecastRepo}
}

func (s *EDIForecastService) CreateNewForecastWithVersion(ctx context.Context,
	headerIn *models.EDI_ForecastResp,
	versionIn *models.EDI_ForecastVersionResp,
) (*models.EDI_ForecastResp, error) {

	now := time.Now().UTC()

	// ------- เตรียม GUID -------
	newGUID := func() mssql.UniqueIdentifier {
		u := uuid.New()
		var id mssql.UniqueIdentifier
		copy(id[:], u[:16])
		return id
	}

	// ถ้า caller ยังไม่ใส่ id มา ให้ gen ใหม่ที่นี่
	if isZeroGUID(headerIn.EDI_ForecastID) {
		headerIn.EDI_ForecastID = newGUID()
	}
	if isZeroGUID(versionIn.EDIForecastVersionID) {
		versionIn.EDIForecastVersionID = newGUID()
	}

	// ------- map models -> domains -------
	headerDom := &domains.EDI_Forecast{
		EDI_ForecastID:        headerIn.EDI_ForecastID,
		NumberForecast:        headerIn.NumberForecast,
		VendorCode:            headerIn.VendorCode,
		FileURL:               headerIn.FileURL,
		ReadForecast:          false,
		StatusForecast:        "New",
		ActiveVersionID:       nil,
		CreatedByExternalID:   headerIn.CreatedByExternalID,
		CreatedBySourceSystem: headerIn.CreatedBySourceSystem,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	// ผูกเวอร์ชันเข้ากับ header + set เป็นเวอร์ชันแรกเสมอ
	versionDom := &domains.EDI_ForecastVersion{
		EDIForecastVersionID:  versionIn.EDIForecastVersionID,
		EDIForecastID:         headerIn.EDI_ForecastID,
		VersionNo:             1,
		PeriodFrom:            versionIn.PeriodFrom,
		PeriodTo:              versionIn.PeriodTo,
		StatusForecast:        versionIn.StatusForecast,
		ReadForecast:          versionIn.ReadForecast,
		Note:                  versionIn.Note,
		SourceFileURL:         versionIn.SourceFileURL,
		CreatedByExternalID:   versionIn.CreatedByExternalID,
		CreatedBySourceSystem: versionIn.CreatedBySourceSystem,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	// 1) สร้าง Header
	if err := s.ediForecastRepo.CreateEDIForecastRepository(headerDom); err != nil {
		return nil, err
	}

	// 2) สร้าง Version (v1)
	if err := s.ediForecastRepo.CreateEDIForecastVersionRepository(versionDom); err != nil {
		return nil, err
	}

	// 3) อัปเดต active_version_id
	if err := s.ediForecastRepo.UpdateActiveVersion(
		uuidToString(headerDom.EDI_ForecastID),
		uuidToString(versionDom.EDIForecastVersionID),
	); err != nil {
		return nil, err
	}

	// ------- map domains -> models response -------
	out := &models.EDI_ForecastResp{
		EDI_ForecastID:        headerDom.EDI_ForecastID,
		NumberForecast:        headerDom.NumberForecast,
		VendorCode:            headerDom.VendorCode,
		ReadForecast:          true,
		ActiveVersionID:       &versionDom.EDIForecastVersionID,
		CreatedByExternalID:   headerIn.CreatedByExternalID,
		CreatedBySourceSystem: headerIn.CreatedBySourceSystem,
		CreatedAt:             headerDom.CreatedAt,
		UpdatedAt:             headerDom.UpdatedAt,
		RowVer:                headerDom.RowVer,
		Versions: []models.EDI_ForecastVersionResp{
			{
				EDIForecastVersionID:  versionDom.EDIForecastVersionID,
				EDIForecastID:         versionDom.EDIForecastID,
				VersionNo:             versionDom.VersionNo,
				PeriodFrom:            versionDom.PeriodFrom,
				PeriodTo:              versionDom.PeriodTo,
				StatusForecast:        versionDom.StatusForecast,
				ReadForecast:          versionDom.ReadForecast,
				Note:                  versionDom.Note,
				SourceFileURL:         versionDom.SourceFileURL,
				CreatedByExternalID:   versionDom.CreatedByExternalID,
				CreatedBySourceSystem: versionDom.CreatedBySourceSystem,
				CreatedAt:             versionDom.CreatedAt,
				RowVer:                versionDom.RowVer,
			},
		},
	}

	return out, nil
}

func (s *EDIForecastService) GenerateRunningNumberService() (string, error) {
	today := time.Now().Format("2006-01-02")

	last, err := s.ediForecastRepo.GetLastForecastRunningForDate(today)
	if err != nil {
		return "", err
	}

	next := last + 1
	return fmt.Sprintf("FC%s-%04d", today, next), nil
}

func (s *EDIForecastService) MarkForecastAsReadService(id mssql.UniqueIdentifier) (*domains.EDI_Forecast, error) {
	return s.ediForecastRepo.MarkForecastAsRead(id)
}

func (s *EDIForecastService) UpdateStatusForecastService(id mssql.UniqueIdentifier, status string) error {
	return s.ediForecastRepo.UpdateStatusForecast(id, status)
}

func (s *EDIForecastService) CreateEDIForecastVersionService(req models.EDI_ForecastVersionResp) error {
	if s.ediForecastRepo == nil {
		return fmt.Errorf("edi forecast repository is not initialized")
	}

	maxVer, err := s.ediForecastRepo.GetMaxVersionNoByForecastID(req.EDIForecastID)
	if err != nil {
		return fmt.Errorf("failed to get current max version: %w", err)
	}

	u := uuid.New()
	var newID mssql.UniqueIdentifier
	copy(newID[:], u[:])

	domainISR := domains.EDI_ForecastVersion{
		EDIForecastVersionID:  newID,
		EDIForecastID:         req.EDIForecastID,
		VersionNo:             maxVer + 1,
		PeriodFrom:            req.PeriodFrom,
		PeriodTo:              req.PeriodTo,
		StatusForecast:        req.StatusForecast,
		ReadForecast:          req.ReadForecast,
		Note:                  req.Note,
		SourceFileURL:         req.SourceFileURL,
		CreatedByExternalID:   req.CreatedByExternalID,
		CreatedBySourceSystem: req.CreatedBySourceSystem,
	}

	if err := s.ediForecastRepo.CreateEDIForecastVersionRepository(&domainISR); err != nil {
		return fmt.Errorf("failed to create forecast version: %w", err)
	}

	if err := s.ediForecastRepo.UpdateActiveForecastVersion(req.EDIForecastID, newID); err != nil {
		return fmt.Errorf("failed to update active version: %w", err)
	}

	if err := s.ediForecastRepo.UpdateStatusForecast(req.EDIForecastID, "Change"); err != nil {
		return fmt.Errorf("failed to update forecast status: %w", err)
	}

	return nil
}

func (s *EDIForecastService) GetEDIForecastWithActiveTopService(limit int, vendorCode string) ([]models.EDIForecastWithActiveReq, error) {
	rows, err := s.ediForecastRepo.GetEDIForecastWithActiveTop(limit, vendorCode)
	if err != nil {
		return nil, err
	}
	out := make([]models.EDIForecastWithActiveReq, 0, len(rows))
	for _, d := range rows {
		m := models.EDIForecastWithActiveReq{
			// header
			EDIForecastID:   d.EDIForecastID,
			NumberForecast:  d.NumberForecast,
			VendorCode:      d.VendorCode,
			ReadForecast:    d.ReadForecast,
			ActiveVersionID: d.ActiveVersionID,
			StatusForecast:  d.StatusForecast,
			FileURL:         d.FileURL,
			CreatedAt:       d.CreatedAt,
			UpdatedAt:       d.UpdatedAt,

			// active version
			AV_ID:            d.AV_ID,
			AV_VersionNo:     d.AV_VersionNo,
			AV_PeriodFrom:    d.AV_PeriodFrom,
			AV_PeriodTo:      d.AV_PeriodTo,
			AV_Status:        d.AV_Status,
			AV_Read:          d.AV_Read,
			AV_Note:          d.AV_Note,
			AV_SourceFileURL: d.AV_SourceFileURL,
			AV_CreatedAt:     d.AV_CreatedAt,

			// latest status log
			LastStatusLogID: d.LastStatusLogID,
			LastOldStatus:   d.LastOldStatus,
			LastNewStatus:   d.LastNewStatus,
			LastStatusNote:  d.LastStatusNote,
			LastFileURL:     d.LastFileURL,
			LastStatusAt:    d.LastStatusAt,
		}

		out = append(out, m)
	}

	return out, nil
}

func (s *EDIForecastService) GetEDIForecastWithActiveByNumberService(number string) (*models.EDIForecastWithActiveReq, error) {
	row, err := s.ediForecastRepo.GetEDIForecastWithActiveByNumber(number)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("ไม่พบข้อมูล Forecast หมายเลข %s", number)
	}

	out := &models.EDIForecastWithActiveReq{
		// header
		EDIForecastID:   row.EDIForecastID,
		NumberForecast:  row.NumberForecast,
		VendorCode:      row.VendorCode,
		ActiveVersionID: row.ActiveVersionID,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,

		// active version
		AV_ID:            row.AV_ID,
		AV_VersionNo:     row.AV_VersionNo,
		AV_PeriodFrom:    row.AV_PeriodFrom,
		AV_PeriodTo:      row.AV_PeriodTo,
		AV_Status:        row.AV_Status,
		AV_Read:          row.AV_Read,
		AV_Note:          row.AV_Note,
		AV_SourceFileURL: row.AV_SourceFileURL,
		AV_CreatedAt:     row.AV_CreatedAt,
	}

	return out, nil
}

func (s *EDIForecastService) GetEDIForecastDetailByNumber(number string) (*models.EDIForecastDetailResp, error) {
	// fmt.Println("number order : ", number)
	h, err := s.ediForecastRepo.GetForecastHeaderByNumber(number)
	if err != nil {
		return nil, err
	}

	vers, err := s.ediForecastRepo.GetForecastVersionsByForecastID(h.EDI_ForecastID)
	if err != nil {
		return nil, err
	}

	out := &models.EDIForecastDetailResp{
		EDI_ForecastID:  h.EDI_ForecastID,
		NumberForecast:  h.NumberForecast,
		VendorCode:      h.VendorCode,
		StatusForecast:  h.StatusForecast,
		FileURL:         h.FileURL,
		ReadForecast:    h.ReadForecast,
		ActiveVersionID: h.ActiveVersionID,
		CreatedAt:       h.CreatedAt,
		UpdatedAt:       h.UpdatedAt,
		Versions:        make([]models.EDIForecastVersionItemResp, 0, len(vers)),
	}

	for _, v := range vers {
		var createdBy *models.CreatedByPrincipal

		if v.CreatedByPrincipal != nil && v.CreatedByPrincipal.ExternalID != "" {
			createdBy = &models.CreatedByPrincipal{
				ExternalID:   v.CreatedByPrincipal.ExternalID,
				Email:        v.CreatedByPrincipal.Email,
				Display_name: v.CreatedByPrincipal.DisplayName,
				Profile:      v.CreatedByPrincipal.Profile,
				Group:        v.CreatedByPrincipal.Group,
				Role:         v.CreatedByPrincipal.Role,
				SourceSystem: v.CreatedByPrincipal.SourceSystem,
				Status:       v.CreatedByPrincipal.Status,
				Username:     v.CreatedByPrincipal.Username,
			}
		}

		out.Versions = append(out.Versions, models.EDIForecastVersionItemResp{
			EDIForecastVersionID: v.EDIForecastVersionID,
			VersionNo:            v.VersionNo,
			PeriodFrom:           v.PeriodFrom,
			PeriodTo:             v.PeriodTo,
			StatusForecast:       v.StatusForecast,
			ReadForecast:         v.ReadForecast,
			Note:                 v.Note,
			SourceFileURL:        v.SourceFileURL,
			CreatedAt:            v.CreatedAt,
			IsActive:             h.ActiveVersionID != nil && *h.ActiveVersionID == v.EDIForecastVersionID,
			CreatedBy:            createdBy,
		})
	}

	return out, nil
}

func (s *EDIForecastService) GetEDIForecastVersionByIDService(ediForecastVersionID string) (domains.EDI_ForecastVersion, error) {
	return s.ediForecastRepo.GetEDIForecastVersionByID(ediForecastVersionID)
}

func (s *EDIForecastService) GetForecastBasicByIDService(ediForecastID string) (domains.ForecastBasicInfo, error) {
	return s.ediForecastRepo.GetForecastBasicByID(ediForecastID)
}

func (s *EDIForecastService) GetEDIVendorNotificationRecipientByCompanyService(company string) ([]models.EDIVendorNotificationRecipientReq, error) {
	recipients, err := s.ediForecastRepo.GetEDIVendorNotificationRecipientByCompany(company)
	if err != nil {
		return nil, err
	}

	res := make([]models.EDIVendorNotificationRecipientReq, 0, len(recipients))

	for _, r := range recipients {
		var principal *models.EDIPrincipalUserEmailReq
		if r.Principal != nil {
			principal = &models.EDIPrincipalUserEmailReq{
				ExternalID: r.Principal.ExternalID,
				Email:      r.Principal.Email,
			}
		}

		res = append(res, models.EDIVendorNotificationRecipientReq{
			VendorNotificationRecipientID: r.VendorNotificationRecipientID,
			Company:                       r.Company,
			NotificationType:              r.NotificationType,
			Principal:                     principal,
		})
	}

	return res, nil
}

func (s *EDIForecastService) GetStatusSummaryByVendorCodeService(vendorCode string) (*models.ForecastStatusSummaryResp, error) {

	domainResult, err := s.ediForecastRepo.GetStatusSummaryByVendorCode(vendorCode)
	if err != nil {
		return nil, err
	}
	if domainResult == nil {
		return nil, nil
	}
	resp := &models.ForecastStatusSummaryResp{
		VendorCode:    domainResult.VendorCode,
		NewCount:      domainResult.NewCount,
		ConfirmCount:  domainResult.ConfirmCount,
		RejectCount:   domainResult.RejectCount,
		ApprovedCount: domainResult.ApprovedCount,
		TotalCount:    domainResult.TotalCount,
	}

	return resp, nil
}

/* -------------------- helpers -------------------- */

func isZeroGUID(id mssql.UniqueIdentifier) bool {
	var zero mssql.UniqueIdentifier
	return id == zero
}

func uuidToString(id mssql.UniqueIdentifier) string {
	// แปลงเป็น uuid string สำหรับเมธอด repository ที่รับ string
	uu, _ := uuid.FromBytes(id[:])
	return uu.String()
}
