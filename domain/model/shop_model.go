package model

import (
	"time"

	"gorm.io/gorm"
)

type Shop struct {
	ID          int            `gorm:"type:int;primaryKey;autoIncrement"`
	NamaToko    string         `gorm:"type:varchar(255);not null"`
	URLToko     string         `gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time      `gorm:"type:timestamp;not null;default:current_timestamp"`
	UpdatedAt   time.Time      `gorm:"type:timestamp"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	IDUser      int            `gorm:"type:int;not null"`

	User        User           `gorm:"foreignKey:IDUser;references:ID"`
	Products    []Product      `gorm:"foreignKey:IDToko;references:ID"`
}

func (Shop) TableName() string {
	return "toko"
}
