package services

import (
	"backend/internal/core/domains"
	"backend/internal/core/models"
	"backend/internal/repositories"
	"fmt"

	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
)

type EDIReadStatService struct {
	repo *repositories.EDIReadStatRepositoryDB
}

func NewEDIReadStatService(r *repositories.EDIReadStatRepositoryDB) *EDIReadStatService {
	return &EDIReadStatService{repo: r}
}

func (s *EDIReadStatService) TrackReadService(req models.EDIReadStatReq) error {
	u := uuid.New()

	var newID mssql.UniqueIdentifier
	copy(newID[:], u[:])

	domainEDIForecastVersionStatusLog := domains.EDIReadStat{
		EDIReadID:  newID,
		Number:     req.Number,
		Type:       req.Type,
		VendorCode: req.VendorCode,
		Read:       req.Read,
		ReadAt:     req.ReadAt,
		CreatedAt:  req.CreatedAt,
	}

	if err := s.repo.TrackRead(&domainEDIForecastVersionStatusLog); err != nil {
		return fmt.Errorf("failed to track read: %w", err)
	}
	return nil
}

func (s *EDIReadStatService) GetReadStatByVendorCodeService(vendorCode string) (models.EDIReadStatResp, error) {
	employee, err := s.repo.GetReadStatByVendorCode(vendorCode)
	if err != nil {
		return models.EDIReadStatResp{}, err
	}

	readStatView := models.EDIReadStatResp{
		EDIReadID:  employee.EDIReadID,
		Number:     employee.Number,
		Type:       employee.Type,
		VendorCode: employee.VendorCode,
		Read:       employee.Read,
		ReadAt:     employee.ReadAt,
		CreatedAt:  employee.CreatedAt,
		UpdatedAt:  employee.UpdatedAt,
		DeletedAt:  employee.DeletedAt,
	}

	return readStatView, nil
}
