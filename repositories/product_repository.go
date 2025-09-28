package repositories

import (
	"github.com/rdsarjito/marketplace-backend/domain/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(id int) (*model.Product, error)
	GetAll() ([]model.Product, error)
	GetByShopID(shopID int) ([]model.Product, error)
	GetByCategoryID(categoryID int) ([]model.Product, error)
	Update(product *model.Product) error
	Delete(id int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id int) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Toko").Preload("Category").Preload("PhotosProduct").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetAll() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Preload("Toko").Preload("Category").Preload("PhotosProduct").Find(&products).Error
	return products, err
}

func (r *productRepository) GetByShopID(shopID int) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Preload("Toko").Preload("Category").Preload("PhotosProduct").Where("id_toko = ?", shopID).Find(&products).Error
	return products, err
}

func (r *productRepository) GetByCategoryID(categoryID int) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Preload("Toko").Preload("Category").Preload("PhotosProduct").Where("id_category = ?", categoryID).Find(&products).Error
	return products, err
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id int) error {
	return r.db.Delete(&model.Product{}, id).Error
}
