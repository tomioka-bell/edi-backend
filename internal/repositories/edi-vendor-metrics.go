package repositories

import (
	"fmt"

	"gorm.io/gorm"

	"backend/internal/core/domains"
	ports "backend/internal/core/ports/repositories"
)

type EDIVendorMetricsRepositoryDB struct {
	db *gorm.DB
}

func NewEDIVendorMetricsRepositoryDB(db *gorm.DB) ports.EDIVendorMetricsRepository {
	// if err := db.AutoMigrate(&domains.EDIVendorMetrics{}); err != nil {
	// 	fmt.Printf("failed to auto migrate: %v", err)
	// }
	return &EDIVendorMetricsRepositoryDB{db: db}
}

func (r *EDIVendorMetricsRepositoryDB) CreateVendorMetrics(vendorMetrics *domains.EDIVendorMetrics) error {
	if err := r.db.Debug().Create(vendorMetrics).Error; err != nil {
		fmt.Printf("CreateVendorMetrics error: %v\n", err)
		return err
	}
	return nil
}

func (r *EDIVendorMetricsRepositoryDB) GetAllVendorMetrics() ([]domains.EDIVendorMetrics, error) {
	var metrics []domains.EDIVendorMetrics

	if err := r.db.
		Select("initials", "company_name").
		Find(&metrics).Error; err != nil {
		return nil, err
	}

	return metrics, nil
}

func (r *EDIVendorMetricsRepositoryDB) GetVendorMetricsByCompany(company string) ([]domains.EDIVendorMetrics, error) {
	var metrics []domains.EDIVendorMetrics
	if err := r.db.Where("company_name = ?", company).Find(&metrics).Error; err != nil {
		return nil, err
	}
	return metrics, nil
}

func (r *EDIVendorMetricsRepositoryDB) GetAllEDIVendorMetricsTop(limit int) ([]domains.EDIVendorMetrics, error) {
	if limit <= 0 || limit > 1000 {
		limit = 50
	}

	const q = `
        SELECT TOP (CAST(? AS INT))
            vendor_metrics_id,
            initials,
            company_name,
            reminder_days,
            active,
            created_at,
            updated_at,
            deleted_at
        FROM edi_vendor_metrics
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
    `

	var out []domains.EDIVendorMetrics
	if err := r.db.Raw(q, limit).Scan(&out).Error; err != nil {
		return nil, err
	}

	return out, nil
}

func (r *EDIVendorMetricsRepositoryDB) UpdateVendorMetricsWithMap(vendorMetricsID string, updates map[string]interface{}) error {
	return r.db.Model(&domains.EDIVendorMetrics{}).
		Where("vendor_metrics_id = ?", vendorMetricsID).
		Updates(updates).
		Error
}

func (r *EDIVendorMetricsRepositoryDB) DeleteVendorMetrics(vendorMetricsID string) error {
	if err := r.db.
		Where("vendor_metrics_id = ?", vendorMetricsID).
		Delete(&domains.EDIVendorMetrics{}).Error; err != nil {

		fmt.Printf("DeleteVendorMetrics error: %v\n", err)
		return err
	}
	return nil
}
