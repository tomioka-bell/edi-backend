package emailclient

import (
	"fmt"
	"os"
	"strings"
)

func (c *Client) SendStatusForecastVendorEmail(
	toEmails []string,
	statusForecast,
	company,
	forecastNumber,
	fileURL string,
	note string,
) error {
	// เตรียม list email ปลายทาง
	cleaned := make([]string, 0, len(toEmails))
	for _, e := range toEmails {
		e = strings.TrimSpace(e)
		if e != "" {
			cleaned = append(cleaned, e)
		}
	}
	if len(cleaned) == 0 {
		return fmt.Errorf("toEmails is empty")
	}

	link := fmt.Sprintf("%s/en/forecast-form/%s", os.Getenv("FRONTEND_BASE_URL"), forecastNumber)

	// Subject
	subject := fmt.Sprintf("Forecast Status Notification | %s | %s", forecastNumber, statusForecast)

	// HTML Email
	htmlBody := fmt.Sprintf(`<!doctype html>
		<html>
		<head>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width,initial-scale=1"/>
		<title>Forecast %s - %s</title>
		</head>
		<body style="margin:0;padding:24px;background: linear-gradient(to bottom, #08a4b8, #000000);font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#111827">
		<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:560px;margin:auto;background:#ffffff;border-radius:12px;box-shadow:0 4px 16px rgba(0,0,0,0.06);overflow:hidden">
			<tr>
			<td style="padding:24px 24px 8px;border-bottom:1px solid #eef1f5">
				<h2 style="margin:0;font-size:20px;line-height:28px">
				Forecast :  %s - %s
				</h2>
				<p style="margin:6px 0 0;color:#6b7280;font-size:12px">
				Company: %s
				</p>
			</td>
			</tr>
			<tr>
			<td style="padding:24px">
				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
				Dear <strong>%s</strong>,
				</p>
				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
				We would like to formally inform you that the forecast for <strong>%s</strong> has been updated in our system.
				</p>

				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
				Please find the key forecast information below:
				</p>
				<table role="presentation" cellspacing="0" cellpadding="0" style="font-size:13px;color:#374151;margin:0 0 16px">
				<tr>
					<td style="padding:2px 0;width:120px"><strong>Company</strong></td>
					<td style="padding:2px 0">: %s</td>
				</tr>
				<tr>
					<td style="padding:2px 0"><strong>Forecast No.</strong></td>
					<td style="padding:2px 0">: %s</td>
				</tr>
				<tr>
					<td style="padding:2px 0"><strong>Status</strong></td>
					<td style="padding:2px 0">: %s</td>
				</tr>
				</table>

				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
				Kindly review the forecast details and provide your confirmation or comments through the link below at your earliest convenience:
				</p>

				<div style="text-align:center;margin:16px 0 20px">
				<a href="%s"
					style="display:inline-block;font-size:14px;font-weight:600;padding:12px 24px;border-radius:12px;background:linear-gradient(135deg, #0284c7 0%%, #0369a1 100%%);color:#ffffff;text-decoration:none;box-shadow:0 4px 6px rgba(2, 132, 199, 0.25);transition:all 0.3s ease;border:none">
					Open Forecast Form
				</a>
				</div>

				<p style="margin:0 0 8px;color:#374151;font-size:13px;line-height:20px">
				Remarks:
				</p>
				<p style="margin:0 0 18px;color:#6b7280;font-size:12px;line-height:20px">
				%s
				</p>

				<p style="margin:0;color:#9ca3af;font-size:12px;line-height:20px">
				The forecast document is attached to this email for your reference.
				</p>

				<p style="font-size:12px;color:#6b7280;line-height:1.6;padding-top:8px;">
									Best regards,<br/>
									Prospira Thailand
								</p>
			</td>
			</tr>
		</table>
		<p style="text-align:center;margin:16px 0 0;color:#9ca3af;font-size:12px">
			This is an automatically generated email. Please do not reply directly to this message.
		</p>
		</body>
		</html>`,
		forecastNumber, // title
		statusForecast, // title
		forecastNumber, // h2
		statusForecast, // h2
		company,        // subtitle
		company,        // paragraph
		company,        // dear
		company,        // table
		forecastNumber, // table
		statusForecast, // table
		link,           // button
		note,           // remarks
	)

	// Create EmailRequest
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

	if attachPath != "" {
		_, err := c.SendMicroserviceWithAttachments(req, attachPath)
		return err
	}
	_, err := c.SendMicroservice(req)
	return err
}

