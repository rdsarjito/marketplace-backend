package model

import (
	"time"

	"gorm.io/gorm"
)

type TRX struct {
	ID               int            `gorm:"type:int;primaryKey;autoIncrement"`
	HargaTotal       int            `gorm:"type:int;not null"`
	KodeInvoice      string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_kode_invoice"`
	MethodBayar      string         `gorm:"type:varchar(255);not null"`
	PaymentStatus    string         `gorm:"type:varchar(50);default:'pending_payment'"`
	PaymentToken     string         `gorm:"type:varchar(255);null"`
	PaymentURL       string         `gorm:"type:text;null"`
	MidtransOrderID  string         `gorm:"type:varchar(255);null;index:idx_midtrans_order_id"`
	PaymentExpiredAt *time.Time     `gorm:"type:timestamp;null"`
	PaymentVANumbers string         `gorm:"type:text;null"`
	PaymentActions   string         `gorm:"type:text;null"`
	PaymentQRString  string         `gorm:"type:text;null"`
	CreatedAt        time.Time      `gorm:"type:timestamp;not null;default:current_timestamp"`
	UpdatedAt        time.Time      `gorm:"type:timestamp"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	IDUser           int            `gorm:"type:int;not null"`
	IDAlamat         int            `gorm:"type:int;not null"`

	User      User        `gorm:"foreignKey:IDUser;references:ID"`
	Address   Address     `gorm:"foreignKey:IDAlamat;references:ID"`
	DetailTRX []DetailTRX `gorm:"foreignKey:IDTRX;references:ID"`
}

type DetailTRX struct {
	ID         int       `gorm:"type:int;primaryKey;autoIncrement"`
	IDTRX      int       `gorm:"type:int;not null"`
	IDProduk   int       `gorm:"type:int;not null"`
	IDToko     int       `gorm:"type:int;not null"`
	Kuantitas  int       `gorm:"type:int;not null"`
	HargaTotal int       `gorm:"type:int;not null"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null;default:current_timestamp"`
	UpdatedAt  time.Time `gorm:"type:timestamp"`

	TRX     TRX     `gorm:"foreignKey:IDTRX;references:ID"`
	Product Product `gorm:"foreignKey:IDProduk;references:ID"`
	Shop    Shop    `gorm:"foreignKey:IDToko;references:ID"`
}

func (TRX) TableName() string {
	return "trx"
}

func (DetailTRX) TableName() string {
	return "detail_trx"
}
