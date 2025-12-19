package repositories

import (
	"fmt"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"

	"backend/internal/core/domains"
	ports "backend/internal/core/ports/repositories"
)

type EDIInvoiceRepositoryDB struct {
	db *gorm.DB
}

func NewEDIInvoiceRepositoryDB(db *gorm.DB) ports.EDIInvoiceRepository {
	// if err := db.AutoMigrate(&domains.EDIInvoice{}, &domains.EDIInvoiceVersion{}, &domains.EDIInvoiceVersionStatusLog{}); err != nil {
	// 	fmt.Printf("failed to auto migrate: %v", err)
	// }
	return &EDIInvoiceRepositoryDB{db: db}
}

func (r *EDIInvoiceRepositoryDB) CreateEDIInvoiceRepository(m *domains.EDIInvoice) error {
	const q = `
		INSERT INTO edi_invoice
			(edi_invoice_id, number_invoice, 
			active_version_id, created_by_external_id, created_by_source_system, status_invoice,
			file_url, vendor_code, number_order, invoice_type, read_invoice, created_at, updated_at)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0,  SYSUTCDATETIME(), SYSUTCDATETIME());
		`

	if err := r.db.Exec(q,
		m.EDIInvoiceID,
		m.NumberInvoice,
		m.ActiveVersionID,
		m.CreatedByExternalID,
		m.CreatedBySourceSystem,
		m.StatusInvoice,
		m.FileURL,
		m.VendorCode,
		m.NumberOrder,
		m.InvoiceType,
	).Error; err != nil {
		fmt.Printf("CreateEDIInvoiceRepository error: %v\n", err)
		return err
	}

	return nil
}

func (r *EDIInvoiceRepositoryDB) CreateEDIInvoiceVersionRepository(v *domains.EDIInvoiceVersion) error {
	const q = `
INSERT INTO edi_invoice_version
    (edi_invoice_version_id, edi_invoice_id, version_no, period_from, period_to,
     status_invoice, note, source_file_url,
     created_by_external_id, created_by_source_system, created_at, updated_at)
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, SYSUTCDATETIME(), SYSUTCDATETIME());
`
	return r.db.Exec(q,
		v.EDIInvoiceVersionID,   // 1  edi_invoice_version_id
		v.EDIInvoiceID,          // 2  edi_invoice_id
		v.VersionNo,             // 3  version_no
		v.PeriodFrom,            // 4  period_from
		v.PeriodTo,              // 5  period_to
		v.StatusInvoice,         // 6  status_invoice
		v.Note,                  // 7  note
		v.SourceFileURL,         // 8  source_file_url
		v.CreatedByExternalID,   // 9  created_by_external_id
		v.CreatedBySourceSystem, // 10 created_by_source_system
	).Error
}

func (r *EDIInvoiceRepositoryDB) GetLastInvoiceRunningForDate(date string) (int, error) {
	prefix := fmt.Sprintf("INV%s-", date)
	var last int

	sql := `
		SELECT 
			COALESCE(MAX(CAST(SUBSTRING(number_invoice, LEN(?) + 1, 10) AS INT)), 0)
		FROM edi_invoice
		WHERE number_invoice LIKE ?;
	`

	err := r.db.Raw(sql, prefix, prefix+"%").Scan(&last).Error
	if err != nil {
		return 0, err
	}

	return last, nil
}

func (r *EDIInvoiceRepositoryDB) UpdateActiveInvoiceVersion(headerID string, activeVersionID string) error {
	const q = `
	UPDATE edi_invoice
	SET active_version_id = ?,
	    updated_at = SYSUTCDATETIME()
	WHERE edi_invoice_id = ?;
	`
	return r.db.Exec(q, activeVersionID, headerID).Error
}

func (r *EDIInvoiceRepositoryDB) MarkInvoiceAsRead(id mssql.UniqueIdentifier) error {
	return r.db.
		Model(&domains.EDIInvoice{}).
		Where("edi_invoice_id = ?", id).
		Update("read_invoice", true).
		Error
}

