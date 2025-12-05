package services

import (
	"errors"
	"time"

	"github.com/rdsarjito/marketplace-backend/constants"
	"github.com/rdsarjito/marketplace-backend/domain/dto/request"
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/domain/model"
	"github.com/rdsarjito/marketplace-backend/repositories"
)

type UserService interface {
	GetMyProfile(userID int) (*response.UserProfile, error)
	UpdateProfile(userID int, req *request.UpdateProfileRequest) (*response.UserProfile, error)
	GetMyAddress(userID int) ([]response.AddressResponse, error)
	GetDetailAddress(userID, addressID int) (*response.AddressResponse, error)
	CreateAddressUser(userID int, req *request.CreateAddressRequest) (*response.AddressResponse, error)
	UpdateAddressUser(userID, addressID int, req *request.UpdateAddressRequest) (*response.AddressResponse, error)
	DeleteAddressUser(userID, addressID int) error
}

type userService struct {
	userRepo    repositories.UserRepository
	addressRepo repositories.AddressRepository
}

func NewUserService(userRepo repositories.UserRepository, addressRepo repositories.AddressRepository) UserService {
	return &userService{
		userRepo:    userRepo,
		addressRepo: addressRepo,
	}
}

func (s *userService) GetMyProfile(userID int) (*response.UserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New(constants.ErrUserNotFound)
	}

	userProfile := &response.UserProfile{
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
		PhotoURL:      user.PhotoURL,
		IsAdmin:       user.IsAdmin,
	}

	return userProfile, nil
}

func (s *userService) UpdateProfile(userID int, req *request.UpdateProfileRequest) (*response.UserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New(constants.ErrUserNotFound)
	}

	// Update only provided fields
	if req.Nama != "" {
		user.Nama = req.Nama
	}
	if req.NoTelp != "" {
		user.NoTelp = req.NoTelp
	}
	if req.TanggalLahir != "" {
		tanggalLahir, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err != nil {
			return nil, errors.New("Invalid date format")
		}
		user.TanggalLahir = tanggalLahir
	}
	if req.JenisKelamin != "" {
		user.JenisKelamin = req.JenisKelamin
	}
	if req.Tentang != "" {
		user.Tentang = req.Tentang
	}
	if req.Pekerjaan != "" {
		user.Pekerjaan = req.Pekerjaan
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.IDProvinsi != "" {
		user.IDProvinsi = req.IDProvinsi
	}
	if req.IDKota != "" {
		user.IDKota = req.IDKota
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	userProfile := &response.UserProfile{
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
		PhotoURL:      user.PhotoURL,
		IsAdmin:       user.IsAdmin,
	}

	return userProfile, nil
}

func (s *userService) GetMyAddress(userID int) ([]response.AddressResponse, error) {
	addresses, err := s.addressRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var addressResponses []response.AddressResponse
	for _, address := range addresses {
		addressResponses = append(addressResponses, response.AddressResponse{
			ID:           address.ID,
			JudulAlamat:  address.JudulAlamat,
			NamaPenerima: address.NamaPenerima,
			NoTelp:       address.NoTelp,
			DetailAlamat: address.DetailAlamat,
			CreatedAt:    address.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    address.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return addressResponses, nil
}

func (s *userService) GetDetailAddress(userID, addressID int) (*response.AddressResponse, error) {
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return nil, errors.New(constants.ErrAddressNotFound)
	}

	if address.IDUser != userID {
		return nil, errors.New(constants.ErrForbidden)
	}

	addressResponse := &response.AddressResponse{
		ID:           address.ID,
		JudulAlamat:  address.JudulAlamat,
		NamaPenerima: address.NamaPenerima,
		NoTelp:       address.NoTelp,
		DetailAlamat: address.DetailAlamat,
		CreatedAt:    address.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    address.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return addressResponse, nil
}

func (s *userService) CreateAddressUser(userID int, req *request.CreateAddressRequest) (*response.AddressResponse, error) {
	address := &model.Address{
		JudulAlamat:  req.JudulAlamat,
		NamaPenerima: req.NamaPenerima,
		NoTelp:       req.NoTelp,
		DetailAlamat: req.DetailAlamat,
		IDUser:       userID,
	}

	if err := s.addressRepo.Create(address); err != nil {
		return nil, err
	}

	addressResponse := &response.AddressResponse{
		ID:           address.ID,
		JudulAlamat:  address.JudulAlamat,
		NamaPenerima: address.NamaPenerima,
		NoTelp:       address.NoTelp,
		DetailAlamat: address.DetailAlamat,
		CreatedAt:    address.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    address.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return addressResponse, nil
}

func (s *userService) UpdateAddressUser(userID, addressID int, req *request.UpdateAddressRequest) (*response.AddressResponse, error) {
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return nil, errors.New(constants.ErrAddressNotFound)
	}

	if address.IDUser != userID {
		return nil, errors.New(constants.ErrForbidden)
	}

	address.JudulAlamat = req.JudulAlamat
	address.NamaPenerima = req.NamaPenerima
	address.NoTelp = req.NoTelp
	address.DetailAlamat = req.DetailAlamat

	if err := s.addressRepo.Update(address); err != nil {
		return nil, err
	}

	addressResponse := &response.AddressResponse{
		ID:           address.ID,
		JudulAlamat:  address.JudulAlamat,
		NamaPenerima: address.NamaPenerima,
		NoTelp:       address.NoTelp,
		DetailAlamat: address.DetailAlamat,
		CreatedAt:    address.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    address.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return addressResponse, nil
}

func (s *userService) DeleteAddressUser(userID, addressID int) error {
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return errors.New(constants.ErrAddressNotFound)
	}

	if address.IDUser != userID {
		return errors.New(constants.ErrForbidden)
	}

	return s.addressRepo.Delete(addressID)
}
