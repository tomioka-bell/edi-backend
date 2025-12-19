package repositories

import (
	"backend/internal/core/domains"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type EDISummaryDataRepositoryDB struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewEDISummaryDataRepositoryDB(db *gorm.DB, redisClient *redis.Client) *EDISummaryDataRepositoryDB {
	return &EDISummaryDataRepositoryDB{db: db, redisClient: redisClient}
}

func (r *EDISummaryDataRepositoryDB) getFromCache(key string, dest interface{}) (bool, error) {
	val, err := r.redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, err
	}

	return true, nil
}

func (r *EDISummaryDataRepositoryDB) setCache(key string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.redisClient.Set(
		context.Background(),
		key,
		jsonData,
		ttl,
	).Err()
}

func (r *EDISummaryDataRepositoryDB) hasCache(key string) (bool, error) {
	exists, err := r.redisClient.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// ================================= Get Status Summary =======================================

func (r *EDISummaryDataRepositoryDB) getStatusSummary(
	tableName string,
	statusField string,
	dest interface{},
) error {

	query := fmt.Sprintf(`
		SELECT
			vendor_code,
			SUM(CASE WHEN %s = 'New'      THEN 1 ELSE 0 END) AS new_count,
			SUM(CASE WHEN %s = 'Confirm'  THEN 1 ELSE 0 END) AS confirm_count,
			SUM(CASE WHEN %s = 'Reject'   THEN 1 ELSE 0 END) AS reject_count,
			SUM(CASE WHEN %s = 'Approved' THEN 1 ELSE 0 END) AS approved_count,
			COUNT(*) AS total_count
		FROM %s
		WHERE deleted_at IS NULL
		GROUP BY vendor_code;
	`, statusField, statusField, statusField, statusField, tableName)

	if err := r.db.Raw(query).Scan(dest).Error; err != nil {
		return err
	}
	return nil
}

func (r *EDISummaryDataRepositoryDB) GetStatusForecastSummaryByVendorCode() ([]domains.ForecastStatusSummary, error) {
	key := "summary:forecast"
	var summary []domains.ForecastStatusSummary

	found, err := r.getFromCache(key, &summary)
	if err != nil {
		fmt.Println("Redis GET error:", err)
	}
	if found {
		fmt.Println("Loaded from Redis!")
		return summary, nil
	}

	err = r.getStatusSummary("edi_forecast", "status_forecast", &summary)
	if err != nil {
		return nil, err
	}

	if err := r.setCache(key, summary, 5*time.Minute); err != nil {
		fmt.Println("Redis SET error:", err)
	} else {
		fmt.Println("Saved to Redis:", key)
	}

	has, err := r.hasCache(key)
	if err != nil {
		fmt.Println("Redis EXISTS error:", err)
	}

	if has {
		fmt.Println("Redis There is information now!")
	}

	return summary, nil
}

func (r *EDISummaryDataRepositoryDB) GetStatusOrderSummaryByVendorCode() ([]domains.OrderStatusSummary, error) {
	key := "summary:order"
	var summary []domains.OrderStatusSummary

	found, err := r.getFromCache(key, &summary)
	if err != nil {
		fmt.Println("Redis GET error:", err)
	}
	if found {
		fmt.Println("Loaded from Redis!")
		return summary, nil
	}

	err = r.getStatusSummary("edi_order", "status_order", &summary)
	if err != nil {
		return nil, err
	}

	if err := r.setCache(key, summary, 5*time.Minute); err != nil {
		fmt.Println("Redis SET error:", err)
	} else {
		fmt.Println("Saved to Redis:", key)
	}

	has, err := r.hasCache(key)
	if err != nil {
		fmt.Println("Redis EXISTS error:", err)
	}

	if has {
		fmt.Println("Redis There is information now!")
	}
	return summary, nil
}

func (r *EDISummaryDataRepositoryDB) GetStatusInvoiceSummaryByVendorCode() ([]domains.OrderStatusSummary, error) {
	var summary []domains.OrderStatusSummary
	err := r.getStatusSummary("edi_invoice", "status_invoice", &summary)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

// ================================= Get Sum Status Summary =======================================
func (r *EDISummaryDataRepositoryDB) getStatusTotal(
	tableName string,
	statusField string,
	dest interface{},
) error {

	query := fmt.Sprintf(`
		SELECT
			SUM(CASE WHEN %[1]s = 'New'      THEN 1 ELSE 0 END) AS new_count,
			SUM(CASE WHEN %[1]s = 'Confirm'  THEN 1 ELSE 0 END) AS confirm_count,
			SUM(CASE WHEN %[1]s = 'Reject'   THEN 1 ELSE 0 END) AS reject_count,
			SUM(CASE WHEN %[1]s = 'Approved' THEN 1 ELSE 0 END) AS approved_count,
			COUNT(*) AS total_count
		FROM %s
		WHERE deleted_at IS NULL;
	`, statusField, tableName)

	if err := r.db.Raw(query).Scan(dest).Error; err != nil {
		return err
	}
	return nil
}

func (r *EDISummaryDataRepositoryDB) GetForecastStatusTotal() (*domains.StatusTotalSummary, error) {
	var summary domains.StatusTotalSummary
	if err := r.getStatusTotal("edi_forecast", "status_forecast", &summary); err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *EDISummaryDataRepositoryDB) GetOrderStatusTotal() (*domains.StatusTotalSummary, error) {
	var summary domains.StatusTotalSummary
	if err := r.getStatusTotal("edi_order", "status_order", &summary); err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *EDISummaryDataRepositoryDB) GetInvoiceStatusTotal() (*domains.StatusTotalSummary, error) {
	var summary domains.StatusTotalSummary
	if err := r.getStatusTotal("edi_invoice", "status_invoice", &summary); err != nil {
		return nil, err
	}
	return &summary, nil
}

// ================================= Get Count By Number =======================================

func (r *EDISummaryDataRepositoryDB) getTotalCount(tableName string) (int64, error) {
	var total int64

	query := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM %s
        WHERE deleted_at IS NULL;
    `, tableName)

	if err := r.db.Raw(query).Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *EDISummaryDataRepositoryDB) GetForecastTotalCount() (int64, error) {
	return r.getTotalCount("edi_forecast")
}

func (r *EDISummaryDataRepositoryDB) GetOrderTotalCount() (int64, error) {
	return r.getTotalCount("edi_order")
}

func (r *EDISummaryDataRepositoryDB) GetInvoiceTotalCount() (int64, error) {
	return r.getTotalCount("edi_invoice")
}

// ================================= Get Status Summary By Month =======================================

func (r *EDISummaryDataRepositoryDB) getMonthlyStatusSummary(
	tableName string,
	statusField string,
	dest interface{},
) error {

	query := fmt.Sprintf(`
		SELECT
			YEAR(created_at)  AS year,
			MONTH(created_at) AS month,
			SUM(CASE WHEN %[1]s = 'New'      THEN 1 ELSE 0 END) AS new_count,
			SUM(CASE WHEN %[1]s = 'Confirm'  THEN 1 ELSE 0 END) AS confirm_count,
			SUM(CASE WHEN %[1]s = 'Reject'   THEN 1 ELSE 0 END) AS reject_count,
			SUM(CASE WHEN %[1]s = 'Approved' THEN 1 ELSE 0 END) AS approved_count,
			COUNT(*) AS total_count
		FROM %s
		WHERE deleted_at IS NULL
		GROUP BY YEAR(created_at), MONTH(created_at)
		ORDER BY YEAR(created_at), MONTH(created_at);
	`, statusField, tableName)

	if err := r.db.Raw(query).Scan(dest).Error; err != nil {
		return err
	}
	return nil
}

func (r *EDISummaryDataRepositoryDB) GetForecastMonthlyStatusSummary() ([]domains.MonthlyStatusSummary, error) {
	var result []domains.MonthlyStatusSummary
	if err := r.getMonthlyStatusSummary("edi_forecast", "status_forecast", &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *EDISummaryDataRepositoryDB) GetOrderMonthlyStatusSummary() ([]domains.MonthlyStatusSummary, error) {
	var result []domains.MonthlyStatusSummary
	if err := r.getMonthlyStatusSummary("edi_order", "status_order", &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *EDISummaryDataRepositoryDB) GetInvoiceMonthlyStatusSummary() ([]domains.MonthlyStatusSummary, error) {
	var result []domains.MonthlyStatusSummary
	if err := r.getMonthlyStatusSummary("edi_invoice", "status_invoice", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ================================= Get Status Summary Vendor =======================================

func (r *EDISummaryDataRepositoryDB) CountUsers() (int64, error) {
	var count int64
	result := r.db.Model(&domains.EDIVendorMetrics{}).Count(&count)
	return count, result.Error
}

func (r *EDISummaryDataRepositoryDB) GetUnreadForecastByVendor(
	ctx context.Context,
	vendorCode string,
) ([]domains.VendorForecastSummary, error) {

	var result []domains.VendorForecastSummary

	query := r.db.WithContext(ctx).
		Model(&domains.EDI_Forecast{}).
		Select("number_forecast, status_forecast, vendor_code, created_at").
		Where("read_forecast = 0 AND deleted_at IS NULL")

	// ถ้าไม่ใช่ Prospira ให้ where vendor code
	if vendorCode != "Prospira (Thailand) Co., Ltd." {
		query = query.Where("vendor_code = ?", vendorCode)
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *EDISummaryDataRepositoryDB) GetUnreadOrderByVendor(
	ctx context.Context,
	vendorCode string,
) ([]domains.VendorOrderSummary, error) {

	var result []domains.VendorOrderSummary

	query := r.db.WithContext(ctx).
		Model(&domains.EDIOrder{}).
		Select("number_order, status_order, vendor_code, created_at").
		Where("read_order = 0 AND deleted_at IS NULL")

	if vendorCode != "Prospira (Thailand) Co., Ltd." {
		query = query.Where("vendor_code = ?", vendorCode)
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *EDISummaryDataRepositoryDB) GetUnreadInvoiceByVendor(
	ctx context.Context,
	vendorCode string,
) ([]domains.VendorInvoiceSummary, error) {

	var result []domains.VendorInvoiceSummary

	query := r.db.WithContext(ctx).
		Model(&domains.EDIInvoice{}).
		Select("number_invoice, status_invoice, vendor_code, created_at").
		Where("read_invoice = 0 AND deleted_at IS NULL")

	if vendorCode != "Prospira (Thailand) Co., Ltd." {
		query = query.Where("vendor_code = ?", vendorCode)
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *EDISummaryDataRepositoryDB) GetForecastPeriodAlerts(ctx context.Context) ([]domains.ForecastPeriodAlert, error) {
	var result []domains.ForecastPeriodAlert

	sql := `
		SELECT
			f.number_forecast,
			f.vendor_code,
			v.version_no,
			v.period_from,

			f.read_forecast,
			ISNULL(m.reminder_days, 3) AS reminder_days,

			DATEADD(DAY, ISNULL(m.reminder_days, 3), f.created_at) AS target_time,

			DATEDIFF(
				SECOND,
				SYSDATETIMEOFFSET(),
				DATEADD(DAY, ISNULL(m.reminder_days, 3), f.created_at)
			) AS seconds_to_target,

			DATEDIFF(
				DAY,
				SYSDATETIMEOFFSET(),
				DATEADD(DAY, ISNULL(m.reminder_days, 3), f.created_at)
			) AS days_to_target,

			p.email AS email
		FROM edi_forecast f

		LEFT JOIN edi_forecast_version v
			ON v.edi_forecast_version_id = f.active_version_id
			AND v.deleted_at IS NULL

		LEFT JOIN edi_vendor_metrics m
			ON m.company_name = f.vendor_code
			AND m.deleted_at IS NULL
			AND m.active = 1

		LEFT JOIN edi_vendor_notification_recipient rcp
			ON rcp.company = f.vendor_code
			AND rcp.deleted_at IS NULL

		LEFT JOIN dbo.edi_principals p
			ON p.edi_principal_id = rcp.edi_principal_id

	WHERE
    f.deleted_at IS NULL
    AND f.read_forecast = 0
    AND SYSDATETIMEOFFSET() >= DATEADD(DAY, 1, f.created_at)
    AND SYSDATETIMEOFFSET() <  DATEADD(DAY, ISNULL(m.reminder_days, 3) + 1, f.created_at)
	`

	if err := r.db.WithContext(ctx).Raw(sql).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *EDISummaryDataRepositoryDB) GetOrderPeriodAlerts(ctx context.Context) ([]domains.OrderPeriodAlert, error) {
	var result []domains.OrderPeriodAlert

	sql := `
		SELECT
			f.number_order,
			f.vendor_code,
			v.version_no,
			v.period_from,

			f.read_order,
			ISNULL(m.reminder_days, 3) AS reminder_days,

			DATEADD(DAY, ISNULL(m.reminder_days, 3), f.created_at) AS target_time,

			DATEDIFF(
				SECOND,
				SYSDATETIMEOFFSET(),
				DATEADD(DAY, ISNULL(m.reminder_days, 3), f.created_at)
			) AS seconds_to_target,

			DATEDIFF(
				DAY,
				SYSDATETIMEOFFSET(),
				DATEADD(DAY, ISNULL(m.reminder_days, 3), f.created_at)
			) AS days_to_target,

			p.email AS email
		FROM edi_order f

		LEFT JOIN edi_order_version v
			ON v.edi_order_version_id = f.active_version_id
			AND v.deleted_at IS NULL

		LEFT JOIN edi_vendor_metrics m
			ON m.company_name = f.vendor_code
			AND m.deleted_at IS NULL
			AND m.active = 1

		LEFT JOIN edi_vendor_notification_recipient rcp
			ON rcp.company = f.vendor_code
			AND rcp.deleted_at IS NULL

		LEFT JOIN dbo.edi_principals p
			ON p.edi_principal_id = rcp.edi_principal_id

	WHERE
    f.deleted_at IS NULL
    AND f.read_order = 0
    AND SYSDATETIMEOFFSET() >= DATEADD(DAY, 1, f.created_at)
    AND SYSDATETIMEOFFSET() <  DATEADD(DAY, ISNULL(m.reminder_days, 3) + 1, f.created_at)
	`

	if err := r.db.WithContext(ctx).Raw(sql).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *EDISummaryDataRepositoryDB) GetInvoicePeriodAlerts(ctx context.Context) ([]domains.InvoicePeriodAlert, error) {
	var result []domains.InvoicePeriodAlert

	sql := `
		SELECT
			f.number_invoice,
			f.vendor_code,
			v.version_no,
			v.period_from,
  
			f.read_invoice,
			ISNULL(m.reminder_days, 3) AS reminder_days,

			DATEADD(DAY, ISNULL(m.reminder_days, 3), f.created_at) AS target_time,

			DATEDIFF(
				SECOND,
				SYSDATETIMEOFFSET(),
				DATEADD(DAY, ISNULL(m.reminder_days, 3), f.created_at)
			) AS seconds_to_target,

			DATEDIFF(
				DAY,
				SYSDATETIMEOFFSET(),
				DATEADD(DAY, ISNULL(m.reminder_days, 3), f.created_at)
			) AS days_to_target,

			p.email AS email
		FROM edi_invoice f

		LEFT JOIN edi_invoice_version v
			ON v.edi_invoice_version_id = f.active_version_id
			AND v.deleted_at IS NULL

		LEFT JOIN edi_vendor_metrics m
			ON m.company_name = f.vendor_code
			AND m.deleted_at IS NULL
			AND m.active = 1

		LEFT JOIN edi_vendor_notification_recipient rcp
			ON rcp.company = f.vendor_code
			AND rcp.deleted_at IS NULL

		LEFT JOIN dbo.edi_principals p
			ON p.edi_principal_id = rcp.edi_principal_id

	WHERE
    f.deleted_at IS NULL
    AND f.read_invoice = 0
    AND SYSDATETIMEOFFSET() >= DATEADD(DAY, 1, f.created_at)
    AND SYSDATETIMEOFFSET() <  DATEADD(DAY, ISNULL(m.reminder_days, 3) + 1, f.created_at)
	`

	if err := r.db.WithContext(ctx).Raw(sql).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