func (r *EDIInvoiceRepositoryDB) GetEDIInvoiceWithActiveTop(limit int, vendorCode string) ([]domains.EDIInvoiceWithActive, error) {
	if limit <= 0 || limit > 1000 {
		limit = 50
	}

	whereVendor := "1=1"
	if vendorCode != "Prospira (Thailand) Co., Ltd." {
		whereVendor = "f.vendor_code = ?"
	}

	query := fmt.Sprintf(`
		SELECT TOP (CAST(? AS INT))
			f.edi_invoice_id,
			f.number_invoice,
			f.number_order,
			f.vendor_code,
			f.read_invoice,
			f.active_version_id,   
			f.status_invoice,
			f.file_url,
			f.created_at,
			f.updated_at,
			f.deleted_at,

			v.edi_invoice_version_id AS av_id,
			v.version_no AS av_version_no,
			v.period_from AS av_period_from,
			v.period_to AS av_period_to,
			v.status_invoice AS av_status,
			v.note AS av_note,
			v.source_file_url AS av_source_file_url,
			v.created_at AS av_created_at,
			v.deleted_at AS av_deleted_at,

			sl.edi_invoice_version_status_log_id AS last_status_log_id, 
			sl.old_status AS last_old_status,
			sl.new_status AS last_new_status,
			sl.note AS last_status_note,
			sl.file_url AS last_file_url,
			sl.created_at AS last_status_at

		FROM edi_invoice f
		LEFT JOIN edi_invoice_version v 
			ON v.edi_invoice_version_id = f.active_version_id

		OUTER APPLY (
			SELECT TOP (1)
				s.edi_invoice_version_status_log_id,
				s.old_status,
				s.new_status,
				s.note,
				s.file_url,
				s.created_at
			FROM edi_invoice_version_status_log s
			WHERE s.edi_invoice_id = f.edi_invoice_id  
			ORDER BY s.created_at DESC, s.edi_invoice_version_status_log_id DESC
		) sl

		WHERE %s AND f.deleted_at IS NULL
		ORDER BY f.created_at DESC;
	`, whereVendor)

	var out []domains.EDIInvoiceWithActive
	return out, r.db.Raw(query, limit, vendorCode).Scan(&out).Error
}

func (r *EDIInvoiceRepositoryDB) GetInvoiceHeaderByNumber(number string) (*domains.EDIInvoice, error) {
	const q = `
		SELECT f.edi_invoice_id, f.number_invoice, f.vendor_code, f.number_order, f.active_version_id,
		       f.status_invoice, f.created_at, f.updated_at, f.file_url, f.vendor_code, f.read_invoice
		FROM edi_invoice f
		WHERE f.deleted_at IS NULL
		  AND f.number_invoice = ?;
	`
	var h domains.EDIInvoice
	if err := r.db.Raw(q, number).Scan(&h).Error; err != nil {
		return nil, err
	}
	if h.EDIInvoiceID == (mssql.UniqueIdentifier{}) {
		return nil, gorm.ErrRecordNotFound
	}
	return &h, nil
}

func (r *EDIInvoiceRepositoryDB) GetInvoiceByNumberOrder(order string) (*domains.EDIInvoice, error) {
	const q = `
		SELECT f.edi_invoice_id, f.number_invoice, f.vendor_code, f.number_order, f.active_version_id,
		       f.status_invoice, f.created_at, f.updated_at, f.file_url, f.vendor_code, f.invoice_type
		FROM edi_invoice f
		WHERE f.deleted_at IS NULL
		  AND f.number_order = ?;
	`
	var h domains.EDIInvoice
	if err := r.db.Debug().Raw(q, order).Scan(&h).Error; err != nil {
		return nil, err
	}
	if h.EDIInvoiceID == (mssql.UniqueIdentifier{}) {
		return nil, gorm.ErrRecordNotFound
	}
	return &h, nil
}

type InvoiceVersionWithPrincipal struct {
	domains.EDIInvoiceVersion
	PrincipalExternalID   *string `gorm:"column:principal_external_id"`
	PrincipalSourceSystem *string `gorm:"column:principal_source_system"`
	PrincipalEmail        *string `gorm:"column:principal_email"`
	PrincipalDisplayName  *string `gorm:"column:principal_display_name"`
	PrincipalProfile      *string `gorm:"column:principal_profile"`
	PrincipalGroup        *string `gorm:"column:principal_group"`
	PrincipalRole         *string `gorm:"column:principal_role"`
	PrincipalStatus       *string `gorm:"column:principal_status"`
	PrincipalUsername     *string `gorm:"column:principal_username"`
}

