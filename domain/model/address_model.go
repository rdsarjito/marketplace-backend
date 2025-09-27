package model

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID           int            `gorm:"type:int;primaryKey;autoIncrement"`
	JudulAlamat  string         `gorm:"type:varchar(255);not null"`
	NamaPenerima string         `gorm:"type:varchar(255);not null"`
	NoTelp       string         `gorm:"type:varchar(255);not null"`
	DetailAlamat string         `gorm:"type:text;not null"`
	CreatedAt    time.Time      `gorm:"type:timestamp;not null;default:current_timestamp"`
	UpdatedAt    time.Time      `gorm:"type:timestamp"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	IDUser       int            `gorm:"type:int;not null"`

	User User `gorm:"foreignKey:IDUser;references:ID"`
}

func (Address) TableName() string {
	return "alamat"
}
