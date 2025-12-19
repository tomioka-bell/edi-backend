package emailclient

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServiceURL string
	FromName   string
	Timeout    int
	Token      string
}

func LoadConfig() (*Config, error) {
	serviceURL := os.Getenv("EMAIL_SERVICE_URL")
	if serviceURL == "" {
		return nil, fmt.Errorf("EMAIL_SERVICE_URL is required")
	}

	timeout := 10
	if t := os.Getenv("EMAIL_SERVICE_TIMEOUT"); t != "" {
		timeout, _ = strconv.Atoi(t)
	}

	return &Config{
		ServiceURL: serviceURL,
		FromName:   os.Getenv("SMTP_SET_FROM_NAME"),
		Timeout:    timeout,
		Token:      os.Getenv("EMAIL_SERVICE_TOKEN"),
	}, nil
}
