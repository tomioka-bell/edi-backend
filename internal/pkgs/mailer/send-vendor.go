package mailer

import emailclient "backend/internal/pkgs/email-client"

func SendStatusForecastVendorEmail(toEmails []string,
	statusForecast,
	company,
	forecastNumber,
	fileURL string,
	note string) error {

	return emailclient.SendStatusForecastVendorEmail(toEmails, statusForecast, company, forecastNumber, fileURL, note)

}

func SendStatusOrderVendorEmail(toEmails []string,
	statusOrder,
	company,
	orderNumber,
	fileURL string,
	note string) error {

	return emailclient.SendStatusOrderVendorEmail(toEmails, statusOrder, company, orderNumber, fileURL, note)

}

func SendStatusInvoiceVendorEmail(toEmails []string,
	statusInvoice,
	company,
	invoiceNumber,
	fileURL string,
	note string) error {

	return emailclient.SendStatusInvoiceVendorEmail(toEmails, statusInvoice, company, invoiceNumber, fileURL, note)

}

// func SendStatusForecastVendorEmail(
// 	toEmails []string,
// 	statusForecast,
// 	company,
// 	forecastNumber,
// 	fileURL string,
// 	note string,
// ) error {
// 	// -----------------------------
// 	// เตรียม list email ปลายทาง
// 	// -----------------------------
// 	cleaned := make([]string, 0, len(toEmails))
// 	for _, e := range toEmails {
// 		e = strings.TrimSpace(e)
// 		if e != "" {
// 			cleaned = append(cleaned, e)
// 		}
// 	}
// 	if len(cleaned) == 0 {
// 		return fmt.Errorf("toEmails is empty")
// 	}

// 	// -----------------------------
// 	// โหลดค่า SMTP จาก ENV
// 	// -----------------------------
// 	host := os.Getenv("SMTP_HOST")
// 	portStr := os.Getenv("SMTP_PORT")
// 	user := os.Getenv("SMTP_USERNAME")
// 	pass := os.Getenv("SMTP_PASSWORD")
// 	from := os.Getenv("SMTP_SET_FROM")
// 	fromName := os.Getenv("SMTP_SET_FROM_NAME")

// 	if host == "" || portStr == "" || user == "" || pass == "" || from == "" {
// 		return fmt.Errorf("smtp env missing (HOST/PORT/USERNAME/PASSWORD/SET_FROM)")
// 	}

// 	port, err := strconv.Atoi(strings.TrimSpace(portStr))
// 	if err != nil {
// 		return fmt.Errorf("invalid SMTP_PORT: %w", err)
// 	}

// 	// -----------------------------
// 	// ตั้งค่า Dialer
// 	// -----------------------------
// 	d := mail.NewDialer(host, port, user, pass)
// 	if port == 465 {
// 		d.SSL = true
// 	}
// 	d.TLSConfig = &tls.Config{
// 		ServerName:         host,
// 		InsecureSkipVerify: true,
// 	}

// 	link := fmt.Sprintf("%s/en/forecast-form/%s", os.Getenv("FRONTEND_BASE_URL"), forecastNumber)

// 	// -----------------------------
// 	// สร้าง Message
// 	// -----------------------------
// 	m := mail.NewMessage()
// 	if fromName != "" {
// 		m.SetHeader("From", m.FormatAddress(from, fromName))
// 	} else {
// 		m.SetHeader("From", from)
// 	}
// 	m.SetHeader("To", cleaned...)

// 	// ทำให้ Subject ดูเป็นทางการมากขึ้น
// 	subject := fmt.Sprintf("Forecast Status Notification | %s | %s", forecastNumber, statusForecast)
// 	m.SetHeader("Subject", subject)

// 	// -----------------------------
// 	// HTML Email (ภาษาทางการมากขึ้น)
// 	// -----------------------------
// 	html := fmt.Sprintf(`<!doctype html>
// 		<html>
// 		<head>
// 		<meta charset="utf-8"/>
// 		<meta name="viewport" content="width=device-width,initial-scale=1"/>
// 		<title>Forecast %s - %s</title>
// 		</head>
// 		<body style="margin:0;padding:24px;background: linear-gradient(to bottom, #08a4b8, #000000);font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#111827">
// 		<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:560px;margin:auto;background:#ffffff;border-radius:12px;box-shadow:0 4px 16px rgba(0,0,0,0.06);overflow:hidden">
// 			<tr>
// 			<td style="padding:24px 24px 8px;border-bottom:1px solid #eef1f5">
// 				<h2 style="margin:0;font-size:20px;line-height:28px">
// 				Forecast :  %s - %s
// 				</h2>
// 				<p style="margin:6px 0 0;color:#6b7280;font-size:12px">
// 				Company: %s
// 				</p>
// 			</td>
// 			</tr>
// 			<tr>
// 			<td style="padding:24px">
// 				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
// 				Dear <strong>%s</strong>,
// 				</p>
// 				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
// 				We would like to formally inform you that the forecast for <strong>%s</strong> has been updated in our system.
// 				</p>

