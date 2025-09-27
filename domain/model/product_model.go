package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID             int            `gorm:"primaryKey;autoIncrement"`
	NamaProduk     string         `gorm:"type:varchar(255);not null"`
	Slug           string         `gorm:"type:varchar(255);not null"`
	HargaReseller  string         `gorm:"type:varchar(255);not null"`
	HargaKonsumen  string         `gorm:"type:varchar(255);not null"`
	Stok           int            `gorm:"type:int;not null;default:0"`
	Deskripsi      string         `gorm:"type:text;not null"`
	CreatedAt      time.Time      `gorm:"type:timestamp;not null;default:current_timestamp"`
	UpdatedAt      time.Time      `gorm:"type:timestamp"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	IDToko         int            `gorm:"type:int;not null"`
	IDCategory     int            `gorm:"type:int;not null"`

	Toko           Shop           `gorm:"foreignKey:IDToko;references:ID"`
	Category       Category       `gorm:"foreignKey:IDCategory;references:ID"`
	PhotosProduct  []PhotoProduct `gorm:"foreignKey:IDProduk;references:ID"`
}

type LogProduct struct {
	ID             int       `gorm:"type:int;primaryKey;autoIncrement"`
	IDProduk       int       `gorm:"type:int;not null"`
	NamaProduk     string    `gorm:"type:varchar(255);not null"`
	Slug           string    `gorm:"type:varchar(255);not null"`
	HargaReseller  string    `gorm:"type:varchar(255);not null"`
	HargaKonsumen  string    `gorm:"type:varchar(255);not null"`
	Stock          int       `gorm:"type:int;not null;default:0"`
	Deskripsi      string    `gorm:"type:text;not null"`
	CreatedAt      time.Time `gorm:"type:timestamp;not null;default:current_timestamp"`
	UpdatedAt      time.Time `gorm:"type:timestamp"`
	IDToko         int       `gorm:"type:int;not null"`
	IDCategory     int       `gorm:"type:int;not null"`

	Toko           Shop           `gorm:"foreignKey:IDToko;references:ID"`
	Produk         Product        `gorm:"foreignKey:IDProduk;references:ID"`
	Category       Category       `gorm:"foreignKey:IDCategory;references:ID"`
}

type PhotoProduct struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	IDProduk  int       `gorm:"type:int;not null"`
	URL       string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`

	Product    Product     `gorm:"foreignKey:IDProduk;references:ID"`
	LogProduct LogProduct  `gorm:"foreignKey:IDProduk;references:IDProduk"`
}

func (Product) TableName() string {
	return "produk"
}

func (LogProduct) TableName() string {
	return "log_produk"
}

func (PhotoProduct) TableName() string {
	return "foto_produk"
}