// Send Order Vendor Email

func (c *Client) SendStatusOrderVendorEmail(
	toEmails []string,
	statusOrder,
	company,
	orderNumber,
	fileURL string,
	note string,
) error {
	// Clean email list
	cleaned := make([]string, 0, len(toEmails))
	for _, e := range toEmails {
		e = strings.TrimSpace(e)
		if e != "" {
			cleaned = append(cleaned, e)
		}
	}
	if len(cleaned) == 0 {
		return fmt.Errorf("toEmails is empty")
	}

	link := fmt.Sprintf("%s/en/order-form/%s", os.Getenv("FRONTEND_BASE_URL"), orderNumber)

	// Subject
	subject := fmt.Sprintf("Order Status Notification | %s | %s", orderNumber, statusOrder)

	// HTML Email
	htmlBody := fmt.Sprintf(`<!doctype html>
		<html>
		<head>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width,initial-scale=1"/>
		<title>Order %s - %s</title>
		</head>
		<body style="margin:0;padding:24px;background: linear-gradient(to bottom, #08a4b8, #000000);font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#111827">
		<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:560px;margin:auto;background:#ffffff;border-radius:12px;box-shadow:0 4px 16px rgba(0,0,0,0.06);overflow:hidden">
			<tr>
			<td style="padding:24px 24px 8px;border-bottom:1px solid #eef1f5">
				<h2 style="margin:0;font-size:20px;line-height:28px">
				Order : %s - %s
				</h2>
				<p style="margin:6px 0 0;color:#6b7280;font-size:12px">
				Company: %s
				</p>
			</td>
			</tr>
			<tr>
			<td style="padding:24px">
				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
				Dear <strong>%s</strong>,
				</p>
				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
				We would like to formally inform you that the order for <strong>%s</strong> has been updated in our system.
				</p>

				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
				Please find the key order information below:
				</p>
				<table role="presentation" cellspacing="0" cellpadding="0" style="font-size:13px;color:#374151;margin:0 0 16px">
				<tr>
					<td style="padding:2px 0;width:120px"><strong>Company</strong></td>
					<td style="padding:2px 0">: %s</td>
				</tr>
				<tr>
					<td style="padding:2px 0"><strong>Order No.</strong></td>
					<td style="padding:2px 0">: %s</td>
				</tr>
				<tr>
					<td style="padding:2px 0"><strong>Status</strong></td>
					<td style="padding:2px 0">: %s</td>
				</tr>
				</table>

				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
				Kindly review the order details and provide your confirmation or comments through the link below at your earliest convenience:
				</p>

				<div style="text-align:center;margin:16px 0 20px">
				<a href="%s"
					style="display:inline-block;font-size:14px;font-weight:600;padding:12px 24px;border-radius:12px;background:linear-gradient(135deg, #0284c7 0%%, #0369a1 100%%);color:#ffffff;text-decoration:none;box-shadow:0 4px 6px rgba(2, 132, 199, 0.25);transition:all 0.3s ease;border:none">
					Open Order Form
				</a>
				</div>

				<p style="margin:0 0 8px;color:#374151;font-size:13px;line-height:20px">
				Remarks:
				</p>
				<p style="margin:0 0 18px;color:#6b7280;font-size:12px;line-height:20px">
				%s
				</p>

				<p style="margin:0;color:#9ca3af;font-size:12px;line-height:20px">
				The order document is attached to this email for your reference.
				</p>

				<p style="font-size:12px;color:#6b7280;line-height:1.6;padding-top:8px;">
									Best regards,<br/>
									Prospira Thailand
								</p>
			</td>
			</tr>
		</table>
		<p style="text-align:center;margin:16px 0 0;color:#9ca3af;font-size:12px">
			This is an automatically generated email. Please do not reply directly to this message.
		</p>
		</body>
		</html>`,
		orderNumber, // title
		statusOrder, // title
		orderNumber, // h2
		statusOrder, // h2
		company,     // subtitle
		company,     // paragraph
		company,     // dear
		company,     // table
		orderNumber, // table
		statusOrder, // table
		link,        // button
		note,        // remarks
	)

	// Create EmailRequest
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

	if attachPath != "" {
		_, err := c.SendMicroserviceWithAttachments(req, attachPath)
		return err
	}
	_, err := c.SendMicroservice(req)
	return err
}