// 				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
// 				Please find the key forecast information below:
// 				</p>
// 				<table role="presentation" cellspacing="0" cellpadding="0" style="font-size:13px;color:#374151;margin:0 0 16px">
// 				<tr>
// 					<td style="padding:2px 0;width:120px"><strong>Company</strong></td>
// 					<td style="padding:2px 0">: %s</td>
// 				</tr>
// 				<tr>
// 					<td style="padding:2px 0"><strong>Forecast No.</strong></td>
// 					<td style="padding:2px 0">: %s</td>
// 				</tr>
// 				<tr>
// 					<td style="padding:2px 0"><strong>Status</strong></td>
// 					<td style="padding:2px 0">: %s</td>
// 				</tr>
// 				</table>

// 				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
// 				Kindly review the forecast details and provide your confirmation or comments through the link below at your earliest convenience:
// 				</p>

// 				<div style="text-align:center;margin:16px 0 20px">
// 				<a href="%s"
// 					style="display:inline-block;font-size:14px;font-weight:600;padding:12px 24px;border-radius:12px;background:linear-gradient(135deg, #0284c7 0%%, #0369a1 100%%);color:#ffffff;text-decoration:none;box-shadow:0 4px 6px rgba(2, 132, 199, 0.25);transition:all 0.3s ease;border:none">
// 					Open Forecast Form
// 				</a>
// 				</div>

// 				<p style="margin:0 0 8px;color:#374151;font-size:13px;line-height:20px">
// 				Remarks:
// 				</p>
// 				<p style="margin:0 0 18px;color:#6b7280;font-size:12px;line-height:20px">
// 				%s
// 				</p>

// 				<p style="margin:0;color:#9ca3af;font-size:12px;line-height:20px">
// 				The forecast document is attached to this email for your reference.
// 				</p>

// 				<p style="font-size:12px;color:#6b7280;line-height:1.6;padding-top:8px;">
// 									Best regards,<br/>
// 									Prospira Thailand
// 								</p>
// 			</td>
// 			</tr>
// 		</table>
// 		<p style="text-align:center;margin:16px 0 0;color:#9ca3af;font-size:12px">
// 			This is an automatically generated email. Please do not reply directly to this message.
// 		</p>
// 		</body>
// 		</html>`,
// 		// map %s ตามลำดับใน template:
// 		forecastNumber, // title Forecast %s - %s
// 		statusForecast, // title
// 		forecastNumber, // h2 Forecast %s - %s
// 		statusForecast, // h2
// 		company,        // Company: %s (subtitle)
// 		company,        // ในประโยค "forecast for <strong>%s</strong>"
// 		company,        // Dear %s,
// 		company,        // ตาราง Company
// 		forecastNumber, // ตาราง Forecast No.
// 		statusForecast, // ตาราง Status
// 		link,           // href ปุ่ม
// 		note,           // Remarks
// 	)

// 	m.SetBody("text/html; charset=UTF-8", html)

// 	// แนบไฟล์ (ถ้ามี)
// 	if fileURL != "" {
// 		filePath := "." + fileURL
// 		m.Attach(filePath)
// 	}

// 	if err := d.DialAndSend(m); err != nil {
// 		return fmt.Errorf("DialAndSend failed: %w", err)
// 	}
// 	return nil
// }

// // =================================== Send Order Vendor Email =======================================================

// func SendStatusOrderVendorEmail(
// 	toEmails []string,
// 	statusOrder,
// 	company,
// 	orderNumber,
// 	fileURL string,
// 	note string,
// ) error {
// 	// -----------------------------
// 	// เตรียม list email ปลายทาง
// 	// -----------------------------
// 	cleaned := make([]string, 0, len(toEmails))
// 	for _, e := range toEmails {
// 		e = strings.TrimSpace(e)
// 		if e != "" {
// 			cleaned = append(cleaned, e)
// 		}
// 	}
// 	if len(cleaned) == 0 {
// 		return fmt.Errorf("toEmails is empty")
// 	}

