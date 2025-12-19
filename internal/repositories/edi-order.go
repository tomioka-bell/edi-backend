package repositories

import (
	"fmt"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"

	"backend/internal/core/domains"
	ports "backend/internal/core/ports/repositories"
)

type EDIOrderRepositoryDB struct {
	db *gorm.DB
}

func NewEDIOrderRepositoryDB(db *gorm.DB) ports.EDIOrderRepository {
	// if err := db.AutoMigrate(&domains.EDIOrder{}, &domains.EDIOrderVersion{}, &domains.EDIOrderVersionStatusLog{}); err != nil {
	// 	fmt.Printf("failed to auto migrate: %v", err)
	// }
	return &EDIOrderRepositoryDB{db: db}
}

func (r *EDIOrderRepositoryDB) CreateEDIOrderRepository(m *domains.EDIOrder) error {
	const q = `
	INSERT INTO edi_order
		(edi_order_id, number_order, vendor_code, number_forecast,
		 active_version_id, created_by_external_id, created_by_source_system, status_order,
		 file_url, read_order, created_at, updated_at)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, 0, SYSUTCDATETIME(), SYSUTCDATETIME());
`

	if err := r.db.Exec(q,
		m.EDIOrderID,            // 1 edi_order_id
		m.NumberOrder,           // 2 number_order
		m.VendorCode,            // 3 vendor_code
		m.NumberForecast,        // 4 number_forecast
		m.ActiveVersionID,       // 5 active_version_id
		m.CreatedByExternalID,   // 6 created_by_external_id
		m.CreatedBySourceSystem, // 7 created_by_source_system
		m.StatusOrder,           // 8 status_order
		m.FileURL,               // 9 file_url
	).Error; err != nil {
		fmt.Printf("CreateEDIOrderRepository error: %v\n", err)
		return err
	}
	return nil
}

func (r *EDIOrderRepositoryDB) CreateEDIOrderVersionRepository(v *domains.EDIOrderVersion) error {
	const q = `
	INSERT INTO edi_order_version
		(edi_order_version_id, edi_order_id, version_no, period_from, period_to,
		 status_order, note, source_file_url,
		 created_by_external_id, created_by_source_system, created_at, updated_at)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, SYSUTCDATETIME(), SYSUTCDATETIME());
	`
	return r.db.Exec(q,
		v.EDIOrderVersionID,
		v.EDIOrderID,
		v.VersionNo,
		v.PeriodFrom,
		v.PeriodTo,
		v.StatusOrder,
		v.Note,
		v.SourceFileURL,
		v.CreatedByExternalID,
		v.CreatedBySourceSystem,
	).Error
}

func (r *EDIOrderRepositoryDB) GetLastOrderRunningForDate(date string) (int, error) {
	prefix := fmt.Sprintf("ORD%s-", date)
	var last int

	sql := `
		SELECT 
			COALESCE(MAX(CAST(SUBSTRING(number_order, LEN(?) + 1, 10) AS INT)), 0)
		FROM edi_order
		WHERE number_order LIKE ?;
	`

	err := r.db.Raw(sql, prefix, prefix+"%").Scan(&last).Error
	if err != nil {
		return 0, err
	}

	return last, nil
}

func (r *EDIOrderRepositoryDB) UpdateActiveOrderVersion(headerID string, activeVersionID string) error {
	const q = `
	UPDATE edi_order
	SET active_version_id = ?,
	    updated_at = SYSUTCDATETIME()
	WHERE edi_order_id = ?;
	`
	return r.db.Exec(q, activeVersionID, headerID).Error
}

func (r *EDIOrderRepositoryDB) MarkOrderAsRead(id mssql.UniqueIdentifier) error {
	return r.db.
		Model(&domains.EDIOrder{}).
		Where("edi_order_id = ?", id).
		Update("read_order", true).
		Error
}

