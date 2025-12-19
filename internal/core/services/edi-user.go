package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"

	"backend/internal/core/domains"
	"backend/internal/core/models"
	ports "backend/internal/core/ports/repositories"
	servicesports "backend/internal/core/ports/services"
	"backend/internal/pkgs/logs"
	"backend/internal/pkgs/mailer"
	"backend/internal/pkgs/utils"
)

type userService struct {
	userisrRepo ports.UserRepository
}

func NewUserService(UserisrRepo ports.UserRepository) servicesports.UserService {
	return &userService{userisrRepo: UserisrRepo}
}

func (s *userService) CreateUserService(req models.UserResp) error {
	u := uuid.New()

	var newID mssql.UniqueIdentifier
	copy(newID[:], u[:])

	encodedPassword := utils.Encode(req.Password)

	domainISR := domains.User{
		UserID:          newID,
		Firstname:       req.Firstname,
		Lastname:        req.Lastname,
		Username:        req.Firstname + "_" + req.Lastname,
		Email:           req.Email,
		Profile:         req.Profile,
		Group:           req.Group,
		Role:            "VENDER",
		Password:        encodedPassword,
		Status:          "active",
		CompanyCode:     req.CompanyCode,
		LoginWithoutOTP: true,
	}

	if s.userisrRepo == nil {
		return fmt.Errorf("user repository is not initialized")
	}
	if err := s.userisrRepo.CreateUserRepository(&domainISR); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *userService) GetEDIUserByGroupService(group string) ([]models.UserReqAll, error) {
	users, err := s.userisrRepo.FindUserByGroup(group)
	if err != nil {
		return nil, err
	}

	res := make([]models.UserReqAll, 0, len(users))

	for _, u := range users {
		res = append(res, models.UserReqAll{
			UserID:          u.UserID,
			Firstname:       u.Firstname,
			Lastname:        u.Lastname,
			Username:        u.Username,
			Email:           u.Email,
			Status:          u.Status,
			Group:           u.Group,
			CompanyCode:     u.CompanyCode,
			LoginWithoutOTP: u.LoginWithoutOTP,
		})
	}

	return res, nil
}

func (s *userService) UpdateUserWithMapService(userID string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if password, ok := updates["password"]; ok {
		encodedPassword := utils.Encode(password.(string))
		updates["password"] = encodedPassword
	}

	return s.userisrRepo.UpdateUserWithMap(userID, updates)
}

func (s *userService) SignIn(dto models.LoginResp) (string, error) {
	userData, err := s.userisrRepo.FindByEmail(dto.Email)
	if err != nil {
		return "", errors.New("อีเมลไม่ถูกต้อง")
	}

	if userData == nil {
		return "", errors.New("ไม่พบผู้ใช้")
	}

	if userData.Status == "disable" {
		return "", errors.New("บัญชีนี้ถูกปิดการใช้งานแล้ว")
	}

	if !utils.Compare(dto.Password, userData.Password) {
		return "", errors.New("รหัสผ่านไม่ถูกต้อง")
	}

	jwtSecretKey := []byte(os.Getenv("TOKEN_SECRET_KEY"))
	claims := jwt.MapClaims{
		"user_id":       userData.UserID,
		"username":      userData.Username,
		"firstname":     userData.Firstname,
		"lastname":      userData.Lastname,
		"profile":       userData.Profile,
		"group":         userData.Group,
		"status":        userData.Status,
		"source_system": "APP_EMPLOYEE",
		"company_code":  userData.CompanyCode,
		"iat":           time.Now().Unix(),
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", errors.New("เกิดข้อผิดพลาดในการเซ็นชื่อ JWT")
	}

	return signedToken, nil
}

func (s *userService) UpdatePasswordByEmailService(email, newPassword string) error {
	FindByEmail, err := s.userisrRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	if FindByEmail == nil {
		return errors.New("ไม่พบผู้ใช้")
	}
	if FindByEmail.Status == "disable" {
		return errors.New("บัญชีนี้ถูกปิดการใช้งาน")
	}
	encodedPassword := utils.Encode(newPassword)
	return s.userisrRepo.UpdatePasswordByEmail(email, encodedPassword)
}

