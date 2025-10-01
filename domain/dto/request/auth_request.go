package request

type RegisterRequest struct {
	Nama         string `json:"nama" validate:"required"`
	KataSandi    string `json:"kata_sandi" validate:"required,min=6"`
	NoTelp       string `json:"no_telp" validate:"required"`
	TanggalLahir string `json:"tanggal_lahir" validate:"required"`
	JenisKelamin string `json:"jenis_kelamin" validate:"required,oneof=L P"`
	Tentang      string `json:"tentang" validate:"required"`
	Pekerjaan    string `json:"pekerjaan" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	IDProvinsi   string `json:"id_provinsi" validate:"required"`
	IDKota       string `json:"id_kota" validate:"required"`
}

type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	KataSandi string `json:"kata_sandi" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	KataSandi   string `json:"kata_sandi" validate:"required,min=6"`
	ConfirmPass string `json:"confirm_password" validate:"required,eqfield=KataSandi"`
}
