package repositories

import (
	"time"

	"github.com/rdsarjito/marketplace-backend/domain/model"
	"gorm.io/gorm"
)

type TRXRepository interface {
	Create(trx *model.TRX) error
	GetByID(id int) (*model.TRX, error)
	GetByUserID(userID int) ([]model.TRX, error)
	GetByInvoiceCode(invoiceCode string) (*model.TRX, error)
	Update(trx *model.TRX) error
	UpdatePaymentStatus(trxID int, paymentStatus string, paymentToken, paymentURL, midtransOrderID string, paymentExpiredAt *time.Time) error
	Delete(id int) error
    CreateDetail(detail *model.DetailTRX) error
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

func (r *trxRepository) GetByInvoiceCode(invoiceCode string) (*model.TRX, error) {
	var trx model.TRX
	err := r.db.Preload("User").Preload("Address").Preload("DetailTRX.Product").Preload("DetailTRX.Shop").Where("kode_invoice = ?", invoiceCode).First(&trx).Error
	if err != nil {
		return nil, err
	}
	return &trx, nil
}

func (r *trxRepository) Update(trx *model.TRX) error {
	return r.db.Save(trx).Error
}

// UpdatePaymentStatus updates only payment-related fields in transaction
func (r *trxRepository) UpdatePaymentStatus(trxID int, paymentStatus string, paymentToken, paymentURL, midtransOrderID string, paymentExpiredAt *time.Time) error {
	updates := map[string]interface{}{
		"payment_status": paymentStatus,
	}

	// Only update fields that are provided (non-empty)
	if paymentToken != "" {
		updates["payment_token"] = paymentToken
	}
	if paymentURL != "" {
		updates["payment_url"] = paymentURL
	}
	if midtransOrderID != "" {
		updates["midtrans_order_id"] = midtransOrderID
	}
	if paymentExpiredAt != nil {
		updates["payment_expired_at"] = paymentExpiredAt
	}

	return r.db.Model(&model.TRX{}).Where("id = ?", trxID).Updates(updates).Error
}

func (r *trxRepository) Delete(id int) error {
	return r.db.Delete(&model.TRX{}, id).Error
}

func (r *trxRepository) CreateDetail(detail *model.DetailTRX) error {
    return r.db.Create(detail).Error
}
