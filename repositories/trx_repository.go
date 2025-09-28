package repositories

import (
	"github.com/rdsarjito/marketplace-backend/domain/model"
	"gorm.io/gorm"
)

type TRXRepository interface {
	Create(trx *model.TRX) error
	GetByID(id int) (*model.TRX, error)
	GetByUserID(userID int) ([]model.TRX, error)
	Update(trx *model.TRX) error
	Delete(id int) error
}

type trxRepository struct {
	db *gorm.DB
}

func NewTRXRepository(db *gorm.DB) TRXRepository {
	return &trxRepository{db: db}
}

func (r *trxRepository) Create(trx *model.TRX) error {
	return r.db.Create(trx).Error
}

func (r *trxRepository) GetByID(id int) (*model.TRX, error) {
	var trx model.TRX
	err := r.db.Preload("User").Preload("Address").Preload("DetailTRX.Product").Preload("DetailTRX.Shop").First(&trx, id).Error
	if err != nil {
		return nil, err
	}
	return &trx, nil
}

func (r *trxRepository) GetByUserID(userID int) ([]model.TRX, error) {
	var trxs []model.TRX
	err := r.db.Preload("User").Preload("Address").Preload("DetailTRX.Product").Preload("DetailTRX.Shop").Where("id_user = ?", userID).Find(&trxs).Error
	return trxs, err
}

func (r *trxRepository) Update(trx *model.TRX) error {
	return r.db.Save(trx).Error
}

func (r *trxRepository) Delete(id int) error {
	return r.db.Delete(&model.TRX{}, id).Error
}
