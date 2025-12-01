# Payment Gateway Testing Guide

Dokumentasi ini menjelaskan cara melakukan testing payment flow dengan Midtrans Sandbox.

## Prerequisites

1. **Akun Midtrans Sandbox**
   - Daftar di [Midtrans Dashboard](https://dashboard.midtrans.com/)
   - Login dan buat akun sandbox (gratis)
   - Dapatkan **Server Key** dan **Client Key** dari dashboard

2. **Environment Setup**
   - Pastikan backend dan frontend sudah running
   - Database sudah ter-setup dengan migration payment fields

## Membuat Akun Midtrans Sandbox (Step by Step)

### Langkah 1: Daftar Akun Midtrans

1. **Buka website Midtrans**
   - Kunjungi: https://dashboard.midtrans.com/
   - Klik tombol **"Sign Up"** atau **"Daftar"** di pojok kanan atas

2. **Isi Form Registrasi**
   - **Email**: Masukkan email yang valid (akan digunakan untuk login)
   - **Password**: Buat password yang kuat (minimal 8 karakter)
   - **Nama Lengkap**: Masukkan nama Anda
   - **Nomor Telepon**: Masukkan nomor telepon yang valid
   - **Perusahaan/Nama Bisnis**: Masukkan nama bisnis atau project Anda (contoh: "Marketplace Project")
   - **Jenis Bisnis**: Pilih kategori yang sesuai (contoh: "E-commerce", "Retail", dll)
   - **Alamat**: Masukkan alamat lengkap

3. **Verifikasi Email**
   - Setelah submit form, cek email Anda
   - Klik link verifikasi yang dikirim oleh Midtrans
   - Jika tidak ada email, cek folder **Spam/Junk**

4. **Login ke Dashboard**
   - Setelah email terverifikasi, kembali ke https://dashboard.midtrans.com/
   - Login dengan email dan password yang sudah dibuat

### Langkah 2: Setup Sandbox Environment

1. **Pilih Environment**
   - Setelah login, di **sidebar kiri** Anda akan melihat dropdown **"Environment"**
   - **PENTING**: Pastikan memilih **"Sandbox"** (bukan "Production")
   - Untuk testing, pilih **"Sandbox"** (gratis, tidak perlu verifikasi)
   - **Production** memerlukan verifikasi dokumen bisnis dan business registration

2. **Cara Switch ke Sandbox:**
   - Di sidebar kiri, cari dropdown **"Environment"** (biasanya di bawah logo Midtrans)
   - Klik dropdown tersebut
   - Pilih **"Sandbox"**
   - Dashboard akan refresh dan masuk ke mode Sandbox

3. **Catatan Penting:**
   - Jika Anda melihat **"Business registration: Not started"** di dashboard, **ABAIKAN** untuk Sandbox
   - Business registration hanya diperlukan untuk **Production**, bukan untuk Sandbox
   - Di Sandbox, Anda bisa langsung menggunakan payment gateway tanpa business registration

### Langkah 3: Dapatkan Access Keys

1. **Buka Settings**
   - Di **sidebar kiri**, scroll ke bawah dan cari menu **"SETTINGS"** (ikon gear ⚙️)
   - Klik **"SETTINGS"** → akan muncul submenu
   - Klik **"Access Keys"** dari submenu

2. **Copy Server Key**
   - Di halaman Access Keys, Anda akan melihat:
     - **Merchant ID**: ID merchant Anda
     - **Server Key**: Key untuk backend (rahasia, jangan share)
     - **Client Key**: Key untuk frontend (bisa di-share)
   - **Server Key** bisa memiliki format berbeda:
     - Format lama: `SB-Mid-server-xxxxxxxxxxxxx` (dengan prefix SB-)
     - Format baru: `Mid-server-xxxxxxxxxxxxx` (tanpa prefix SB-)
   - **PENTING**: Yang penting adalah **Environment sudah di Sandbox**, bukan format prefix key-nya!
   - Jika key tersembunyi, klik tombol **"Show"** atau icon **mata** 👁️ untuk melihat
   - Klik icon **copy** 📋 untuk copy Server Key
   - **PENTING**: Simpan Server Key dengan aman (jangan commit ke Git!)

3. **Copy Client Key**
   - **Client Key** juga bisa memiliki format berbeda:
     - Format lama: `SB-Mid-client-xxxxxxxxxxxxx` (dengan prefix SB-)
     - Format baru: `Mid-client-xxxxxxxxxxxxx` (tanpa prefix SB-)
   - Copy Client Key juga (untuk frontend, jika diperlukan)
   - Client Key bisa di-share (tidak se-rahasia Server Key)

4. **Verifikasi Keys:**
   - **Yang PENTING**: Pastikan **Environment dropdown menunjukkan "Sandbox"** (bukan "Production")
   - Format keys bisa berbeda (dengan atau tanpa prefix SB-), yang penting Environment-nya Sandbox
   - Jika Environment sudah "Sandbox", keys tersebut adalah Sandbox keys dan aman untuk testing
   - Contoh keys yang valid di Sandbox:
     - `Mid-server-xxxxxxxxxxxxx` ✅ (Sandbox, format baru)
     - `SB-Mid-server-xxxxxxxxxxxxx` ✅ (Sandbox, format lama)
     - `Mid-client-xxxxxxxxxxxxx` ✅ (Sandbox, format baru)
     - `SB-Mid-client-xxxxxxxxxxxxx` ✅ (Sandbox, format lama)

5. **Contoh Format Keys:**
   ```
   # Format Baru (tanpa prefix SB-)
   Server Key: Mid-server-xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   Client Key: Mid-client-xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   
   # Format Lama (dengan prefix SB-)
   Server Key: SB-Mid-server-xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   Client Key: SB-Mid-client-xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   
   # Keduanya valid untuk Sandbox!
   ```

### Langkah 4: Setup Webhook URL (Opsional untuk Testing)

**Pertanyaan: Apakah perlu dilakukan?**

**Jawaban: TIDAK WAJIB untuk testing awal, tapi DIREKOMENDASIKAN untuk testing lengkap.**

#### Kapan Webhook Diperlukan:
- ✅ **Testing automatic payment status update**: Webhook akan otomatis update status payment setelah user bayar di Midtrans
- ✅ **Testing real-time notification**: Untuk simulasi production environment
- ✅ **Testing email notification**: Email success/expired akan terkirim otomatis via webhook
- ✅ **Production**: WAJIB untuk production environment

#### Kapan Bisa Di-Skip:
- ✅ **Testing awal**: Jika hanya ingin test payment creation dan redirect ke Midtrans
- ✅ **Development**: Bisa skip dulu, test manual check payment status
- ✅ **Local testing tanpa ngrok**: Jika tidak ingin setup ngrok

#### Alternatif Tanpa Webhook:
Jika tidak setup webhook, Anda bisa:
1. **Manual Check Payment Status**: 
   - Gunakan endpoint `POST /api/v1/trx/:id/check-payment`
   - Atau klik tombol "Periksa Status Pembayaran" di payment status page
2. **Test Payment Flow**:
   - Create transaction → Redirect ke Midtrans → Complete payment
   - Kembali ke payment status page → Klik "Periksa Status Pembayaran"
   - Status akan ter-update setelah manual check

#### Cara Setup Webhook (Jika Ingin):

1. **Install & Setup Ngrok** (untuk local testing):
   ```bash
   # Download ngrok dari: https://ngrok.com/
   # Atau install via Homebrew (Mac):
   brew install ngrok
   
   # Start ngrok tunnel (pastikan backend running di port 8080)
   ngrok http 8080
   ```
   - Ngrok akan memberikan URL seperti: `https://abc123.ngrok-free.app`
   - Copy URL tersebut

2. **Buka Settings → Configuration di Midtrans Dashboard**
   - Di sidebar, klik **"Settings"** → **"Configuration"**
   - Scroll ke bagian **"Payment Notification URL"**

3. **Set Webhook URL**
   - Masukkan webhook URL: `https://abc123.ngrok-free.app/api/v1/payment/webhook`
   - Ganti `abc123.ngrok-free.app` dengan URL ngrok Anda
   - Klik **"Save"**

4. **Catatan Penting:**
   - Webhook URL harus accessible dari internet (makanya perlu ngrok untuk local)
   - Setiap kali restart ngrok, URL akan berubah (perlu update di dashboard)
   - Untuk production, gunakan domain yang valid dan permanent
   - Sandbox webhook akan mengirim notifikasi ke URL ini setiap ada perubahan status payment

#### Testing Webhook (Setelah Setup):
1. Create transaction dengan payment method non-COD
2. Complete payment di Midtrans (gunakan testing tools)
3. **Webhook akan otomatis terkirim** ke backend
4. Status payment akan ter-update otomatis tanpa perlu manual check
5. Email notification akan terkirim (jika email service configured)

**Kesimpulan:**
- Untuk **testing awal**: Bisa skip, gunakan manual check payment status
- Untuk **testing lengkap**: Setup webhook untuk simulasi production environment
- Untuk **production**: WAJIB setup webhook dengan domain yang valid

### Langkah 5: Test Connection (Verifikasi Setup)

1. **Cek API Status**
   - Di dashboard, buka **"Settings"** → **"Access Keys"**
   - Pastikan status menunjukkan **"Active"** atau **"Connected"**

2. **Test dengan API Call** (Opsional)
   ```bash
   # Test dengan curl (ganti YOUR_SERVER_KEY dengan Server Key Anda)
   curl -X GET \
     https://api.sandbox.midtrans.com/v2/ping \
     -H "Authorization: Basic $(echo -n 'YOUR_SERVER_KEY:' | base64)"
   ```
   - Jika berhasil, akan return status OK

### Langkah 6: Tambahkan Keys ke Environment

1. **Buka file `.env` di backend**
   ```bash
   cd marketplace-backend
   nano .env  # atau gunakan editor lain
   ```

2. **Tambahkan konfigurasi Midtrans**
   ```env
   # Midtrans Configuration
   # Copy Server Key dan Client Key dari dashboard Midtrans
   # Format bisa dengan atau tanpa prefix SB- (keduanya valid)
   MIDTRANS_SERVER_KEY=Mid-server-xxxxxxxxxxxxx
   MIDTRANS_CLIENT_KEY=Mid-client-xxxxxxxxxxxxx
   MIDTRANS_IS_PRODUCTION=false
   FRONTEND_URL=http://localhost:5173
   ```
   
   **Catatan:**
   - Copy **Server Key** dan **Client Key** langsung dari dashboard Midtrans
   - Format keys bisa berbeda (dengan atau tanpa prefix SB-), keduanya valid
   - Yang penting: pastikan Environment di dashboard adalah **"Sandbox"**

3. **Restart Backend**
   ```bash
   # Stop server (Ctrl+C)
   # Start lagi
   go run main.go
   ```

### Troubleshooting Akun Midtrans

#### Problem: Email verifikasi tidak terkirim
**Solution:**
- Cek folder **Spam/Junk** email
- Pastikan email yang digunakan valid
- Request ulang verifikasi email dari dashboard
- Gunakan email provider yang reliable (Gmail, Outlook, dll)

#### Problem: Tidak bisa login setelah verifikasi
**Solution:**
- Pastikan menggunakan email dan password yang benar
- Coba reset password jika lupa
- Pastikan browser tidak memblokir cookies
- Coba gunakan browser lain atau incognito mode

#### Problem: Server Key tidak muncul
**Solution:**
- Pastikan sudah login ke **Sandbox** environment (cek dropdown "Environment" di sidebar)
- Pastikan akun sudah terverifikasi
- Refresh halaman atau logout/login lagi
- Hubungi support Midtrans jika masih bermasalah

#### Problem: Keys tidak ada prefix "SB-"
**Solution:**
- **INI NORMAL!** Format keys bisa berbeda-beda
- Yang penting adalah **Environment sudah di "Sandbox"** (cek dropdown di sidebar)
- Keys dengan format `Mid-server-xxx` atau `Mid-client-xxx` tetap valid untuk Sandbox
- Format lama: `SB-Mid-server-xxx` (dengan prefix SB-)
- Format baru: `Mid-server-xxx` (tanpa prefix SB-)
- Keduanya valid selama Environment di Sandbox!

#### Problem: Webhook tidak diterima
**Solution:**
- Pastikan webhook URL accessible dari internet (gunakan ngrok)
- Pastikan URL format benar: `https://your-url/api/v1/payment/webhook`
- Cek backend logs untuk melihat apakah webhook diterima
- Test webhook manual dengan curl (lihat section Testing)

### Informasi Penting

1. **Sandbox vs Production**
   - **Sandbox**: Gratis, untuk testing, tidak perlu verifikasi
   - **Production**: Perlu verifikasi dokumen, untuk transaksi real

2. **Sandbox Limitations**
   - Tidak ada transaksi real (semua test)
   - Tidak ada biaya
   - Data bisa di-reset kapan saja
   - Perfect untuk development dan testing

3. **Security**
   - **JANGAN** commit Server Key ke Git
   - Gunakan `.env` file dan tambahkan ke `.gitignore`
   - Jangan share Server Key dengan siapapun
   - Client Key bisa di-share (untuk frontend)

4. **Support**
   - Dokumentasi: https://docs.midtrans.com/
   - Support: support@midtrans.com
   - Community: https://github.com/Midtrans

### Next Steps

Setelah akun Midtrans sandbox sudah setup:
1. ✅ Tambahkan keys ke `.env` file
2. ✅ Restart backend server
3. ✅ Lanjut ke section **"Setup Environment Variables"** di bawah
4. ✅ Ikuti **"Testing Checklist"** untuk test payment flow

## Setup Environment Variables

Tambahkan konfigurasi berikut ke file `.env` di `marketplace-backend`:

```env
# Midtrans Configuration
MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxxxxxxxxxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxxxxxxxxxx
MIDTRANS_IS_PRODUCTION=false

# Frontend URL (untuk redirect setelah payment)
FRONTEND_URL=http://localhost:5173
```

**Cara mendapatkan Server Key dan Client Key:**
- Lihat section **"Membuat Akun Midtrans Sandbox (Step by Step)"** di atas untuk panduan lengkap
- Atau langsung:
  1. Login ke [Midtrans Dashboard](https://dashboard.midtrans.com/)
  2. Pilih **Settings** → **Access Keys**
  3. Copy **Server Key** (untuk backend) dan **Client Key** (untuk frontend, jika diperlukan)
  4. Pastikan menggunakan **Sandbox** keys (bukan Production)

## Testing Checklist

### 1. Setup & Configuration Testing

- [ ] Environment variables sudah di-set dengan benar
- [ ] Backend server bisa start tanpa error
- [ ] Database migration sudah dijalankan (kolom payment sudah ada di tabel `trx`)
- [ ] Frontend bisa connect ke backend API

### 2. Payment Method Testing

#### A. Virtual Account (Bank Transfer)

**Test Case:**
1. Login sebagai user
2. Tambah produk ke cart
3. Checkout dengan payment method: `virtual_account` atau `va`
4. Submit checkout
5. **Expected:**
   - Redirect ke payment status page
   - Auto-redirect ke Midtrans payment page
   - Tampil Virtual Account number (BCA, BNI, atau Mandiri)
   - Status payment: `pending_payment`

**Testing Payment:**
- Gunakan **Midtrans Testing Tools** di dashboard untuk simulate payment
- Atau transfer ke VA number yang ditampilkan (sandbox mode)
- Cek status payment di payment status page

#### B. E-Wallet (GoPay, OVO, DANA, LinkAja)

**Test Case:**
1. Checkout dengan payment method: `e_wallet`, `gopay`, `ovo`, `dana`, atau `linkaja`
2. **Expected:**
   - Redirect ke Midtrans payment page
   - Tampil QR code atau deep link untuk e-wallet
   - Status payment: `pending_payment`

**Testing Payment:**
- Scan QR code dengan aplikasi e-wallet (sandbox mode)
- Atau gunakan Midtrans testing tools untuk simulate payment

#### C. Bank Transfer

**Test Case:**
1. Checkout dengan payment method: `bank_transfer`, `bank_transfer_bca`, `bank_transfer_bni`, atau `bank_transfer_mandiri`
2. **Expected:**
   - Redirect ke Midtrans payment page
   - Tampil instruksi transfer bank
   - Status payment: `pending_payment`

#### D. Credit Card

**Test Case:**
1. Checkout dengan payment method: `credit_card` atau `cc`
2. **Expected:**
   - Redirect ke Midtrans payment page
   - Tampil form credit card
   - Status payment: `pending_payment`

**Testing Payment:**
- Gunakan test card dari Midtrans:
  - **Success:** `4811111111111114`
  - **3DS Challenge:** `4811111111111114` (akan trigger 3DS)
  - **Decline:** `4911111111111113`

#### E. COD (Cash on Delivery)

**Test Case:**
1. Checkout dengan payment method: `COD` atau `cod`
2. **Expected:**
   - **TIDAK** redirect ke Midtrans
   - Langsung ke orders page
   - Status payment: `pending_payment` (atau sesuai logic COD)
   - Cart cleared

### 3. Payment Status Flow Testing

#### A. Payment Success Flow

**Test Case:**
1. Create transaction dengan payment method non-COD
2. Complete payment di Midtrans (gunakan testing tools)
3. **Expected:**
   - Webhook diterima oleh backend (`POST /api/v1/payment/webhook`)
   - Status payment berubah menjadi `paid`
   - Email notification terkirim (jika email service configured)
   - Payment status page menampilkan status "Pembayaran Berhasil"

**Verification:**
- Cek database: `trx.payment_status = 'paid'`
- Cek email inbox (jika configured)
- Cek payment status page UI

#### B. Payment Expired Flow

**Test Case:**
1. Create transaction dengan payment method non-COD
2. Tunggu sampai payment expired (default: 24 jam, bisa di-test dengan simulate)
3. **Expected:**
   - Status payment berubah menjadi `expired`
   - Email notification terkirim (jika configured)
   - Payment status page menampilkan status "Pembayaran Kadaluarsa"

**Verification:**
- Cek database: `trx.payment_status = 'expired'`
- Cek email inbox (jika configured)

#### C. Payment Failed/Cancelled Flow

**Test Case:**
1. Create transaction dengan payment method non-COD
2. Cancel atau fail payment di Midtrans
3. **Expected:**
   - Status payment berubah menjadi `failed` atau `cancelled`
   - Payment status page menampilkan status sesuai

### 4. Webhook Testing

#### Manual Webhook Test

**Test Case:**
1. Dapatkan `order_id` dari transaction yang dibuat
2. Call webhook endpoint:
   ```bash
   curl -X POST http://localhost:8080/api/v1/payment/webhook \
     -H "Content-Type: application/json" \
     -d '{
       "order_id": "INV-1234567890",
       "transaction_status": "settlement"
     }'
   ```
3. **Expected:**
   - Response: `{"status": true, "message": "Webhook processed successfully"}`
   - Transaction status ter-update di database

#### Webhook dari Midtrans

**Setup:**
1. Login ke Midtrans Dashboard
2. Go to **Settings** → **Configuration** → **Payment Notification URL**
3. Set URL: `http://your-backend-url/api/v1/payment/webhook`
   - Untuk local testing, gunakan ngrok atau similar tool:
     ```bash
     ngrok http 8080
     # Gunakan URL dari ngrok: https://xxxxx.ngrok.io/api/v1/payment/webhook
     ```

**Test Case:**
1. Complete payment di Midtrans
2. **Expected:**
   - Midtrans mengirim webhook ke URL yang di-set
   - Backend memproses webhook dan update status

### 5. Manual Payment Status Check

**Test Case:**
1. Login sebagai user
2. Call endpoint check payment:
   ```bash
   curl -X POST http://localhost:8080/api/v1/trx/{transaction_id}/check-payment \
     -H "Authorization: Bearer {token}"
   ```
3. **Expected:**
   - Response berisi transaction dengan status terbaru
   - Status ter-update di database

**Via Frontend:**
- Buka payment status page: `/payment/{transaction_id}`
- Klik tombol "Periksa Status Pembayaran"
- Status akan ter-update

### 6. Frontend Flow Testing

#### A. Checkout Flow

**Test Case:**
1. User login
2. Tambah produk ke cart
3. Go to checkout page
4. Pilih payment method
5. Fill form checkout
6. Submit checkout
7. **Expected:**
   - Redirect ke `/payment/{transaction_id}`
   - Payment status page menampilkan loading
   - Auto-redirect ke Midtrans payment page

#### B. Payment Status Page

**Test Case:**
1. Buka payment status page: `/payment/{transaction_id}`
2. **Expected:**
   - Menampilkan status payment yang benar
   - Menampilkan informasi transaksi (invoice, total, dll)
   - Tombol aksi sesuai status:
     - Pending: "Lanjutkan Pembayaran", "Periksa Status"
     - Paid: "Lihat Pesanan"
     - Expired/Failed: "Belanja Lagi"

#### C. Redirect dari Midtrans

**Test Case:**
1. Complete payment di Midtrans
2. Midtrans redirect ke `finish_url` (yang sudah di-set: `/payment/{transaction_id}`)
3. **Expected:**
   - Payment status page auto-check status
   - Status ter-update dan UI menampilkan status terbaru

### 7. Error Handling Testing

#### A. Invalid Payment Method

**Test Case:**
1. Submit checkout dengan payment method yang tidak valid
2. **Expected:**
   - Validation error
   - Transaction tidak dibuat

#### B. Payment Creation Failed

**Test Case:**
1. Set invalid Midtrans Server Key
2. Submit checkout dengan payment method non-COD
3. **Expected:**
   - Error message ditampilkan
   - Transaction dibuat dengan status `failed`

#### C. Webhook dengan Invalid Data

**Test Case:**
1. Call webhook dengan `order_id` yang tidak ada
2. **Expected:**
   - Response error: "transaction not found"
   - Status code: 200 (Midtrans expects 200 OK)

## Testing Tools

### 1. Midtrans Dashboard Testing Tools

- Login ke [Midtrans Dashboard](https://dashboard.midtrans.com/)
- Go to **Transactions** → **Sandbox Transactions**
- Bisa simulate payment status changes

### 2. Midtrans Testing Cards

**Credit Card Testing:**
- Success: `4811111111111114`
- 3DS Challenge: `4811111111111114` (akan trigger 3DS)
- Decline: `4911111111111113`

**Expiry Date:** Any future date (e.g., `12/25`)
**CVV:** Any 3 digits (e.g., `123`)

### 3. Virtual Account Testing

- Gunakan VA number yang ditampilkan di payment page
- Di sandbox mode, tidak perlu transfer real money
- Gunakan Midtrans testing tools untuk simulate payment

### 4. E-Wallet Testing

- Scan QR code dengan aplikasi e-wallet (sandbox mode)
- Atau gunakan deep link yang ditampilkan
- Di sandbox mode, tidak perlu real payment

## Common Issues & Solutions

### Issue: Payment URL tidak ter-generate

**Solution:**
- Cek Midtrans Server Key sudah benar
- Cek `MIDTRANS_IS_PRODUCTION=false` untuk sandbox
- Cek logs backend untuk error message

### Issue: Webhook tidak diterima

**Solution:**
- Pastikan webhook URL accessible dari internet (gunakan ngrok untuk local)
- Cek webhook URL di Midtrans dashboard sudah benar
- Cek backend logs untuk melihat apakah webhook diterima

### Issue: Status payment tidak ter-update

**Solution:**
- Cek webhook endpoint bisa diakses
- Cek database untuk melihat apakah status ter-update
- Manual check payment status via API

### Issue: Redirect tidak bekerja

**Solution:**
- Cek `FRONTEND_URL` sudah benar di `.env`
- Cek payment status page route sudah terdaftar
- Cek browser console untuk error

## Production Checklist

Sebelum deploy ke production:

- [ ] Ganti ke Production Midtrans keys
- [ ] Set `MIDTRANS_IS_PRODUCTION=true`
- [ ] Setup webhook URL yang accessible dari internet
- [ ] Test semua payment methods
- [ ] Setup monitoring untuk webhook
- [ ] Setup error logging dan alerting
- [ ] Test email notifications
- [ ] Review security (HTTPS, webhook verification, dll)

## Additional Resources

- [Midtrans Documentation](https://docs.midtrans.com/)
- [Midtrans Sandbox Testing](https://docs.midtrans.com/docs/core-api/overview)
- [Midtrans Dashboard](https://dashboard.midtrans.com/)


