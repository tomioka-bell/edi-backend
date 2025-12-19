package emailclient

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

func (c *Client) SendModifyInvoiceVendorEmail(
	vendorCompany string,
	invoiceNumber string,
	fileURL string,
	note string,
	notificationType string,
) error {

	// recipients (internal Prospira team)
	toEmails, err := GetEDIEmployeeNotificationRecipientByCompany(notificationType)
	if err != nil {
		return err
	}

	seen := map[string]struct{}{}
	cleaned := make([]string, 0, len(toEmails))
	for _, e := range toEmails {
		addr := strings.TrimSpace(e.Principal.Email)
		if addr == "" {
			continue
		}
		key := strings.ToLower(addr)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		cleaned = append(cleaned, addr)
	}

	if len(cleaned) == 0 {
		return fmt.Errorf("recipient emails is empty")
	}

	// Link
	baseFE := strings.TrimRight(strings.TrimSpace(os.Getenv("FRONTEND_BASE_URL")), "/")
	if baseFE == "" {
		baseFE = "http://localhost:5173"
	}
	link := fmt.Sprintf("%s/en/invoice-form/%s", baseFE, url.PathEscape(invoiceNumber))

	// Formal text helpers
	remark := strings.TrimSpace(note)
	if remark == "" {
		remark = "-"
	}
	nowUTC := time.Now().UTC().Format("2006-01-02 15:04:05 UTC")

	// Subject (formal / internal)
	subject := fmt.Sprintf("[Prospira] Vendor Invoice Response Received – %s", invoiceNumber)

	// HTML (formal tone)
	htmlBody := fmt.Sprintf(`<!doctype html>
		<html>
		<head>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width,initial-scale=1"/>
		<title>Vendor Invoice Response - %[1]s</title>
		</head>

		<body style="margin:0;padding:24px;background: linear-gradient(to bottom, #08a4b8, #000000);
		font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#111827">

		<table role="presentation" width="100%%" cellspacing="0" cellpadding="0"
			style="max-width:680px;margin:auto;background:#ffffff;border-radius:14px;
			box-shadow:0 10px 24px rgba(0,0,0,0.12);overflow:hidden">

			<!-- Header -->
			<tr>
			<td style="padding:22px 26px;background:#ffffff;border-bottom:1px solid #e5e7eb">
				<div style="font-size:18px;font-weight:700;line-height:1.2;color:#111827">
				Vendor Invoice Response Notification
				</div>
				<div style="margin-top:6px;font-size:13px;color:#111827;opacity:0.85">
				Prospira (Thailand) Co., Ltd.
				</div>
			</td>
			</tr>

			<!-- Body -->
			<tr>
			<td style="padding:26px">
				<p style="margin:0 0 12px;font-size:14px;line-height:1.7;color:#111827">
				Dear Team,
				</p>

				<p style="margin:0 0 14px;font-size:14px;line-height:1.7;color:#374151">
				This is to notify you that the vendor has submitted or updated an invoice response in the system.
				Please review and proceed with the next steps as applicable.
				</p>

				<table role="presentation" cellspacing="0" cellpadding="0"
				style="width:100%%;margin:16px 0 18px;border:1px solid #e5e7eb;border-radius:12px;overflow:hidden">
				<tr>
					<td style="padding:10px 12px;background:#f9fafb;width:190px;font-size:13px;color:#111827"><strong>Vendor / Company</strong></td>
					<td style="padding:10px 12px;font-size:13px;color:#374151">%[2]s</td>
				</tr>
				<tr>
					<td style="padding:10px 12px;background:#f9fafb;font-size:13px;color:#111827"><strong>Invoice No.</strong></td>
					<td style="padding:10px 12px;font-size:13px;color:#374151">%[1]s</td>
				</tr>
				<tr>
					<td style="padding:10px 12px;background:#f9fafb;font-size:13px;color:#111827"><strong>Notification Type</strong></td>
					<td style="padding:10px 12px;font-size:13px;color:#374151">%[3]s</td>
				</tr>
				<tr>
					<td style="padding:10px 12px;background:#f9fafb;font-size:13px;color:#111827"><strong>Received At</strong></td>
					<td style="padding:10px 12px;font-size:13px;color:#374151">%[6]s</td>
				</tr>
				</table>

				<p style="margin:0 0 14px;font-size:14px;line-height:1.7;color:#374151">
				You may open the invoice record in the system using the button below:
				</p>

				<div style="text-align:center;margin:18px 0 18px">
				<a href="%[4]s" target="_blank" rel="noopener noreferrer"
					style="display:inline-block;padding:12px 24px;border-radius:12px;
					background:#08a4b8;color:#ffffff;text-decoration:none;font-size:14px;font-weight:700;
					box-shadow:0 6px 14px rgba(0,0,0,0.18)">
					Open Invoice in System
				</a>
				</div>

				<p style="margin:0 0 6px;font-size:12px;color:#6b7280;line-height:1.7">
				If the button does not work, please copy and paste this link into your browser:<br/>
				<span style="word-break:break-all;color:#0369a1">%[4]s</span>
				</p>

				<div style="margin-top:14px">
				<div style="margin:0 0 6px;font-size:13px;color:#111827;font-weight:600">Vendor Remarks</div>
				<div style="margin:0;font-size:12px;color:#6b7280;line-height:1.7;white-space:pre-line">%[5]s</div>
				</div>

				<p style="margin-top:22px;font-size:12px;color:#6b7280;line-height:1.7">
				Best regards,<br/>
				<strong>Prospira (Thailand) Co., Ltd.</strong>
				</p>
			</td>
			</tr>
		</table>

		<p style="text-align:center;margin:14px 0 0;font-size:11px;color:#9ca3af;line-height:1.6">
			This is an automatically generated email for internal use. Please do not reply to this message.
		</p>
		</body>
		</html>`,
		invoiceNumber,    // %[1]s
		vendorCompany,    // %[2]s
		notificationType, // %[3]s
		link,             // %[4]s
		remark,           // %[5]s
		nowUTC,           // %[6]s
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

	if attachPath != "" {
		// Send with attachments using form-data endpoint
		_, err = c.SendMicroserviceWithAttachments(req, attachPath)
	} else {
		// Send without attachments using JSON endpoint
		_, err = c.SendMicroservice(req)
	}
	return err
}
