package services

import (
	"errors"
	"time"

	"github.com/rdsarjito/marketplace-backend/constants"
	"github.com/rdsarjito/marketplace-backend/domain/dto/request"
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/domain/model"
	"github.com/rdsarjito/marketplace-backend/repositories"
	"github.com/rdsarjito/marketplace-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(req *request.RegisterRequest) (*response.AuthResponse, error)
	LoginUser(req *request.LoginRequest) (*response.AuthResponse, error)
}

type authService struct {
	userRepo           repositories.UserRepository
	shopRepo           repositories.ShopRepository
	provinceCityRepo   repositories.ProvinceCityRepository
}

func NewAuthService(userRepo repositories.UserRepository, shopRepo repositories.ShopRepository, provinceCityRepo repositories.ProvinceCityRepository) AuthService {
	return &authService{
		userRepo:         userRepo,
		shopRepo:         shopRepo,
		provinceCityRepo: provinceCityRepo,
	}
}

func (s *authService) RegisterUser(req *request.RegisterRequest) (*response.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New(constants.ErrUserAlreadyExists)
	}

	existingUser, _ = s.userRepo.GetByPhone(req.NoTelp)
	if existingUser != nil {
		return nil, errors.New(constants.ErrUserAlreadyExists)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.KataSandi), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Parse tanggal lahir
	tanggalLahir, err := time.Parse("2006-01-02", req.TanggalLahir)
	if err != nil {
		return nil, errors.New("Invalid date format")
	}

	// Create user
	user := &model.User{
		Nama:         req.Nama,
		KataSandi:    string(hashedPassword),
		NoTelp:       req.NoTelp,
		TanggalLahir: tanggalLahir,
		JenisKelamin: req.JenisKelamin,
		Tentang:      req.Tentang,
		Pekerjaan:    req.Pekerjaan,
		Email:        req.Email,
		IDProvinsi:   req.IDProvinsi,
		IDKota:       req.IDKota,
		IsAdmin:      false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Create shop for user
	shop := &model.Shop{
		NamaToko: req.Nama + "'s Shop",
		URLToko:  "shop-" + req.Email,
		IDUser:   user.ID,
	}

	if err := s.shopRepo.Create(shop); err != nil {
		return nil, err
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	// Return response
	userProfile := response.UserProfile{
		ID:            user.ID,
		Nama:          user.Nama,
		NoTelp:        user.NoTelp,
		TanggalLahir:  user.TanggalLahir.Format("2006-01-02"),
		JenisKelamin:  user.JenisKelamin,
		Tentang:       user.Tentang,
		Pekerjaan:     user.Pekerjaan,
		Email:         user.Email,
		IDProvinsi:    user.IDProvinsi,
		IDKota:        user.IDKota,
		IsAdmin:       user.IsAdmin,
	}

	return &response.AuthResponse{
		Token: token,
		User:  userProfile,
	}, nil
}

func (s *authService) LoginUser(req *request.LoginRequest) (*response.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New(constants.ErrInvalidCredentials)
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(req.KataSandi)); err != nil {
		return nil, errors.New(constants.ErrInvalidCredentials)
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	// Return response
	userProfile := response.UserProfile{
		ID:            user.ID,
		Nama:          user.Nama,
		NoTelp:        user.NoTelp,
		TanggalLahir:  user.TanggalLahir.Format("2006-01-02"),
		JenisKelamin:  user.JenisKelamin,
		Tentang:       user.Tentang,
		Pekerjaan:     user.Pekerjaan,
		Email:         user.Email,
		IDProvinsi:    user.IDProvinsi,
		IDKota:        user.IDKota,
		IsAdmin:       user.IsAdmin,
	}

	return &response.AuthResponse{
		Token: token,
		User:  userProfile,
	}, nil
}
