package services

import (
	"backend/internal/core/domains"
	"backend/internal/core/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

func (s *EDIForecastService) CreateEDIForecastVersionStatusLogService(req models.EDI_ForecastVersionStatusLogResp) error {
	u := uuid.New()

	var newID mssql.UniqueIdentifier
	copy(newID[:], u[:])

	domainEDIForecastVersionStatusLog := domains.EDI_ForecastVersionStatusLog{
		EDIForecastVersionStatusLogID: newID,
		EDIForecastID:                 req.EDIForecastID,
		OldStatus:                     req.OldStatus,
		NewStatus:                     req.NewStatus,
		Note:                          req.Note,
		ChangedByExternalID:           req.ChangedByExternalID,
		ChangedBySourceSystem:         req.ChangedBySourceSystem,
		FileURL:                       req.FileURL,
		CreatedAt:                     req.CreatedAt,
	}

	if s.ediForecastRepo == nil {
		return fmt.Errorf("edi forecast repository is not initialized")
	}
	if err := s.ediForecastRepo.CreateEDIForecastVersionStatusLog(&domainEDIForecastVersionStatusLog); err != nil {
		return fmt.Errorf("failed to create forecast version status log: %w", err)
	}
	return nil
}

func (s *EDIForecastService) GetForecastVersionStatusLogByForecastVersionIDService(
	forecastVersionID string,
) ([]models.EDI_ForecastVersionStatusLogReq, error) {

	logs, err := s.ediForecastRepo.GetForecastVersionStatusLogByForecastVersionID(forecastVersionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	resp := make([]models.EDI_ForecastVersionStatusLogReq, 0, len(logs))

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

		resp = append(resp, models.EDI_ForecastVersionStatusLogReq{
			EDIForecastID:         log.EDIForecastID,
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

func (s *EDIForecastService) GetForecastVersionStatusLogByForecastVersionIDAndApprovedService(
	forecastVersionID string,
) ([]models.EDI_ForecastVersionStatusLogReq, error) {

	logs, err := s.ediForecastRepo.GetForecastVersionStatusLogByForecastVersionIDAndApproved(forecastVersionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	resp := make([]models.EDI_ForecastVersionStatusLogReq, 0, len(logs))

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

		resp = append(resp, models.EDI_ForecastVersionStatusLogReq{
			EDIForecastID:         log.EDIForecastID,
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
