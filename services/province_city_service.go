package services

import (
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/repositories"
)

type ProvinceCityService interface {
	GetListProvince() ([]response.ProvinceResponse, error)
	GetDetailProvince(provID string) (*response.ProvinceResponse, error)
	GetListCity(provID string) ([]response.CityResponse, error)
	GetDetailCity(cityID string) (*response.CityResponse, error)
}

type provinceCityService struct {
	provinceCityRepo repositories.ProvinceCityRepository
}

func NewProvinceCityService(provinceCityRepo repositories.ProvinceCityRepository) ProvinceCityService {
	return &provinceCityService{
		provinceCityRepo: provinceCityRepo,
	}
}

func (s *provinceCityService) GetListProvince() ([]response.ProvinceResponse, error) {
	provinces, err := s.provinceCityRepo.GetListProvince()
	if err != nil {
		return nil, err
	}

	var provinceResponses []response.ProvinceResponse
	for _, province := range provinces {
		provinceResponses = append(provinceResponses, response.ProvinceResponse{
			ID:   province["id"].(string),
			Name: province["name"].(string),
		})
	}

	return provinceResponses, nil
}

func (s *provinceCityService) GetDetailProvince(provID string) (*response.ProvinceResponse, error) {
	province, err := s.provinceCityRepo.GetDetailProvince(provID)
	if err != nil {
		return nil, err
	}

	provinceResponse := &response.ProvinceResponse{
		ID:   province["id"].(string),
		Name: province["name"].(string),
	}

	return provinceResponse, nil
}

func (s *provinceCityService) GetListCity(provID string) ([]response.CityResponse, error) {
	cities, err := s.provinceCityRepo.GetListCity(provID)
	if err != nil {
		return nil, err
	}

	var cityResponses []response.CityResponse
	for _, city := range cities {
		cityResponses = append(cityResponses, response.CityResponse{
			ID:     city["id"].(string),
			ProvID: city["province_id"].(string),
			Name:   city["name"].(string),
		})
	}

	return cityResponses, nil
}

func (s *provinceCityService) GetDetailCity(cityID string) (*response.CityResponse, error) {
	city, err := s.provinceCityRepo.GetDetailCity(cityID)
	if err != nil {
		return nil, err
	}

	cityResponse := &response.CityResponse{
		ID:     city["id"].(string),
		ProvID: city["province_id"].(string),
		Name:   city["name"].(string),
	}

	return cityResponse, nil
}
