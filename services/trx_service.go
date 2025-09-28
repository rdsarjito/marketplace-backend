package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/rdsarjito/marketplace-backend/constants"
	"github.com/rdsarjito/marketplace-backend/domain/dto/request"
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/domain/model"
	"github.com/rdsarjito/marketplace-backend/repositories"
)

type TRXService interface {
	GetListTRX(userID int) ([]response.TRXResponse, error)
	GetDetailTRX(userID, trxID int) (*response.TRXResponse, error)
	CreateTRX(userID int, req *request.CreateTRXRequest) (*response.TRXResponse, error)
}

type trxService struct {
	trxRepo       repositories.TRXRepository
	productRepo   repositories.ProductRepository
	addressRepo   repositories.AddressRepository
	shopRepo      repositories.ShopRepository
	categoryRepo  repositories.CategoryRepository
}

func NewTRXService(trxRepo repositories.TRXRepository, productRepo repositories.ProductRepository, addressRepo repositories.AddressRepository, shopRepo repositories.ShopRepository, categoryRepo repositories.CategoryRepository) TRXService {
	return &trxService{
		trxRepo:      trxRepo,
		productRepo:  productRepo,
		addressRepo:  addressRepo,
		shopRepo:     shopRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *trxService) GetListTRX(userID int) ([]response.TRXResponse, error) {
	trxs, err := s.trxRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var trxResponses []response.TRXResponse
	for _, trx := range trxs {
		trxResponses = append(trxResponses, s.mapTRXToResponse(trx))
	}

	return trxResponses, nil
}

func (s *trxService) GetDetailTRX(userID, trxID int) (*response.TRXResponse, error) {
	trx, err := s.trxRepo.GetByID(trxID)
	if err != nil {
		return nil, errors.New("Transaction not found")
	}

	if trx.IDUser != userID {
		return nil, errors.New(constants.ErrForbidden)
	}

	trxResponse := s.mapTRXToResponse(*trx)
	return &trxResponse, nil
}

func (s *trxService) CreateTRX(userID int, req *request.CreateTRXRequest) (*response.TRXResponse, error) {
	// Validate address belongs to user
	address, err := s.addressRepo.GetByID(req.IDAlamat)
	if err != nil {
		return nil, errors.New(constants.ErrAddressNotFound)
	}

	if address.IDUser != userID {
		return nil, errors.New(constants.ErrForbidden)
	}

	// Validate products and calculate total
	totalHarga := 0
	for _, detail := range req.DetailTRX {
		product, err := s.productRepo.GetByID(detail.IDProduk)
		if err != nil {
			return nil, errors.New(constants.ErrProductNotFound)
		}

		// Check stock
		if product.Stok < detail.Kuantitas {
			return nil, errors.New(constants.ErrInsufficientStock)
		}

		// Validate shop
		shop, err := s.shopRepo.GetByID(detail.IDToko)
		if err != nil {
			return nil, errors.New(constants.ErrShopNotFound)
		}

		// Calculate price (using harga konsumen)
		hargaKonsumen, err := strconv.Atoi(product.HargaKonsumen)
		if err != nil {
			return nil, errors.New("Invalid price format")
		}

		detailHarga := hargaKonsumen * detail.Kuantitas
		if detailHarga != detail.HargaTotal {
			return nil, errors.New("Price calculation mismatch")
		}

		totalHarga += detailHarga
	}

	// Validate total price
	if totalHarga != req.HargaTotal {
		return nil, errors.New("Total price mismatch")
	}

	// Generate invoice code
	kodeInvoice := s.generateInvoiceCode()

	// Create transaction
	trx := &model.TRX{
		HargaTotal:  req.HargaTotal,
		KodeInvoice: kodeInvoice,
		MethodBayar: req.MethodBayar,
		IDUser:      userID,
		IDAlamat:    req.IDAlamat,
	}

	if err := s.trxRepo.Create(trx); err != nil {
		return nil, err
	}

	// Create detail transactions and update stock
	for _, detailReq := range req.DetailTRX {
		detail := &model.DetailTRX{
			IDTRX:      trx.ID,
			IDProduk:   detailReq.IDProduk,
			IDToko:     detailReq.IDToko,
			Kuantitas:  detailReq.Kuantitas,
			HargaTotal: detailReq.HargaTotal,
		}

		// Update product stock
		product, _ := s.productRepo.GetByID(detailReq.IDProduk)
		product.Stok -= detailReq.Kuantitas
		s.productRepo.Update(product)

		// Note: In a real application, you would also create the detail record
		// For now, we'll skip the detail creation as we don't have a detail repository
	}

	// Get created transaction with relations
	createdTRX, err := s.trxRepo.GetByID(trx.ID)
	if err != nil {
		return nil, err
	}

	trxResponse := s.mapTRXToResponse(*createdTRX)
	return &trxResponse, nil
}

func (s *trxService) generateInvoiceCode() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("INV-%d", timestamp)
}