func (s *userService) SendInitialPasswordSetup(email, lang string) error {
	frontend := os.Getenv("FRONTEND_BASE_URL")
	if frontend == "" {
		frontend = "http://localhost:5173"
	}

	link, err := s.userisrRepo.CreatePasswordResetLink(email, frontend, lang)
	if err != nil {
		return fmt.Errorf("create reset link: %w", err)
	}
	if link == "" {
		return nil
	}

	if err := mailer.SendCreatePasswordRequest(email, link); err != nil {
		return fmt.Errorf("send reset email: %w", err)
	}

	return nil
}

func (s *userService) RequestPasswordReset(email, lang string) error {
	frontend := os.Getenv("FRONTEND_BASE_URL")
	if frontend == "" {
		frontend = "http://localhost:5173"
	}

	link, err := s.userisrRepo.CreatePasswordResetLink(email, frontend, lang)
	if err != nil {
		return fmt.Errorf("create reset link: %w", err)
	}
	if link == "" {
		return nil
	}

	if err := mailer.SendPasswordResetEmail(email, link); err != nil {
		return fmt.Errorf("send reset email: %w", err)
	}

	return nil
}

func (s *userService) PerformPasswordResetByToken(token, newPassword string) error {
	encoded := utils.Encode(newPassword)
	return s.userisrRepo.ResetPasswordByToken(token, encoded)
}

func (s *userService) RequestPasswordResetMany(emails []string, lang string) error {
	frontend := os.Getenv("FRONTEND_BASE_URL")
	if frontend == "" {
		frontend = "http://localhost:5173"
	}

	uniq := make([]string, 0, len(emails))
	seen := map[string]struct{}{}
	for _, e := range emails {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		key := strings.ToLower(e)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		uniq = append(uniq, e)
	}

	for _, email := range uniq {
		link, err := s.userisrRepo.CreatePasswordResetLink(email, frontend, lang)
		if err != nil {
			return fmt.Errorf("create reset link (%s): %w", email, err)
		}
		if link == "" {
			continue
		}

		if err := mailer.SendPasswordResetEmail(email, link); err != nil {
			return fmt.Errorf("send reset email (%s): %w", email, err)
		}
	}

	return nil
}

func (s *userService) GetProfileByCookieId(userID string) (models.UserReq, error) {
	// ดึงข้อมูล User จาก repository
	user, err := s.userisrRepo.GetUserByID(userID)
	if err != nil {
		return models.UserReq{}, err
	}

	// สร้างโครงสร้าง UserReq
	userReq := models.UserReq{
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Username:    user.Username,
		Email:       user.Email,
		Status:      user.Status,
		RoleName:    user.Role,
		Profile:     user.Profile,
		Group:       user.Group,
		CompanyCode: user.CompanyCode,
	}

	// เตรียมข้อมูล RoleName
	// var roleInfoList []models.RoleInfo
	// for _, role := range user.Role {
	// 	roleInfo := models.RoleInfo{
	// 		Name:        role.RoleName,
	// 		Description: role.RoleDescription,
	// 	}
	// 	roleInfoList = append(roleInfoList, roleInfo)
	// }
	// userReq.RoleName = roleInfoList

	return userReq, nil
}

func (s *userService) GetUserByID(userID string) (models.UserAdminReq, error) {
	// ดึงข้อมูล User จาก repository
	user, err := s.userisrRepo.GetUserByID(userID)
	if err != nil {
		return models.UserAdminReq{}, err
	}

	// สร้างโครงสร้าง UserReq
	userReq := models.UserAdminReq{
		UserID:      user.UserID,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		Status:      user.Status,
		RoleName:    user.Role,
		CompanyCode: user.CompanyCode,
	}

	// เตรียมข้อมูล RoleName
	// var roleInfoList []models.RoleInfo
	// for _, role := range user.Role {
	// 	roleInfo := models.RoleInfo{
	// 		Name: role.RoleName,
	// 	}
	// 	roleInfoList = append(roleInfoList, roleInfo)
	// }
	// userReq.RoleName = roleInfoList

	return userReq, nil
}

