package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Init() {
	initEnvLoader()
	initConfigLoader()
	initTimeZone()
}

func initEnvLoader() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables only")
	} else {
		fmt.Println(".env file loaded successfully")
	}
}

func initConfigLoader() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Cannot read config.yml: %v\n", err)
	}
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Printf("Error loading timezone: %v\n", err)
		os.Exit(1)
	}
	time.Local = ict
}
