package emailclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Client struct {
	config     *Config
	httpClient *http.Client
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

func NewClientFromEnv() (*Client, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return NewClient(config), nil
}

// Send ส่งอีเมลแบบเต็มรูปแบบ
func (c *Client) SendMicroservice(req *EmailRequest) (*EmailResponse, error) {
	if req.FromName == "" && c.config.FromName != "" {
		req.FromName = c.config.FromName
	}

	if len(req.ToEmails) == 0 {
		return nil, fmt.Errorf("to_emails is required")
	}

	emailArray := make([]MicroserviceEmailRequest, len(req.ToEmails))
	for i, toEmail := range req.ToEmails {
		emailArray[i] = MicroserviceEmailRequest{
			Subject:  req.Subject,
			MailFrom: os.Getenv("SMTP_SET_FROM"),
			MailTo:   toEmail,
			HTMLBody: req.HTMLBody,
		}
	}

	jsonData, err := json.Marshal(emailArray)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request:  %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.config.ServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.config.Token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.Token)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var emailResp EmailResponse
	if err := json.Unmarshal(body, &emailResp); err != nil {
		return nil, fmt.Errorf("failed to parse response:  %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return &emailResp, fmt.Errorf("email service returned status %d: %s", resp.StatusCode, emailResp.Error)
	}

	return &emailResp, nil
}

// SendMicroserviceWithAttachments ส่งอีเมลพร้อมไฟล์แนบโดยใช้ form-data
func (c *Client) SendMicroserviceWithAttachments(req *EmailRequest, attachmentPath string) (*EmailResponse, error) {

	fmt.Println("📧 Sending email with attachments to", req.ToEmails)
	if len(req.ToEmails) == 0 {
		return nil, fmt.Errorf("to_emails is required")
	}

	mailFrom := os.Getenv("SMTP_SET_FROM")
	if req.FromName != "" {
		mailFrom = req.FromName
	}
	if mailFrom == "" {
		mailFrom = "Prospira <noreply@prospira-th.com>"
	}

	// Create multipart form for each recipient
	for _, toEmail := range req.ToEmails {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add form fields
		if err := writer.WriteField("subject", req.Subject); err != nil {
			return nil, fmt.Errorf("failed to write subject field: %w", err)
		}
		if err := writer.WriteField("mail_from", mailFrom); err != nil {
			return nil, fmt.Errorf("failed to write mail_from field: %w", err)
		}
		if err := writer.WriteField("mail_to", toEmail); err != nil {
			return nil, fmt.Errorf("failed to write mail_to field: %w", err)
		}
		if err := writer.WriteField("html_body", req.HTMLBody); err != nil {
			return nil, fmt.Errorf("failed to write html_body field: %w", err)
		}

		// Add file attachment if provided
		if attachmentPath != "" {
			fileData, err := os.ReadFile(attachmentPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read attachment file: %w", err)
			}

			part, err := writer.CreateFormFile("files", getFileName(attachmentPath))
			if err != nil {
				return nil, fmt.Errorf("failed to create form file: %w", err)
			}

			if _, err := part.Write(fileData); err != nil {
				return nil, fmt.Errorf("failed to write file data: %w", err)
			}
		}

		writer.Close()

		// Determine service URL - use attachment endpoint
		serviceURL := c.config.ServiceURL
		if !isAttachmentEndpoint(serviceURL) {
			// Replace the standard endpoint with the attachment endpoint
			serviceURL = os.Getenv("EMAIL_SERVICE_URL_ATTACHED")
		}

		fmt.Printf("📧 Sending email with attachments to %s\n", toEmail)
		fmt.Printf("📄 Using endpoint: %s\n", serviceURL)

		httpReq, err := http.NewRequest("POST", serviceURL, body)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		httpReq.Header.Set("Content-Type", writer.FormDataContentType())
		if c.config.Token != "" {
			httpReq.Header.Set("Authorization", "Bearer "+c.config.Token)
		}

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		var emailResp EmailResponse
		if err := json.Unmarshal(respBody, &emailResp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			return &emailResp, fmt.Errorf("email service returned status %d: %s", resp.StatusCode, emailResp.Error)
		}
	}

	return &EmailResponse{Status: "sent"}, nil
}

// getFileName extracts filename from path
func getFileName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}

// isAttachmentEndpoint checks if URL is the attachment endpoint
func isAttachmentEndpoint(url string) bool {
	return len(url) > len("/api/send-mail/attached") &&
		url[len(url)-len("/api/send-mail/attached"):] == "/api/send-mail/attached"
}

// SendEmailOnly ส่งอีเมลคนเดียว (รับ string แปลงเป็น array อัตโนมัติ)
func (c *Client) SendEmailOnly(toEmail, subject, htmlBody string) error {

	fmt.Println("📧 Sending single email to", toEmail)
	req := &EmailRequest{
		ToEmails: []string{toEmail},
		Subject:  subject,
		HTMLBody: htmlBody,
	}

	_, err := c.SendMicroservice(req)
	return err
}

// SendBulk ส่งอีเมลหลายคนพร้อมกัน (รับ array โดยตรง)
func (c *Client) SendMicroserviceBulk(toEmails []string, subject, htmlBody string) error {
	req := &EmailRequest{
		ToEmails: toEmails,
		Subject:  subject,
		HTMLBody: htmlBody,
	}

	_, err := c.SendMicroservice(req)
	return err
}

// UnmarshalJSON handles both string and boolean error fields from the microservice
func (e *EmailResponse) UnmarshalJSON(data []byte) error {
	type Alias EmailResponse
	aux := &struct {
		Error interface{} `json:"error"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Error != nil {
		switch v := aux.Error.(type) {
		case string:
			e.Error = v
		case bool:
			if v {
				e.Error = "true"
			} else {
				e.Error = ""
			}
		}
	}
	return nil
}
