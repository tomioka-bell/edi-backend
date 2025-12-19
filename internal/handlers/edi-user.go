package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	database "backend/external/db"
	"backend/internal/core/models"
	services "backend/internal/core/ports/services"
	"backend/internal/pkgs/utils"
)

type UserHandler struct {
	UserSrv services.UserService
}

func NewUserHandler(insSrv services.UserService) *UserHandler {
	return &UserHandler{UserSrv: insSrv}
}

func (h *UserHandler) LoginDBHandler(c *fiber.Ctx) error {
	var loginData models.LoginResp
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลไม่ถูกต้อง",
		})
	}

	token, err := h.UserSrv.SignIn(loginData)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "PPR_",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		// SameSite: "Lax",
		Secure:   false,
		SameSite: "none",
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Login": "เข้าสู่ระบบสำเร็จ",
	})
}

func (h *UserHandler) LoginHandler(c *fiber.Ctx) error {
	var loginData models.LoginResp
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลไม่ถูกต้อง",
		})
	}

	token, err := h.UserSrv.SignIn(loginData)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "เข้าสู่ระบบสำเร็จ",
		"token":   token,
	})
}

func (h *UserHandler) LogoutDBHandler(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "PPR_",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "none",
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Logout": "ออกจากระบบสำเร็จ",
	})
}

func (h *UserHandler) GetProfileCookie(c *fiber.Ctx) error {
	cookie := c.Cookies("ACG_")

	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Cookie not found",
		})
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token user")
		}
		return []byte(viper.GetString("token.secret_key")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token " + err.Error(),
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)

		result, err := h.UserSrv.GetProfileByCookieId(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"error":  err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"result": result,
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Get profile failed",
	})
}

// func (h *UserHandler) GetProfileHandler(c *fiber.Ctx) error {
// 	authHeader := c.Get("Authorization")
// 	if authHeader == "" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Unauthorized access",
// 		})
// 	}

// 	tokenString := ""
// 	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
// 		tokenString = authHeader[7:]
// 	} else {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Invalid authorization header",
// 		})
// 	}

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("invalid token")
// 		}
// 		return []byte(viper.GetString("token.secret_key")), nil
// 	})

// 	if err != nil || !token.Valid {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Invalid token",
// 		})
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Unauthorized access",
// 		})
// 	}

// 	userID := claims["user_id"].(string)

// 	result, err := h.UserSrv.GetProfileByCookieId(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status": "error",
// 			"error":  err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"status": "success",
// 		"result": result,
// 	})
// }

func (h *UserHandler) GetProfileHandler(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "token not found",
		})
	}

	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(viper.GetString("token.secret_key")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "StatusUnauthorized access",
		})
	}

	userID := claims["user_id"].(string)

	result, err := h.UserSrv.GetEDIPrincipalUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"result": result,
	})
}

func (h *UserHandler) GetUserByIDHandler(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "กรุณาระบุ user_id"})
	}

	companyContents, err := h.UserSrv.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(companyContents)
}

func (h *UserHandler) CreateUserHandler(c *fiber.Ctx) error {
	var req models.UserResp

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if h.UserSrv == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Service is not available",
		})
	}

	profileURL, err := utils.UploadFileFromForm(
		c,
		"profile",
		"uploads/profile",
		"/uploads/profile",
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to upload profile image: %v", err),
		})
	}

	if profileURL != nil {
		req.Profile = *profileURL
	}

	if err := h.UserSrv.SendInitialPasswordSetup(req.Email, "en"); err != nil {
		fmt.Printf("password reset failed for %s: %v\n", req.Email, err)
	}

	if err := h.UserSrv.CreateUserService(req); err != nil {
		log.Println("Error creating User:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create User",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

// func (h *UserHandler) CreateUserHandler(c *fiber.Ctx) error {
// 	var req models.UserResp

// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request payload",
// 		})
// 	}

// 	if h.UserSrv == nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Service is not available",
// 		})
// 	}

// 	err := h.UserSrv.CreateUserService(req)
// 	if err != nil {
// 		log.Println("Error creating  User:", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to create  User",
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "User  created successfully",
// 	})
// }

func (h *UserHandler) GetAllUserHandler(c *fiber.Ctx) error {
	Companies, err := h.UserSrv.GetAllUserSevice()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to user ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Companies)
}

func (h *UserHandler) CheckAuth(c *fiber.Ctx) error {
	cookie := c.Cookies("EDI-System")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized access",
		})
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(viper.GetString("token.secret_key")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"]
		c.Locals("user_id", userID)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized access",
	})
}