func (r *EDIInvoiceRepositoryDB) GetInvoiceVersionsByInvoiceID(
	InvoiceID mssql.UniqueIdentifier,
) ([]domains.EDIInvoiceVersion, error) {

	var results []InvoiceVersionWithPrincipal

	err := r.db.
		Model(&domains.EDIInvoiceVersion{}).
		Select("edi_invoice_version.*, "+
			"p.external_id as principal_external_id, "+
			"p.source_system as principal_source_system, "+
			"p.email as principal_email, "+
			"p.display_name as principal_display_name, "+
			"p.profile as principal_profile, "+
			"p.[group] as principal_group, "+
			"p.role as principal_role, "+
			"p.status as principal_status, "+
			"p.username as principal_username").
		Joins("LEFT JOIN edi_principals p ON "+
			"edi_invoice_version.created_by_external_id = p.external_id AND "+
			"edi_invoice_version.created_by_source_system = p.source_system AND "+
			"p.deleted_at IS NULL").
		Where("edi_invoice_version.edi_invoice_id = ? AND edi_invoice_version.deleted_at IS NULL", InvoiceID).
		Order("version_no DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	vers := make([]domains.EDIInvoiceVersion, len(results))
	for i, r := range results {
		vers[i] = r.EDIInvoiceVersion

		if r.PrincipalExternalID != nil && *r.PrincipalExternalID != "" {
			vers[i].CreatedByPrincipal = &domains.EDI_Principal{
				ExternalID:   *r.PrincipalExternalID,
				SourceSystem: *r.PrincipalSourceSystem,
				Email:        getStringValue(r.PrincipalEmail),
				DisplayName:  getStringValue(r.PrincipalDisplayName),
				Profile:      getStringValue(r.PrincipalProfile),
				Group:        getStringValue(r.PrincipalGroup),
				Role:         getStringValue(r.PrincipalRole),
				Status:       getStringValue(r.PrincipalStatus),
				Username:     getStringValue(r.PrincipalUsername),
			}
		}
	}

	return vers, nil
}

func (r *EDIInvoiceRepositoryDB) UpdateStatusInvoice(id mssql.UniqueIdentifier, status string) error {
	const q = `
	UPDATE edi_invoice
	SET status_invoice = ?,
	    updated_at = SYSUTCDATETIME()
	WHERE edi_invoice_id = ?;
	`
	return r.db.Exec(q, status, id).Error
}

func (r *EDIInvoiceRepositoryDB) GetMaxVersionNoByInvoiceID(ediInvoiceID mssql.UniqueIdentifier) (int, error) {
	const q = `
        SELECT ISNULL(MAX(version_no), 0)
        FROM edi_invoice_version
        WHERE edi_invoice_id = ?;
    `
	var maxVer int
	if err := r.db.Raw(q, ediInvoiceID).Scan(&maxVer).Error; err != nil {
		return 0, err
	}
	return maxVer, nil
}

func (r *EDIInvoiceRepositoryDB) GetStatusInvoiceSummaryByVendorCode(vendorCode string) (*domains.InvoiceStatusSummary, error) {
	var (
		q    string
		args []interface{}
	)

	if vendorCode == "Prospira (Thailand) Co., Ltd." {
		q = `
			SELECT
				'All' AS vendor_code,
				SUM(CASE WHEN status_invoice = 'New'  THEN 1 ELSE 0 END) AS new_count,
				SUM(CASE WHEN status_invoice = 'Confirm'  THEN 1 ELSE 0 END) AS confirm_count,
				SUM(CASE WHEN status_invoice = 'Reject'   THEN 1 ELSE 0 END) AS reject_count,
				SUM(CASE WHEN status_invoice = 'Approved' THEN 1 ELSE 0 END) AS approved_count,
				COUNT(*) AS total_count
			FROM edi_invoice
			WHERE deleted_at IS NULL;
		`
	} else {
		q = `
			SELECT
				vendor_code,
				SUM(CASE WHEN status_invoice = 'New'  THEN 1 ELSE 0 END) AS new_count,
				SUM(CASE WHEN status_invoice = 'Confirm'  THEN 1 ELSE 0 END) AS confirm_count,
				SUM(CASE WHEN status_invoice = 'Reject'   THEN 1 ELSE 0 END) AS reject_count,
				SUM(CASE WHEN status_invoice = 'Approved' THEN 1 ELSE 0 END) AS approved_count,
				COUNT(*) AS total_count
			FROM edi_invoice
			WHERE deleted_at IS NULL
			  AND vendor_code = ?
			GROUP BY vendor_code;
		`
		args = append(args, vendorCode)
	}

	var summary domains.InvoiceStatusSummary
	if err := r.db.Raw(q, args...).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}

func (r *EDIInvoiceRepositoryDB) GetEDIInvoiceVersionByID(ediInvoiceVersionID string) (domains.EDIInvoiceVersion, error) {

	const q = `
	SELECT TOP (1)
		v.edi_invoice_version_id,
		v.edi_invoice_id,
		v.version_no,
		v.period_from,
		v.period_to,
		v.status_invoice,
		v.note,
		v.source_file_url,
		v.created_at,
		v.deleted_at,
		h.vendor_code,
		h.number_invoice,
		h.number_order
	FROM edi_invoice_version v
	INNER JOIN edi_invoice h
		ON h.edi_invoice_id = v.edi_invoice_id
	WHERE v.edi_invoice_version_id = ?;
`

	var out domains.EDIInvoiceVersion

	if err := r.db.Debug().Raw(q, ediInvoiceVersionID).Scan(&out).Error; err != nil {
		return domains.EDIInvoiceVersion{}, err
	}

	if out.EDIInvoiceVersionID == (domains.EDIInvoiceVersion{}.EDIInvoiceVersionID) {
		return domains.EDIInvoiceVersion{}, gorm.ErrRecordNotFound
	}

	return out, nil
}

func (r *EDIInvoiceRepositoryDB) GetEDIVendorNotificationRecipientByCompany(company string) ([]domains.EDIVendorNotificationRecipient, error) {

	rows, err := r.db.Raw(`
    SELECT 
        r.vendor_notification_recipient_id,
        r.company,
        r.notification_type,
        r.edi_principal_id,
        r.created_at,
        r.updated_at,
        r.deleted_at,

        ISNULL(p.edi_principal_id, '00000000-0000-0000-0000-000000000000') AS p_edi_principal_id,
        p.external_id AS p_external_id,
        p.source_system AS p_source_system,
        p.email AS p_email,
        p.display_name AS p_display_name,
        p.profile AS p_profile,
        p.[group] AS p_group,
        p.role AS p_role,
        p.status AS p_status,
        p.username AS p_username,
        p.created_at AS p_created_at,
        p.updated_at AS p_updated_at,
        p.deleted_at AS p_deleted_at
    FROM edi_vendor_notification_recipient r
    LEFT JOIN edi_principals p 
        ON r.edi_principal_id = p.edi_principal_id
    WHERE r.company = ?
      AND r.deleted_at IS NULL
`, company).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipients := []domains.EDIVendorNotificationRecipient{}

	for rows.Next() {
		var rec domains.EDIVendorNotificationRecipient
		var principal domains.EDI_Principal

		err := rows.Scan(
			&rec.VendorNotificationRecipientID,
			&rec.Company,
			&rec.NotificationType,
			&rec.EDI_PrincipalID,
			&rec.CreatedAt,
			&rec.UpdatedAt,
			&rec.DeletedAt,

			&principal.EDI_PrincipalID,
			&principal.ExternalID,
			&principal.SourceSystem,
			&principal.Email,
			&principal.DisplayName,
			&principal.Profile,
			&principal.Group,
			&principal.Role,
			&principal.Status,
			&principal.Username,
			&principal.CreatedAt,
			&principal.UpdatedAt,
			&principal.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		if principal.EDI_PrincipalID != (mssql.UniqueIdentifier{}) {
			rec.Principal = &principal
		}

		recipients = append(recipients, rec)
	}

	return recipients, nil
}

func (r *EDIInvoiceRepositoryDB) GetInvoiceBasicByID(ediInvoiceID string) (domains.InvoiceBasicInfo, error) {

	fmt.Println("ediInvoiceID :", ediInvoiceID)

	const q = `
		SELECT
			number_invoice,
			vendor_code
		FROM edi_invoice
		WHERE edi_invoice_id = ?;
	`

	var out domains.InvoiceBasicInfo

	if err := r.db.Debug().Raw(q, ediInvoiceID).Scan(&out).Error; err != nil {
		return domains.InvoiceBasicInfo{}, err
	}

	if out.NumberInvoice == "" && out.VendorCode == "" {
		return domains.InvoiceBasicInfo{}, gorm.ErrRecordNotFound
	}

	return out, nil
}
