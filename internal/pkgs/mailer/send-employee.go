package mailer

import (
	emailclient "backend/internal/pkgs/email-client"
)

func SendStatusForecastEmployeeEmail(vendorCompany string,
	statusForecast,
	forecastNumber,
	fileURL string,
	note string,
	notificationType string) error {

	return emailclient.SendStatusForecastEmployeeEmail(vendorCompany, statusForecast, forecastNumber, fileURL, note, notificationType)

}

func SendStatusOrderEmployeeEmail(vendorCompany string,
	statusForecast,
	forecastNumber,
	fileURL string,
	note string,
	notificationType string) error {
	return emailclient.SendStatusOrderEmployeeEmail(vendorCompany, statusForecast, forecastNumber, fileURL, note, notificationType)
}

func SendStatusInvoiceEmployeeEmail(vendorCompany string,
	statusForecast,
	forecastNumber,
	fileURL string,
	note string,
	notificationType string) error {
	return emailclient.SendStatusInvoiceEmployeeEmail(vendorCompany, statusForecast, forecastNumber, fileURL, note, notificationType)
}

// func GetEDIEmployeeNotificationRecipientByCompany(notification_type string) ([]domains.EDIVendorNotificationRecipient, error) {
// 	db := database.InitDataBase()
// 	company := "Prospira (Thailand) Co., Ltd."

// 	rows, err := db.Raw(`
//     SELECT
//         r.vendor_notification_recipient_id,
//         r.company,
//         r.notification_type,
//         r.edi_principal_id,
//         r.created_at,
//         r.updated_at,
//         r.deleted_at,

//         ISNULL(p.edi_principal_id, '00000000-0000-0000-0000-000000000000') AS p_edi_principal_id,
//         p.external_id AS p_external_id,
//         p.source_system AS p_source_system,
//         p.email AS p_email,
//         p.display_name AS p_display_name,
//         p.profile AS p_profile,
//         p.[group] AS p_group,
//         p.role AS p_role,
//         p.status AS p_status,
//         p.username AS p_username,
//         p.created_at AS p_created_at,
//         p.updated_at AS p_updated_at,
//         p.deleted_at AS p_deleted_at
//     FROM edi_vendor_notification_recipient r
//     LEFT JOIN edi_principals p
//         ON r.edi_principal_id = p.edi_principal_id
//     WHERE r.company = ?
//       AND r.notification_type = ?
//       AND r.deleted_at IS NULL
// 	  AND p.source_system = 'APP_EMPLOYEE'
// `, company, notification_type).Rows()

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	recipients := []domains.EDIVendorNotificationRecipient{}

// 	for rows.Next() {
// 		var rec domains.EDIVendorNotificationRecipient
// 		var principal domains.EDI_Principal

// 		err := rows.Scan(
// 			&rec.VendorNotificationRecipientID,
// 			&rec.Company,
// 			&rec.NotificationType,
// 			&rec.EDI_PrincipalID,
// 			&rec.CreatedAt,
// 			&rec.UpdatedAt,
// 			&rec.DeletedAt,

// 			&principal.EDI_PrincipalID,
// 			&principal.ExternalID,
// 			&principal.SourceSystem,
// 			&principal.Email,
// 			&principal.DisplayName,
// 			&principal.Profile,
// 			&principal.Group,
// 			&principal.Role,
// 			&principal.Status,
// 			&principal.Username,
// 			&principal.CreatedAt,
// 			&principal.UpdatedAt,
// 			&principal.DeletedAt,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		if principal.EDI_PrincipalID != (mssql.UniqueIdentifier{}) {
// 			rec.Principal = &principal
// 		}

// 		recipients = append(recipients, rec)
// 	}

// 	return recipients, nil
// }

// // Vendor ส่งอีเมลหาบริษัท Prospira (Thailand) Co., Ltd.
// func SendStatusForecastEmployeeEmail(
// 	vendorCompany string,
// 	statusForecast,
// 	forecastNumber,
// 	fileURL string,
// 	note string,
// 	notificationType string,
// ) error {

