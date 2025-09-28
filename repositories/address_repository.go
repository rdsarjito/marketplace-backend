package repositories

import (
	"github.com/rdsarjito/marketplace-backend/domain/model"
	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(address *model.Address) error
	GetByID(id int) (*model.Address, error)
	GetByUserID(userID int) ([]model.Address, error)
	Update(address *model.Address) error
	Delete(id int) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(address *model.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepository) GetByID(id int) (*model.Address, error) {
	var address model.Address
	err := r.db.First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) GetByUserID(userID int) ([]model.Address, error) {
	var addresses []model.Address
	err := r.db.Where("id_user = ?", userID).Find(&addresses).Error
	return addresses, err
}

func (r *addressRepository) Update(address *model.Address) error {
	return r.db.Save(address).Error
}

func (r *addressRepository) Delete(id int) error {
	return r.db.Delete(&model.Address{}, id).Error
}
