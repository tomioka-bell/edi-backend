package repositories

import (
	"backend/internal/core/domains"
	"fmt"
)

func (r *EDIForecastRepositoryDB) CreateEDIForecastVersionStatusLog(m *domains.EDI_ForecastVersionStatusLog) error {
	const q = `
	INSERT INTO edi_forecast_version_status_log
		(edi_forecast_id, old_status, new_status, note, changed_by_external_id, changed_by_source_system, file_url, created_at)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, SYSUTCDATETIME());
`
	if err := r.db.Exec(q,
		m.EDIForecastID,
		m.OldStatus,
		m.NewStatus,
		m.Note,
		m.ChangedByExternalID,
		m.ChangedBySourceSystem,
		m.FileURL,
	).Error; err != nil {
		fmt.Printf("CreateEDIForecastRepository error: %v\n", err)
		return err
	}
	return nil
}

func (r *EDIForecastRepositoryDB) GetForecastVersionStatusLogByForecastVersionID(
	forecastID string,
) ([]domains.EDI_ForecastVersionStatusLog, error) {

	var logs []domains.EDI_ForecastVersionStatusLog

	if err := r.db.
		Where("edi_forecast_id = ?", forecastID).
		Order("created_at ASC").
		Find(&logs).Error; err != nil {
		return nil, err
	}

	for i := range logs {
		log := &logs[i]
		if log.ChangedByExternalID != "" && log.ChangedBySourceSystem != "" {
			var principal domains.EDI_Principal
			if err := r.db.
				Where("external_id = ? AND source_system = ?", log.ChangedByExternalID, log.ChangedBySourceSystem).
				First(&principal).Error; err == nil {
				log.ChangedByPrincipal = &principal
			}
		}
	}

	return logs, nil
}

func (r *EDIForecastRepositoryDB) GetForecastVersionStatusLogByForecastVersionIDAndApproved(
	forecastVersionID string,
) ([]domains.EDI_ForecastVersionStatusLog, error) {

	var logs []domains.EDI_ForecastVersionStatusLog

	if err := r.db.
		Where("edi_forecast_id = ?", forecastVersionID).
		Order("created_at ASC").
		Find(&logs).Error; err != nil {
		return nil, err
	}

	for i := range logs {
		log := &logs[i]
		if log.ChangedByExternalID != "" && log.ChangedBySourceSystem != "" {
			var principal domains.EDI_Principal
			if err := r.db.
				Where("external_id = ? AND source_system = ?", log.ChangedByExternalID, log.ChangedBySourceSystem).
				First(&principal).Error; err == nil {
				log.ChangedByPrincipal = &principal
			}
		}
	}

	return logs, nil
}