// 	fmt.Printf("ข้อมูลบริษัทที่ส่งมา : %q\n", vendorCompany)
// 	// -----------------------------
// 	// ดึงอีเมลผู้รับภายใน
// 	// -----------------------------
// 	toEmails, err := GetEDIEmployeeNotificationRecipientByCompany(notificationType)
// 	if err != nil {
// 		return err
// 	}

// 	cleaned := make([]string, 0, len(toEmails))
// 	for _, e := range toEmails {
// 		if e.Principal.Email != "" {
// 			cleaned = append(cleaned, e.Principal.Email)
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

// 	// -----------------------------
// 	// เตรียมลิงก์เข้าเว็บ (เปลี่ยน localhost เป็น ENV ถ้ามี)
// 	// -----------------------------
// 	baseFE := strings.TrimRight(os.Getenv("FRONTEND_BASE_URL"), "/")
// 	if baseFE == "" {
// 		baseFE = "http://localhost:5173"
// 	}
// 	link := fmt.Sprintf("%s/en/forecast-form/%s", baseFE, forecastNumber)

// 	// -----------------------------
// 	// สร้าง Message
// 	// -----------------------------
// 	m := mail.NewMessage()
// 	m.SetHeader("From", from)
// 	m.SetHeader("To", cleaned...)

// 	subject := fmt.Sprintf("Vendor Forecast Response | %s | %s", forecastNumber, statusForecast)
// 	m.SetHeader("Subject", subject)

// 	// -----------------------------
// 	// HTML Email (ปรับให้เป็น vendor แจ้งมาหาทีมเรา)
// 	// -----------------------------
// 	html := fmt.Sprintf(`<!doctype html>
// 		<html>
// 		<head>
// 		<meta charset="utf-8"/>
// 		<meta name="viewport" content="width=device-width,initial-scale=1"/>
// 		<title>Vendor Forecast Response - %[1]s</title>
// 		</head>
// 		<body style="
// 			margin:0;
// 			padding:0;
// 			font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;
// 			background: linear-gradient(to bottom, #08a4b8, #000000);
// 		">
// 		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:24px 0;">
// 			<tr>
// 			<td align="center">
// 				<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:640px;background-color:#ffffff;border-radius:8px;padding:24px 28px;border:1px solid #e5e7eb;">
// 				<tr>
// 					<td style="font-size:18px;font-weight:600;color:#111827;padding-bottom:12px;">
// 					Vendor Forecast Response - %[1]s (%[2]s)
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:8px;">
// 					Dear Team,
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:12px;">
// 					Vendor <strong>%[3]s</strong> has submitted or updated a response for forecast <strong>%[1]s</strong>.
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:12px;">
// 					<table cellpadding="0" cellspacing="0" style="font-size:13px;color:#374151;">
// 						<tr>
// 						<td style="padding:2px 0;width:160px;"><strong>Vendor / Company</strong></td>
// 						<td style="padding:2px 0;">: %[3]s</td>
// 						</tr>
// 						<tr>
// 						<td style="padding:2px 0;"><strong>Forecast No.</strong></td>
// 						<td style="padding:2px 0;">: %[1]s</td>
// 						</tr>
// 						<tr>
// 						<td style="padding:2px 0;"><strong>Forecast Status</strong></td>
// 						<td style="padding:2px 0;">: %[2]s</td>
// 						</tr>
// 						<tr>
// 						<td style="padding:2px 0;"><strong>Notification Type</strong></td>
// 						<td style="padding:2px 0;">: %[4]s</td>
// 						</tr>
// 					</table>
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:18px;">
// 					You can review the full details in the system via the link below:
// 					</td>
// 				</tr>
// 				<tr>
// 					<td align="center" style="padding-bottom:20px;">
// 					<a href="%[5]s"
// 						style="display:inline-block;padding:10px 20px;border-radius:6px;
// 						background-color:#0284c7;color:#ffffff;text-decoration:none;
// 						font-size:14px;font-weight:500;">
// 						Open Forecast in System
// 					</a>
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-bottom:16px;">
// 					If the button above does not work, you can copy and paste the following link into your browser:<br/>
// 					<span style="word-break:break-all;color:#0369a1;">%[5]s</span>
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-bottom:12px;">
// 					Additional remarks from vendor:<br/>
// 					<span style="white-space:pre-line;">%[6]s</span>
// 					</td>
// 				</tr>
// 				</table>
// 				<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:640px;margin-top:12px;">
// 				<tr>
// 					<td style="font-size:11px;color:#9ca3af;text-align:center;line-height:1.4;">
// 					This is an automatically generated email for internal use. Please do not reply to this message.
// 					</td>
// 				</tr>
// 				</table>
// 			</td>
// 			</tr>
// 		</table>
// 		</body>
// 		</html>`,
// 		forecastNumber,   // %[1]s
// 		statusForecast,   // %[2]s
// 		vendorCompany,    // %[3]s
// 		notificationType, // %[4]s
// 		link,             // %[5]s (button + plain link)
// 		note,             // %[6]s
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

