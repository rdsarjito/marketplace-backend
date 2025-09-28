package request

type CreateCategoryRequest struct {
	Nama string `json:"nama" validate:"required"`
}

type UpdateCategoryRequest struct {
	Nama string `json:"nama" validate:"required"`
}
