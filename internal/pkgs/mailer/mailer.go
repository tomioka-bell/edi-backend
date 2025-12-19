package mailer

import (
	emailclient "backend/internal/pkgs/email-client"
)

func SendPasswordResetEmail(toEmail, resetLink string) error {
	return emailclient.SendPasswordResetEmailPersonalized(toEmail, resetLink)
}

func SendLoginOTPEmail(toEmail, otp string) error {
	return emailclient.SendLoginOTPEmail(toEmail, otp)
}

// func SendLoginOTPEmail(toEmail, otp string) error {
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

// 	d := mail.NewDialer(host, port, user, pass)
// 	if port == 465 {
// 		d.SSL = true
// 	}
// 	d.TLSConfig = &tls.Config{
// 		ServerName:         host,
// 		InsecureSkipVerify: true,
// 	}

// 	m := mail.NewMessage()
// 	if fromName != "" {
// 		m.SetHeader("From", m.FormatAddress(from, fromName))
// 	} else {
// 		m.SetHeader("From", from)
// 	}
// 	m.SetHeader("To", toEmail)
// 	m.SetHeader("Subject", "Your verification code")

// 	expireMin := 1
// 	expiryTextEN := "1 minute"

// 	loc, _ := time.LoadLocation("Asia/Bangkok")
// 	nowTH := time.Now().In(loc)
// 	expiresAtTH := nowTH.Add(time.Duration(expireMin) * time.Minute)

// 	generatedAtText := nowTH.Format(time.RFC1123)
// 	expiresAtText := expiresAtTH.Format(time.RFC1123)

// 	dear := utils.DisplayName(toEmail)

// 	html := fmt.Sprintf(`<!doctype html>
// 	<html>
// 	<head><meta charset="utf-8"/><meta name="viewport" content="width=device-width,initial-scale=1"/></head>
// 	<body style="margin:0;padding:24px;background: linear-gradient(to bottom, #08a4b8, #000000);font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#111827">
// 	<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:560px;margin:auto;background:#ffffff;border-radius:12px;box-shadow:0 4px 16px rgba(0,0,0,0.06);overflow:hidden">
// 		<tr><td style="padding:24px 24px 8px;border-bottom:1px solid #eef1f5">
// 		<h2 style="margin:0;font-size:20px;line-height:28px">Verification code</h2>
// 		<p style="margin:6px 0 0;color:#6b7280;font-size:12px">Generated at %s</p>
// 		</td></tr>
// 		<tr><td style="padding:24px">
// 		<p style="margin:0 0 8px;font-size:14px;line-height:22px">Dear %s,</p>
// 		<p style="margin:0 0 12px;font-size:14px;line-height:22px">Use this code to finish signing in:</p>
// 		<div style="text-align:center;margin:16px 0 20px">
// 			<div style="display:inline-block;font-size:28px;letter-spacing:6px;font-weight:700;padding:12px 16px;border:1px solid #e5e7eb;border-radius:10px;background:#fafafa">%s</div>
// 		</div>
// 		<p style="margin:0 0 6px;color:#374151;font-size:13px;line-height:20px">
// 			This code expires in <strong>%s</strong>.
// 		</p>
// 		<p style="margin:0 0 18px;color:#6b7280;font-size:12px;line-height:20px">
// 			Expires at: %s
// 		</p>
// 		<p style="margin:0;color:#9ca3af;font-size:12px;line-height:20px">
// 			If you didn’t request it, you can safely ignore this email.
// 		</p>
// 		</td></tr>
// 	</table>
// 	<p style="text-align:center;margin:16px 0 0;color:#9ca3af;font-size:12px">Do not share this code with anyone.</p>
// 	</body>
// 	</html>`, generatedAtText, dear, otp, expiryTextEN, expiresAtText)

// 	m.SetBody("text/html; charset=UTF-8", html)

// 	if err := d.DialAndSend(m); err != nil {
// 		return fmt.Errorf("DialAndSend failed: %w", err)
// 	}
// 	return nil
// }

//==============================================================================================================

func SendTestEmail(toEmail string) error {
	return emailclient.SendTestEmail(toEmail)
}

// func SendTestEmail(toEmail string) error {
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

// 	d := mail.NewDialer(host, port, user, pass)
// 	if port == 465 {
// 		d.SSL = true
// 	}
// 	d.TLSConfig = &tls.Config{
// 		ServerName:         host,
// 		InsecureSkipVerify: true,
// 	}

// 	// สร้างอีเมล
// 	m := mail.NewMessage()
// 	if fromName != "" {
// 		m.SetHeader("From", m.FormatAddress(from, fromName))
// 	} else {
// 		m.SetHeader("From", from)
// 	}
// 	m.SetHeader("To", toEmail)
// 	m.SetHeader("Subject", "Test Email from Backend")

// 	htmlBody := `
// 	<!doctype html>
// 	<html>
// 	<head><meta charset="utf-8"/></head>
// 	<body style="font-family:Arial, sans-serif;padding:20px;">
// 		<h2>Test Email Successful 🎉</h2>
// 		<p>This is a test email sent from your backend SMTP settings.</p>
// 		<p>If you received this email, your SMTP configuration is working correctly.</p>
// 	</body>
// 	</html>
// 	`

// 	m.SetBody("text/html; charset=UTF-8", htmlBody)

// 	// ส่งอีเมล
// 	if err := d.DialAndSend(m); err != nil {
// 		return fmt.Errorf("DialAndSend failed: %w", err)
// 	}
// 	return nil
// }
