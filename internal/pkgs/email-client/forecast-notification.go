package emailclient

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func (c *Client) SendForecastReadReminder(
	toEmail string,
	numberForecast string,
	vendorName string,
	periodFrom time.Time,
	targetTime time.Time,
) error {

	// Sanitize email
	toEmail = strings.TrimSpace(toEmail)
	if toEmail == "" {
		return fmt.Errorf("toEmail is empty")
	}

	// สร้าง link
	baseFE := strings.TrimRight(strings.TrimSpace(os.Getenv("FRONTEND_BASE_URL")), "/")
	if baseFE == "" {
		baseFE = "http://localhost:5173"
	}
	link := fmt.Sprintf("%s/en/forecast-form/%s", baseFE, numberForecast)

	// Subject
	subject := fmt.Sprintf("Forecast %s reminder – Prospira Thailand", numberForecast)

	// Format dates
	periodStr := periodFrom.Format("02 Jan 2006 15:04")
	targetStr := targetTime.Format("02 Jan 2006 15:04")

	// สร้าง HTML
	htmlBody := fmt.Sprintf(`
		<!doctype html>
		<html>
		<head>
			<meta charset="utf-8" />
			<title>Forecast reminder</title>
		</head>
		<body style="
			margin:0;
			padding:0;
			font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;
			background: linear-gradient(to bottom, #08a4b8, #000000);
		">
			<table width="100%%" cellpadding="0" cellspacing="0" style="padding:24px 0;">
				<tr>
					<td align="center">
						<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:520px;background-color:#ffffff;border-radius:8px;padding:24px 28px;border:1px solid #e5e7eb;">
							<tr>
								<td style="font-size:18px;font-weight:600;color:#111827;padding-bottom:12px;">
									Forecast reminder – %s
								</td>
							</tr>
							<tr>
								<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:8px;">
									Dear vendor,
								</td>
							</tr>
							<tr>
								<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:8px;">
									This is a reminder for your forecast document from <strong>%s</strong>.
								</td>
							</tr>
							<tr>
								<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:16px;">
									Please review and confirm the following forecast in the Prospira Thailand system:
								</td>
							</tr>
							<tr>
								<td style="font-size:13px;color:#111827;line-height:1.6;padding-bottom:12px;">
									<strong>Forecast No.:</strong> %s<br/>
									<strong>Vendor:</strong> %s<br/>
									<strong>Period From:</strong> %s<br/>
									<strong>Reminder Time:</strong> %s
								</td>
							</tr>
							<tr>
							<td align="center" style="padding-bottom:20px;">
							<a href="%s"
								style="display:inline-block;padding:10px 20px;border-radius:6px;
								background-color:#0284c7;color:#ffffff;text-decoration:none;
								font-size:14px;font-weight:500;">
								Open Forecast in System
							</a>
							</td>
						</tr>
							<tr>
								<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-top:8px;">
									Please log in to the Prospira Thailand system to view and acknowledge this forecast.
								</td>
							</tr>
							<tr>
								<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-top:8px;">
									Best regards,<br/>
									Prospira Thailand
								</td>
							</tr>
						</table>
						<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:520px;margin-top:12px;">
							<tr>
								<td style="font-size:11px;color:#9ca3af;text-align:center;line-height:1.4;">
									This is an automatically generated email. Please do not reply.
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</body>
		</html>`,
		numberForecast, // title
		vendorName,     // intro
		numberForecast, // Forecast No.
		vendorName,     // Vendor
		periodStr,      // Period From
		targetStr,      // Reminder Time
		link,
	)

	// สร้าง request
	req := &EmailRequest{
		ToEmails: []string{toEmail},
		Subject:  subject,
		HTMLBody: htmlBody,
	}

	_, err := c.SendMicroservice(req)
	return err
}
