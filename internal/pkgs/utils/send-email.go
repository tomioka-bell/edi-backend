package utils

import (
	"fmt"
	"net/smtp"

	"github.com/spf13/viper"
)

func SendTestMail(to []string) {
	from := viper.GetString("ai_key.email_key")
	password := viper.GetString("ai_key.email_pass")

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	message := []byte("Subject: Test Email\n\nสวัสดี นี่คือข้อความทดสอบ!")

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("ส่งอีเมลสำเร็จ")
}
