package repositories

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"backend/internal/core/domains"
	ports "backend/internal/core/ports/repositories"
	"backend/internal/pkgs/utils"
)

type UserRepositoryDB struct {
	db *gorm.DB
}

func NewUserRepositoryDB(db *gorm.DB) ports.UserRepository {
	// if err := db.AutoMigrate(&domains.User{}, &domains.LoginVerification{}); err != nil {
	// 	fmt.Printf("failed to auto migrate: %v", err)
	// }
	return &UserRepositoryDB{db: db}
}

func (r *UserRepositoryDB) CreateUserRepository(User *domains.User) error {
	if err := r.db.Create(User).Error; err != nil {
		fmt.Printf("CreateUserRepository error: %v\n", err)
		return err
	}
	return nil
}

func (r *UserRepositoryDB) FindByUsername(username string) (*domains.User, error) {
	var user domains.User

	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryDB) FindUserByGroup(group string) ([]domains.User, error) {
	var users []domains.User

	if err := r.db.Where("[group] = ?", group).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryDB) FindByEmail(email string) (*domains.User, error) {
	var user domains.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryDB) UpdatePasswordByEmail(email, newPassword string) error {
	if err := r.db.Model(&domains.User{}).
		Where("email = ?", email).
		Update("password", newPassword).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryDB) GetUserByID(userID string) (domains.User, error) {
	var user domains.User
	if err := r.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return domains.User{}, err
	}
	return user, nil
}

func (r UserRepositoryDB) GetAllUser() ([]domains.User, error) {
	var reviews []domains.User
	return reviews, r.db.Find(&reviews).Error
}

func (r UserRepositoryDB) UpdateUserWithMap(userID string, updates map[string]interface{}) error {
	return r.db.Model(&domains.User{}).
		Where("user_id = ?", userID).
		Updates(updates).
		Error
}

func (r UserRepositoryDB) GetUserCount() (int64, error) {
	var count int64
	return count, r.db.Model(&domains.User{}).Count(&count).Error
}

func (r *UserRepositoryDB) CreatePasswordResetLink(email, frontendBaseURL, lang string) (string, error) {
	var u domains.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	if strings.EqualFold(u.Status, "disable") {
		return "", nil
	}

	plain, hash, err := utils.NewResetToken()
	if err != nil {
		return "", err
	}
	expires := time.Now().Add(30 * time.Minute)

	returnLink := ""
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			DELETE pr
			FROM password_resets pr
			WHERE pr.user_id IN (SELECT u.user_id FROM users u WHERE u.email = ?)
		`, email).Error; err != nil {
			return fmt.Errorf("cleanup old tokens: %w", err)
		}

		res := tx.Exec(`
			INSERT INTO password_resets (id, user_id, token_hash, expires_at, created_at)
			SELECT NEWID(), u.user_id, ?, ?, SYSUTCDATETIME()
			FROM users u
			WHERE u.email = ?
		`, hash, expires, email)
		if res.Error != nil {
			return fmt.Errorf("insert token: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("insert token: no row inserted")
		}

		if lang == "" {
			lang = "en"
		}
		returnLink = strings.TrimRight(frontendBaseURL, "/") + "/" + lang + "/reset-password?token=" + plain
		return nil
	}); err != nil {
		return "", err
	}

	return returnLink, nil
}

// ยืนยันโทเคน + ตั้งรหัสผ่านใหม่ (แนะนำส่งมาจาก handler)
func (r *UserRepositoryDB) ResetPasswordByToken(plainToken, newHashedPassword string) error {
	hash := utils.HashToken(plainToken)

	var pr domains.PasswordReset
	if err := r.db.Where("token_hash = ?", hash).First(&pr).Error; err != nil {
		return errors.New("invalid token")
	}
	if pr.UsedAt != nil || time.Now().After(pr.ExpiresAt) {
		return errors.New("token expired or used")
	}

	// หา user

	fmt.Println("user_id: ", pr.UserID)
	var u domains.User
	if err := r.db.Where("user_id = ?", pr.UserID).First(&u).Error; err != nil {
		return errors.New("user not found")
	}
	if u.Status == "disable" {
		return errors.New("account disabled")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domains.User{}).
			Where("user_id = ?", u.UserID).
			Update("password", newHashedPassword).Error; err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&domains.PasswordReset{}).
			Where("id = ?", pr.ID).
			Update("used_at", &now).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *UserRepositoryDB) StartLoginWithEmailOTP(email string) (*domains.User, error) {
	var u domains.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if strings.EqualFold(u.Status, "disable") {
		return nil, nil
	}

	plain, hash, err := utils.New6DigitCode()
	if err != nil {
		return nil, err
	}

	fmt.Println("plain : ", plain)

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DELETE lv FROM login_verifications AS lv WHERE lv.user_id = ?`, u.UserID).Error; err != nil {
			return fmt.Errorf("cleanup old login codes: %w", err)
		}

		if err := tx.Exec(`
			INSERT INTO login_verifications (id, user_id, code_hash, expires_at, created_at, attempt_count)
			SELECT NEWID(), ?, ?, DATEADD(MINUTE, 1, t.created_at), t.created_at, 0
			FROM (SELECT SYSUTCDATETIME() AS created_at) AS t
		`, u.UserID, hash).Error; err != nil {
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

func (r *UserRepositoryDB) VerifyLoginCode(email, plainCode string) (*domains.User, error) {
	var u domains.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, errors.New("ไม่พบผู้ใช้")
	}

	if strings.EqualFold(u.Status, "disable") {
		return nil, errors.New("บัญชีนี้ถูกปิดการใช้งาน")
	}

	codeHash := utils.HashCode(plainCode)
	var lv domains.LoginVerification

	if err := r.db.
		Where("user_id = ? AND code_hash = ? AND consumed_at IS NULL AND expires_at > ?",
			u.UserID, codeHash, time.Now()).
		Order("created_at DESC").
		First(&lv).Error; err != nil {
		return nil, errors.New("โค้ดไม่ถูกต้อง หรือหมดอายุ")
	}

	if lv.AttemptCount >= 5 {
		return nil, errors.New("พยายามเกินจำนวนครั้งที่กำหนด")
	}

	now := time.Now()
	if err := r.db.Model(&domains.LoginVerification{}).
		Where("id = ?", lv.ID).
		Updates(map[string]any{
			"consumed_at":   &now,
			"attempt_count": gorm.Expr("attempt_count + 1"),
		}).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepositoryDB) FindActiveEmailsByGroup(group string) ([]string, error) {
	var emails []string

	err := r.db.
		Model(&domains.User{}).
		Select("DISTINCT email").
		Where("[group] = ?", group).
		Where("status <> ?", "disable").
		Where("email IS NOT NULL AND LTRIM(RTRIM(email)) <> ''").
		Where("deleted_at IS NULL").
		Pluck("email", &emails).Error

	if err != nil {
		return nil, err
	}

	return emails, nil
}
