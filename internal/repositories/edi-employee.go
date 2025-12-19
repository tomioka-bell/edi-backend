package repositories

import (
	"backend/internal/core/domains"
	ports "backend/internal/core/ports/repositories"
	"backend/internal/pkgs/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type EmployeeRepositoryDB struct {
	db *gorm.DB
}

func NewEmployeeRepositoryDB(db *gorm.DB) ports.EmployeeRepository {
	// if err := db.AutoMigrate(&domains.LoginVerificationEmployee{}); err != nil {
	// 	fmt.Printf("failed to auto migrate: %v", err)
	// }
	return &EmployeeRepositoryDB{db: db}
}

func (r *EmployeeRepositoryDB) CreateEmployeeRepository(Employee *domains.EmployeeView) error {
	if err := r.db.Create(Employee).Error; err != nil {
		fmt.Printf("CreateEmployeeRepository error: %v\n", err)
		return err
	}
	return nil
}

func (r *EmployeeRepositoryDB) GetEmployeeByADLogon(adUserLogon string) (*domains.EmployeeView, error) {
	var emp domains.EmployeeView
	if err := r.db.
		Where("AD_UserLogon = ?", adUserLogon).
		First(&emp).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepositoryDB) GetEmployeeByFullNameEn(fullNameEn string) (*domains.EmployeeView, error) {
	var emp domains.EmployeeView
	if err := r.db.
		Where("UHR_FullNameEn = ?", fullNameEn).
		First(&emp).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepositoryDB) FindEmployeeByAccount(account string) (*domains.EmployeeView, error) {
	var user domains.EmployeeView

	if err := r.db.Where("AD_UserLogon = ?", account).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *EmployeeRepositoryDB) FindPrincipalByEmail(email string) (*domains.EDI_Principal, error) {
	var p domains.EDI_Principal
	if err := r.db.Where("email = ?", email).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *EmployeeRepositoryDB) FindPrincipalByUsername(username string) (*domains.EDI_Principal, error) {
	var p domains.EDI_Principal
	if err := r.db.Where("username = ?", username).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *EmployeeRepositoryDB) CreateEDIPrincipalRepository(EDIUser *domains.EDI_Principal) error {
	if err := r.db.Create(EDIUser).Error; err != nil {
		fmt.Printf("CreateEDIPrincipalRepository error: %v\n", err)
		return err
	}
	return nil
}

func (r *EmployeeRepositoryDB) StartLoginWithEmailEmployeeOTP(email string) (*domains.EmployeeView, error) {
	var u domains.EmployeeView
	if err := r.db.Where("AD_Mail = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if strings.EqualFold(u.AD_AccountStatus, "DISABLE") {
		return nil, nil
	}

	plain, hash, err := utils.New6DigitCode()
	if err != nil {
		return nil, err
	}

	fmt.Println("plain", plain)

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			`DELETE FROM login_verifications_employee WHERE user_login = ?`,
			u.AD_UserLogon,
		).Error; err != nil {
			return fmt.Errorf("cleanup old login codes: %w", err)
		}

		if err := tx.Exec(`
			INSERT INTO login_verifications_employee (
				id, user_login, emp_code, email, code_hash, 
				expires_at, created_at, attempt_count
			)
			SELECT 
				NEWID(), ?, ?, ?, ?, 
				DATEADD(MINUTE, 1, t.created_at), t.created_at, 0
			FROM (SELECT SYSUTCDATETIME() AS created_at) AS t
		`,
			u.AD_UserLogon,
			u.UHR_EmpCode,
			u.AD_Mail,
			hash,
		).Error; err != nil {
			return fmt.Errorf("insert login code: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	u.TempOTP = plain
	return &u, nil
}

func (r *EmployeeRepositoryDB) VerifyLoginCodeEmployee(
	login string,
	plainCode string,
) (*domains.LoginVerificationEmployee, error) {

	codeHash := utils.HashCode(plainCode)

	var lv domains.LoginVerificationEmployee

	// หาโค้ดที่:
	// - user_login ตรงกับ login
	// - code_hash ตรงกับโค้ดที่กรอก
	// - ยังไม่ถูกใช้ (consumed_at IS NULL)
	// - ยังไม่หมดอายุ (expires_at > now)
	if err := r.db.
		Where(`
            user_login = ? 
            AND code_hash = ? 
            AND consumed_at IS NULL 
            AND expires_at > ?
        `, login, codeHash, time.Now()).
		Order("created_at DESC").
		First(&lv).Error; err != nil {

		return nil, errors.New("โค้ดไม่ถูกต้อง หรือหมดอายุ")
	}

	// เช็คจำนวนความพยายาม
	if lv.AttemptCount >= 5 {
		return nil, errors.New("พยายามเกินจำนวนครั้งที่กำหนด")
	}

	// mark ว่าโค้ดนี้ถูกใช้แล้ว + เพิ่ม attempt_count
	now := time.Now()
	if err := r.db.Model(&domains.LoginVerificationEmployee{}).
		Where("id = ?", lv.ID).
		Updates(map[string]any{
			"consumed_at":   &now,
			"attempt_count": gorm.Expr("attempt_count + 1"),
		}).Error; err != nil {
		return nil, err
	}

	return &lv, nil
}

// func (r *EmployeeRepositoryDB) VerifyLoginCodeEmployee(email, plainCode string) (*domains.EmployeeView, error) {
// 	var u domains.EmployeeView
// 	if err := r.db.Where("AD_UserLogon = ?", email).First(&u).Error; err != nil {
// 		return nil, errors.New("ไม่พบผู้ใช้")
// 	}

// 	if strings.EqualFold(u.AD_AccountStatus, "DISABLE") {
// 		return nil, errors.New("บัญชีนี้ถูกปิดการใช้งาน")
// 	}

// 	codeHash := utils.HashCode(plainCode)
// 	var lv domains.LoginVerificationEmployee

// 	if err := r.db.Where("user_login = ? AND code_hash = ? AND consumed_at IS NULL AND expires_at > ?",
// 		u.AD_UserLogon, codeHash, time.Now()).
// 		Order("created_at DESC").
// 		First(&lv).Error; err != nil {
// 		return nil, errors.New("โค้ดไม่ถูกต้อง หรือหมดอายุ")
// 	}

// 	if lv.AttemptCount >= 5 {
// 		return nil, errors.New("พยายามเกินจำนวนครั้งที่กำหนด")
// 	}

// 	now := time.Now()
// 	if err := r.db.Model(&domains.LoginVerificationEmployee{}).
// 		Where("id = ?", lv.ID).
// 		Updates(map[string]any{
// 			"consumed_at":   &now,
// 			"attempt_count": gorm.Expr("attempt_count + 1"),
// 		}).Error; err != nil {
// 		return nil, err
// 	}

// 	return &u, nil
// }

func (r *EmployeeRepositoryDB) UpdatePrincipalWithMap(principalID string, updates map[string]interface{}) error {
	return r.db.Model(&domains.EDI_Principal{}).
		Where("edi_principal_id = ?", principalID).
		Updates(updates).
		Error
}
