package response

type ShopResponse struct {
	ID        int    `json:"id"`
	NamaToko  string `json:"nama_toko"`
	URLToko   string `json:"url_toko"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IDUser    int    `json:"id_user"`
}
