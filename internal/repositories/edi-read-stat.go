package repositories

import (
	"backend/internal/core/domains"
	"fmt"

	"gorm.io/gorm"
)

type EDIReadStatRepositoryDB struct {
	db *gorm.DB
}

func NewEDIReadStatRepositoryDB(db *gorm.DB) *EDIReadStatRepositoryDB {
	// if err := db.AutoMigrate(&domains.EDIReadStat{}); err != nil {
	// 	fmt.Printf("failed to auto migrate: %v", err)
	// }
	return &EDIReadStatRepositoryDB{db: db}
}

func (r *EDIReadStatRepositoryDB) TrackRead(m *domains.EDIReadStat) error {
	const q = `
	INSERT INTO edi_read_stat
		([edi_read_id], [number], [type], [vendor_code], [read], [read_at], [created_at])
	VALUES
		(?, ?, ?, ?, ?, ?, SYSUTCDATETIME());
`
	if err := r.db.Debug().Exec(q,
		m.EDIReadID,
		m.Number,
		m.Type,
		m.VendorCode,
		m.Read,
		m.ReadAt,
	).Error; err != nil {
		fmt.Printf("TrackRead error: %v\n", err)
		return err
	}
	return nil
}

func (r *EDIReadStatRepositoryDB) GetReadStatByVendorCode(vendorCode string) (*domains.EDIReadStat, error) {
	var m domains.EDIReadStat
	if err := r.db.Where("vendor_code = ?", vendorCode).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