// // =================================== Send Order Employee Email =======================================================

// func SendStatusOrderEmployeeEmail(
// 	vendorCompany string,
// 	statusOrder,
// 	orderNumber,
// 	fileURL string,
// 	note string,
// 	notificationType string,
// ) error {

// 	fmt.Printf("ข้อมูลบริษัทที่ส่งมา : %q\n", vendorCompany)
// 	// -----------------------------
// 	// ดึงอีเมลผู้รับภายใน
// 	// -----------------------------
// 	toEmails, err := GetEDIEmployeeNotificationRecipientByCompany(notificationType)
// 	if err != nil {
// 		return err
// 	}

// 	cleaned := make([]string, 0, len(toEmails))
// 	for _, e := range toEmails {
// 		if e.Principal.Email != "" {
// 			cleaned = append(cleaned, e.Principal.Email)
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

// 	// -----------------------------
// 	// เตรียมลิงก์เข้าเว็บ (เปลี่ยน localhost เป็น ENV ถ้ามี)
// 	// -----------------------------
// 	baseFE := strings.TrimRight(os.Getenv("FRONTEND_BASE_URL"), "/")
// 	if baseFE == "" {
// 		baseFE = "http://localhost:5173"
// 	}
// 	link := fmt.Sprintf("%s/en/order-form/%s", baseFE, orderNumber)

// 	// -----------------------------
// 	// สร้าง Message
// 	// -----------------------------
// 	m := mail.NewMessage()
// 	m.SetHeader("From", from)
// 	m.SetHeader("To", cleaned...)

// 	subject := fmt.Sprintf("Vendor Order Response | %s | %s", orderNumber, statusOrder)
// 	m.SetHeader("Subject", subject)

