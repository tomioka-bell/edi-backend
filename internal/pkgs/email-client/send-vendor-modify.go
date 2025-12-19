package emailclient

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (c *Client) SendModifyForecastVendorEmail(
	toEmails []string,
	company,
	forecastNumber,
	fileURL string,
	note string,
) error {

	fmt.Println("📧 Preparing to send Modify Forecast Vendor Email to", toEmails)
	// Sanitize emails
	cleaned := sanitizeAndDeduplicateEmails(toEmails)
	if len(cleaned) == 0 {
		return fmt.Errorf("toEmails is empty")
	}

	// สร้าง link
	frontend := strings.TrimRight(strings.TrimSpace(os.Getenv("FRONTEND_BASE_URL")), "/")
	if frontend == "" {
		frontend = "http://localhost:5173"
	}
	link := fmt.Sprintf("%s/en/forecast-form/%s", frontend, url.PathEscape(forecastNumber))

	// จัดการ remark
	remark := strings.TrimSpace(note)
	if remark == "" {
		remark = "-"
	}

	// Subject
	subject := fmt.Sprintf("[Prospira] Forecast Update Notification – %s", forecastNumber)

	// เวลา
	nowUTC := time.Now().UTC().Format("2006-01-02 15:04:05 UTC")

	// สร้าง HTML
	htmlBody := fmt.Sprintf(`<!doctype html>
		<html>
		<head>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width,initial-scale=1"/>
		<title>Forecast Update - %s</title>
		</head>

		<body style="margin:0;padding:24px;background:linear-gradient(to bottom, #08a4b8, #000000);
		font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#111827">

		<table role="presentation" width="100%%" cellspacing="0" cellpadding="0"
			style="max-width:640px;margin:auto;background:#ffffff;border-radius:14px;
			box-shadow:0 10px 24px rgba(0,0,0,0.12);overflow:hidden">

			<tr>
			<td style="padding:24px 26px;background:#F5F5F5;color:#ffffff;border-bottom:1px solid rgba(255,255,255,0.15)">
			<div style="font-size:18px;font-weight:700;margin:0;line-height:1.2;color:#111827;">
			Forecast Update Notification
			</div>
			<div style="margin-top:6px;font-size:13px;opacity:0.95;color:#111827;">
			Prospira (Thailand) Co., Ltd.
			</div>

			</td>
			</tr>

			<tr>
			<td style="padding:  26px">
				<p style="margin:0 0 14px;font-size:14px;line-height:1.7;color:#111827">
				Dear <strong>%s</strong>,
				</p>

				<p style="margin:0 0 14px;font-size:14px;line-height: 1.7;color:#374151">
				This message is to inform you that the forecast listed below has been <strong>updated</strong> in our system.
				Kindly review the details at your earliest convenience. 
				</p>

				<table role="presentation" cellspacing="0" cellpadding="0"
				style="width:100%%;margin:16px 0 18px;border:1px solid #e5e7eb;border-radius:12px;overflow:hidden">
				<tr>
					<td style="padding: 10px 12px;background:#f9fafb;width: 160px;font-size:13px;color:#111827"><strong>Company</strong></td>
					<td style="padding:10px 12px;font-size:13px;color:#374151">%s</td>
				</tr>
				<tr>
					<td style="padding: 10px 12px;background:#f9fafb;font-size: 13px;color:#111827"><strong>Forecast No. </strong></td>
					<td style="padding:10px 12px;font-size:13px;color:#374151">%s</td>
				</tr>
				<tr>
					<td style="padding:10px 12px;background:#f9fafb;font-size:13px;color:#111827"><strong>Updated At</strong></td>
					<td style="padding: 10px 12px;font-size:13px;color:#374151">%s</td>
				</tr>
				</table>

				<div style="text-align:center;margin:18px 0 22px">
				<a href="%s" target="_blank" rel="noopener noreferrer"
					style="display:inline-block;font-size:14px;font-weight:700;
					padding:  12px 28px;border-radius:12px;background:#08a4b8;color:#ffffff;
					text-decoration: none;box-shadow:0 6px 14px rgba(0,0,0,0.18)">
					View Forecast Details
				</a>
				</div>

				<div style="margin-top:8px">
				<div style="margin:0 0 6px;font-size:13px;color:#111827;font-weight:600">
					Remarks
				</div>
				<div style="margin: 0;font-size: 12px;color:#6b7280;line-height:1.7;white-space:pre-line">
					%s
				</div>
				</div>

				<p style="margin: 22px 0 0;font-size:12px;color:#6b7280;line-height:1.7">
				Best regards,<br/>
				<strong>Prospira (Thailand) Co., Ltd.</strong>
				</p>
			</td>
			</tr>
		</table>

		<p style="text-align:center;margin: 14px 0 0;font-size:11px;color:#9ca3af;line-height:1.6">
			This is an automatically generated email. Please do not reply to this message.
		</p>
		</body>
		</html>`,
		forecastNumber, // <title>
		company,        // Dear %s
		company,        // Company (table)
		forecastNumber, // Forecast No.
		nowUTC,         // Updated At
		link,           // Button link
		remark,         // Remarks
	)

	req := &EmailRequest{
		ToEmails: cleaned,
		Subject:  subject,
		HTMLBody: htmlBody,
	}

	// Handle attachment
	attachPath := ""
	if strings.TrimSpace(fileURL) != "" {
		attachPath = "." + strings.TrimSpace(fileURL)
	}

	var err error
	if attachPath != "" {
		// Send with attachments using form-data endpoint
		_, err = c.SendMicroserviceWithAttachments(req, attachPath)
	} else {
		// Send without attachments using JSON endpoint
		_, err = c.SendMicroservice(req)
	}
	return err
}

