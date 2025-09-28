package response

type CategoryResponse struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