// 	// -----------------------------
// 	// HTML Email (ปรับให้เป็น vendor แจ้งมาหาทีมเรา)
// 	// -----------------------------
// 	html := fmt.Sprintf(`<!doctype html>
// 		<html>
// 		<head>
// 		<meta charset="utf-8"/>
// 		<meta name="viewport" content="width=device-width,initial-scale=1"/>
// 		<title>Vendor Order Response - %[1]s</title>
// 		</head>
// 		<body style="
// 			margin:0;
// 			padding:0;
// 			font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;
// 			background: linear-gradient(to bottom, #08a4b8, #000000);
// 		">
// 		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:24px 0;">
// 			<tr>
// 			<td align="center">
// 				<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:640px;background-color:#ffffff;border-radius:8px;padding:24px 28px;border:1px solid #e5e7eb;">
// 				<tr>
// 					<td style="font-size:18px;font-weight:600;color:#111827;padding-bottom:12px;">
// 					Vendor Order Response - %[1]s (%[2]s)
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:8px;">
// 					Dear Team,
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:12px;">
// 					Vendor <strong>%[3]s</strong> has submitted or updated a response for order <strong>%[1]s</strong>.
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:12px;">
// 					<table cellpadding="0" cellspacing="0" style="font-size:13px;color:#374151;">
// 						<tr>
// 						<td style="padding:2px 0;width:160px;"><strong>Vendor / Company</strong></td>
// 						<td style="padding:2px 0;">: %[3]s</td>
// 						</tr>
// 						<tr>
// 						<td style="padding:2px 0;"><strong>Order No.</strong></td>
// 						<td style="padding:2px 0;">: %[1]s</td>
// 						</tr>
// 						<tr>
// 						<td style="padding:2px 0;"><strong>Order Status</strong></td>
// 						<td style="padding:2px 0;">: %[2]s</td>
// 						</tr>
// 						<tr>
// 						<td style="padding:2px 0;"><strong>Notification Type</strong></td>
// 						<td style="padding:2px 0;">: %[4]s</td>
// 						</tr>
// 					</table>
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:18px;">
// 					You can review the full details in the system via the link below:
// 					</td>
// 				</tr>
// 				<tr>
// 					<td align="center" style="padding-bottom:20px;">
// 					<a href="%[5]s"
// 						style="display:inline-block;padding:10px 20px;border-radius:6px;
// 						background-color:#0284c7;color:#ffffff;text-decoration:none;
// 						font-size:14px;font-weight:500;">
// 						Open Order in System
// 					</a>
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-bottom:16px;">
// 					If the button above does not work, you can copy and paste the following link into your browser:<br/>
// 					<span style="word-break:break-all;color:#0369a1;">%[5]s</span>
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-bottom:12px;">
// 					Additional remarks from vendor:<br/>
// 					<span style="white-space:pre-line;">%[6]s</span>
// 					</td>
// 				</tr>
// 				</table>
// 				<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:640px;margin-top:12px;">
// 				<tr>
// 					<td style="font-size:11px;color:#9ca3af;text-align:center;line-height:1.4;">
// 					This is an automatically generated email for internal use. Please do not reply to this message.
// 					</td>
// 				</tr>
// 				</table>
// 			</td>
// 			</tr>
// 		</table>
// 		</body>
// 		</html>`,
// 		orderNumber,      // %[1]s
// 		statusOrder,      // %[2]s
// 		vendorCompany,    // %[3]s
// 		notificationType, // %[4]s
// 		link,             // %[5]s (button + plain link)
// 		note,             // %[6]s
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

// // =================================== Send Invoice Employee Email =======================================================

// func SendStatusInvoiceEmployeeEmail(
// 	vendorCompany string,
// 	statusInvoice,
// 	invoiceNumber,
// 	fileURL string,
// 	note string,
// 	notificationType string,
// ) error {

// 	fmt.Printf("ข้อมูลบริษัทที่ส่งมา : %q\n", vendorCompany)
// 	// -----------------------------
// 	// ดึงอีเมลผู้รับภายใน
// 	// -----------------------------
// 	toEmails, err := GetEDIEmployeeNotificationRecipientByCompany(notificationType)
// 	if err != nil {
// 		return err
// 	}

// 	cleaned := make([]string, 0, len(toEmails))
// 	for _, e := range toEmails {
// 		if e.Principal.Email != "" {
// 			cleaned = append(cleaned, e.Principal.Email)
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

// 	// -----------------------------
// 	// เตรียมลิงก์เข้าเว็บ (เปลี่ยน localhost เป็น ENV ถ้ามี)
// 	// -----------------------------
// 	baseFE := strings.TrimRight(os.Getenv("FRONTEND_BASE_URL"), "/")
// 	if baseFE == "" {
// 		baseFE = "http://localhost:5173"
// 	}
// 	link := fmt.Sprintf("%s/en/invoice-form/%s", baseFE, invoiceNumber)