// 	// -----------------------------
// 	// โหลดค่า SMTP จาก ENV
// 	// -----------------------------
// 	host := os.Getenv("SMTP_HOST")
// 	portStr := os.Getenv("SMTP_PORT")
// 	user := os.Getenv("SMTP_USERNAME")
// 	pass := os.Getenv("SMTP_PASSWORD")
// 	from := os.Getenv("SMTP_SET_FROM")
// 	fromName := os.Getenv("SMTP_SET_FROM_NAME")

// 	if host == "" || portStr == "" || user == "" || pass == "" || from == "" {
// 		return fmt.Errorf("smtp env missing (HOST/PORT/USERNAME/PASSWORD/SET_FROM)")
// 	}

// 	port, err := strconv.Atoi(strings.TrimSpace(portStr))
// 	if err != nil {
// 		return fmt.Errorf("invalid SMTP_PORT: %w", err)
// 	}

// 	// -----------------------------
// 	// ตั้งค่า Dialer
// 	// -----------------------------
// 	d := mail.NewDialer(host, port, user, pass)
// 	if port == 465 {
// 		d.SSL = true
// 	}
// 	d.TLSConfig = &tls.Config{
// 		ServerName:         host,
// 		InsecureSkipVerify: true,
// 	}

// 	link := fmt.Sprintf("%s/en/order-form/%s", os.Getenv("FRONTEND_BASE_URL"), orderNumber)

// 	// -----------------------------
// 	// สร้าง Message
// 	// -----------------------------
// 	m := mail.NewMessage()
// 	if fromName != "" {
// 		m.SetHeader("From", m.FormatAddress(from, fromName))
// 	} else {
// 		m.SetHeader("From", from)
// 	}
// 	m.SetHeader("To", cleaned...)

// 	// ทำให้ Subject ดูเป็นทางการมากขึ้น
// 	subject := fmt.Sprintf("Order Status Notification | %s | %s", orderNumber, statusOrder)
// 	m.SetHeader("Subject", subject)

// 	// -----------------------------
// 	// HTML Email (ภาษาทางการมากขึ้น)
// 	// -----------------------------
// 	html := fmt.Sprintf(`<!doctype html>
// 		<html>
// 		<head>
// 		<meta charset="utf-8"/>
// 		<meta name="viewport" content="width=device-width,initial-scale=1"/>
// 		<title>Order %s - %s</title>
// 		</head>
// 		<body style="margin:0;padding:24px;background: linear-gradient(to bottom, #08a4b8, #000000);font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#111827">
// 		<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:560px;margin:auto;background:#ffffff;border-radius:12px;box-shadow:0 4px 16px rgba(0,0,0,0.06);overflow:hidden">
// 			<tr>
// 			<td style="padding:24px 24px 8px;border-bottom:1px solid #eef1f5">
// 				<h2 style="margin:0;font-size:20px;line-height:28px">
// 				Order : %s - %s
// 				</h2>
// 				<p style="margin:6px 0 0;color:#6b7280;font-size:12px">
// 				Company: %s
// 				</p>
// 			</td>
// 			</tr>
// 			<tr>
// 			<td style="padding:24px">
// 				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
// 				Dear <strong>%s</strong>,
// 				</p>
// 				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
// 				We would like to formally inform you that the order for <strong>%s</strong> has been updated in our system.
// 				</p>

// 				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
// 				Please find the key order information below:
// 				</p>
// 				<table role="presentation" cellspacing="0" cellpadding="0" style="font-size:13px;color:#374151;margin:0 0 16px">
// 				<tr>
// 					<td style="padding:2px 0;width:120px"><strong>Company</strong></td>
// 					<td style="padding:2px 0">: %s</td>
// 				</tr>
// 				<tr>
// 					<td style="padding:2px 0"><strong>Order No.</strong></td>
// 					<td style="padding:2px 0">: %s</td>
// 				</tr>
// 				<tr>
// 					<td style="padding:2px 0"><strong>Status</strong></td>
// 					<td style="padding:2px 0">: %s</td>
// 				</tr>
// 				</table>

// 				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
// 				Kindly review the order details and provide your confirmation or comments through the link below at your earliest convenience:
// 				</p>

// 				<div style="text-align:center;margin:16px 0 20px">
// 				<a href="%s"
// 					style="display:inline-block;font-size:14px;font-weight:600;padding:12px 24px;border-radius:12px;background:linear-gradient(135deg, #0284c7 0%%, #0369a1 100%%);color:#ffffff;text-decoration:none;box-shadow:0 4px 6px rgba(2, 132, 199, 0.25);transition:all 0.3s ease;border:none">
// 					Open Order Form
// 				</a>
// 				</div>

// 				<p style="margin:0 0 8px;color:#374151;font-size:13px;line-height:20px">
// 				Remarks:
// 				</p>
// 				<p style="margin:0 0 18px;color:#6b7280;font-size:12px;line-height:20px">
// 				%s
// 				</p>

