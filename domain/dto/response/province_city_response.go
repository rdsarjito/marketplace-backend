package response

type ProvinceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CityResponse struct {
	ID       string `json:"id"`
	ProvID   string `json:"prov_id"`
	Name     string `json:"name"`
}

type ProvinceDetailResponse struct {
	ID   string         `json:"id"`
	Name string         `json:"name"`
	Cities []CityResponse `json:"cities"`
}