// 	// -----------------------------
// 	// สร้าง Message
// 	// -----------------------------
// 	m := mail.NewMessage()
// 	m.SetHeader("From", from)
// 	m.SetHeader("To", cleaned...)

// 	subject := fmt.Sprintf("Vendor Invoice Response | %s | %s", invoiceNumber, statusInvoice)
// 	m.SetHeader("Subject", subject)

// 	// -----------------------------
// 	// HTML Email (ปรับให้เป็น vendor แจ้งมาหาทีมเรา)
// 	// -----------------------------
// 	html := fmt.Sprintf(`<!doctype html>
// 		<html>
// 		<head>
// 		<meta charset="utf-8"/>
// 		<meta name="viewport" content="width=device-width,initial-scale=1"/>
// 		<title>Vendor Invoice Response - %[1]s</title>
// 		</head>
// 		<body style="
// 			margin:0;
// 			padding:0;
// 			font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;
// 			background: linear-gradient(to bottom, #08a4b8, #000000);
// 		">
// 		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:24px 0;">
// 			<tr>
// 			<td align="center">
// 				<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:640px;background-color:#ffffff;border-radius:8px;padding:24px 28px;border:1px solid #e5e7eb;">
// 				<tr>
// 					<td style="font-size:18px;font-weight:600;color:#111827;padding-bottom:12px;">
// 					Vendor Invoice Response - %[1]s (%[2]s)
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:8px;">
// 					Dear Team,
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:12px;">
// 					Vendor <strong>%[3]s</strong> has submitted or updated a response for invoice <strong>%[1]s</strong>.
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:12px;">
// 					<table cellpadding="0" cellspacing="0" style="font-size:13px;color:#374151;">
// 						<tr>
// 						<td style="padding:2px 0;width:160px;"><strong>Vendor / Company</strong></td>
// 						<td style="padding:2px 0;">: %[3]s</td>
// 						</tr>
// 						<tr>
// 						<td style="padding:2px 0;"><strong>Invoice No.</strong></td>
// 						<td style="padding:2px 0;">: %[1]s</td>
// 						</tr>
// 						<tr>
// 						<td style="padding:2px 0;"><strong>Invoice Status</strong></td>
// 						<td style="padding:2px 0;">: %[2]s</td>
// 						</tr>
// 						<tr>
// 						<td style="padding:2px 0;"><strong>Notification Type</strong></td>
// 						<td style="padding:2px 0;">: %[4]s</td>
// 						</tr>
// 					</table>
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:18px;">
// 					You can review the full details in the system via the link below:
// 					</td>
// 				</tr>
// 				<tr>
// 					<td align="center" style="padding-bottom:20px;">
// 					<a href="%[5]s"
// 						style="display:inline-block;padding:10px 20px;border-radius:6px;
// 						background-color:#0284c7;color:#ffffff;text-decoration:none;
// 						font-size:14px;font-weight:500;">
// 						Open Invoice in System
// 					</a>
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-bottom:16px;">
// 					If the button above does not work, you can copy and paste the following link into your browser:<br/>
// 					<span style="word-break:break-all;color:#0369a1;">%[5]s</span>
// 					</td>
// 				</tr>
// 				<tr>
// 					<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-bottom:12px;">
// 					Additional remarks from vendor:<br/>
// 					<span style="white-space:pre-line;">%[6]s</span>
// 					</td>
// 				</tr>
// 				</table>
// 				<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:640px;margin-top:12px;">
// 				<tr>
// 					<td style="font-size:11px;color:#9ca3af;text-align:center;line-height:1.4;">
// 					This is an automatically generated email for internal use. Please do not reply to this message.
// 					</td>
// 				</tr>
// 				</table>
// 			</td>
// 			</tr>
// 		</table>
// 		</body>
// 		</html>`,
// 		invoiceNumber,    // %[1]s
// 		statusInvoice,    // %[2]s
// 		vendorCompany,    // %[3]s
// 		notificationType, // %[4]s
// 		link,             // %[5]s (button + plain link)
// 		note,             // %[6]s
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
