package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          int            `gorm:"type:int;primaryKey;autoIncrement"`
	Nama        string         `gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time      `gorm:"type:timestamp;not null;default:current_timestamp"`
	UpdatedAt   time.Time      `gorm:"type:timestamp"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Products    []Product      `gorm:"foreignKey:IDCategory;references:ID"`
}

func (Category) TableName() string {
	return "kategori"
}
