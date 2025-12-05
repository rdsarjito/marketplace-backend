package request

type UpdateProfileRequest struct {
	Nama         string `json:"nama" validate:"omitempty"`
	NoTelp       string `json:"no_telp" validate:"omitempty"`
	TanggalLahir string `json:"tanggal_lahir" validate:"omitempty"`
	JenisKelamin string `json:"jenis_kelamin" validate:"omitempty,oneof=L P"`
	Tentang      string `json:"tentang" validate:"omitempty"`
	Pekerjaan    string `json:"pekerjaan" validate:"omitempty"`
	Email        string `json:"email" validate:"omitempty,email"`
	IDProvinsi   string `json:"id_provinsi" validate:"omitempty"`
	IDKota       string `json:"id_kota" validate:"omitempty"`
}
