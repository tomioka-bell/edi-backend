package mailer

import emailclient "backend/internal/pkgs/email-client"

func SendCreatePasswordRequest(toEmail, resetLink string) error {
	return emailclient.SendCreatePasswordRequestPersonalized(toEmail, resetLink)
}

// func SendCreatePasswordRequest(toEmail, resetLink string) error {
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

// 	d := mail.NewDialer(host, port, user, pass)
// 	d.TLSConfig = &tls.Config{
// 		ServerName:         host,
// 		InsecureSkipVerify: true,
// 	}

// 	m := mail.NewMessage()
// 	m.SetHeader("From", from)
// 	m.SetHeader("To", toEmail)
// 	m.SetHeader("Subject", "Welcome to Prospira Thailand – Set up your password")

// 	html := fmt.Sprintf(`
// 	<!doctype html>
// 	<html>
// 	<head>
// 		<meta charset="utf-8" />
// 		<title>Set up your password</title>
// 	</head>
// 	<body style="
// 		margin:0;
// 		padding:0;
// 		font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;
// 		background: linear-gradient(to bottom, #08a4b8, #000000);
// 	">
// 		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:24px 0;">
// 			<tr>
// 				<td align="center">
// 					<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:520px;background-color:#ffffff;border-radius:8px;padding:24px 28px;border:1px solid #e5e7eb;">
// 						<tr>
// 							<td style="font-size:18px;font-weight:600;color:#111827;padding-bottom:12px;">
// 								Welcome to Prospira Thailand
// 							</td>
// 						</tr>
// 						<tr>
// 							<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:12px;">
// 								Dear user,
// 							</td>
// 						</tr>
// 						<tr>
// 							<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:12px;">
// 								An account has been created for you in the Prospira Thailand system using this email address.
// 							</td>
// 						</tr>
// 						<tr>
// 							<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:18px;">
// 								To activate and secure your account, please create your password by clicking the button below.
// 								For security reasons, this link will remain valid for 30 minutes.
// 							</td>
// 						</tr>
// 						<tr>
// 							<td align="center" style="padding-bottom:20px;">
// 								<a href="%s"
// 								style="display:inline-block;padding:10px 20px;border-radius:6px;
// 								background-color:#0284c7;color:#ffffff;text-decoration:none;
// 								font-size:14px;font-weight:500;">
// 									Set up password
// 								</a>
// 							</td>
// 						</tr>
// 						<tr>
// 							<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-bottom:16px;">
// 								If the button above does not work, you can copy and paste the following link into your browser:<br/>
// 								<span style="word-break:break-all;color:#0369a1;">%s</span>
// 							</td>
// 						</tr>
// 						<tr>
// 							<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-bottom:12px;">
// 								If you did not expect to receive this email or believe it was sent to you in error,
// 								you can safely ignore it. Your account will not be activated until a password is set.
// 							</td>
// 						</tr>
// 						<tr>
// 							<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-top:8px;">
// 								Best regards,<br/>
// 								Prospira Thailand
// 							</td>
// 						</tr>
// 					</table>
// 					<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:520px;margin-top:12px;">
// 						<tr>
// 							<td style="font-size:11px;color:#9ca3af;text-align:center;line-height:1.4;">
// 								This is an automatically generated email. Please do not reply.
// 							</td>
// 						</tr>
// 					</table>
// 				</td>
// 			</tr>
// 		</table>
// 	</body>
// 	</html>`, resetLink, resetLink)

// 	plain := fmt.Sprintf(
// 		"Welcome to Prospira Thailand\n\n"+
// 			"An account has been created for you in the Prospira Thailand system using this email address.\n\n"+
// 			"To activate and secure your account, please create your password using the link below (valid for 30 minutes):\n%s\n\n"+
// 			"If you did not expect this email or believe it was sent to you in error, you can ignore it. "+
// 			"Your account will not be activated until a password is set.\n\n"+
// 			"Best regards,\nProspira Thailand",
// 		resetLink,
// 	)

// 	m.SetBody("text/html; charset=UTF-8", html)
// 	m.AddAlternative("text/plain; charset=UTF-8", plain)

// 	if err := d.DialAndSend(m); err != nil {
// 		return fmt.Errorf("DialAndSend failed: %w", err)
// 	}
// 	return nil
// }
