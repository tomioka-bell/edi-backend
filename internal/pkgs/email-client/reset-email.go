package emailclient

import (
	"fmt"
	"strings"
)

func (c *Client) SendCreatePasswordRequest(toEmail, resetLink string) error {
	// Sanitize email
	toEmail = strings.TrimSpace(toEmail)
	if toEmail == "" {
		return fmt.Errorf("toEmail is empty")
	}

	// Subject
	subject := "Welcome to Prospira Thailand – Set up your password"

	// สร้าง HTML
	htmlBody := fmt.Sprintf(`
	<!doctype html>
	<html>
	<head>
		<meta charset="utf-8" />
		<title>Set up your password</title>
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
								Welcome to Prospira Thailand
							</td>
						</tr>
						<tr>
							<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:12px;">
								Dear user,
							</td>
						</tr>
						<tr>
							<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:12px;">
								An account has been created for you in the Prospira Thailand system using this email address.
							</td>
						</tr>
						<tr>
							<td style="font-size:14px;color:#374151;line-height:1.6;padding-bottom:18px;">
								To activate and secure your account, please create your password by clicking the button below.
								For security reasons, this link will remain valid for 30 minutes.
							</td>
						</tr>
						<tr>
							<td align="center" style="padding-bottom:20px;">
								<a href="%s"
								style="display:inline-block;padding:10px 20px;border-radius:6px;
								background-color:#0284c7;color:#ffffff;text-decoration:none;
								font-size:14px;font-weight:500;">
									Set up password
								</a>
							</td>
						</tr>
						<tr>
							<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-bottom:16px;">
								If the button above does not work, you can copy and paste the following link into your browser:<br/>
								<span style="word-break:break-all;color:#0369a1;">%s</span>
							</td>
						</tr>
						<tr>
							<td style="font-size:12px;color:#6b7280;line-height:1.6;padding-bottom:12px;">
								If you did not expect to receive this email or believe it was sent to you in error,
								you can safely ignore it. Your account will not be activated until a password is set.
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
	</html>`, resetLink, resetLink)

	// สร้าง request
	req := &EmailRequest{
		ToEmails: []string{toEmail},
		Subject:  subject,
		HTMLBody: htmlBody,
	}

	// ✅ ส่งผ่าน microservice
	_, err := c.SendMicroservice(req)
	return err
}
