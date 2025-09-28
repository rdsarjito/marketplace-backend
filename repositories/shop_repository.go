package repositories

import (
	"github.com/rdsarjito/marketplace-backend/domain/model"
	"gorm.io/gorm"
)

type ShopRepository interface {
	Create(shop *model.Shop) error
	GetByID(id int) (*model.Shop, error)
	GetByUserID(userID int) (*model.Shop, error)
	GetAll() ([]model.Shop, error)
	Update(shop *model.Shop) error
	Delete(id int) error
}

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) ShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) Create(shop *model.Shop) error {
	return r.db.Create(shop).Error
}

func (r *shopRepository) GetByID(id int) (*model.Shop, error) {
	var shop model.Shop
	err := r.db.Preload("User").First(&shop, id).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *shopRepository) GetByUserID(userID int) (*model.Shop, error) {
	var shop model.Shop
	err := r.db.Preload("User").Where("id_user = ?", userID).First(&shop).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *shopRepository) GetAll() ([]model.Shop, error) {
	var shops []model.Shop
	err := r.db.Preload("User").Find(&shops).Error
	return shops, err
}

func (r *shopRepository) Update(shop *model.Shop) error {
	return r.db.Save(shop).Error
}

func (r *shopRepository) Delete(id int) error {
	return r.db.Delete(&model.Shop{}, id).Error
}
