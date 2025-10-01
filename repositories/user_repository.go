package repositories

import (
	"github.com/rdsarjito/marketplace-backend/domain/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByPhone(phone string) (*model.User, error)
	Update(user *model.User) error
	Delete(id int) error
	CreatePasswordResetToken(token *model.PasswordResetToken) error
	GetPasswordResetToken(token string) (*model.PasswordResetToken, error)
	MarkTokenAsUsed(token string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id int) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("notelp = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id int) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) CreatePasswordResetToken(token *model.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *userRepository) GetPasswordResetToken(token string) (*model.PasswordResetToken, error) {
	var resetToken model.PasswordResetToken
	err := r.db.Where("token = ? AND used = ?", token, false).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

func (r *userRepository) MarkTokenAsUsed(token string) error {
	return r.db.Model(&model.PasswordResetToken{}).Where("token = ?", token).Update("used", true).Error
}
