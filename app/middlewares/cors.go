package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var NewCorsMiddleware = cors.New(cors.Config{
	AllowOrigins: strings.Join([]string{
		// local dev
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:5173",
		"http://localhost:5174",
		"http://localhost:13002",

		// IP/port ภายใน
		"http://10.0.98.208",
		"http://10.0.98.208:13002",
		"http://10.0.98.208:5757",

		// โดเมนจริง (ไม่มี path)
		"http://external.psth.com",
		"https://external.psth.com",
		"http://prospira.th.com",
		"https://prospira.th.com",
		"http://10.144.1.103",
		"https://edi.prospira.co.th",
	}, ","),

	AllowMethods: strings.Join([]string{
		fiber.MethodGet,
		fiber.MethodPost,
		fiber.MethodPut,
		fiber.MethodPatch,
		fiber.MethodDelete,
		fiber.MethodHead,
		fiber.MethodOptions,
	}, ","),

	AllowHeaders: strings.Join([]string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
	}, ","),

	AllowCredentials: true,
})
