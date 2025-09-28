package services

import (
	"errors"

	"github.com/rdsarjito/marketplace-backend/constants"
	"github.com/rdsarjito/marketplace-backend/domain/dto/request"
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/repositories"
)

type ShopService interface {
	MyShop(userID int) (*response.ShopResponse, error)
	GetListShop() ([]response.ShopResponse, error)
	GetDetailShop(shopID int) (*response.ShopResponse, error)
	UpdateProfileShop(userID, shopID int, req *request.UpdateShopRequest) (*response.ShopResponse, error)
}

type shopService struct {
	shopRepo repositories.ShopRepository
}

func NewShopService(shopRepo repositories.ShopRepository) ShopService {
	return &shopService{
		shopRepo: shopRepo,
	}
}

func (s *shopService) MyShop(userID int) (*response.ShopResponse, error) {
	shop, err := s.shopRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New(constants.ErrShopNotFound)
	}

	shopResponse := &response.ShopResponse{
		ID:        shop.ID,
		NamaToko:  shop.NamaToko,
		URLToko:   shop.URLToko,
		CreatedAt: shop.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: shop.UpdatedAt.Format("2006-01-02 15:04:05"),
		IDUser:    shop.IDUser,
	}

	return shopResponse, nil
}

func (s *shopService) GetListShop() ([]response.ShopResponse, error) {
	shops, err := s.shopRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var shopResponses []response.ShopResponse
	for _, shop := range shops {
		shopResponses = append(shopResponses, response.ShopResponse{
			ID:        shop.ID,
			NamaToko:  shop.NamaToko,
			URLToko:   shop.URLToko,
			CreatedAt: shop.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: shop.UpdatedAt.Format("2006-01-02 15:04:05"),
			IDUser:    shop.IDUser,
		})
	}

	return shopResponses, nil
}

func (s *shopService) GetDetailShop(shopID int) (*response.ShopResponse, error) {
	shop, err := s.shopRepo.GetByID(shopID)
	if err != nil {
		return nil, errors.New(constants.ErrShopNotFound)
	}

	shopResponse := &response.ShopResponse{
		ID:        shop.ID,
		NamaToko:  shop.NamaToko,
		URLToko:   shop.URLToko,
		CreatedAt: shop.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: shop.UpdatedAt.Format("2006-01-02 15:04:05"),
		IDUser:    shop.IDUser,
	}

	return shopResponse, nil
}

func (s *shopService) UpdateProfileShop(userID, shopID int, req *request.UpdateShopRequest) (*response.ShopResponse, error) {
	shop, err := s.shopRepo.GetByID(shopID)
	if err != nil {
		return nil, errors.New(constants.ErrShopNotFound)
	}

	if shop.IDUser != userID {
		return nil, errors.New(constants.ErrForbidden)
	}

	shop.NamaToko = req.NamaToko
	shop.URLToko = req.URLToko

	if err := s.shopRepo.Update(shop); err != nil {
		return nil, err
	}

	shopResponse := &response.ShopResponse{
		ID:        shop.ID,
		NamaToko:  shop.NamaToko,
		URLToko:   shop.URLToko,
		CreatedAt: shop.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: shop.UpdatedAt.Format("2006-01-02 15:04:05"),
		IDUser:    shop.IDUser,
	}

	return shopResponse, nil
}
