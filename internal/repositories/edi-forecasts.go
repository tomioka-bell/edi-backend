package repositories

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"

	"backend/internal/core/domains"
	ports "backend/internal/core/ports/repositories"
)

type EDIForecastRepositoryDB struct {
	db *gorm.DB
}

func NewEDIForecastRepositoryDB(db *gorm.DB) ports.EDIForecastRepository {
	// if err := db.AutoMigrate(&domains.EDI_Forecast{}, &domains.EDI_ForecastVersion{}, &domains.EDI_ForecastVersionStatusLog{}); err != nil {
	// 	fmt.Printf("failed to auto migrate: %v", err)
	// }
	return &EDIForecastRepositoryDB{db: db}
}

func (r *EDIForecastRepositoryDB) CreateNewForecastWithVersion(
	db *gorm.DB,
	header *domains.EDI_Forecast,
	version *domains.EDI_ForecastVersion,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1 สร้าง Header ก่อน
		if err := tx.Create(header).Error; err != nil {
			return fmt.Errorf("create header failed: %w", err)
		}

		// 2️ ผูกเวอร์ชันกับ header ที่เพิ่งสร้าง
		version.EDIForecastID = header.EDI_ForecastID
		version.VersionNo = 1 // เวอร์ชันแรกเสมอ

		if err := tx.Create(version).Error; err != nil {
			return fmt.Errorf("create version failed: %w", err)
		}

		// 3️ อัปเดต header ให้รู้ว่าเวอร์ชันไหนคือ active version
		if err := tx.Model(&domains.EDI_Forecast{}).
			Where("edi_forecast_id = ?", header.EDI_ForecastID).
			Update("active_version_id", version.EDIForecastVersionID).Error; err != nil {
			return fmt.Errorf("update active version failed: %w", err)
		}

		return nil
	})
}

func (r *EDIForecastRepositoryDB) CreateEDIForecastRepository(m *domains.EDI_Forecast) error {
	const q = `
	INSERT INTO edi_forecast
		(edi_forecast_id, number_forecast, vendor_code,
		 active_version_id, created_by_external_id, created_by_source_system,
		 status_forecast, file_url, read_forecast, created_at, updated_at)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, 0, SYSUTCDATETIME(), SYSUTCDATETIME());
	`
	if err := r.db.Debug().Exec(q,
		m.EDI_ForecastID,
		m.NumberForecast,
		m.VendorCode,
		m.ActiveVersionID,
		m.CreatedByExternalID,
		m.CreatedBySourceSystem,
		m.StatusForecast,
		m.FileURL,
	).Error; err != nil {
		fmt.Printf("CreateEDIForecastRepository error: %v\n", err)
		return err
	}
	return nil
}

func (r *EDIForecastRepositoryDB) CreateEDIForecastVersionRepository(v *domains.EDI_ForecastVersion) error {
	const q = `
	INSERT INTO edi_forecast_version
		(edi_forecast_version_id, edi_forecast_id, version_no, period_from, period_to,
		 status_forecast, read_forecast, note, source_file_url,
		 created_by_external_id, created_by_source_system, created_at, updated_at)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, SYSUTCDATETIME(), SYSUTCDATETIME());
	`
	return r.db.Debug().Exec(q,
		v.EDIForecastVersionID,
		v.EDIForecastID,
		v.VersionNo,
		v.PeriodFrom,
		v.PeriodTo,
		v.StatusForecast,
		v.ReadForecast,
		v.Note,
		v.SourceFileURL,
		v.CreatedByExternalID,
		v.CreatedBySourceSystem,
	).Error
}

func (r *EDIForecastRepositoryDB) GetLastForecastRunningForDate(date string) (int, error) {
	prefix := fmt.Sprintf("FC%s-", date)
	var last int

	sql := `
		SELECT 
			COALESCE(MAX(CAST(SUBSTRING(number_forecast, LEN(?) + 1, 10) AS INT)), 0)
		FROM edi_forecast
		WHERE number_forecast LIKE ?;
	`

	err := r.db.Raw(sql, prefix, prefix+"%").Scan(&last).Error
	if err != nil {
		return 0, err
	}

	return last, nil
}

