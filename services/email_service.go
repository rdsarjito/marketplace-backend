package services

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendPasswordResetEmail(email, token string) error
	SendPaymentSuccessEmail(email, invoiceCode string, totalAmount int) error
	SendPaymentExpiredEmail(email, invoiceCode string, totalAmount int) error
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
		fromEmail:    getEnv("FROM_EMAIL", "noreply@warungbudehramah.com"),
		fromName:     getEnv("FROM_NAME", "Warung Budeh Ramah"),
	}
}

func (s *emailService) SendPasswordResetEmail(email, token string) error {
	// Jika tidak ada konfigurasi SMTP, log ke console (untuk development)
	if s.smtpUsername == "" || s.smtpPassword == "" {
		fmt.Printf("=== EMAIL RESET PASSWORD ===\n")
		fmt.Printf("To: %s\n", email)
		fmt.Printf("Subject: Reset Password - Warung Budeh Ramah\n")
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
		<title>Reset Password - Warung Budeh Ramah</title>
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
				<h1>Warung Budeh Ramah</h1>
			</div>
			<div class="content">
				<h2>Reset Password</h2>
				<p>Halo,</p>
				<p>Kami menerima permintaan untuk mereset password akun Warung Budeh Ramah Anda.</p>
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
				<p>&copy; 2024 Warung Budeh Ramah. All rights reserved.</p>
			</div>
		</div>
	</body>
	</html>
	`, resetURL, resetURL)

	// Template email plain text
	textBody := fmt.Sprintf(`
Reset Password - Warung Budeh Ramah

Halo,

Kami menerima permintaan untuk mereset password akun Warung Budeh Ramah Anda.

Klik link berikut untuk mereset password Anda:
%s

Catatan penting:
- Link ini hanya berlaku selama 24 jam
- Link ini hanya bisa digunakan sekali
- Jika Anda tidak meminta reset password, abaikan email ini

Email ini dikirim secara otomatis, mohon tidak membalas email ini.

© 2024 Warung Budeh Ramah. All rights reserved.
	`, resetURL)

	// Buat email message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Reset Password - Warung Budeh Ramah")
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

func (s *emailService) SendPaymentSuccessEmail(email, invoiceCode string, totalAmount int) error {
	// Jika tidak ada konfigurasi SMTP, log ke console (untuk development)
	if s.smtpUsername == "" || s.smtpPassword == "" {
		fmt.Printf("=== EMAIL PAYMENT SUCCESS ===\n")
		fmt.Printf("To: %s\n", email)
		fmt.Printf("Subject: Pembayaran Berhasil - Warung Budeh Ramah\n")
		fmt.Printf("Invoice: %s\n", invoiceCode)
		fmt.Printf("Total: Rp %d\n", totalAmount)
		fmt.Printf("=============================\n")
		return nil
	}

	// Format total amount
	totalAmountStr := fmt.Sprintf("Rp %s", formatCurrency(totalAmount))

	// Template email HTML
	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Pembayaran Berhasil - Warung Budeh Ramah</title>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background-color: #03AC0E; color: white; padding: 20px; text-align: center; }
			.content { padding: 30px; background-color: #f9f9f9; }
			.success-badge { background-color: #03AC0E; color: white; padding: 10px 20px; border-radius: 5px; display: inline-block; margin: 20px 0; }
			.info-box { background-color: #fff; border: 1px solid #ddd; border-radius: 5px; padding: 15px; margin: 20px 0; }
			.info-row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #eee; }
			.info-row:last-child { border-bottom: none; }
			.footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Warung Budeh Ramah</h1>
			</div>
			<div class="content">
				<div style="text-align: center;">
					<span class="success-badge">✓ Pembayaran Berhasil</span>
				</div>
				<h2>Terima Kasih!</h2>
				<p>Halo,</p>
				<p>Pembayaran Anda telah berhasil diproses. Pesanan Anda sedang dipersiapkan.</p>
				<div class="info-box">
					<div class="info-row">
						<strong>Nomor Invoice:</strong>
						<span>%s</span>
					</div>
					<div class="info-row">
						<strong>Total Pembayaran:</strong>
						<span>%s</span>
					</div>
				</div>
				<p>Kami akan mengirimkan update lebih lanjut mengenai status pengiriman pesanan Anda.</p>
				<p>Jika Anda memiliki pertanyaan, silakan hubungi customer service kami.</p>
			</div>
			<div class="footer">
				<p>Email ini dikirim secara otomatis, mohon tidak membalas email ini.</p>
				<p>&copy; 2024 Warung Budeh Ramah. All rights reserved.</p>
			</div>
		</div>
	</body>
	</html>
	`, invoiceCode, totalAmountStr)

	// Template email plain text
	textBody := fmt.Sprintf(`
Pembayaran Berhasil - Warung Budeh Ramah

Halo,

Pembayaran Anda telah berhasil diproses. Pesanan Anda sedang dipersiapkan.

Nomor Invoice: %s
Total Pembayaran: %s

Kami akan mengirimkan update lebih lanjut mengenai status pengiriman pesanan Anda.

Jika Anda memiliki pertanyaan, silakan hubungi customer service kami.

Email ini dikirim secara otomatis, mohon tidak membalas email ini.

© 2024 Warung Budeh Ramah. All rights reserved.
	`, invoiceCode, totalAmountStr)

	// Buat email message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Pembayaran Berhasil - Warung Budeh Ramah")
	m.SetBody("text/plain", textBody)
	m.AddAlternative("text/html", htmlBody)

	// Kirim email
	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUsername, s.smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (s *emailService) SendPaymentExpiredEmail(email, invoiceCode string, totalAmount int) error {
	// Jika tidak ada konfigurasi SMTP, log ke console (untuk development)
	if s.smtpUsername == "" || s.smtpPassword == "" {
		fmt.Printf("=== EMAIL PAYMENT EXPIRED ===\n")
		fmt.Printf("To: %s\n", email)
		fmt.Printf("Subject: Pembayaran Kadaluarsa - Warung Budeh Ramah\n")
		fmt.Printf("Invoice: %s\n", invoiceCode)
		fmt.Printf("Total: Rp %d\n", totalAmount)
		fmt.Printf("=============================\n")
		return nil
	}

	// Format total amount
	totalAmountStr := fmt.Sprintf("Rp %s", formatCurrency(totalAmount))

	// Template email HTML
	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Pembayaran Kadaluarsa - Warung Budeh Ramah</title>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background-color: #FF6B6B; color: white; padding: 20px; text-align: center; }
			.content { padding: 30px; background-color: #f9f9f9; }
			.warning-badge { background-color: #FF6B6B; color: white; padding: 10px 20px; border-radius: 5px; display: inline-block; margin: 20px 0; }
			.info-box { background-color: #fff; border: 1px solid #ddd; border-radius: 5px; padding: 15px; margin: 20px 0; }
			.info-row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #eee; }
			.info-row:last-child { border-bottom: none; }
			.footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Warung Budeh Ramah</h1>
			</div>
			<div class="content">
				<div style="text-align: center;">
					<span class="warning-badge">⚠ Pembayaran Kadaluarsa</span>
				</div>
				<h2>Pembayaran Tidak Ditemukan</h2>
				<p>Halo,</p>
				<p>Kami ingin memberitahu bahwa waktu pembayaran untuk pesanan Anda telah kadaluarsa.</p>
				<div class="info-box">
					<div class="info-row">
						<strong>Nomor Invoice:</strong>
						<span>%s</span>
					</div>
					<div class="info-row">
						<strong>Total Pembayaran:</strong>
						<span>%s</span>
					</div>
				</div>
				<p>Pesanan Anda telah dibatalkan karena pembayaran tidak dilakukan dalam waktu yang ditentukan.</p>
				<p>Jika Anda masih ingin melakukan pembayaran, silakan buat pesanan baru.</p>
				<p>Jika Anda memiliki pertanyaan, silakan hubungi customer service kami.</p>
			</div>
			<div class="footer">
				<p>Email ini dikirim secara otomatis, mohon tidak membalas email ini.</p>
				<p>&copy; 2024 Warung Budeh Ramah. All rights reserved.</p>
			</div>
		</div>
	</body>
	</html>
	`, invoiceCode, totalAmountStr)

	// Template email plain text
	textBody := fmt.Sprintf(`
Pembayaran Kadaluarsa - Warung Budeh Ramah

Halo,

Kami ingin memberitahu bahwa waktu pembayaran untuk pesanan Anda telah kadaluarsa.

Nomor Invoice: %s
Total Pembayaran: %s

Pesanan Anda telah dibatalkan karena pembayaran tidak dilakukan dalam waktu yang ditentukan.

Jika Anda masih ingin melakukan pembayaran, silakan buat pesanan baru.

Jika Anda memiliki pertanyaan, silakan hubungi customer service kami.

Email ini dikirim secara otomatis, mohon tidak membalas email ini.

© 2024 Warung Budeh Ramah. All rights reserved.
	`, invoiceCode, totalAmountStr)

	// Buat email message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Pembayaran Kadaluarsa - Warung Budeh Ramah")
	m.SetBody("text/plain", textBody)
	m.AddAlternative("text/html", htmlBody)

	// Kirim email
	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUsername, s.smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// formatCurrency formats number to Indonesian currency format
func formatCurrency(amount int) string {
	amountStr := fmt.Sprintf("%d", amount)
	result := ""
	for i, char := range amountStr {
		if i > 0 && (len(amountStr)-i)%3 == 0 {
			result += "."
		}
		result += string(char)
	}
	return result
}
