package request

type CreateProductRequest struct {
	NamaProduk    string `json:"nama_produk" validate:"required"`
	HargaReseller string `json:"harga_reseller" validate:"required"`
	HargaKonsumen string `json:"harga_konsumen" validate:"required"`
	Stok          int    `json:"stok" validate:"required,min=0"`
	Deskripsi     string `json:"deskripsi" validate:"required"`
	IDToko        int    `json:"id_toko" validate:"required"`
	IDCategory    int    `json:"id_category" validate:"required"`
}

type UpdateProductRequest struct {
	NamaProduk    string `json:"nama_produk" validate:"required"`
	HargaReseller string `json:"harga_reseller" validate:"required"`
	HargaKonsumen string `json:"harga_konsumen" validate:"required"`
	Stok          int    `json:"stok" validate:"required,min=0"`
	Deskripsi     string `json:"deskripsi" validate:"required"`
	IDCategory    int    `json:"id_category" validate:"required"`
}
