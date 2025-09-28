package response

type ProductResponse struct {
	ID             int                `json:"id"`
	NamaProduk     string             `json:"nama_produk"`
	Slug           string             `json:"slug"`
	HargaReseller  string             `json:"harga_reseller"`
	HargaKonsumen  string             `json:"harga_konsumen"`
	Stok           int                `json:"stok"`
	Deskripsi      string             `json:"deskripsi"`
	CreatedAt      string             `json:"created_at"`
	UpdatedAt      string             `json:"updated_at"`
	IDToko         int                `json:"id_toko"`
	IDCategory     int                `json:"id_category"`
	Toko           ShopResponse       `json:"toko"`
	Category       CategoryResponse   `json:"category"`
	PhotosProduct  []PhotoProductResponse `json:"photos_product"`
}

type PhotoProductResponse struct {
	ID        int    `json:"id"`
	IDProduk  int    `json:"id_produk"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
