package services

import (
	"backend/internal/core/domains"
	"backend/internal/core/models"
	ports "backend/internal/core/ports/repositories"
	servicesports "backend/internal/core/ports/services"
	"fmt"
	"time"

	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
)

type EDIVendorMetricsService struct {
	repo ports.EDIVendorMetricsRepository
}

func NewEDIVendorMetricsService(repo ports.EDIVendorMetricsRepository) servicesports.EDIVendorMetricsService {
	return &EDIVendorMetricsService{repo: repo}
}

func (s *EDIVendorMetricsService) CreateVendorMetricsService(req models.VendorMetricsResp) error {
	u := uuid.New()

	var newID mssql.UniqueIdentifier
	copy(newID[:], u[:])

	domainISR := domains.EDIVendorMetrics{
		VendorMetricsID: newID,
		Initials:        req.Initials,
		CompanyName:     req.CompanyName,
		ReminderDays:    req.ReminderDays,
		Active:          true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if s.repo == nil {
		return fmt.Errorf("vendor metrics is not initialized")
	}
	if err := s.repo.CreateVendorMetrics(&domainISR); err != nil {
		return fmt.Errorf("failed to create vendor metrics: %w", err)
	}
	return nil
}

func (s *EDIVendorMetricsService) GetVendorMetricsByCompanyService(company string) ([]models.VendorMetricsResp, error) {
	metrics, err := s.repo.GetVendorMetricsByCompany(company)
	if err != nil {
		return nil, err
	}

	res := make([]models.VendorMetricsResp, 0, len(metrics))
	const layout = "2006-01-02 15:04:05"
	for _, m := range metrics {
		res = append(res, models.VendorMetricsResp{
			VendorMetricsID: m.VendorMetricsID,
			Initials:        m.Initials,
			CompanyName:     m.CompanyName,
			ReminderDays:    m.ReminderDays,
			Active:          m.Active,
			CreatedAt:       m.CreatedAt.Format(layout),
			UpdatedAt:       m.UpdatedAt.Format(layout),
		})
	}

	return res, nil
}

func (s *EDIVendorMetricsService) GetAllVendorMetricservice() ([]models.VendorMetricsCompanyResp, error) {
	metrics, err := s.repo.GetAllVendorMetrics()
	if err != nil {
		return nil, err
	}

	res := make([]models.VendorMetricsCompanyResp, 0, len(metrics))
	for _, m := range metrics {
		res = append(res, models.VendorMetricsCompanyResp{
			Initials:    m.Initials,
			CompanyName: m.CompanyName,
		})
	}

	return res, nil
}

func (s *EDIVendorMetricsService) GetAllEDIVendorMetricsTopService(limit int) ([]models.VendorMetricsResp, error) {
	metrics, err := s.repo.GetAllEDIVendorMetricsTop(limit)
	if err != nil {
		return nil, err
	}

	res := make([]models.VendorMetricsResp, 0, len(metrics))
	const layout = "2006-01-02 15:04:05"
	for _, m := range metrics {
		res = append(res, models.VendorMetricsResp{
			VendorMetricsID: m.VendorMetricsID,
			Initials:        m.Initials,
			CompanyName:     m.CompanyName,
			ReminderDays:    m.ReminderDays,
			Active:          m.Active,
			CreatedAt:       m.CreatedAt.Format(layout),
			UpdatedAt:       m.UpdatedAt.Format(layout),
		})
	}

	return res, nil
}

func (s *EDIVendorMetricsService) UpdateVendorMetricsService(id mssql.UniqueIdentifier, updates map[string]interface{}) error {
	if s.repo == nil {
		return fmt.Errorf("vendor metrics is not initialized")
	}
	if err := s.repo.UpdateVendorMetricsWithMap(id.String(), updates); err != nil {
		return fmt.Errorf("failed to update vendor metrics: %w", err)
	}
	return nil
}

func (s *EDIVendorMetricsService) DeleteVendorMetricsService(vendorMetricsID mssql.UniqueIdentifier) error {
	if s.repo == nil {
		return fmt.Errorf("vendor metrics is not initialized")
	}
	if err := s.repo.DeleteVendorMetrics(vendorMetricsID.String()); err != nil {
		return fmt.Errorf("failed to delete vendor metrics: %w", err)
	}
	return nil
}