func (r *EDIOrderRepositoryDB) GetEDIOrderWithActiveTop(limit int, vendorCode string) ([]domains.EDIOrderWithActive, error) {
	if limit <= 0 || limit > 1000 {
		limit = 50
	}

	whereVendor := "1=1"
	if vendorCode != "Prospira (Thailand) Co., Ltd." {
		whereVendor = "f.vendor_code = ?"
	}

	query := fmt.Sprintf(`
		SELECT TOP (CAST(? AS INT))
			f.edi_order_id,
			f.number_order,
			f.vendor_code,
			f.read_order,
			f.active_version_id,   
			f.status_order,
			f.file_url,
			f.created_at,
			f.updated_at,
			f.deleted_at,

			v.edi_order_version_id AS av_id,
			v.version_no AS av_version_no,
			v.period_from AS av_period_from,
			v.period_to AS av_period_to,
			v.status_order AS av_status,
			v.note AS av_note,
			v.source_file_url AS av_source_file_url,
			v.created_at AS av_created_at,
			v.deleted_at AS av_deleted_at,

			sl.edi_order_version_status_log_id AS last_status_log_id, 
			sl.old_status AS last_old_status,
			sl.new_status AS last_new_status,
			sl.note AS last_status_note,
			sl.file_url AS last_file_url,
			sl.created_at AS last_status_at

		FROM edi_order f
		LEFT JOIN edi_order_version v 
			ON v.edi_order_version_id = f.active_version_id

		OUTER APPLY (
			SELECT TOP (1)
				s.edi_order_version_status_log_id,
				s.old_status,
				s.new_status,
				s.note,
				s.file_url,
				s.created_at
			FROM edi_order_version_status_log s
			WHERE s.edi_order_id = f.edi_order_id  
			ORDER BY s.created_at DESC, s.edi_order_version_status_log_id DESC
		) sl

		WHERE %s AND f.deleted_at IS NULL
		ORDER BY f.created_at DESC;
	`, whereVendor)

	var out []domains.EDIOrderWithActive
	return out, r.db.Raw(query, limit, vendorCode).Scan(&out).Error
}

func (r *EDIOrderRepositoryDB) GetOrderHeaderByNumber(number string) (*domains.EDIOrder, error) {
	const q = `
		SELECT f.edi_order_id, f.number_order, f.vendor_code, f.number_forecast, f.active_version_id,
		       f.status_order, f.created_at, f.updated_at, f.file_url, f.read_order
		FROM edi_order f
		WHERE f.deleted_at IS NULL
		  AND f.number_order = ?;
	`
	var h domains.EDIOrder
	if err := r.db.Raw(q, number).Scan(&h).Error; err != nil {
		return nil, err
	}
	if h.EDIOrderID == (mssql.UniqueIdentifier{}) {
		return nil, gorm.ErrRecordNotFound
	}
	return &h, nil
}