func (s *userService) GetAllUserSevice() ([]models.UserReqAll, error) {
	domainRequests, err := s.userisrRepo.GetAllUser()
	if err != nil {
		logs.Error(err)
		return nil, errors.New("failed to get all user")
	}

	var requests []models.UserReqAll

	for _, domainRequest := range domainRequests {
		// if domainRequest.Status == "Disable" {
		// 	continue
		// }

		request := models.UserReqAll{
			UserID:      domainRequest.UserID,
			Firstname:   domainRequest.Firstname,
			Lastname:    domainRequest.Lastname,
			Username:    domainRequest.Username,
			Email:       domainRequest.Email,
			Status:      domainRequest.Status,
			Profile:     domainRequest.Profile,
			Group:       domainRequest.Group,
			CompanyCode: domainRequest.CompanyCode,
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (s *userService) GetUserCountService() (int64, error) {
	count, err := s.userisrRepo.GetUserCount()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// services/user_service_otp.go
func (s *userService) StartLogin(email, password string) error {
	// 1) ตรวจอีเมล/รหัสผ่านก่อน (ยังไม่ออก JWT)
	userData, err := s.userisrRepo.FindByEmail(email)
	if err != nil || userData == nil {
		return errors.New("No user found in the system.")
	}
	if userData.Status == "disable" {
		return errors.New("This account has been deactivated.")
	}
	if !utils.Compare(password, userData.Password) {
		return errors.New("The email or password is incorrect.")
	}

	// 2) ออกโค้ด 6 หลัก + ส่งเมล
	u, err := s.userisrRepo.StartLoginWithEmailOTP(email)
	if err != nil || u == nil {
		// เพื่อความปลอดภัยไม่บอกว่า user ไม่มี
		return errors.New("ไม่สามารถเริ่มการยืนยันได้")
	}

	if err := mailer.SendLoginOTPEmail(email, u.TempOTP); err != nil {
		return fmt.Errorf("send login otp email: %w", err)
	}
	return nil
}

// func (s *userService) CompleteLogin(email, code string) (string, error) {
// 	u, err := s.userisrRepo.VerifyLoginCode(email, code)
// 	if err != nil {
// 		return "", err
// 	}

// 	jwtSecretKey := []byte(os.Getenv("TOKEN_SECRET_KEY"))
// 	claims := jwt.MapClaims{
// 		"user_id":   u.UserID,
// 		"username":  u.Username,
// 		"firstname": u.Firstname,
// 		"lastname":  u.Lastname,
// 		"profile":   u.Profile,
// 		"group":     u.Group,
// 		"status":    u.Status,
// 		"iat":       time.Now().Unix(),
// 		"exp":       time.Now().Add(time.Hour * 24).Unix(),
// 	}
// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	signedToken, err := jwtToken.SignedString(jwtSecretKey)
// 	if err != nil {
// 		return "", errors.New("เกิดข้อผิดพลาดในการเซ็นชื่อ JWT")
// 	}
// 	return signedToken, nil
// }

func isUniqueViolation(err error) bool {
	es := strings.ToLower(err.Error())
	return strings.Contains(es, "unique") || strings.Contains(es, "duplicate")
}

func (s *userService) CompleteLogin(email, code string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	// 1) verify OTP
	user, err := s.userisrRepo.VerifyLoginCode(email, code)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("verify code สำเร็จ แต่ไม่พบข้อมูลผู้ใช้")
	}

	// 2) find principal by email
	userPrincipal, err := s.userisrRepo.FindPrincipalByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// error จริง (ไม่ใช่ not found)
		return "", err
	}

	// 3) create principal if not found or nil
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) || userPrincipal == nil {
		displayName := strings.TrimSpace(user.Firstname + " " + user.Lastname)

		username := user.Username
		if username == "" {
			if at := strings.Index(email, "@"); at > 0 {
				username = email[:at]
			} else {
				username = email
			}
		}

		newPrincipal := &domains.EDI_Principal{
			ExternalID:      strings.ToLower(user.UserID.String()),
			Email:           email,
			Username:        username,
			DisplayName:     displayName,
			Profile:         user.Profile,
			Role:            user.Role,
			Group:           user.Group,
			CompanyCode:     user.CompanyCode,
			SourceSystem:    "APP_USER",
			Status:          "active",
			LoginWithoutOTP: true,
		}

		if err2 := s.userisrRepo.CreateEDIPrincipalRepository(newPrincipal); err2 != nil {
			if isUniqueViolation(err2) {
				newPrincipal.Username = fmt.Sprintf("%s_%d", username, time.Now().Unix())
				if err3 := s.userisrRepo.CreateEDIPrincipalRepository(newPrincipal); err3 != nil {
					return "", fmt.Errorf("create principal (retry): %w", err3)
				}
			} else {
				return "", fmt.Errorf("create principal: %w", err2)
			}
		}

		userPrincipal = newPrincipal
	} else {
		// --------- กรณีพบ principal: ตรวจว่าข้อมูลเปลี่ยนไหม ถ้าเปลี่ยนให้ update ---------
		updates := map[string]interface{}{}
		displayName := strings.TrimSpace(user.Firstname + " " + user.Lastname)

		if userPrincipal.ExternalID != user.UserID.String() {
			updates["external_id"] = user.UserID.String()
			userPrincipal.ExternalID = user.UserID.String()
		}
		if userPrincipal.Email != user.Email {
			updates["email"] = user.Email
			userPrincipal.Email = user.Email
		}
		if userPrincipal.Username != user.Username {
			updates["username"] = user.Username
			userPrincipal.Username = user.Username
		}
		if userPrincipal.DisplayName != displayName {
			updates["display_name"] = displayName
			userPrincipal.DisplayName = displayName
		}
		if userPrincipal.Group != user.Group {
			updates["group"] = user.Group
			userPrincipal.Group = user.Group
		}
		if userPrincipal.Status != user.Status {
			updates["status"] = user.Status
			userPrincipal.Status = user.Status
		}
		if userPrincipal.Profile != user.Profile {
			updates["profile"] = user.Profile
			userPrincipal.Profile = user.Profile
		}
		if userPrincipal.LoginWithoutOTP != user.LoginWithoutOTP {
			updates["login_without_otp"] = user.LoginWithoutOTP
			userPrincipal.LoginWithoutOTP = user.LoginWithoutOTP
		}
		if userPrincipal.CompanyCode != user.CompanyCode {
			updates["company_code"] = user.CompanyCode
			userPrincipal.CompanyCode = user.CompanyCode
		}

		if len(updates) > 0 {
			if err := s.userisrRepo.UpdatePrincipalWithMap(userPrincipal.EDI_PrincipalID.String(), updates); err != nil {
				return "", fmt.Errorf("update principal: %w", err)
			}
		}
	}

	// safety: ถ้ายัง nil อยู่ ให้ error ชัดเจน
	if userPrincipal == nil {
		return "", errors.New("principal is nil after ensure")
	}

	// 4) issue JWT from principal
	secret := os.Getenv("TOKEN_SECRET_KEY")
	if secret == "" {
		return "", errors.New("missing TOKEN_SECRET_KEY")
	}

	jwtSecretKey := []byte(secret)

	claims := jwt.MapClaims{
		"user_id":      userPrincipal.ExternalID,
		"username":     userPrincipal.Username,
		"display_name": userPrincipal.DisplayName,
		"profile":      userPrincipal.Profile,
		"group":        userPrincipal.Group,
		"status":       userPrincipal.Status,
		"company_code": userPrincipal.CompanyCode,
		"iat":          time.Now().Unix(),
		"exp":          time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", errors.New("เกิดข้อผิดพลาดในการเซ็นชื่อ JWT")
	}
	return signedToken, nil
}

func (s *userService) GetActiveEmailsByGroup(ctx context.Context, group string) ([]string, error) {
	if group == "" {
		return nil, errors.New("group is required")
	}

	emails, err := s.userisrRepo.FindActiveEmailsByGroup(group)
	if err != nil {
		return nil, err
	}

	if emails == nil {
		emails = []string{}
	}

	return emails, nil
}
