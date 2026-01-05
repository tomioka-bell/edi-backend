package repositories

import (
	"fmt"

	"backend/internal/core/domains"
)

func (r *EDIOrderRepositoryDB) CreateEDIOrderVersionStatusLog(m *domains.EDIOrderVersionStatusLog) error {
	const q = `
	INSERT INTO edi_order_version_status_log
		(edi_order_id, old_status, new_status, note, changed_by_external_id, changed_by_source_system, file_url, created_at)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, SYSUTCDATETIME());
`
	if err := r.db.Exec(q,
		m.EDIOrderID,
		m.OldStatus,
		m.NewStatus,
		m.Note,
		m.ChangedByExternalID,
		m.ChangedBySourceSystem,
		m.FileURL,
	).Error; err != nil {
		fmt.Printf("CreateEDIOrderVersionStatusLog error: %v\n", err)
		return err
	}
	return nil
}

func (r *EDIOrderRepositoryDB) GetOrderVersionStatusLogByOrderVersionIDAndApproved(
	orderVersionID string,
) ([]domains.EDIOrderVersionStatusLog, error) {

	var logs []domains.EDIOrderVersionStatusLog

	if err := r.db.
		Where("edi_order_id = ?", orderVersionID).
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

func (r *EDIOrderRepositoryDB) GetOrderVersionStatusLogByOrderVersionID(
	orderVersionID string,
) ([]domains.EDIOrderVersionStatusLog, error) {

	var logs []domains.EDIOrderVersionStatusLog

	if err := r.db.
		Where("edi_order_id = ?", orderVersionID).
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

func (r *EDIOrderRepositoryDB) GetEDIOrderByNumberOrderData(
	orderNumber string,
) ([]domains.EDIOrderVersionStatusLog, error) {

	var logs []domains.EDIOrderVersionStatusLog

	// 1) Find order header
	order, err := r.GetOrderHeaderByNumber(orderNumber)
	if err != nil {
		return nil, err
	}

	// 2) Find versions of this order
	var versions []domains.EDIOrderVersion
	if err := r.db.
		Where("edi_order_id = ?", order.EDIOrderID).
		Find(&versions).Error; err != nil {
		return nil, err
	}

	// 3) Loop and fetch logs from each version
	for _, version := range versions {

		var versionLogs []domains.EDIOrderVersionStatusLog

		if err := r.db.
			Where("edi_order_version_id = ?", version.EDIOrderVersionID).
			Order("created_at ASC").
			Find(&versionLogs).Error; err != nil {
			return nil, err
		}

		logs = append(logs, versionLogs...)
	}

	// 4) Load user principal
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

func (r *EDIOrderRepositoryDB) GetOrderVersionStatusLogByOrderNumberAndApproved(
	orderNumber string,
) ([]domains.EDIOrderVersionStatusLog, error) {

	var logs []domains.EDIOrderVersionStatusLog

	// Query using JOIN to get status logs by order number in one query
	if err := r.db.
		Joins("JOIN edi_order ON edi_order_version_status_log.edi_order_id = edi_order.edi_order_id").
		Where("edi_order.number_order = ?", orderNumber).
		Order("edi_order_version_status_log.created_at ASC").
		Find(&logs).Error; err != nil {
		return nil, err
	}

	// Load ChangedByPrincipal manually for each log
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
