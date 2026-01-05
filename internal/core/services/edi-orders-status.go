package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"

	"backend/internal/core/domains"
	"backend/internal/core/models"
)

func (s *EDIOrderService) CreateEDIOrderVersionStatusLogService(req models.EDIOrderVersionStatusLogResp) error {
	u := uuid.New()

	var newID mssql.UniqueIdentifier
	copy(newID[:], u[:])

	domainEDIForecastVersionStatusLog := domains.EDIOrderVersionStatusLog{
		EDIOrderVersionStatusLogID: newID,
		EDIOrderID:                 req.EDIOrderID,
		OldStatus:                  req.OldStatus,
		NewStatus:                  req.NewStatus,
		Note:                       req.Note,
		ChangedByExternalID:        req.ChangedByExternalID,
		ChangedBySourceSystem:      req.ChangedBySourceSystem,
		FileURL:                    req.FileURL,
		CreatedAt:                  req.CreatedAt,
	}

	if s.ediOrderRepo == nil {
		return fmt.Errorf("edi order repository is not initialized")
	}
	if err := s.ediOrderRepo.CreateEDIOrderVersionStatusLog(&domainEDIForecastVersionStatusLog); err != nil {
		return fmt.Errorf("failed to create order version status log: %w", err)
	}
	return nil
}

func (s *EDIOrderService) GetOrderVersionStatusLogByOrderVersionIDService(
	orderVersionID string,
) ([]models.EDIOrderVersionStatusLogReq, error) {

	logs, err := s.ediOrderRepo.GetOrderVersionStatusLogByOrderVersionID(orderVersionID)
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

func (s *EDIOrderService) GetOrderVersionStatusLogByOrderVersionIDAndApprovedService(
	orderVersionID string,
) ([]models.EDIOrderVersionStatusLogReq, error) {

	logs, err := s.ediOrderRepo.GetOrderVersionStatusLogByOrderVersionIDAndApproved(orderVersionID)
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

func (s *EDIOrderService) GetOrderVersionStatusLogByOrderNumberAndApprovedService(
	orderNumber string,
) ([]models.EDIOrderVersionStatusLogReq, error) {

	logs, err := s.ediOrderRepo.GetOrderVersionStatusLogByOrderNumberAndApproved(orderNumber)
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
