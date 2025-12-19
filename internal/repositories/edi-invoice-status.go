package repositories

import (
	"backend/internal/core/domains"
	"fmt"
)

func (r *EDIInvoiceRepositoryDB) CreateEDIInvoiceVersionStatusLog(m *domains.EDIInvoiceVersionStatusLog) error {
	const q = `
	INSERT INTO edi_invoice_version_status_log
		(edi_invoice_id, old_status, new_status, note, changed_by_external_id, changed_by_source_system, file_url, created_at)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, SYSUTCDATETIME());
`
	if err := r.db.Exec(q,
		m.EDIInvoiceID,
		m.OldStatus,
		m.NewStatus,
		m.Note,
		m.ChangedByExternalID,
		m.ChangedBySourceSystem,
		m.FileURL,
	).Error; err != nil {
		fmt.Printf("CreateEDIInvoiceVersionStatusLog error: %v\n", err)
		return err
	}
	return nil
}

func (r *EDIInvoiceRepositoryDB) GetInvoiceVersionStatusLogByInvoiceVersionID(
	InvoiceVersionID string,
) ([]domains.EDIInvoiceVersionStatusLog, error) {
	var logs []domains.EDIInvoiceVersionStatusLog

	if err := r.db.
		Where("edi_invoice_id = ?", InvoiceVersionID).
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

func (r *EDIInvoiceRepositoryDB) GetInvoiceVersionStatusLogByInvoiceVersionIDAndApproved(
	InvoiceVersionID string,
) ([]domains.EDIInvoiceVersionStatusLog, error) {

	var logs []domains.EDIInvoiceVersionStatusLog

	if err := r.db.
		Where("edi_invoice_id = ? AND new_status = ?", InvoiceVersionID).
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
