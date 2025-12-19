package middlewares

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

var NewLoggerMiddleware = logger.New(logger.Config{
	Next:          nil,
	Done:          nil,
	Format:        "[${time}] | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	TimeFormat:    "2006-01-02 15:04:05",
	TimeZone:      "Local",
	TimeInterval:  500 * time.Millisecond,
	Output:        os.Stdout,
	DisableColors: false,
})
