package request

type CreateShopRequest struct {
	NamaToko string `json:"nama_toko" validate:"required"`
	URLToko  string `json:"url_toko" validate:"required"`
}

type UpdateShopRequest struct {
	NamaToko string `json:"nama_toko" validate:"required"`
	URLToko  string `json:"url_toko" validate:"required"`
}