type orderVersionWithPrincipal struct {
	domains.EDIOrderVersion
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

func (r *EDIOrderRepositoryDB) GetOrderVersionsByOrderID(
	orderID mssql.UniqueIdentifier,
) ([]domains.EDIOrderVersion, error) {

	var results []orderVersionWithPrincipal

	err := r.db.
		Model(&domains.EDIOrderVersion{}).
		Select("edi_order_version.*, "+
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
			"edi_order_version.created_by_external_id = p.external_id AND "+
			"edi_order_version.created_by_source_system = p.source_system AND "+
			"p.deleted_at IS NULL").
		Where("edi_order_version.edi_order_id = ? AND edi_order_version.deleted_at IS NULL", orderID).
		Order("version_no DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	vers := make([]domains.EDIOrderVersion, len(results))
	for i, r := range results {
		vers[i] = r.EDIOrderVersion

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

func (r *EDIOrderRepositoryDB) UpdateStatusOrder(id mssql.UniqueIdentifier, status string) error {
	const q = `
	UPDATE edi_order
	SET status_order = ?,
	    updated_at = SYSUTCDATETIME()
	WHERE edi_order_id = ?;
	`
	return r.db.Exec(q, status, id).Error
}

func (r *EDIOrderRepositoryDB) GetMaxVersionNoByOrderID(ediOrderID mssql.UniqueIdentifier) (int, error) {
	const q = `
        SELECT ISNULL(MAX(version_no), 0)
        FROM edi_order_version
        WHERE edi_order_id = ?;
    `
	var maxVer int
	if err := r.db.Raw(q, ediOrderID).Scan(&maxVer).Error; err != nil {
		return 0, err
	}
	return maxVer, nil
}

func (r *EDIOrderRepositoryDB) GetStatusOrderSummaryByVendorCode(vendorCode string) (*domains.OrderStatusSummary, error) {
	var (
		q    string
		args []interface{}
	)

	if vendorCode == "Prospira (Thailand) Co., Ltd." {
		q = `
			SELECT
				'All' AS vendor_code,
				SUM(CASE WHEN status_order = 'New'  THEN 1 ELSE 0 END) AS new_count,
				SUM(CASE WHEN status_order = 'Confirm'  THEN 1 ELSE 0 END) AS confirm_count,
				SUM(CASE WHEN status_order = 'Reject'   THEN 1 ELSE 0 END) AS reject_count,
				SUM(CASE WHEN status_order = 'Approved' THEN 1 ELSE 0 END) AS approved_count,
				COUNT(*) AS total_count
			FROM edi_order
			WHERE deleted_at IS NULL;
		`
	} else {
		q = `
			SELECT
				vendor_code,
				SUM(CASE WHEN status_order = 'New'  THEN 1 ELSE 0 END) AS new_count,
				SUM(CASE WHEN status_order = 'Confirm'  THEN 1 ELSE 0 END) AS confirm_count,
				SUM(CASE WHEN status_order = 'Reject'   THEN 1 ELSE 0 END) AS reject_count,
				SUM(CASE WHEN status_order = 'Approved' THEN 1 ELSE 0 END) AS approved_count,
				COUNT(*) AS total_count
			FROM edi_order
			WHERE deleted_at IS NULL
			  AND vendor_code = ?
			GROUP BY vendor_code;
		`
		args = append(args, vendorCode)
	}

	var summary domains.OrderStatusSummary
	if err := r.db.Raw(q, args...).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}

func (r *EDIOrderRepositoryDB) GetOrderHeaderByVendorCode(VendorCode string) (*domains.EDIOrder, error) {
	const q = `
		SELECT f.edi_order_id, f.number_order, f.vendor_code, f.number_forecast
		FROM edi_order f
		WHERE f.deleted_at IS NULL
		  AND f.vendor_code = ?;
	`
	var h domains.EDIOrder
	if err := r.db.Raw(q, VendorCode).Scan(&h).Error; err != nil {
		return nil, err
	}
	if h.EDIOrderID == (mssql.UniqueIdentifier{}) {
		return nil, gorm.ErrRecordNotFound
	}
	return &h, nil
}

func (r *EDIOrderRepositoryDB) GetOrderByNumberForecast(
	numberForecast string,
) ([]domains.EDIOrderHeaderWithPeriod, error) {

	const q = `
		SELECT
		    o.edi_order_id,
		    o.number_order,
			o.number_forecast,	
		    o.status_order,
			o.vendor_code,
		    o.created_at,
		    v.period_to
		FROM edi_order o
		LEFT JOIN edi_order_version v
		    ON v.edi_order_version_id = o.active_version_id
		   AND v.deleted_at IS NULL
		WHERE o.deleted_at IS NULL
		  AND o.number_forecast = ?;
	`

	var list []domains.EDIOrderHeaderWithPeriod
	if err := r.db.Raw(q, numberForecast).Scan(&list).Error; err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return list, nil
}

func (r *EDIOrderRepositoryDB) GetEDIOrderVersionByID(ediOrderVersionID string) (domains.EDIOrderVersion, error) {

	const q = `
	SELECT TOP (1)
		v.edi_order_version_id,
		v.edi_order_id,
		v.version_no,
		v.period_from,
		v.period_to,
		v.status_order,
		v.note,
		v.source_file_url,
		v.created_at,
		v.deleted_at,
		h.vendor_code,
		h.number_order
	FROM edi_order_version v
	INNER JOIN edi_order h
		ON h.edi_order_id = v.edi_order_id
	WHERE v.edi_order_version_id = ?;
`

	var out domains.EDIOrderVersion

	if err := r.db.Debug().Raw(q, ediOrderVersionID).Scan(&out).Error; err != nil {
		return domains.EDIOrderVersion{}, err
	}

	if out.EDIOrderVersionID == (domains.EDIOrderVersion{}.EDIOrderVersionID) {
		return domains.EDIOrderVersion{}, gorm.ErrRecordNotFound
	}

	return out, nil
}

func (r *EDIOrderRepositoryDB) GetEDIVendorNotificationRecipientByCompany(company string) ([]domains.EDIVendorNotificationRecipient, error) {

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

func (r *EDIOrderRepositoryDB) GetOrderBasicByID(ediOrderID string) (domains.OrderBasicInfo, error) {

	const q = `
		SELECT
			number_order,
			vendor_code
		FROM edi_order
		WHERE edi_order_id = ?;
	`

	var out domains.OrderBasicInfo

	if err := r.db.Debug().Raw(q, ediOrderID).Scan(&out).Error; err != nil {
		return domains.OrderBasicInfo{}, err
	}

	if out.NumberOrder == "" && out.VendorCode == "" {
		return domains.OrderBasicInfo{}, gorm.ErrRecordNotFound
	}

	return out, nil
}
