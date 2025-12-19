package services

import (
	"backend/internal/core/domains"
	"backend/internal/core/models"
	ports "backend/internal/core/ports/repositories"
	servicesports "backend/internal/core/ports/services"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

type EDIOrderService struct {
	ediOrderRepo ports.EDIOrderRepository
}

func NewEDIOrderService(ediOrderRepo ports.EDIOrderRepository) servicesports.EDIOrderService {

	return &EDIOrderService{ediOrderRepo: ediOrderRepo}
}

func (s *EDIOrderService) CreateNewOrderWithVersion(ctx context.Context,
	headerIn *models.EDIOrderResp,
	versionIn *models.EDIOrderVersionResp,
) (*models.EDIOrderResp, error) {

	now := time.Now().UTC()

	newGUID := func() mssql.UniqueIdentifier {
		u := uuid.New()
		var id mssql.UniqueIdentifier
		copy(id[:], u[:16])
		return id
	}

	// ถ้า caller ยังไม่ใส่ id มา ให้ gen ใหม่ที่นี่
	if isZeroGUID(headerIn.EDIOrderID) {
		headerIn.EDIOrderID = newGUID()
	}
	if isZeroGUID(versionIn.EDIOrderVersionID) {
		versionIn.EDIOrderVersionID = newGUID()
	}

	// ------- map models -> domains -------
	headerDom := &domains.EDIOrder{
		EDIOrderID:            headerIn.EDIOrderID,
		NumberOrder:           headerIn.NumberOrder,
		VendorCode:            headerIn.VendorCode,
		NumberForecast:        headerIn.NumberForecast,
		FileURL:               headerIn.FileURL,
		ReadOrder:             false,
		StatusOrder:           "New",
		ActiveVersionID:       nil,
		CreatedByExternalID:   headerIn.CreatedByExternalID,
		CreatedBySourceSystem: headerIn.CreatedBySourceSystem,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	// ผูกเวอร์ชันเข้ากับ header + set เป็นเวอร์ชันแรกเสมอ
	versionDom := &domains.EDIOrderVersion{
		EDIOrderVersionID:     versionIn.EDIOrderVersionID,
		EDIOrderID:            headerIn.EDIOrderID,
		VersionNo:             1,
		PeriodFrom:            versionIn.PeriodFrom,
		PeriodTo:              versionIn.PeriodTo,
		StatusOrder:           versionIn.StatusOrder,
		Note:                  versionIn.Note,
		SourceFileURL:         versionIn.SourceFileURL,
		CreatedByExternalID:   versionIn.CreatedByExternalID,
		CreatedBySourceSystem: versionIn.CreatedBySourceSystem,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	// 1) สร้าง Header
	if err := s.ediOrderRepo.CreateEDIOrderRepository(headerDom); err != nil {
		return nil, err
	}

	// 2) สร้าง Version (v1)
	if err := s.ediOrderRepo.CreateEDIOrderVersionRepository(versionDom); err != nil {
		return nil, err
	}

	// 3) อัปเดต active_version_id
	if err := s.ediOrderRepo.UpdateActiveOrderVersion(
		uuidToString(headerDom.EDIOrderID),
		uuidToString(versionDom.EDIOrderVersionID),
	); err != nil {
		return nil, err
	}

	// ------- map domains -> models response -------
	out := &models.EDIOrderResp{
		EDIOrderID:            headerDom.EDIOrderID,
		NumberOrder:           headerDom.NumberOrder,
		VendorCode:            headerDom.VendorCode,
		NumberForecast:        headerDom.NumberForecast,
		ReadOrder:             true,
		ActiveVersionID:       &versionDom.EDIOrderVersionID,
		CreatedByExternalID:   headerIn.CreatedByExternalID,
		CreatedBySourceSystem: headerIn.CreatedBySourceSystem,
		CreatedAt:             headerDom.CreatedAt,
		UpdatedAt:             headerDom.UpdatedAt,
		RowVer:                headerDom.RowVer,
		Versions: []models.EDIOrderVersionResp{
			{
				EDIOrderVersionID:     versionDom.EDIOrderVersionID,
				EDIOrderID:            versionDom.EDIOrderID,
				VersionNo:             versionDom.VersionNo,
				PeriodFrom:            versionDom.PeriodFrom,
				PeriodTo:              versionDom.PeriodTo,
				StatusOrder:           versionDom.StatusOrder,
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

func (s *EDIOrderService) GenerateRunningNumberService() (string, error) {
	today := time.Now().Format("2006-01-02")

	last, err := s.ediOrderRepo.GetLastOrderRunningForDate(today)
	if err != nil {
		return "", err
	}

	next := last + 1
	return fmt.Sprintf("ORD%s-%04d", today, next), nil
}

func (s *EDIOrderService) CreateEDIOrderVersionService(req models.EDIOrderVersionResp) error {
	if s.ediOrderRepo == nil {
		return fmt.Errorf("edi order repository is not initialized")
	}

	maxVer, err := s.ediOrderRepo.GetMaxVersionNoByOrderID(req.EDIOrderID)
	if err != nil {
		return fmt.Errorf("failed to get current max version: %w", err)
	}

	u := uuid.New()
	var newID mssql.UniqueIdentifier
	copy(newID[:], u[:])

	domainISR := domains.EDIOrderVersion{
		EDIOrderVersionID:     newID,
		EDIOrderID:            req.EDIOrderID,
		VersionNo:             maxVer + 1,
		PeriodFrom:            req.PeriodFrom,
		PeriodTo:              req.PeriodTo,
		StatusOrder:           req.StatusOrder,
		Note:                  req.Note,
		SourceFileURL:         req.SourceFileURL,
		CreatedByExternalID:   req.CreatedByExternalID,
		CreatedBySourceSystem: req.CreatedBySourceSystem,
	}

	if err := s.ediOrderRepo.CreateEDIOrderVersionRepository(&domainISR); err != nil {
		return fmt.Errorf("failed to create order version: %w", err)
	}
	if err := s.ediOrderRepo.UpdateActiveOrderVersion(req.EDIOrderID.String(), newID.String()); err != nil {
		return fmt.Errorf("failed to update active version: %w", err)
	}

	if err := s.ediOrderRepo.UpdateStatusOrder(req.EDIOrderID, "Change"); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (s *EDIOrderService) MarkOrderAsReadService(id mssql.UniqueIdentifier) error {
	return s.ediOrderRepo.MarkOrderAsRead(id)
}

func (s *EDIOrderService) UpdateStatusOrderService(id mssql.UniqueIdentifier, status string) error {
	return s.ediOrderRepo.UpdateStatusOrder(id, status)
}

func (s *EDIOrderService) GetEDIOrderWithActiveTopService(limit int, vendorCode string) ([]models.EDIOrderWithActiveReq, error) {
	rows, err := s.ediOrderRepo.GetEDIOrderWithActiveTop(limit, vendorCode)
	if err != nil {
		return nil, err
	}
	out := make([]models.EDIOrderWithActiveReq, 0, len(rows))
	for _, d := range rows {
		m := models.EDIOrderWithActiveReq{
			// header
			EDIOrderID:      d.EDIOrderID,
			NumberOrder:     d.NumberOrder,
			VendorCode:      d.VendorCode,
			NumberForecast:  d.NumberForecast,
			ReadOrder:       d.ReadOrder,
			ActiveVersionID: d.ActiveVersionID,
			StatusOrder:     d.StatusOrder,
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

func (s *EDIOrderService) GetEDIOrderDetailByNumber(number string) (*models.EDIOrderDetailResp, error) {
	h, err := s.ediOrderRepo.GetOrderHeaderByNumber(number)
	if err != nil {
		return nil, err
	}

	vers, err := s.ediOrderRepo.GetOrderVersionsByOrderID(h.EDIOrderID)
	if err != nil {
		return nil, err
	}

	out := &models.EDIOrderDetailResp{
		EDIOrderID:      h.EDIOrderID,
		NumberOrder:     h.NumberOrder,
		NumberForecast:  h.NumberForecast,
		ReadOrder:       h.ReadOrder,
		VendorCode:      h.VendorCode,
		StatusOrder:     h.StatusOrder,
		FileURL:         h.FileURL,
		ActiveVersionID: h.ActiveVersionID,
		CreatedAt:       h.CreatedAt,
		UpdatedAt:       h.UpdatedAt,
		Versions:        make([]models.EDIOrderVersionItemResp, 0, len(vers)),
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

		out.Versions = append(out.Versions, models.EDIOrderVersionItemResp{
			EDIOrderVersionID: v.EDIOrderVersionID,
			VersionNo:         v.VersionNo,
			PeriodFrom:        v.PeriodFrom,
			PeriodTo:          v.PeriodTo,
			StatusOrder:       v.StatusOrder,
			Note:              v.Note,
			SourceFileURL:     v.SourceFileURL,
			CreatedAt:         v.CreatedAt,
			IsActive:          h.ActiveVersionID != nil && *h.ActiveVersionID == v.EDIOrderVersionID,
			CreatedBy:         createdBy,
		})
	}

	return out, nil
}

func (s *EDIOrderService) GetEDIOrderVersionByIDService(ediOrderVersionID string) (domains.EDIOrderVersion, error) {
	return s.ediOrderRepo.GetEDIOrderVersionByID(ediOrderVersionID)
}

func (s *EDIOrderService) GetOrderBasicByIDService(ediOrderID string) (domains.OrderBasicInfo, error) {
	return s.ediOrderRepo.GetOrderBasicByID(ediOrderID)
}

func (s *EDIOrderService) GetEDIVendorNotificationRecipientByCompanyService(company string) ([]models.EDIVendorNotificationRecipientReq, error) {
	recipients, err := s.ediOrderRepo.GetEDIVendorNotificationRecipientByCompany(company)
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

func (s *EDIOrderService) GetEDIOrderByNumberOrderDataService(
	numberOrder string,
) ([]models.EDIOrderVersionStatusLogReq, error) {

	logs, err := s.ediOrderRepo.GetEDIOrderByNumberOrderData(numberOrder)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	resp := make([]models.EDIOrderVersionStatusLogReq, 0, len(logs))

	for _, log := range logs {
		var principal *models.PrincipalResp
		if log.ChangedByPrincipal != nil {
			principal = &models.PrincipalResp{
				ExternalID:   log.ChangedByPrincipal.ExternalID,
				SourceSystem: log.ChangedByPrincipal.SourceSystem,
				Email:        log.ChangedByPrincipal.Email,
				DisplayName:  log.ChangedByPrincipal.DisplayName,
				Username:     log.ChangedByPrincipal.Username,
				Profile:      log.ChangedByPrincipal.Profile,
				Group:        log.ChangedByPrincipal.Group,
				Role:         log.ChangedByPrincipal.Role,
				Status:       log.ChangedByPrincipal.Status,
			}
		}

		resp = append(resp, models.EDIOrderVersionStatusLogReq{
			EDIOrderID:            log.EDIOrderID,
			OldStatus:             log.OldStatus,
			NewStatus:             log.NewStatus,
			Note:                  log.Note,
			ChangedByExternalID:   log.ChangedByExternalID,
			ChangedBySourceSystem: log.ChangedBySourceSystem,
			FileURL:               log.FileURL,
			CreatedAt:             log.CreatedAt,
			ChangedByUser:         principal,
		})
	}

	return resp, nil
}

func (s *EDIOrderService) GetOrderHeaderByVendorCodeService(VendorCode string) ([]models.EDIOrderByvendorReq, error) {

	log, err := s.ediOrderRepo.GetOrderHeaderByVendorCode(VendorCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	resp := []models.EDIOrderByvendorReq{
		{
			NumberForecast: log.NumberForecast,
		},
	}
	return resp, nil
}

func (s *EDIOrderService) GetStatusOrderSummaryByVendorCodeService(vendorCode string) (*models.StatusOrderSummaryResp, error) {

	domainResult, err := s.ediOrderRepo.GetStatusOrderSummaryByVendorCode(vendorCode)
	if err != nil {
		return nil, err
	}
	if domainResult == nil {
		return nil, nil
	}
	resp := &models.StatusOrderSummaryResp{
		VendorCode:    domainResult.VendorCode,
		NewCount:      domainResult.NewCount,
		ConfirmCount:  domainResult.ConfirmCount,
		RejectCount:   domainResult.RejectCount,
		ApprovedCount: domainResult.ApprovedCount,
		TotalCount:    domainResult.TotalCount,
	}

	return resp, nil
}

func (s *EDIOrderService) GetOrderHeaderByNumberForecastService(
	numberForecast string,
) ([]models.EDIOrderByNumberForecastResp, error) {

	domainResults, err := s.ediOrderRepo.GetOrderByNumberForecast(numberForecast)
	if err != nil {
		return nil, err
	}

	if len(domainResults) == 0 {
		return []models.EDIOrderByNumberForecastResp{}, nil
	}

	resp := make([]models.EDIOrderByNumberForecastResp, 0, len(domainResults))

	for _, d := range domainResults {
		resp = append(resp, models.EDIOrderByNumberForecastResp{
			EDIOrderID:     d.EDIOrderID,
			NumberOrder:    d.NumberOrder,
			NumberForecast: d.NumberForecast,
			VendorCode:     d.VendorCode,
			StatusOrder:    d.StatusOrder,
			CreatedAt:      d.CreatedAt,
			PeriodTo:       d.PeriodTo,
		})
	}

	return resp, nil
}