// readFileAsAttachment อ่านไฟล์และแปลงเป็น base64
func readFileAsAttachment(filePath string) (EmailAttachment, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return EmailAttachment{}, fmt.Errorf("failed to read file: %w", err)
	}

	filename := filepath.Base(filePath)
	base64Data := base64.StdEncoding.EncodeToString(data)

	contentType := "application/octet-stream"
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".pdf":
		contentType = "application/pdf"
	case ".xlsx", ".xls":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".docx", ".doc":
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".csv":
		contentType = "text/csv"
	}

	return EmailAttachment{
		Filename:    filename,
		ContentType: contentType,
		Base64Data:  base64Data,
	}, nil
}

func (c *Client) SendModifyOrderVendorEmail(
	toEmails []string,
	company,
	orderNumber,
	fileURL string,
	note string,
) error {
	// Sanitize emails
	cleaned := sanitizeAndDeduplicateEmails(toEmails)
	if len(cleaned) == 0 {
		return fmt.Errorf("toEmails is empty")
	}

	frontend := strings.TrimRight(strings.TrimSpace(os.Getenv("FRONTEND_BASE_URL")), "/")
	if frontend == "" {
		frontend = "http://localhost:5173"
	}
	link := fmt.Sprintf("%s/en/order-form/%s", frontend, url.PathEscape(orderNumber))

	// จัดการ remark
	remark := strings.TrimSpace(note)
	if remark == "" {
		remark = "-"
	}

	// Subject
	subject := fmt.Sprintf("[Prospira] Order Update Notification – %s", orderNumber)

	// เวลา
	nowUTC := time.Now().UTC().Format("2006-01-02 15:04:05 UTC")

	// สร้าง HTML
	htmlBody := fmt.Sprintf(`<!doctype html>
		<html>
		<head>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width,initial-scale=1"/>
		<title>Order Update - %s</title>
		</head>

		<body style="margin:0;padding:24px;background:linear-gradient(to bottom, #08a4b8, #000000);
		font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#111827">

		<table role="presentation" width="100%%" cellspacing="0" cellpadding="0"
			style="max-width:640px;margin:auto;background:#ffffff;border-radius:14px;
			box-shadow:0 10px 24px rgba(0,0,0,0.12);overflow:hidden">

			<tr>
			<td style="padding:24px 26px;background:#F5F5F5;color:#ffffff;border-bottom:1px solid rgba(255,255,255,0.15)">
			<div style="font-size:18px;font-weight:700;margin:0;line-height:1.2;color:#111827;">
			Order Update Notification
			</div>
			<div style="margin-top:6px;font-size:13px;opacity:0.95;color:#111827;">
			Prospira (Thailand) Co., Ltd.
			</div>

			</td>
			</tr>

			<tr>
			<td style="padding:26px">
				<p style="margin:0 0 14px;font-size:14px;line-height:1.7;color:#111827">
				Dear <strong>%s</strong>,
				</p>

				<p style="margin:0 0 14px;font-size:14px;line-height:1.7;color:#374151">
				This message is to inform you that the order listed below has been <strong>updated</strong> in our system.
				Kindly review the details at your earliest convenience.
				</p>

				<table role="presentation" cellspacing="0" cellpadding="0"
				style="width:100%%;margin:16px 0 18px;border:1px solid #e5e7eb;border-radius:12px;overflow:hidden">
				<tr>
					<td style="padding:10px 12px;background:#f9fafb;width:160px;font-size:13px;color:#111827"><strong>Company</strong></td>
					<td style="padding:10px 12px;font-size:13px;color:#374151">%s</td>
				</tr>
				<tr>
					<td style="padding:10px 12px;background:#f9fafb;font-size:13px;color:#111827"><strong>Order No.</strong></td>
					<td style="padding:10px 12px;font-size:13px;color:#374151">%s</td>
				</tr>
				<tr>
					<td style="padding:10px 12px;background:#f9fafb;font-size:13px;color:#111827"><strong>Updated At</strong></td>
					<td style="padding:10px 12px;font-size:13px;color:#374151">%s</td>
				</tr>
				</table>

				<div style="text-align:center;margin:18px 0 22px">
				<a href="%s" target="_blank" rel="noopener noreferrer"
					style="display:inline-block;font-size:14px;font-weight:700;
					padding:12px 28px;border-radius:12px;background:#08a4b8;color:#ffffff;
					text-decoration:none;box-shadow:0 6px 14px rgba(0,0,0,0.18)">
					View Order Details
				</a>
				</div>

				<div style="margin-top:8px">
				<div style="margin:0 0 6px;font-size:13px;color:#111827;font-weight:600">
					Remarks
				</div>
				<div style="margin:0;font-size:12px;color:#6b7280;line-height:1.7;white-space:pre-line">
					%s
				</div>
				</div>

				<p style="margin:22px 0 0;font-size:12px;color:#6b7280;line-height:1.7">
				Best regards,<br/>
				<strong>Prospira (Thailand) Co., Ltd.</strong>
				</p>
			</td>
			</tr>
		</table>

		<p style="text-align:center;margin:14px 0 0;font-size:11px;color:#9ca3af;line-height:1.6">
			This is an automatically generated email. Please do not reply to this message.
		</p>
		</body>
		</html>`,
		orderNumber, // <title>
		company,     // Dear %s
		company,     // Company (table)
		orderNumber, // Order No.
		nowUTC,      // Updated At
		link,        // Button link
		remark,      // Remarks
	)

	// สร้าง request
	req := &EmailRequest{
		ToEmails: cleaned,
		Subject:  subject,
		HTMLBody: htmlBody,
	}

	// Handle attachment
	attachPath := ""
	if strings.TrimSpace(fileURL) != "" {
		attachPath = "." + strings.TrimSpace(fileURL)
	}

	var err error
	if attachPath != "" {
		// Send with attachments using form-data endpoint
		_, err = c.SendMicroserviceWithAttachments(req, attachPath)
	} else {
		// Send without attachments using JSON endpoint
		_, err = c.SendMicroservice(req)
	}
	return err
}