// 				<p style="margin:0;color:#9ca3af;font-size:12px;line-height:20px">
// 				The order document is attached to this email for your reference.
// 				</p>

// 				<p style="font-size:12px;color:#6b7280;line-height:1.6;padding-top:8px;">
// 									Best regards,<br/>
// 									Prospira Thailand
// 								</p>
// 			</td>
// 			</tr>
// 		</table>
// 		<p style="text-align:center;margin:16px 0 0;color:#9ca3af;font-size:12px">
// 			This is an automatically generated email. Please do not reply directly to this message.
// 		</p>
// 		</body>
// 		</html>`,
// 		// map %s ตามลำดับใน template:
// 		orderNumber, // title Order %s - %s
// 		statusOrder, // title
// 		orderNumber, // h2 Order %s - %s
// 		statusOrder, // h2
// 		company,     // Company: %s (subtitle)
// 		company,     // ในประโยค "order for <strong>%s</strong>"
// 		company,     // Dear %s,
// 		company,     // ตาราง Company
// 		orderNumber, // ตาราง Order No.
// 		statusOrder, // ตาราง Status
// 		link,        // href ปุ่ม
// 		note,        // Remarks
// 	)

// 	m.SetBody("text/html; charset=UTF-8", html)

// 	// แนบไฟล์ (ถ้ามี)
// 	if fileURL != "" {
// 		filePath := "." + fileURL
// 		m.Attach(filePath)
// 	}

// 	if err := d.DialAndSend(m); err != nil {
// 		return fmt.Errorf("DialAndSend failed: %w", err)
// 	}
// 	return nil
// }

// // =================================== Send Invoice Vendor Email =======================================================

// func SendStatusInvoiceVendorEmail(
// 	toEmails []string,
// 	statusInvoice,
// 	company,
// 	invoiceNumber,
// 	fileURL string,
// 	note string,
// ) error {
// 	// -----------------------------
// 	// เตรียม list email ปลายทาง
// 	// -----------------------------
// 	cleaned := make([]string, 0, len(toEmails))
// 	for _, e := range toEmails {
// 		e = strings.TrimSpace(e)
// 		if e != "" {
// 			cleaned = append(cleaned, e)
// 		}
// 	}
// 	if len(cleaned) == 0 {
// 		return fmt.Errorf("toEmails is empty")
// 	}

// 	// -----------------------------
// 	// โหลดค่า SMTP จาก ENV
// 	// -----------------------------
// 	host := os.Getenv("SMTP_HOST")
// 	portStr := os.Getenv("SMTP_PORT")
// 	user := os.Getenv("SMTP_USERNAME")
// 	pass := os.Getenv("SMTP_PASSWORD")
// 	from := os.Getenv("SMTP_SET_FROM")
// 	fromName := os.Getenv("SMTP_SET_FROM_NAME")

// 	if host == "" || portStr == "" || user == "" || pass == "" || from == "" {
// 		return fmt.Errorf("smtp env missing (HOST/PORT/USERNAME/PASSWORD/SET_FROM)")
// 	}

// 	port, err := strconv.Atoi(strings.TrimSpace(portStr))
// 	if err != nil {
// 		return fmt.Errorf("invalid SMTP_PORT: %w", err)
// 	}

// 	// -----------------------------
// 	// ตั้งค่า Dialer
// 	// -----------------------------
// 	d := mail.NewDialer(host, port, user, pass)
// 	if port == 465 {
// 		d.SSL = true
// 	}
// 	d.TLSConfig = &tls.Config{
// 		ServerName:         host,
// 		InsecureSkipVerify: true,
// 	}

// 	link := fmt.Sprintf("%s/en/invoice-form/%s", os.Getenv("FRONTEND_BASE_URL"), invoiceNumber)

// 	// -----------------------------
// 	// สร้าง Message
// 	// -----------------------------
// 	m := mail.NewMessage()
// 	if fromName != "" {
// 		m.SetHeader("From", m.FormatAddress(from, fromName))
// 	} else {
// 		m.SetHeader("From", from)
// 	}
// 	m.SetHeader("To", cleaned...)

// 	// ทำให้ Subject ดูเป็นทางการมากขึ้น
// 	subject := fmt.Sprintf("Invoice Status Notification | %s | %s", invoiceNumber, statusInvoice)
// 	m.SetHeader("Subject", subject)

