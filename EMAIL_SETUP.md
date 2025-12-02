# Email Configuration Setup

## Environment Variables

Tambahkan konfigurasi email berikut ke file `.env`:

```bash
# Email Configuration (Optional - jika tidak diisi akan log ke console)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password
FROM_EMAIL=noreply@warungbudehramah.com
FROM_NAME=Warung Budeh Ramah
```

## Gmail Setup

### 1. Enable 2-Factor Authentication
- Buka Google Account settings
- Aktifkan 2-Factor Authentication

### 2. Generate App Password
- Buka Google Account → Security → App passwords
- Generate password untuk "Mail"
- Gunakan password ini sebagai `SMTP_PASSWORD`

### 3. Gmail SMTP Settings
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_gmail@gmail.com
SMTP_PASSWORD=your_16_char_app_password
```

## Other Email Providers

### Outlook/Hotmail
```
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
```

### Yahoo
```
SMTP_HOST=smtp.mail.yahoo.com
SMTP_PORT=587
```

## Development Mode

Jika tidak ada konfigurasi email, sistem akan:
- Log email content ke console
- Menampilkan token reset di terminal
- Tidak mengirim email sesungguhnya

## Testing

1. Set environment variables
2. Restart backend server
3. Test forgot password dengan email yang valid
4. Cek email inbox untuk link reset password