func (h *UserHandler) UpdateUserHandler(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	updates := map[string]interface{}{}

	// ===== TEXT FIELDS =====
	if v := c.FormValue("firstname"); v != "" {
		updates["firstname"] = v
	}
	if v := c.FormValue("lastname"); v != "" {
		updates["lastname"] = v
	}
	if v := c.FormValue("username"); v != "" {
		updates["username"] = v
	}
	if v := c.FormValue("email"); v != "" {
		updates["email"] = v
	}
	if v := c.FormValue("password"); v != "" {
		updates["password"] = utils.Encode(v)
	}
	if v := c.FormValue("status"); v != "" {
		updates["status"] = v
	}
	if v := c.FormValue("role"); v != "" {
		updates["role"] = v
	}
	if v := c.FormValue("company_code"); v != "" {
		updates["company_code"] = v
	}

	// ===== BOOLEAN =====
	if v := c.FormValue("login_without_otp"); v != "" {
		updates["login_without_otp"] = v == "true" || v == "1"
	}

	// ===== FILE UPLOAD =====
	profileURL, err := utils.UploadFileFromForm(
		c,
		"profile",          // key in form-data
		"uploads/profile",  // disk path
		"/uploads/profile", // public path
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to upload profile image: %v", err),
		})
	}

	if profileURL != nil {
		updates["profile"] = *profileURL
	}

	// ===== VALIDATION =====
	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "no fields to update",
		})
	}

	// ===== UPDATE =====
	if err := h.UserSrv.UpdateUserWithMapService(userID, updates); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func (h *UserHandler) GetUserCountHandler(c *fiber.Ctx) error {
	count, err := h.UserSrv.GetUserCountService()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve user count",
		})
	}

	return c.JSON(fiber.Map{
		"total_users": count,
	})
}

// POST http://127.0.0.1:1970/api/user/request-password-reset  { "email": "user@example.com" }
func (h *UserHandler) RequestPasswordReset(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid email"})
	}

	if err := h.UserSrv.RequestPasswordReset(req.Email, "en"); err != nil {
		fmt.Printf("password reset failed for %s: %v\n", req.Email, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If your email exists, you will receive a reset link.",
	})
}

// POST http://127.0.0.1:1970/api/user/reset-password  { "token": "...", "new_password": "..." }
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid json"})
	}
	if req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "missing token"})
	}
	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "password must be at least 8 characters"})
	}

	if err := h.UserSrv.PerformPasswordResetByToken(req.Token, req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "password updated"})
}

// เริ่มล็อกอิน: ตรวจ email/password แล้ว "ส่งโค้ด 6 หลักทางอีเมล"
func (h *UserHandler) StartLogin(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.UserSrv.StartLogin(req.Email, req.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "We have sent a 6-digit verification code to your email.",
	})
}

// ยืนยันโค้ด: รับ email + code แล้วออก JWT
func (h *UserHandler) VerifyLoginCode(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := c.BodyParser(&req); err != nil || req.Email == "" || len(req.Code) != 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	token, err := h.UserSrv.CompleteLogin(req.Email, req.Code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ล็อกอินสำเร็จ",
		"token":   token,
	})
}

func (h *UserHandler) LoginLDAPHandler(c *fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request payload",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "username or password is required",
		})
	}

	ok, msg := database.AuthenticateUserDomainLogin(req.Username, req.Password)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.LoginResponse{
		Success: true,
		Message: "authenticated via LDAP",
		// Token:   token,
	})
}

func GetUserInfo(l *ldap.Conn, targetUsername string) (*ldap.Entry, error) {

	baseDN := os.Getenv("LDAP_BASE_DN") // เช่น "DC=prospira,DC=com"
	if baseDN == "" {
		return nil, fmt.Errorf("LDAP_BASE_DN not defined")
	}

	// Filter ที่ใช้ค้นหา user
	filter := fmt.Sprintf("(sAMAccountName=%s)", ldap.EscapeFilter(targetUsername))

	searchReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter,
		[]string{"dn", "sAMAccountName", "displayName", "mail", "userPrincipalName", "memberOf"},
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("LDAP search failed: %v", err)
	}

	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("user not found: %s", targetUsername)
	}

	return sr.Entries[0], nil
}

func (h *UserHandler) GetEDIPrincipalUserByGroupHandler(c *fiber.Ctx) error {
	group := c.Query("group")
	user, err := h.UserSrv.GetEDIPrincipalUserByGroupService(group)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve user by group",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) GetActiveEmailsByGroupHandler(c *fiber.Ctx) error {
	group := c.Query("group")
	user, err := h.UserSrv.GetActiveEmailsByGroup(c.Context(), group)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve user by group",
		})
	}

	if err := h.UserSrv.RequestPasswordResetMany(user, "en"); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve user by group",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdatePrincipalWithMapHandler(c *fiber.Ctx) error {
	idStr := c.Params("edi_principal_id")

	if _, err := uuid.Parse(idStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}

	var req models.EDIPrincipalUserUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	updates := map[string]interface{}{}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	updates["login_without_otp"] = req.LoginWithoutOTP

	if err := h.UserSrv.UpdatePrincipalWithMapService(idStr, updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Principal user updated successfully",
	})
}

func (h *UserHandler) GetEDIPrincipalUserByCompanyHandler(c *fiber.Ctx) error {
	company := c.Query("company")
	user, err := h.UserSrv.GetEDIPrincipalUserByCompanyService(company)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve user by company",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) GetEDIUserByGroupHandler(c *fiber.Ctx) error {
	group := c.Query("company")
	user, err := h.UserSrv.GetEDIUserByGroupService(group)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve user by group",
		})
	}

	return c.JSON(user)
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) LoginUserHandler(c *fiber.Ctx) error {
	var req LoginUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request payload",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email and password are required",
		})
	}

	token, err := h.UserSrv.StartLoginWithEmailEmployeeOTPService(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to start login flow: " + err.Error(),
		})
	}

	if token != "" {
		return c.JSON(fiber.Map{
			"success":    true,
			"message":    "Login success without OTP.",
			"need_otp":   false,
			"token":      token,
			"login_type": "direct",
		})
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    "We have sent a 6-digit verification code to your email.",
		"need_otp":   true,
		"login_type": "otp",
	})
}