func (s *trxService) mapTRXToResponse(trx model.TRX) response.TRXResponse {
	// Map user
	userResponse := response.UserProfile{
		ID:            trx.User.ID,
		Nama:          trx.User.Nama,
		NoTelp:        trx.User.NoTelp,
		TanggalLahir:  trx.User.TanggalLahir.Format("2006-01-02"),
		JenisKelamin:  trx.User.JenisKelamin,
		Tentang:       trx.User.Tentang,
		Pekerjaan:     trx.User.Pekerjaan,
		Email:         trx.User.Email,
		IDProvinsi:    trx.User.IDProvinsi,
		IDKota:        trx.User.IDKota,
		IsAdmin:       trx.User.IsAdmin,
	}

	// Map address
	addressResponse := response.AddressResponse{
		ID:           trx.Address.ID,
		JudulAlamat:  trx.Address.JudulAlamat,
		NamaPenerima: trx.Address.NamaPenerima,
		NoTelp:       trx.Address.NoTelp,
		DetailAlamat: trx.Address.DetailAlamat,
		CreatedAt:    trx.Address.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    trx.Address.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Map detail transactions
	var detailResponses []response.DetailTRXResponse
	for _, detail := range trx.DetailTRX {
		// Map product
		productResponse := response.ProductResponse{
			ID:             detail.Product.ID,
			NamaProduk:     detail.Product.NamaProduk,
			Slug:           detail.Product.Slug,
			HargaReseller:  detail.Product.HargaReseller,
			HargaKonsumen:  detail.Product.HargaKonsumen,
			Stok:           detail.Product.Stok,
			Deskripsi:      detail.Product.Deskripsi,
			CreatedAt:      detail.Product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      detail.Product.UpdatedAt.Format("2006-01-02 15:04:05"),
			IDToko:         detail.Product.IDToko,
			IDCategory:     detail.Product.IDCategory,
		}

		// Map shop
		shopResponse := response.ShopResponse{
			ID:        detail.Shop.ID,
			NamaToko:  detail.Shop.NamaToko,
			URLToko:   detail.Shop.URLToko,
			CreatedAt: detail.Shop.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: detail.Shop.UpdatedAt.Format("2006-01-02 15:04:05"),
			IDUser:    detail.Shop.IDUser,
		}

		detailResponses = append(detailResponses, response.DetailTRXResponse{
			ID:         detail.ID,
			IDTRX:      detail.IDTRX,
			IDProduk:   detail.IDProduk,
			IDToko:     detail.IDToko,
			Kuantitas:  detail.Kuantitas,
			HargaTotal: detail.HargaTotal,
			CreatedAt:  detail.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  detail.UpdatedAt.Format("2006-01-02 15:04:05"),
			Product:    productResponse,
			Shop:       shopResponse,
		})
	}

	return response.TRXResponse{
		ID:          trx.ID,
		HargaTotal:  trx.HargaTotal,
		KodeInvoice: trx.KodeInvoice,
		MethodBayar: trx.MethodBayar,
		CreatedAt:   trx.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   trx.UpdatedAt.Format("2006-01-02 15:04:05"),
		IDUser:      trx.IDUser,
		IDAlamat:    trx.IDAlamat,
		User:        userResponse,
		Address:     addressResponse,
		DetailTRX:   detailResponses,
	}
}