// อัปเดต header บางฟิลด์ (ตัวอย่าง)
func (r *EDIForecastRepositoryDB) UpdateActiveVersion(headerID string, activeVersionID string) error {
	const q = `
	UPDATE edi_forecast
	SET active_version_id = ?,
	    updated_at = SYSUTCDATETIME()
	WHERE edi_forecast_id = ?;
	`
	return r.db.Exec(q, activeVersionID, headerID).Error
}

func (r *EDIForecastRepositoryDB) MarkForecastAsRead(id mssql.UniqueIdentifier) (*domains.EDI_Forecast, error) {
	var f domains.EDI_Forecast
	if err := r.db.
		Where("edi_forecast_id = ?", id).
		First(&f).Error; err != nil {
		return nil, err
	}

	if f.ReadForecast && f.ReadAt != nil {
		return &f, nil
	}

	now := time.Now().UTC()

	if err := r.db.
		Model(&domains.EDI_Forecast{}).
		Where("edi_forecast_id = ?", id).
		Updates(map[string]any{
			"read_forecast": true,
			"read_at":       now,
		}).
		Error; err != nil {
		return nil, err
	}

	f.ReadForecast = true
	f.ReadAt = &now

	return &f, nil
}

func (r *EDIForecastRepositoryDB) UpdateStatusForecast(id mssql.UniqueIdentifier, status string) error {
	const q = `
	UPDATE edi_forecast
	SET status_forecast = ?,
	    updated_at = SYSUTCDATETIME()
	WHERE edi_forecast_id = ?;
	`
	return r.db.Exec(q, status, id).Error
}

/* -------------------------------- READ ----------------------------------- */

func (r *EDIForecastRepositoryDB) GetEDIForecastWithActiveTop(limit int, vendorCode string) ([]domains.EDIForecastWithActive, error) {
	if limit <= 0 || limit > 1000 {
		limit = 50
	}

	// ถ้าเป็น Prospira → vendor filter = 1=1 (ดึงทั้งหมด)
	// ถ้าไม่ใช่ → ใช้ vendor_code = ?
	whereVendor := "1=1"
	if vendorCode != "Prospira (Thailand) Co., Ltd." {
		whereVendor = "f.vendor_code = ?"
	}

	query := fmt.Sprintf(`
		SELECT TOP (CAST(? AS INT))
			f.edi_forecast_id,
			f.number_forecast,
			f.vendor_code,
			f.read_forecast,
			f.active_version_id,   
			f.status_forecast,
			f.file_url,
			f.created_at,
			f.updated_at,
			f.deleted_at,

			v.edi_forecast_version_id AS av_id,
			v.version_no AS av_version_no,
			v.period_from AS av_period_from,
			v.period_to AS av_period_to,
			v.status_forecast AS av_status,
			v.read_forecast AS av_read,
			v.note AS av_note,
			v.source_file_url AS av_source_file_url,
			v.created_at AS av_created_at,
			v.deleted_at AS av_deleted_at,

			sl.edi_forecast_version_status_log_id AS last_status_log_id,
			sl.old_status AS last_old_status,
			sl.new_status AS last_new_status,
			sl.note AS last_status_note,
			sl.file_url AS last_file_url,
			sl.created_at AS last_status_at

		FROM edi_forecast f
		LEFT JOIN edi_forecast_version v 
			ON v.edi_forecast_version_id = f.active_version_id

		OUTER APPLY (
			SELECT TOP (1)
				s.edi_forecast_version_status_log_id,
				s.old_status,
				s.new_status,
				s.note,
				s.file_url,
				s.created_at
			FROM edi_forecast_version_status_log s
			WHERE s.edi_forecast_id = f.edi_forecast_id  
			ORDER BY s.created_at DESC, s.edi_forecast_version_status_log_id DESC
		) sl

		WHERE %s AND f.deleted_at IS NULL
		ORDER BY f.created_at DESC;
	`, whereVendor)

	var out []domains.EDIForecastWithActive

	if vendorCode == "Prospira (Thailand) Co., Ltd." {
		return out, r.db.Raw(query, limit).Scan(&out).Error
	}

	return out, r.db.Raw(query, limit, vendorCode).Scan(&out).Error
}