// Send Invoice Vendor Email

func (c *Client) SendStatusInvoiceVendorEmail(
	toEmails []string,
	statusInvoice,
	company,
	invoiceNumber,
	fileURL string,
	note string,
) error {
	// Clean email list
	cleaned := make([]string, 0, len(toEmails))
	for _, e := range toEmails {
		e = strings.TrimSpace(e)
		if e != "" {
			cleaned = append(cleaned, e)
		}
	}
	if len(cleaned) == 0 {
		return fmt.Errorf("toEmails is empty")
	}

	link := fmt.Sprintf("%s/en/invoice-form/%s", os.Getenv("FRONTEND_BASE_URL"), invoiceNumber)

	// Subject
	subject := fmt.Sprintf("Invoice Status Notification | %s | %s", invoiceNumber, statusInvoice)

	// HTML Email
	htmlBody := fmt.Sprintf(`<!doctype html>
		<html>
		<head>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width,initial-scale=1"/>
		<title>Invoice %s - %s</title>
		</head>
		<body style="margin:0;padding:24px;background: linear-gradient(to bottom, #08a4b8, #000000);font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#111827">
		<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:560px;margin:auto;background:#ffffff;border-radius:12px;box-shadow:0 4px 16px rgba(0,0,0,0.06);overflow:hidden">
			<tr>
			<td style="padding:24px 24px 8px;border-bottom:1px solid #eef1f5">
				<h2 style="margin:0;font-size:20px;line-height:28px">
				Invoice : %s - %s
				</h2>
				<p style="margin:6px 0 0;color:#6b7280;font-size:12px">
				Company: %s
				</p>
			</td>
			</tr>
			<tr>
			<td style="padding:24px">
				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
				Dear <strong>%s</strong>,
				</p>
				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
				We would like to formally inform you that the invoice for <strong>%s</strong> has been updated in our system.
				</p>

				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
				Please find the key invoice information below:
				</p>
				<table role="presentation" cellspacing="0" cellpadding="0" style="font-size:13px;color:#374151;margin:0 0 16px">
				<tr>
					<td style="padding:2px 0;width:120px"><strong>Company</strong></td>
					<td style="padding:2px 0">: %s</td>
				</tr>
				<tr>
					<td style="padding:2px 0"><strong>Invoice No.</strong></td>
					<td style="padding:2px 0">: %s</td>
				</tr>
				<tr>
					<td style="padding:2px 0"><strong>Status</strong></td>
					<td style="padding:2px 0">: %s</td>
				</tr>
				</table>

				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
				Kindly review the invoice details and provide your confirmation or comments through the link below at your earliest convenience:
				</p>

				<div style="text-align:center;margin:16px 0 20px">
				<a href="%s"
					style="display:inline-block;font-size:14px;font-weight:600;padding:12px 24px;border-radius:12px;background:linear-gradient(135deg, #0284c7 0%%, #0369a1 100%%);color:#ffffff;text-decoration:none;box-shadow:0 4px 6px rgba(2, 132, 199, 0.25);transition:all 0.3s ease;border:none">
					Open Invoice Form
				</a>
				</div>

				<p style="margin:0 0 8px;color:#374151;font-size:13px;line-height:20px">
				Remarks:
				</p>
				<p style="margin:0 0 18px;color:#6b7280;font-size:12px;line-height:20px">
				%s
				</p>

				<p style="margin:0;color:#9ca3af;font-size:12px;line-height:20px">
				The invoice document is attached to this email for your reference.
				</p>

				<p style="font-size:12px;color:#6b7280;line-height:1.6;padding-top:8px;">
									Best regards,<br/>
									Prospira Thailand
								</p>
			</td>
			</tr>
		</table>
		<p style="text-align:center;margin:16px 0 0;color:#9ca3af;font-size:12px">
			This is an automatically generated email. Please do not reply directly to this message.
		</p>
		</body>
		</html>`,
		invoiceNumber, // title
		statusInvoice, // title
		invoiceNumber, // h2
		statusInvoice, // h2
		company,       // subtitle
		company,       // paragraph
		company,       // dear
		company,       // table
		invoiceNumber, // table
		statusInvoice, // table
		link,          // button
		note,          // remarks
	)

	// Create EmailRequest
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

	if attachPath != "" {
		_, err := c.SendMicroserviceWithAttachments(req, attachPath)
		return err
	}
	_, err := c.SendMicroservice(req)
	return err
}
