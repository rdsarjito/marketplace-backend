package response

type AuthResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type UserProfile struct {
	ID            int    `json:"id"`
	Nama          string `json:"nama"`
	NoTelp        string `json:"no_telp"`
	TanggalLahir  string `json:"tanggal_lahir"`
	JenisKelamin  string `json:"jenis_kelamin"`
	Tentang       string `json:"tentang"`
	Pekerjaan     string `json:"pekerjaan"`
	Email         string `json:"email"`
	IDProvinsi    string `json:"id_provinsi"`
	IDKota        string `json:"id_kota"`
	IsAdmin       bool   `json:"is_admin"`
}
