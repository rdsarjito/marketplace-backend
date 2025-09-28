package request

type CreateTRXRequest struct {
	HargaTotal  int                    `json:"harga_total" validate:"required"`
	MethodBayar string                 `json:"method_bayar" validate:"required"`
	IDAlamat    int                    `json:"id_alamat" validate:"required"`
	DetailTRX   []CreateDetailTRXRequest `json:"detail_trx" validate:"required"`
}

type CreateDetailTRXRequest struct {
	IDProduk  int `json:"id_produk" validate:"required"`
	IDToko    int `json:"id_toko" validate:"required"`
	Kuantitas int `json:"kuantitas" validate:"required,min=1"`
	HargaTotal int `json:"harga_total" validate:"required"`
}