// 	// -----------------------------
// 	// HTML Email (ภาษาทางการมากขึ้น)
// 	// -----------------------------
// 	html := fmt.Sprintf(`<!doctype html>
// 		<html>
// 		<head>
// 		<meta charset="utf-8"/>
// 		<meta name="viewport" content="width=device-width,initial-scale=1"/>
// 		<title>Invoice %s - %s</title>
// 		</head>
// 		<body style="margin:0;padding:24px;background: linear-gradient(to bottom, #08a4b8, #000000);font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#111827">
// 		<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:560px;margin:auto;background:#ffffff;border-radius:12px;box-shadow:0 4px 16px rgba(0,0,0,0.06);overflow:hidden">
// 			<tr>
// 			<td style="padding:24px 24px 8px;border-bottom:1px solid #eef1f5">
// 				<h2 style="margin:0;font-size:20px;line-height:28px">
// 				Invoice : %s - %s
// 				</h2>
// 				<p style="margin:6px 0 0;color:#6b7280;font-size:12px">
// 				Company: %s
// 				</p>
// 			</td>
// 			</tr>
// 			<tr>
// 			<td style="padding:24px">
// 				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
// 				Dear <strong>%s</strong>,
// 				</p>
// 				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
// 				We would like to formally inform you that the invoice for <strong>%s</strong> has been updated in our system.
// 				</p>

// 				<p style="margin:0 0 8px;font-size:14px;line-height:22px">
// 				Please find the key invoice information below:
// 				</p>
// 				<table role="presentation" cellspacing="0" cellpadding="0" style="font-size:13px;color:#374151;margin:0 0 16px">
// 				<tr>
// 					<td style="padding:2px 0;width:120px"><strong>Company</strong></td>
// 					<td style="padding:2px 0">: %s</td>
// 				</tr>
// 				<tr>
// 					<td style="padding:2px 0"><strong>Invoice No.</strong></td>
// 					<td style="padding:2px 0">: %s</td>
// 				</tr>
// 				<tr>
// 					<td style="padding:2px 0"><strong>Status</strong></td>
// 					<td style="padding:2px 0">: %s</td>
// 				</tr>
// 				</table>

// 				<p style="margin:0 0 12px;font-size:14px;line-height:22px">
// 				Kindly review the invoice details and provide your confirmation or comments through the link below at your earliest convenience:
// 				</p>

// 				<div style="text-align:center;margin:16px 0 20px">
// 				<a href="%s"
// 					style="display:inline-block;font-size:14px;font-weight:600;padding:12px 24px;border-radius:12px;background:linear-gradient(135deg, #0284c7 0%%, #0369a1 100%%);color:#ffffff;text-decoration:none;box-shadow:0 4px 6px rgba(2, 132, 199, 0.25);transition:all 0.3s ease;border:none">
// 					Open Invoice Form
// 				</a>
// 				</div>

// 				<p style="margin:0 0 8px;color:#374151;font-size:13px;line-height:20px">
// 				Remarks:
// 				</p>
// 				<p style="margin:0 0 18px;color:#6b7280;font-size:12px;line-height:20px">
// 				%s
// 				</p>

// 				<p style="margin:0;color:#9ca3af;font-size:12px;line-height:20px">
// 				The invoice document is attached to this email for your reference.
// 				</p>

// 				<p style="font-size:12px;color:#6b7280;line-height:1.6;padding-top:8px;">
// 									Best regards,<br/>
// 									Prospira Thailand
// 								</p>
// 			</td>
// 			</tr>
// 		</table>
// 		<p style="text-align:center;margin:16px 0 0;color:#9ca3af;font-size:12px">
// 			This is an automatically generated email. Please do not reply directly to this message.
// 		</p>
// 		</body>
// 		</html>`,
// 		// map %s ตามลำดับใน template:
// 		invoiceNumber, // title Invoice %s - %s
// 		statusInvoice, // title
// 		invoiceNumber, // h2 Invoice %s - %s
// 		statusInvoice, // h2
// 		company,       // Company: %s (subtitle)
// 		company,       // ในประโยค "invoice for <strong>%s</strong>"
// 		company,       // Dear %s,
// 		company,       // ตาราง Company
// 		invoiceNumber, // ตาราง Invoice No.
// 		statusInvoice, // ตาราง Status
// 		link,          // href ปุ่ม
// 		note,          // Remarks
// 	)

// 	m.SetBody("text/html; charset=UTF-8", html)

// 	// แนบไฟล์ (ถ้ามี)
// 	if fileURL != "" {
// 		filePath := "." + fileURL
// 		m.Attach(filePath)
// 	}

// 	if err := d.DialAndSend(m); err != nil {
// 		return fmt.Errorf("DialAndSend failed: %w", err)
// 	}
// 	return nil
// }
