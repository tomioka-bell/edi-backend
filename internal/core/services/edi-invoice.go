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

type EDIInvoiceService struct {
	ediInvoiceRepo ports.EDIInvoiceRepository
}

func NewEDIInvoiceService(ediInvoiceRepo ports.EDIInvoiceRepository) servicesports.EDIInvoiceService {
	return &EDIInvoiceService{ediInvoiceRepo: ediInvoiceRepo}
}

func (s *EDIInvoiceService) GenerateRunningNumberService() (string, error) {
	today := time.Now().Format("2006-01-02")

	last, err := s.ediInvoiceRepo.GetLastInvoiceRunningForDate(today)
	if err != nil {
		return "", err
	}

	next := last + 1
	return fmt.Sprintf("INV%s-%04d", today, next), nil
}

func (s *EDIInvoiceService) CreateNewInvoiceWithVersion(ctx context.Context,
	headerIn *models.EDIInvoiceResp,
	versionIn *models.EDIInvoiceVersionResp,
) (*models.EDIInvoiceResp, error) {

	now := time.Now().UTC()

	newGUID := func() mssql.UniqueIdentifier {
		u := uuid.New()
		var id mssql.UniqueIdentifier
		copy(id[:], u[:16])
		return id
	}

	// ถ้า caller ยังไม่ใส่ id มา ให้ gen ใหม่ที่นี่
	if isZeroGUID(headerIn.EDIInvoiceID) {
		headerIn.EDIInvoiceID = newGUID()
	}
	if isZeroGUID(versionIn.EDIInvoiceVersionID) {
		versionIn.EDIInvoiceVersionID = newGUID()
	}

	// ------- map models -> domains -------
	headerDom := &domains.EDIInvoice{
		EDIInvoiceID:          headerIn.EDIInvoiceID,
		NumberInvoice:         headerIn.NumberInvoice,
		NumberOrder:           headerIn.NumberOrder,
		InvoiceType:           headerIn.InvoiceType,
		FileURL:               headerIn.FileURL,
		VendorCode:            headerIn.VendorCode,
		ReadInvoice:           false,
		StatusInvoice:         "New",
		ActiveVersionID:       nil,
		CreatedByExternalID:   headerIn.CreatedByExternalID,
		CreatedBySourceSystem: headerIn.CreatedBySourceSystem,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	// ผูกเวอร์ชันเข้ากับ header + set เป็นเวอร์ชันแรกเสมอ
	versionDom := &domains.EDIInvoiceVersion{
		EDIInvoiceVersionID:   versionIn.EDIInvoiceVersionID,
		EDIInvoiceID:          headerIn.EDIInvoiceID,
		VersionNo:             1,
		PeriodFrom:            versionIn.PeriodFrom,
		PeriodTo:              versionIn.PeriodTo,
		StatusInvoice:         versionIn.StatusInvoice,
		Note:                  versionIn.Note,
		SourceFileURL:         versionIn.SourceFileURL,
		CreatedByExternalID:   versionIn.CreatedByExternalID,
		CreatedBySourceSystem: versionIn.CreatedBySourceSystem,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	// 1) สร้าง Header
	if err := s.ediInvoiceRepo.CreateEDIInvoiceRepository(headerDom); err != nil {
		return nil, err
	}

	// 2) สร้าง Version (v1)
	if err := s.ediInvoiceRepo.CreateEDIInvoiceVersionRepository(versionDom); err != nil {
		return nil, err
	}

	// 3) อัปเดต active_version_id
	if err := s.ediInvoiceRepo.UpdateActiveInvoiceVersion(
		uuidToString(headerDom.EDIInvoiceID),
		uuidToString(versionDom.EDIInvoiceVersionID),
	); err != nil {
		return nil, err
	}

	// ------- map domains -> models response -------
	out := &models.EDIInvoiceResp{
		EDIInvoiceID:          headerDom.EDIInvoiceID,
		NumberInvoice:         headerDom.NumberInvoice,
		NumberOrder:           headerDom.NumberOrder,
		InvoiceType:           headerDom.InvoiceType,
		VendorCode:            headerDom.VendorCode,
		ReadInvoice:           true,
		ActiveVersionID:       &versionDom.EDIInvoiceVersionID,
		CreatedByExternalID:   headerIn.CreatedByExternalID,
		CreatedBySourceSystem: headerIn.CreatedBySourceSystem,
		CreatedAt:             headerDom.CreatedAt,
		UpdatedAt:             headerDom.UpdatedAt,
		RowVer:                headerDom.RowVer,
		Versions: []models.EDIInvoiceVersionResp{
			{
				EDIInvoiceVersionID:   versionDom.EDIInvoiceVersionID,
				EDIInvoiceID:          versionDom.EDIInvoiceID,
				VersionNo:             versionDom.VersionNo,
				PeriodFrom:            versionDom.PeriodFrom,
				PeriodTo:              versionDom.PeriodTo,
				StatusInvoice:         versionDom.StatusInvoice,
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

func (s *EDIInvoiceService) CreateEDIInvoiceVersionService(req models.EDIInvoiceVersionResp) error {
	if s.ediInvoiceRepo == nil {
		return fmt.Errorf("edi Invoice repository is not initialized")
	}

	maxVer, err := s.ediInvoiceRepo.GetMaxVersionNoByInvoiceID(req.EDIInvoiceID)
	if err != nil {
		return fmt.Errorf("failed to get current max version: %w", err)
	}

	u := uuid.New()
	var newID mssql.UniqueIdentifier
	copy(newID[:], u[:])

	domainISR := domains.EDIInvoiceVersion{
		EDIInvoiceVersionID:   newID,
		EDIInvoiceID:          req.EDIInvoiceID,
		VersionNo:             maxVer + 1,
		PeriodFrom:            req.PeriodFrom,
		PeriodTo:              req.PeriodTo,
		StatusInvoice:         req.StatusInvoice,
		Note:                  req.Note,
		SourceFileURL:         req.SourceFileURL,
		CreatedByExternalID:   req.CreatedByExternalID,
		CreatedBySourceSystem: req.CreatedBySourceSystem,
	}

	if err := s.ediInvoiceRepo.CreateEDIInvoiceVersionRepository(&domainISR); err != nil {
		return fmt.Errorf("failed to create Invoice version: %w", err)
	}
	if err := s.ediInvoiceRepo.UpdateActiveInvoiceVersion(req.EDIInvoiceID.String(), newID.String()); err != nil {
		return fmt.Errorf("failed to update active version: %w", err)
	}

	if err := s.ediInvoiceRepo.UpdateStatusInvoice(req.EDIInvoiceID, "Change"); err != nil {
		return fmt.Errorf("failed to update Invoice status: %w", err)
	}

	return nil
}

func (s *EDIInvoiceService) MarkInvoiceAsReadService(id mssql.UniqueIdentifier) error {
	return s.ediInvoiceRepo.MarkInvoiceAsRead(id)
}

func (s *EDIInvoiceService) UpdateStatusInvoiceService(id mssql.UniqueIdentifier, status string) error {
	return s.ediInvoiceRepo.UpdateStatusInvoice(id, status)
}

func (s *EDIInvoiceService) GetEDIInvoiceWithActiveTopService(limit int, vendorCode string) ([]models.EDIInvoiceWithActiveReq, error) {
	rows, err := s.ediInvoiceRepo.GetEDIInvoiceWithActiveTop(limit, vendorCode)
	if err != nil {
		return nil, err
	}
	out := make([]models.EDIInvoiceWithActiveReq, 0, len(rows))
	for _, d := range rows {
		m := models.EDIInvoiceWithActiveReq{
			// header
			EDIInvoiceID:    d.EDIInvoiceID,
			NumberInvoice:   d.NumberInvoice,
			NumberOrder:     d.NumberOrder,
			ReadInvoice:     d.ReadInvoice,
			ActiveVersionID: d.ActiveVersionID,
			StatusInvoice:   d.StatusInvoice,
			VendorCode:      d.VendorCode,
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

func (s *EDIInvoiceService) GetEDIInvoiceDetailByNumber(number string) (*models.EDIInvoiceDetailResp, error) {
	h, err := s.ediInvoiceRepo.GetInvoiceHeaderByNumber(number)
	if err != nil {
		return nil, err
	}

	vers, err := s.ediInvoiceRepo.GetInvoiceVersionsByInvoiceID(h.EDIInvoiceID)
	if err != nil {
		return nil, err
	}

	out := &models.EDIInvoiceDetailResp{
		EDIInvoiceID:  h.EDIInvoiceID,
		NumberInvoice: h.NumberInvoice,
		ReadInvoice:   h.ReadInvoice,
		NumberOrder:   h.NumberOrder,
		InvoiceType:   h.InvoiceType,
		StatusInvoice: h.StatusInvoice,
		VanderCode:    h.VendorCode,
		FileURL:       h.FileURL,

		ActiveVersionID: h.ActiveVersionID,
		CreatedAt:       h.CreatedAt,
		UpdatedAt:       h.UpdatedAt,
		Versions:        make([]models.EDIInvoiceVersionItemResp, 0, len(vers)),
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

		out.Versions = append(out.Versions, models.EDIInvoiceVersionItemResp{
			EDIInvoiceVersionID: v.EDIInvoiceVersionID,
			VersionNo:           v.VersionNo,
			PeriodFrom:          v.PeriodFrom,
			PeriodTo:            v.PeriodTo,
			StatusInvoice:       v.StatusInvoice,
			Note:                v.Note,
			SourceFileURL:       v.SourceFileURL,
			CreatedAt:           v.CreatedAt,
			IsActive:            h.ActiveVersionID != nil && *h.ActiveVersionID == v.EDIInvoiceVersionID,
			CreatedBy:           createdBy,
		})
	}

	return out, nil
}

func (s *EDIInvoiceService) GetInvoiceDetailByNumberOrderService(number string) (*models.EDIInvoiceDetailResp, error) {
	h, err := s.ediInvoiceRepo.GetInvoiceByNumberOrder(number)
	if err != nil {
		return nil, err
	}

	vers, err := s.ediInvoiceRepo.GetInvoiceVersionsByInvoiceID(h.EDIInvoiceID)
	if err != nil {
		return nil, err
	}

	out := &models.EDIInvoiceDetailResp{
		EDIInvoiceID:    h.EDIInvoiceID,
		NumberInvoice:   h.NumberInvoice,
		NumberOrder:     h.NumberOrder,
		InvoiceType:     h.InvoiceType,
		StatusInvoice:   h.StatusInvoice,
		VanderCode:      h.VendorCode,
		FileURL:         h.FileURL,
		ActiveVersionID: h.ActiveVersionID,
		CreatedAt:       h.CreatedAt,
		UpdatedAt:       h.UpdatedAt,
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

		out.Versions = append(out.Versions, models.EDIInvoiceVersionItemResp{
			EDIInvoiceVersionID: v.EDIInvoiceVersionID,
			VersionNo:           v.VersionNo,
			PeriodFrom:          v.PeriodFrom,
			PeriodTo:            v.PeriodTo,
			StatusInvoice:       v.StatusInvoice,
			Note:                v.Note,
			SourceFileURL:       v.SourceFileURL,
			CreatedAt:           v.CreatedAt,
			IsActive:            h.ActiveVersionID != nil && *h.ActiveVersionID == v.EDIInvoiceVersionID,
			CreatedBy:           createdBy,
		})
	}

	return out, nil
}

func (s *EDIInvoiceService) GetStatusInvoiceSummaryByVendorCodeService(vendorCode string) (*models.StatusInvoiceSummaryResp, error) {

	domainResult, err := s.ediInvoiceRepo.GetStatusInvoiceSummaryByVendorCode(vendorCode)
	if err != nil {
		return nil, err
	}
	if domainResult == nil {
		return nil, nil
	}

	resp := &models.StatusInvoiceSummaryResp{
		VendorCode:    domainResult.VendorCode,
		NewCount:      domainResult.NewCount,
		ConfirmCount:  domainResult.ConfirmCount,
		RejectCount:   domainResult.RejectCount,
		ApprovedCount: domainResult.ApprovedCount,
		TotalCount:    domainResult.TotalCount,
	}

	return resp, nil
}

func (s *EDIInvoiceService) GetEDIInvoiceVersionByIDService(ediInvoiceVersionID string) (domains.EDIInvoiceVersion, error) {
	return s.ediInvoiceRepo.GetEDIInvoiceVersionByID(ediInvoiceVersionID)
}

func (s *EDIInvoiceService) GetInvoiceBasicByIDService(ediInvoiceID string) (domains.InvoiceBasicInfo, error) {
	return s.ediInvoiceRepo.GetInvoiceBasicByID(ediInvoiceID)
}

func (s *EDIInvoiceService) GetEDIVendorNotificationRecipientByCompanyService(company string) ([]models.EDIVendorNotificationRecipientReq, error) {
	recipients, err := s.ediInvoiceRepo.GetEDIVendorNotificationRecipientByCompany(company)
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
