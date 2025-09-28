package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ProvinceCityRepository interface {
	GetListProvince() ([]map[string]interface{}, error)
	GetDetailProvince(provID string) (map[string]interface{}, error)
	GetListCity(provID string) ([]map[string]interface{}, error)
	GetDetailCity(cityID string) (map[string]interface{}, error)
}

type provinceCityRepository struct {
	apiURL string
	client *http.Client
}

func NewProvinceCityRepository(apiURL string) ProvinceCityRepository {
	return &provinceCityRepository{
		apiURL: apiURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *provinceCityRepository) GetListProvince() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/provinces.json", r.apiURL)
	return r.makeRequest(url)
}

func (r *provinceCityRepository) GetDetailProvince(provID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/provinces/%s.json", r.apiURL, provID)
	result, err := r.makeRequest(url)
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, fmt.Errorf("province not found")
}

func (r *provinceCityRepository) GetListCity(provID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/regencies/%s.json", r.apiURL, provID)
	return r.makeRequest(url)
}

func (r *provinceCityRepository) GetDetailCity(cityID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/regencies/%s.json", r.apiURL, cityID)
	result, err := r.makeRequest(url)
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, fmt.Errorf("city not found")
}

func (r *provinceCityRepository) makeRequest(url string) ([]map[string]interface{}, error) {
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
