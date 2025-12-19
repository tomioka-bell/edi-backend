package handlers

import (
	database "backend/external/db"
	services "backend/internal/core/ports/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	UserSrv     services.UserService
	EmployeeSrv services.EmployeeService
}

type StartLoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func NewAuthHandler(insSrv services.UserService, empSrv services.EmployeeService) *AuthHandler {
	return &AuthHandler{UserSrv: insSrv, EmployeeSrv: empSrv}
}

func (h *AuthHandler) StartLoginHandler(c *fiber.Ctx) error {
	var req StartLoginRequest
	if err := c.BodyParser(&req); err != nil || req.Identifier == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	id := strings.TrimSpace(req.Identifier)

	// CASE 1: email → user ปกติ
	if strings.Contains(id, "@") {
		if err := h.UserSrv.StartLogin(id, req.Password); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"message":   "We sent you a 6-digit code to your email.",
			"loginType": "user",
		})
	}

	// CASE 2: ไม่มี @ → LDAP employee
	ok, msg := database.AuthenticateUserDomainLogin(id, req.Password)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	// if err := h.EmployeeSrv.StartLoginWithEmailEmployeeOTPService(id); err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"success": false,
	// 		"message": "LDAP ok but failed to start login flow: " + err.Error(),
	// 	})
	// }

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   "We sent you a 6-digit code to your corporate email.",
		"loginType": "employee",
	})
}

type VerifyLoginRequest struct {
	Identifier string `json:"identifier"`
	Code       string `json:"code"`
	LoginType  string `json:"loginType"`
}

func (h *AuthHandler) VerifyLoginHandler(c *fiber.Ctx) error {
	var req VerifyLoginRequest
	if err := c.BodyParser(&req); err != nil || len(req.Code) != 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	id := strings.TrimSpace(req.Identifier)
	loginType := strings.ToLower(strings.TrimSpace(req.LoginType))

	switch loginType {
	case "user":
		token, err := h.UserSrv.CompleteLogin(id, req.Code)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"message": "ล็อกอินสำเร็จ",
			"token":   token,
		})

	case "employee":
		token, err := h.EmployeeSrv.CompleteLoginEmployee(id, req.Code)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"message": "ล็อกอินสำเร็จ",
			"token":   token,
		})

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid login type",
		})
	}
}
