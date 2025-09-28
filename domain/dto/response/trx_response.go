package response

type TRXResponse struct {
	ID          int                `json:"id"`
	HargaTotal  int                `json:"harga_total"`
	KodeInvoice string             `json:"kode_invoice"`
	MethodBayar string             `json:"method_bayar"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	IDUser      int                `json:"id_user"`
	IDAlamat    int                `json:"id_alamat"`
	User        UserProfile        `json:"user"`
	Address     AddressResponse    `json:"address"`
	DetailTRX   []DetailTRXResponse `json:"detail_trx"`
}

type DetailTRXResponse struct {
	ID         int              `json:"id"`
	IDTRX      int              `json:"id_trx"`
	IDProduk   int              `json:"id_produk"`
	IDToko     int              `json:"id_toko"`
	Kuantitas  int              `json:"kuantitas"`
	HargaTotal int              `json:"harga_total"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
	Product    ProductResponse  `json:"product"`
	Shop       ShopResponse     `json:"shop"`
}
