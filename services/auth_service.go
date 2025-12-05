package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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
	ForgotPassword(req *request.ForgotPasswordRequest) (*response.ForgotPasswordResponse, error)
	ResetPassword(req *request.ResetPasswordRequest) (*response.ResetPasswordResponse, error)
    LoginWithGoogle(email, name string) (*response.AuthResponse, error)
}

type authService struct {
	userRepo           repositories.UserRepository
	shopRepo           repositories.ShopRepository
	provinceCityRepo   repositories.ProvinceCityRepository
	emailService       EmailService
}

func NewAuthService(userRepo repositories.UserRepository, shopRepo repositories.ShopRepository, provinceCityRepo repositories.ProvinceCityRepository, emailService EmailService) AuthService {
	return &authService{
		userRepo:         userRepo,
		shopRepo:         shopRepo,
		provinceCityRepo: provinceCityRepo,
		emailService:     emailService,
	}
}

func (s *authService) RegisterUser(req *request.RegisterRequest) (*response.AuthResponse, error) {
	// Check if email already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New(constants.ErrEmailAlreadyExists)
	}

	// Check if phone number already exists
	existingUser, _ = s.userRepo.GetByPhone(req.NoTelp)
	if existingUser != nil {
		return nil, errors.New(constants.ErrPhoneAlreadyExists)
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

// LoginWithGoogle logs the user in using Google account email. If the user
// does not exist, it creates a minimal user profile and a shop entry.
func (s *authService) LoginWithGoogle(email, name string) (*response.AuthResponse, error) {
    if email == "" {
        return nil, errors.New("Invalid Google account: email missing")
    }

    // Try to find by email
    user, err := s.userRepo.GetByEmail(email)
    if err != nil || user == nil {
        // Create minimal user (fill required non-null fields with defaults)
        now := time.Now()
        newUser := &model.User{
            Nama:         name,
            KataSandi:    "google-oauth", // not used for Google login
            NoTelp:       "google-" + email,
            TanggalLahir: now,
            JenisKelamin: "",
            Tentang:      "",
            Pekerjaan:    "",
            Email:        email,
            IDProvinsi:   "",
            IDKota:       "",
            IsAdmin:      false,
        }
        if err := s.userRepo.Create(newUser); err != nil {
            return nil, err
        }
        // Create default shop
        shop := &model.Shop{
            NamaToko: name + "'s Shop",
            URLToko:  "shop-" + email,
            IDUser:   newUser.ID,
        }
        if err := s.shopRepo.Create(shop); err != nil {
            return nil, err
        }
        user = newUser
    }

    // Generate token
    token, err := utils.GenerateToken(user.ID, user.IsAdmin)
    if err != nil {
        return nil, err
    }

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

func (s *authService) ForgotPassword(req *request.ForgotPasswordRequest) (*response.ForgotPasswordResponse, error) {
	// Check if user exists
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		// For security, don't reveal if email exists or not
		return &response.ForgotPasswordResponse{
			Message: "Jika email terdaftar, link reset password telah dikirim ke email Anda",
		}, nil
	}

	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	// Create password reset token
	resetToken := &model.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Token expires in 24 hours
		Used:      false,
	}

	if err := s.userRepo.CreatePasswordResetToken(resetToken); err != nil {
		return nil, err
	}

	// Send email with reset link
	if err := s.emailService.SendPasswordResetEmail(user.Email, token); err != nil {
		// Log error but don't fail the request (for security)
		fmt.Printf("Failed to send email to %s: %v\n", user.Email, err)
	}

	return &response.ForgotPasswordResponse{
		Message: "Jika email terdaftar, link reset password telah dikirim ke email Anda",
	}, nil
}

func (s *authService) ResetPassword(req *request.ResetPasswordRequest) (*response.ResetPasswordResponse, error) {
	// Get and validate token
	resetToken, err := s.userRepo.GetPasswordResetToken(req.Token)
	if err != nil {
		return nil, errors.New("Token tidak valid atau sudah expired")
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		return nil, errors.New("Token sudah expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(resetToken.UserID)
	if err != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.KataSandi), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Update user password
	user.KataSandi = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Mark token as used
	if err := s.userRepo.MarkTokenAsUsed(req.Token); err != nil {
		return nil, err
	}

	return &response.ResetPasswordResponse{
		Message: "Password berhasil direset",
	}, nil
}
