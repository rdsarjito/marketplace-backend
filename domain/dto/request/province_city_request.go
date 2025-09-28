package request

type GetProvinceRequest struct {
	ProvID string `json:"prov_id" validate:"required"`
}

type GetCityRequest struct {
	CityID string `json:"city_id" validate:"required"`
}
