package model

import "time"

type User struct {
	ID            int       `gorm:"type:int;primaryKey;autoIncrement"`
	Nama          string    `gorm:"type:varchar(255);not null"`
	KataSandi     string    `gorm:"type:varchar(255);not null"`
	NoTelp        string    `gorm:"column:notelp;type:varchar(255);not null;uniqueIndex:idx_notelp"`
	TanggalLahir  time.Time `gorm:"type:date;not null"`
	JenisKelamin  string    `gorm:"type:varchar(255);not null"`
	Tentang       string    `gorm:"type:text;not null"`
	Pekerjaan     string    `gorm:"type:varchar(255);not null"`
	Email         string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_email"`
	IDProvinsi    string    `gorm:"type:varchar(255);not null"`
	IDKota        string    `gorm:"type:varchar(255);not null"`
	IsAdmin       bool      `gorm:"column:isAdmin;default:false"`
	CreatedAt     time.Time `gorm:"type:timestamp"`
	UpdatedAt     time.Time `gorm:"type:timestamp"`
}

func (User) TableName() string {
	return "user"
}
