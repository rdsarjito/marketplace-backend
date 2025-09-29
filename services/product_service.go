package services

import (
	"errors"

	"github.com/rdsarjito/marketplace-backend/constants"
	"github.com/rdsarjito/marketplace-backend/domain/dto/request"
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/domain/model"
	"github.com/rdsarjito/marketplace-backend/repositories"
	"github.com/gosimple/slug"
)

type ProductService interface {
	GetListProduct() ([]response.ProductResponse, error)
	GetDetailProduct(id int) (*response.ProductResponse, error)
	CreateProduct(userID int, req *request.CreateProductRequest) (*response.ProductResponse, error)
	UpdateProduct(userID, id int, req *request.UpdateProductRequest) (*response.ProductResponse, error)
	DeleteProduct(userID, id int) error
    AddProductPhoto(userID, productID int, url string) (*response.ProductResponse, error)
}

type productService struct {
	productRepo  repositories.ProductRepository
	shopRepo     repositories.ShopRepository
	categoryRepo repositories.CategoryRepository
}

func NewProductService(productRepo repositories.ProductRepository, shopRepo repositories.ShopRepository, categoryRepo repositories.CategoryRepository) ProductService {
	return &productService{
		productRepo:  productRepo,
		shopRepo:     shopRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) GetListProduct() ([]response.ProductResponse, error) {
	products, err := s.productRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var productResponses []response.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, s.mapProductToResponse(product))
	}

	return productResponses, nil
}

func (s *productService) GetDetailProduct(id int) (*response.ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(constants.ErrProductNotFound)
	}

	productResponse := s.mapProductToResponse(*product)
	return &productResponse, nil
}

func (s *productService) CreateProduct(userID int, req *request.CreateProductRequest) (*response.ProductResponse, error) {
	// Check if shop belongs to user
	shop, err := s.shopRepo.GetByID(req.IDToko)
	if err != nil {
		return nil, errors.New(constants.ErrShopNotFound)
	}

	if shop.IDUser != userID {
		return nil, errors.New(constants.ErrForbidden)
	}

	// Check if category exists
	_, err = s.categoryRepo.GetByID(req.IDCategory)
	if err != nil {
		return nil, errors.New(constants.ErrCategoryNotFound)
	}

	// Generate slug
	productSlug := slug.Make(req.NamaProduk)

	product := &model.Product{
		NamaProduk:    req.NamaProduk,
		Slug:          productSlug,
		HargaReseller: req.HargaReseller,
		HargaKonsumen: req.HargaKonsumen,
		Stok:          req.Stok,
		Deskripsi:     req.Deskripsi,
		IDToko:        req.IDToko,
		IDCategory:    req.IDCategory,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	// Get created product with relations
	createdProduct, err := s.productRepo.GetByID(product.ID)
	if err != nil {
		return nil, err
	}

	productResponse := s.mapProductToResponse(*createdProduct)
	return &productResponse, nil
}

func (s *productService) UpdateProduct(userID, id int, req *request.UpdateProductRequest) (*response.ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(constants.ErrProductNotFound)
	}

	// Check if shop belongs to user
	shop, err := s.shopRepo.GetByID(product.IDToko)
	if err != nil {
		return nil, errors.New(constants.ErrShopNotFound)
	}

	if shop.IDUser != userID {
		return nil, errors.New(constants.ErrForbidden)
	}

	// Check if category exists
	_, err = s.categoryRepo.GetByID(req.IDCategory)
	if err != nil {
		return nil, errors.New(constants.ErrCategoryNotFound)
	}

	// Generate new slug if product name changed
	if product.NamaProduk != req.NamaProduk {
		product.Slug = slug.Make(req.NamaProduk)
	}

	product.NamaProduk = req.NamaProduk
	product.HargaReseller = req.HargaReseller
	product.HargaKonsumen = req.HargaKonsumen
	product.Stok = req.Stok
	product.Deskripsi = req.Deskripsi
	product.IDCategory = req.IDCategory

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	// Get updated product with relations
	updatedProduct, err := s.productRepo.GetByID(product.ID)
	if err != nil {
		return nil, err
	}

	productResponse := s.mapProductToResponse(*updatedProduct)
	return &productResponse, nil
}

func (s *productService) DeleteProduct(userID, id int) error {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return errors.New(constants.ErrProductNotFound)
	}

	// Check if shop belongs to user
	shop, err := s.shopRepo.GetByID(product.IDToko)
	if err != nil {
		return errors.New(constants.ErrShopNotFound)
	}

	if shop.IDUser != userID {
		return errors.New(constants.ErrForbidden)
	}

	return s.productRepo.Delete(id)
}

func (s *productService) AddProductPhoto(userID, productID int, url string) (*response.ProductResponse, error) {
    product, err := s.productRepo.GetByID(productID)
    if err != nil {
        return nil, errors.New(constants.ErrProductNotFound)
    }
    shop, err := s.shopRepo.GetByID(product.IDToko)
    if err != nil {
        return nil, errors.New(constants.ErrShopNotFound)
    }
    if shop.IDUser != userID {
        return nil, errors.New(constants.ErrForbidden)
    }
    photo := &model.PhotoProduct{ IDProduk: productID, URL: url }
    if err := s.productRepo.AddPhoto(photo); err != nil {
        return nil, err
    }
    updated, err := s.productRepo.GetByID(productID)
    if err != nil {
        return nil, err
    }
    resp := s.mapProductToResponse(*updated)
    return &resp, nil
}

func (s *productService) mapProductToResponse(product model.Product) response.ProductResponse {
	// Map shop
	shopResponse := response.ShopResponse{
		ID:        product.Toko.ID,
		NamaToko:  product.Toko.NamaToko,
		URLToko:   product.Toko.URLToko,
		CreatedAt: product.Toko.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.Toko.UpdatedAt.Format("2006-01-02 15:04:05"),
		IDUser:    product.Toko.IDUser,
	}

	// Map category
	categoryResponse := response.CategoryResponse{
		ID:        product.Category.ID,
		Nama:      product.Category.Nama,
		CreatedAt: product.Category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.Category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Map photos
	var photoResponses []response.PhotoProductResponse
	for _, photo := range product.PhotosProduct {
		photoResponses = append(photoResponses, response.PhotoProductResponse{
			ID:        photo.ID,
			IDProduk:  photo.IDProduk,
			URL:       photo.URL,
			CreatedAt: photo.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: photo.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response.ProductResponse{
		ID:             product.ID,
		NamaProduk:     product.NamaProduk,
		Slug:           product.Slug,
		HargaReseller:  product.HargaReseller,
		HargaKonsumen:  product.HargaKonsumen,
		Stok:           product.Stok,
		Deskripsi:      product.Deskripsi,
		CreatedAt:      product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      product.UpdatedAt.Format("2006-01-02 15:04:05"),
		IDToko:         product.IDToko,
		IDCategory:     product.IDCategory,
		Toko:           shopResponse,
		Category:       categoryResponse,
		PhotosProduct:  photoResponses,
	}
}