func (r *EDIForecastRepositoryDB) GetEDIForecastWithActiveByNumber(number string) (*domains.EDIForecastWithActive, error) {
	const q = `
		SELECT
			f.edi_forecast_id, f.number_forecast, f.vendor_code, f.active_version_id,
			f.status_forecast,
			f.created_at, f.updated_at, f.deleted_at, f.row_ver,

			v.edi_forecast_version_id AS av_id,
			v.version_no              AS av_version_no,
			v.period_from             AS av_period_from,
			v.period_to               AS av_period_to,
			v.status_forecast         AS av_status,
			v.read_forecast           AS av_read,
			v.note                    AS av_note,
			v.source_file_url         AS av_source_file_url,
			v.created_at              AS av_created_at,
			v.deleted_at              AS av_deleted_at,
			v.row_ver                 AS av_row_ver
		FROM edi_forecast AS f
		LEFT JOIN edi_forecast_version AS v
			ON v.edi_forecast_version_id = f.active_version_id
		WHERE f.deleted_at IS NULL
			AND f.number_forecast = ?
	`
	var out domains.EDIForecastWithActive
	if err := r.db.Raw(q, number).Scan(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

// ดึง Header ตาม Number
func (r *EDIForecastRepositoryDB) GetForecastHeaderByNumber(number string) (*domains.EDI_Forecast, error) {
	const q = `
		SELECT f.edi_forecast_id, f.number_forecast, f.vendor_code, f.active_version_id,
		       f.status_forecast, f.created_at, f.updated_at, f.file_url, f.read_forecast
		FROM edi_forecast f
		WHERE f.deleted_at IS NULL
		  AND f.number_forecast = ?;
	`
	var h domains.EDI_Forecast
	if err := r.db.Raw(q, number).Scan(&h).Error; err != nil {
		return nil, err
	}
	if h.EDI_ForecastID == (mssql.UniqueIdentifier{}) {
		return nil, gorm.ErrRecordNotFound
	}
	return &h, nil
}

type versionWithPrincipal struct {
	domains.EDI_ForecastVersion
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

func (r *EDIForecastRepositoryDB) GetForecastVersionsByForecastID(
	forecastID mssql.UniqueIdentifier,
) ([]domains.EDI_ForecastVersion, error) {

	var results []versionWithPrincipal

	err := r.db.
		Model(&domains.EDI_ForecastVersion{}).
		Select("edi_forecast_version.*, "+
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
			"edi_forecast_version.created_by_external_id = p.external_id AND "+
			"edi_forecast_version.created_by_source_system = p.source_system AND "+
			"p.deleted_at IS NULL").
		Where("edi_forecast_version.edi_forecast_id = ? AND edi_forecast_version.deleted_at IS NULL", forecastID).
		Order("version_no DESC").
		Scan(&results).Error

	if err != nil {
		log.Printf("[GetForecastVersionsByForecastID] error: %v", err)
		return nil, err
	}

	vers := make([]domains.EDI_ForecastVersion, len(results))
	for i, r := range results {
		vers[i] = r.EDI_ForecastVersion

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

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ดึง header + ผูก Versions ให้ครบในหน่วยความจำ
func (r *EDIForecastRepositoryDB) GetEDIForecastWithAllVersionsTop(limit int) ([]domains.EDI_Forecast, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}

	// 1) ดึงหัวเอกสารตามลำดับเวลา (จำกัดจำนวน)
	const headQ = `
		SELECT
		edi_forecast_id, number_forecast, vendor_code, active_version_id,
		created_at, updated_at, deleted_at, row_ver
		FROM edi_forecast
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		OFFSET 0 ROWS FETCH NEXT ? ROWS ONLY;
		`
	var headers []domains.EDI_Forecast
	if err := r.db.Raw(headQ, limit).Scan(&headers).Error; err != nil {
		return nil, err
	}
	if len(headers) == 0 {
		return headers, nil
	}

	ids := make([]string, 0, len(headers))
	toStr := func(id mssql.UniqueIdentifier) string {
		u, _ := uuid.FromBytes(id[:])
		return u.String()
	}
	idx := make(map[string]int) // map headerID(string) -> index ใน headers
	for i, h := range headers {
		s := toStr(h.EDI_ForecastID)
		ids = append(ids, s)
		idx[s] = i
	}

	const verQ = `
		SELECT
		edi_forecast_version_id, edi_forecast_id, version_no, period_from, period_to,
		status_forecast, read_forecast, note, source_file_url,
		created_at, deleted_at, row_ver
		FROM edi_forecast_version
		WHERE deleted_at IS NULL
		AND edi_forecast_id IN ?
		ORDER BY edi_forecast_id, version_no DESC;
		`
	var versions []domains.EDI_ForecastVersion
	if err := r.db.Raw(verQ, ids).Scan(&versions).Error; err != nil {
		return nil, err
	}

	// 3) ประกอบ versions กลับเข้าไปในแต่ละ header
	for _, v := range versions {
		key := toStr(v.EDIForecastID)
		if i, ok := idx[key]; ok {
			headers[i].Versions = append(headers[i].Versions, v)
		}
	}

	return headers, nil
}

func (r *EDIForecastRepositoryDB) GetAllEDIForecast() ([]domains.EDI_Forecast, error) {
	const q = `
	SELECT
		edi_forecast_id, number_forecast, vendor_code, active_version_id,
		created_at, updated_at, deleted_at, row_ver
	FROM edi_forecast
	WHERE deleted_at IS NULL
	ORDER BY created_at DESC;
	`
	var list []domains.EDI_Forecast
	return list, r.db.Raw(q).Scan(&list).Error
}

func (r *EDIForecastRepositoryDB) GetAllEDIForecastVersion() ([]domains.EDI_ForecastVersion, error) {
	const q = `
	SELECT
		edi_forecast_version_id, edi_forecast_id, version_no, period_from, period_to,
		status_forecast, read_forecast, note, source_file_url,
		created_at, deleted_at, row_ver
	FROM edi_forecast_version
	WHERE deleted_at IS NULL
	ORDER BY created_at DESC;
	`
	var list []domains.EDI_ForecastVersion
	return list, r.db.Raw(q).Scan(&list).Error
}

func (r *EDIForecastRepositoryDB) GetEDIForecastTop(limit int) ([]domains.EDI_Forecast, error) {
	if limit <= 0 || limit > 1000 {
		limit = 50
	}
	const q = `
		SELECT
			edi_forecast_id, number_forecast, vendor_code, active_version_id,
			created_at, updated_at, deleted_at, row_ver
		FROM edi_forecast
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		OFFSET 0 ROWS FETCH NEXT ? ROWS ONLY;
	`
	var list []domains.EDI_Forecast
	return list, r.db.Raw(q, limit).Scan(&list).Error
}

func (r *EDIForecastRepositoryDB) GetEDIForecastByID(ediForecastID string) (domains.EDI_Forecast, error) {
	const q = `
	SELECT TOP (1)
		edi_forecast_id, number_forecast, vendor_code, active_version_id,
		created_at, updated_at, deleted_at, row_ver
	FROM edi_forecast
	WHERE edi_forecast_id = ?;
	`
	var out domains.EDI_Forecast
	if err := r.db.Raw(q, ediForecastID).Scan(&out).Error; err != nil {
		return domains.EDI_Forecast{}, err
	}
	if out.EDI_ForecastID == (domains.EDI_Forecast{}.EDI_ForecastID) {
		return domains.EDI_Forecast{}, gorm.ErrRecordNotFound
	}
	return out, nil
}

func (r *EDIForecastRepositoryDB) GetEDIForecastVersionByID(ediForecastVersionID string) (domains.EDI_ForecastVersion, error) {

	const q = `
	SELECT TOP (1)
		v.edi_forecast_version_id,
		v.edi_forecast_id,
		v.version_no,
		v.period_from,
		v.period_to,
		v.status_forecast,
		v.read_forecast,
		v.note,
		v.source_file_url,
		v.created_at,
		v.deleted_at,
		h.vendor_code,
		h.number_forecast
	FROM edi_forecast_version v
	INNER JOIN edi_forecast h
		ON h.edi_forecast_id = v.edi_forecast_id
	WHERE v.edi_forecast_version_id = ?;
`

	var out domains.EDI_ForecastVersion

	if err := r.db.Raw(q, ediForecastVersionID).Scan(&out).Error; err != nil {
		return domains.EDI_ForecastVersion{}, err
	}

	if out.EDIForecastVersionID == (domains.EDI_ForecastVersion{}.EDIForecastVersionID) {
		return domains.EDI_ForecastVersion{}, gorm.ErrRecordNotFound
	}

	return out, nil
}

func (r *EDIForecastRepositoryDB) GetEDIForecastVersionByVersionNo(versionNo string) (domains.EDI_ForecastVersion, error) {
	const q = `
	SELECT TOP (1)
		edi_forecast_version_id, edi_forecast_id, version_no, period_from, period_to,
		status_forecast, read_forecast, note, source_file_url,
		created_at, deleted_at, row_ver
	FROM edi_forecast_version
	WHERE version_no = @p1
	ORDER BY created_at DESC; 
	`
	var out domains.EDI_ForecastVersion
	if err := r.db.Raw(q, versionNo).Scan(&out).Error; err != nil {
		return domains.EDI_ForecastVersion{}, err
	}
	if out.EDIForecastVersionID == (domains.EDI_ForecastVersion{}.EDIForecastVersionID) {
		return domains.EDI_ForecastVersion{}, gorm.ErrRecordNotFound
	}
	return out, nil
}

func (r *EDIForecastRepositoryDB) GetEDIForecastVersionByDocAndVersion(ediForecastID string, versionNo int) (domains.EDI_ForecastVersion, error) {
	const q = `
	SELECT TOP (1)
		edi_forecast_version_id, edi_forecast_id, version_no, period_from, period_to,
		status_forecast, read_forecast, note, source_file_url,
		created_at, deleted_at, row_ver
	FROM edi_forecast_version
	WHERE edi_forecast_id = @p1 AND version_no = @p2;
	`
	var out domains.EDI_ForecastVersion
	if err := r.db.Raw(q, ediForecastID, versionNo).Scan(&out).Error; err != nil {
		return domains.EDI_ForecastVersion{}, err
	}
	if out.EDIForecastVersionID == (domains.EDI_ForecastVersion{}.EDIForecastVersionID) {
		return domains.EDI_ForecastVersion{}, gorm.ErrRecordNotFound
	}
	return out, nil
}

// ลบแบบ soft delete
func (r *EDIForecastRepositoryDB) SoftDeleteEDIForecast(headerID string) error {
	const q = `
	UPDATE edi_forecast
	SET deleted_at = SYSUTCDATETIME()
	WHERE edi_forecast_id = ? AND deleted_at IS NULL;
	`
	res := r.db.Exec(q, headerID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found or already deleted")
	}
	return nil
}

func (r *EDIForecastRepositoryDB) GetMaxVersionNoByForecastID(ediForecastID mssql.UniqueIdentifier) (int, error) {
	const q = `
        SELECT ISNULL(MAX(version_no), 0)
        FROM edi_forecast_version
        WHERE edi_forecast_id = ?;
    `
	var maxVer int
	if err := r.db.Raw(q, ediForecastID).Scan(&maxVer).Error; err != nil {
		return 0, err
	}
	return maxVer, nil
}

func (r *EDIForecastRepositoryDB) UpdateForecastVersionWithMap(forecastVersionID string, updates map[string]any) error {
	return r.db.Model(&domains.EDI_ForecastVersion{}).
		Where("edi_forecast_version_id = ?", forecastVersionID).
		Updates(updates).
		Error
}

func (r *EDIForecastRepositoryDB) UpdateActiveForecastVersion(
	forecastID mssql.UniqueIdentifier,
	versionID mssql.UniqueIdentifier,
) error {
	const q = `
        UPDATE edi_forecast
        SET active_version_id = ?
        WHERE edi_forecast_id = ?;
    `
	return r.db.Exec(q, versionID, forecastID).Error
}

func (r *EDIForecastRepositoryDB) GetEDIVendorNotificationRecipientByCompany(company string) ([]domains.EDIVendorNotificationRecipient, error) {

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

// ================================================= DataStatus =========================================================

func (r *EDIForecastRepositoryDB) GetStatusSummaryByVendorCode(vendorCode string) (*domains.ForecastStatusSummary, error) {
	var (
		q    string
		args []interface{}
	)

	if vendorCode == "Prospira (Thailand) Co., Ltd." {
		q = `
			SELECT
				'All' AS vendor_code,
				SUM(CASE WHEN status_forecast = 'New'  THEN 1 ELSE 0 END) AS new_count,
				SUM(CASE WHEN status_forecast = 'Confirm'  THEN 1 ELSE 0 END) AS confirm_count,
				SUM(CASE WHEN status_forecast = 'Reject'   THEN 1 ELSE 0 END) AS reject_count,
				SUM(CASE WHEN status_forecast = 'Approved' THEN 1 ELSE 0 END) AS approved_count,
				COUNT(*) AS total_count
			FROM edi_forecast
			WHERE deleted_at IS NULL;
		`
	} else {
		q = `
			SELECT
				vendor_code,
				SUM(CASE WHEN status_forecast = 'New'  THEN 1 ELSE 0 END) AS new_count,
				SUM(CASE WHEN status_forecast = 'Confirm'  THEN 1 ELSE 0 END) AS confirm_count,
				SUM(CASE WHEN status_forecast = 'Reject'   THEN 1 ELSE 0 END) AS reject_count,
				SUM(CASE WHEN status_forecast = 'Approved' THEN 1 ELSE 0 END) AS approved_count,
				COUNT(*) AS total_count
			FROM edi_forecast
			WHERE deleted_at IS NULL
			  AND vendor_code = ?
			GROUP BY vendor_code;
		`
		args = append(args, vendorCode)
	}

	var summary domains.ForecastStatusSummary
	if err := r.db.Raw(q, args...).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}

func (r *EDIForecastRepositoryDB) GetForecastBasicByID(ediForecastID string) (domains.ForecastBasicInfo, error) {

	const q = `
		SELECT
			number_forecast,
			vendor_code
		FROM edi_forecast
		WHERE edi_forecast_id = ?;
	`

	var out domains.ForecastBasicInfo

	if err := r.db.Debug().Raw(q, ediForecastID).Scan(&out).Error; err != nil {
		return domains.ForecastBasicInfo{}, err
	}

	if out.NumberForecast == "" && out.VendorCode == "" {
		return domains.ForecastBasicInfo{}, gorm.ErrRecordNotFound
	}

	return out, nil
}
