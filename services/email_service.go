package services

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendPasswordResetEmail(email, token string) error
}

type emailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
	fromName     string
}

func NewEmailService() EmailService {
	return &emailService{
		smtpHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort:     getEnvInt("SMTP_PORT", 587),
		smtpUsername: getEnv("SMTP_USERNAME", ""),
		smtpPassword: getEnv("SMTP_PASSWORD", ""),
		fromEmail:    getEnv("FROM_EMAIL", "noreply@evermos.com"),
		fromName:     getEnv("FROM_NAME", "Evermos"),
	}
}

func (s *emailService) SendPasswordResetEmail(email, token string) error {
	// Jika tidak ada konfigurasi SMTP, log ke console (untuk development)
	if s.smtpUsername == "" || s.smtpPassword == "" {
		fmt.Printf("=== EMAIL RESET PASSWORD ===\n")
		fmt.Printf("To: %s\n", email)
		fmt.Printf("Subject: Reset Password - Evermos\n")
		fmt.Printf("Token: %s\n", token)
		fmt.Printf("Reset URL: http://localhost:5173/reset-password?token=%s\n", token)
		fmt.Printf("=============================\n")
		return nil
	}

	// Buat link reset password
	resetURL := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", token)

	// Template email HTML
	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Reset Password - Evermos</title>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background-color: #03AC0E; color: white; padding: 20px; text-align: center; }
			.content { padding: 30px; background-color: #f9f9f9; }
			.button { display: inline-block; background-color: #03AC0E; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; margin: 20px 0; }
			.footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Evermos</h1>
			</div>
			<div class="content">
				<h2>Reset Password</h2>
				<p>Halo,</p>
				<p>Kami menerima permintaan untuk mereset password akun Evermos Anda.</p>
				<p>Klik tombol di bawah ini untuk mereset password Anda:</p>
				<p style="text-align: center;">
					<a href="%s" class="button">Reset Password</a>
				</p>
				<p>Atau copy dan paste link berikut ke browser Anda:</p>
				<p style="word-break: break-all; background-color: #eee; padding: 10px; border-radius: 3px;">
					%s
				</p>
				<p><strong>Catatan penting:</strong></p>
				<ul>
					<li>Link ini hanya berlaku selama 24 jam</li>
					<li>Link ini hanya bisa digunakan sekali</li>
					<li>Jika Anda tidak meminta reset password, abaikan email ini</li>
				</ul>
			</div>
			<div class="footer">
				<p>Email ini dikirim secara otomatis, mohon tidak membalas email ini.</p>
				<p>&copy; 2024 Evermos. All rights reserved.</p>
			</div>
		</div>
	</body>
	</html>
	`, resetURL, resetURL)

	// Template email plain text
	textBody := fmt.Sprintf(`
Reset Password - Evermos

Halo,

Kami menerima permintaan untuk mereset password akun Evermos Anda.

Klik link berikut untuk mereset password Anda:
%s

Catatan penting:
- Link ini hanya berlaku selama 24 jam
- Link ini hanya bisa digunakan sekali
- Jika Anda tidak meminta reset password, abaikan email ini

Email ini dikirim secara otomatis, mohon tidak membalas email ini.

© 2024 Evermos. All rights reserved.
	`, resetURL)

	// Buat email message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Reset Password - Evermos")
	m.SetBody("text/plain", textBody)
	m.AddAlternative("text/html", htmlBody)

	// Kirim email
	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUsername, s.smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := fmt.Sscanf(value, "%d", &defaultValue); err == nil && intValue == 1 {
			return defaultValue
		}
	}
	return defaultValue
}
