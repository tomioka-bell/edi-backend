package emailclient

type EmailRequest struct {
	ToEmails    []string `json:"to_emails"`
	FromName    string   `json:"from_name,omitempty"`
	Subject     string   `json:"subject"`
	HTMLBody    string   `json:"html_body"`
	TextBody    string   `json:"text_body,omitempty"`
	CC          []string `json:"cc,omitempty"`
	BCC         []string `json:"bcc,omitempty"`
	ReplyTo     string   `json:"reply_to,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type EmailAttachment struct {
	Filename    string `json:"filename"`               // ชื่อไฟล์
	ContentType string `json:"content_type,omitempty"` // MIME type
	URL         string `json:"url,omitempty"`          // URL ของไฟล์ (ให้ microservice ดาวน์โหลด)
	Base64Data  string `json:"base64_data,omitempty"`  // หรือส่ง base64 โดยตรง
}

// EmailResponse คือ response จาก email service
type EmailResponse struct {
	Status    string `json:"status"`
	MessageID string `json:"message_id,omitempty"`
	Error     string `json:"error,omitempty"`
}

type MicroserviceEmailRequest struct {
	Subject  string `json:"subject"`
	MailFrom string `json:"mail_from"`
	MailTo   string `json:"mail_to"`
	HTMLBody string `json:"html_body"`
}
