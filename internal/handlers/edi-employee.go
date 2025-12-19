package handlers

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/clients"
	services "backend/internal/core/ports/services"
)

type EmployeeHandler struct {
	EmployeeSrv services.EmployeeService
}

func NewEmployeeHandler(insSrv services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{EmployeeSrv: insSrv}
}

func (h *EmployeeHandler) GetEmployeeByADLogonHandler(c *fiber.Ctx) error {
	username := c.Query("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "กรุณาระบุ username"})
	}

	employeeView, err := h.EmployeeSrv.GetEmployeeByADLogonService(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(employeeView)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Token     string `json:"token,omitempty"`
	NeedOTP   bool   `json:"need_otp"`
	LoginType string `json:"login_type"`
}

func (h *EmployeeHandler) LoginWithLdapHandler(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request payload",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "username and password are required",
		})
	}

	// 1) LDAP auth + ดึง user_info
	ldapUser, ok, msg := clients.LdapAuthenticate(req.Username, req.Password)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	// 2) เริ่ม flow login / OTP ด้วยข้อมูลจาก LDAP
	token, err := h.EmployeeSrv.StartLoginWithEmailEmployeeOTPService(ldapUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "LDAP ok but failed to start login flow: " + err.Error(),
		})
	}

	// ถ้าไม่ต้องใช้ OTP
	if token != "" {
		return c.JSON(fiber.Map{
			"success":    true,
			"message":    "Login success without OTP.",
			"need_otp":   false,
			"token":      token,
			"login_type": "direct",
		})
	}

	// ต้องกรอก OTP ต่อ
	return c.JSON(fiber.Map{
		"success":    true,
		"message":    "We have sent a 6-digit verification code to your email.",
		"need_otp":   true,
		"login_type": "otp",
	})
}

func (h *EmployeeHandler) VerifyLoginCodeEmployeeHandler(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := c.BodyParser(&req); err != nil || req.Email == "" || len(req.Code) != 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	token, err := h.EmployeeSrv.CompleteLoginEmployee(req.Email, req.Code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ล็อกอินสำเร็จ",
		"token":   token,
	})
}
