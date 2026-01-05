package emailclient

import (
	"fmt"
	"log"
	"time"
)

var DefaultClient *Client

func InitDefaultClient() error {
	var err error
	DefaultClient, err = NewClientFromEnv()
	if err != nil {
		return err
	}
	log.Println("Email client initialized successfully")
	return nil
}

// SendLoginOTPEmail ส่งอีเมลรหัส OTP สำหรับการล็อกอิน
func SendLoginOTPEmail(toEmail, otp string) error {
	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendLoginOTPEmail(toEmail, otp)
}

func SendModifyForecastVendorEmail(toEmails []string, company, forecastNumber, fileURL, note string) error {
	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendModifyForecastVendorEmail(toEmails, company, forecastNumber, fileURL, note)
}

// SendPasswordResetEmailMany ส่งอีเมลรีเซ็ตรหัสผ่านไปยังหลายอีเมล
func SendPasswordResetEmailMany(toEmails []string, resetLink string) error {
	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendPasswordResetEmailManyPersonalized(toEmails, resetLink)
}

// SendModifyOrderVendorEmail ส่งอีเมลแจ้งแก้ไขคำสั่งซื้อไปยังผู้ขาย
func SendModifyOrderVendorEmail(toEmails []string, company, orderNumber, fileURL, note string) error {
	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendModifyOrderVendorEmail(toEmails, company, orderNumber, fileURL, note)
}

// SendForecastReadReminder ส่งอีเมลแจ้งเตือนให้อ่านสัญญา Forecast
func SendForecastReadReminder(toEmail, numberForecast, vendorName string, periodFrom, targetTime time.Time) error {
	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendForecastReadReminder(toEmail, numberForecast, vendorName, periodFrom, targetTime)
}

// SendPasswordResetEmailPersonalized ส่งอีเมลรีเซ็ตรหัสผ่านแบบระบุชื่อผู้รับ
func SendPasswordResetEmailPersonalized(toEmail, resetLink string) error {
	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendPasswordResetEmail(toEmail, resetLink)
}

// SendCreatePasswordRequestPersonalized ส่งอีเมลสร้างรหัสผ่านใหม่แบบระบุชื่อผู้รับ
func SendCreatePasswordRequestPersonalized(toEmails string, resetLink string) error {
	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendCreatePasswordRequest(toEmails, resetLink)
}

func SendStatusForecastEmployeeEmail(vendorCompany string,
	statusForecast,
	forecastNumber,
	fileURL string,
	note string,
	notificationType string) error {

	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendStatusForecastEmployeeEmail(vendorCompany, statusForecast, forecastNumber, fileURL, note, notificationType)

}

func SendStatusOrderEmployeeEmail(vendorCompany string,
	statusForecast,
	forecastNumber,
	fileURL string,
	note string,
	notificationType string) error {

	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendStatusOrderEmployeeEmail(vendorCompany, statusForecast, forecastNumber, fileURL, note, notificationType)

}

func SendStatusInvoiceEmployeeEmail(vendorCompany string,
	statusForecast,
	forecastNumber,
	fileURL string,
	note string,
	notificationType string) error {

	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendStatusInvoiceEmployeeEmail(vendorCompany, statusForecast, forecastNumber, fileURL, note, notificationType)

}

func SendModifyInvoiceVendorEmail(
	vendorCompany string,
	invoiceNumber string,
	fileURL string,
	note string,
	notificationType string,
) error {

	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendModifyInvoiceVendorEmail(vendorCompany, invoiceNumber, fileURL, note, notificationType)
}

// SendStatusForecastVendorEmail ส่งอีเมลสถานะ Forecast ไปยังผู้ขาย
func SendStatusForecastVendorEmail(toEmails []string, statusForecast, company, forecastNumber, fileURL, note string) error {
	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendStatusForecastVendorEmail(toEmails, statusForecast, company, forecastNumber, fileURL, note)
}

// SendStatusOrderVendorEmail ส่งอีเมลสถานะ Order ไปยังผู้ขาย
func SendStatusOrderVendorEmail(toEmails []string, statusOrder, company, orderNumber, fileURL, note string) error {
	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendStatusOrderVendorEmail(toEmails, statusOrder, company, orderNumber, fileURL, note)
}

// SendStatusInvoiceVendorEmail ส่งอีเมลสถานะ Invoice ไปยังผู้ขาย
func SendStatusInvoiceVendorEmail(toEmails []string, statusInvoice, company, invoiceNumber, fileURL, note string) error {
	if DefaultClient == nil {
		return ErrClientNotInitialized
	}
	return DefaultClient.SendStatusInvoiceVendorEmail(toEmails, statusInvoice, company, invoiceNumber, fileURL, note)
}

var ErrClientNotInitialized = fmt.Errorf("email client not initialized, call InitDefaultClient() first")
