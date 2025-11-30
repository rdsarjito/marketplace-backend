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
	trxRepo         repositories.TRXRepository
	productRepo     repositories.ProductRepository
	addressRepo     repositories.AddressRepository
	shopRepo        repositories.ShopRepository
	categoryRepo    repositories.CategoryRepository
	userRepo        repositories.UserRepository
	midtransService MidtransService
}

func NewTRXService(trxRepo repositories.TRXRepository, productRepo repositories.ProductRepository, addressRepo repositories.AddressRepository, shopRepo repositories.ShopRepository, categoryRepo repositories.CategoryRepository, userRepo repositories.UserRepository, midtransService MidtransService) TRXService {
	return &trxService{
		trxRepo:         trxRepo,
		productRepo:     productRepo,
		addressRepo:     addressRepo,
		shopRepo:        shopRepo,
		categoryRepo:    categoryRepo,
		userRepo:        userRepo,
		midtransService: midtransService,
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
		_, err = s.shopRepo.GetByID(detail.IDToko)
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

	// Get user data for customer details
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("User not found")
	}

	// Create transaction
	trx := &model.TRX{
		HargaTotal:    req.HargaTotal,
		KodeInvoice:   kodeInvoice,
		MethodBayar:   req.MethodBayar,
		PaymentStatus: "pending_payment",
		IDUser:        userID,
		IDAlamat:      req.IDAlamat,
	}

	if err := s.trxRepo.Create(trx); err != nil {
		return nil, err
	}

	// Create detail transactions and update stock
	var itemDetails []map[string]interface{}
	for _, detailReq := range req.DetailTRX {
		// Update product stock
		product, _ := s.productRepo.GetByID(detailReq.IDProduk)
		product.Stok -= detailReq.Kuantitas
		s.productRepo.Update(product)

		// Create detail record
		detail := &model.DetailTRX{
			IDTRX:      trx.ID,
			IDProduk:   detailReq.IDProduk,
			IDToko:     detailReq.IDToko,
			Kuantitas:  detailReq.Kuantitas,
			HargaTotal: detailReq.HargaTotal,
		}
		_ = s.trxRepo.CreateDetail(detail)

		// Build item details for Midtrans
		itemDetails = append(itemDetails, map[string]interface{}{
			"id":       fmt.Sprintf("product-%d", product.ID),
			"price":    product.HargaKonsumen,
			"quantity": detailReq.Kuantitas,
			"name":     product.NamaProduk,
		})
	}

	// If payment method is not COD, create payment via Midtrans
	if req.MethodBayar != "COD" {
		// Map payment method to Midtrans payment type
		paymentType := s.mapPaymentMethodToMidtransType(req.MethodBayar)

		// Build customer details
		customerDetails := map[string]interface{}{
			"first_name": user.Nama,
			"email":      user.Email,
			"phone":      user.NoTelp,
		}

		// Build billing address from user address
		billingAddress := map[string]interface{}{
			"first_name": address.NamaPenerima,
			"phone":      address.NoTelp,
			"address":    address.DetailAlamat,
		}
		customerDetails["billing_address"] = billingAddress

		// Create payment request
		midtransReq := &CreatePaymentRequest{
			OrderID:         kodeInvoice,
			GrossAmount:     req.HargaTotal,
			PaymentType:     paymentType,
			CustomerDetails: customerDetails,
			ItemDetails:     itemDetails,
			CustomExpiry: &CustomExpiry{
				ExpiryDuration: 24 * 60, // 24 hours in minutes
				Unit:           "minute",
			},
		}

		// Call Midtrans service
		paymentResp, err := s.midtransService.CreatePayment(midtransReq)
		if err != nil {
			// If payment creation fails, still return transaction but with error status
			trx.PaymentStatus = "failed"
			s.trxRepo.Update(trx)
			return nil, fmt.Errorf("failed to create payment: %w", err)
		}

		// Parse expiry time (Midtrans returns in various formats)
		var expiryTime *time.Time
		if paymentResp.ExpiryTime != "" {
			// Try multiple date formats
			formats := []string{
				"2006-01-02 15:04:05",
				time.RFC3339,
				"2006-01-02T15:04:05Z07:00",
				"2006-01-02T15:04:05",
			}
			for _, format := range formats {
				if parsed, err := time.Parse(format, paymentResp.ExpiryTime); err == nil {
					expiryTime = &parsed
					break
				}
			}
		}

		// Update transaction with payment info
		trx.PaymentToken = paymentResp.Token
		trx.PaymentURL = paymentResp.RedirectURL
		trx.MidtransOrderID = paymentResp.OrderID
		trx.PaymentExpiredAt = expiryTime
		trx.PaymentStatus = "pending_payment"

		if err := s.trxRepo.Update(trx); err != nil {
			return nil, fmt.Errorf("failed to update transaction with payment info: %w", err)
		}
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

// mapPaymentMethodToMidtransType maps our payment method to Midtrans payment type
func (s *trxService) mapPaymentMethodToMidtransType(methodBayar string) string {
	switch methodBayar {
	case "virtual_account", "va":
		return "bank_transfer" // Will be configured as VA in the service
	case "e_wallet", "ewallet", "gopay", "ovo", "dana", "linkaja":
		return "e_wallet"
	case "bank_transfer", "bank_transfer_bca", "bank_transfer_bni", "bank_transfer_mandiri":
		return "bank_transfer"
	case "credit_card", "cc":
		return "credit_card"
	default:
		return "bank_transfer" // Default fallback
	}
}

func (s *trxService) mapTRXToResponse(trx model.TRX) response.TRXResponse {
	// Map user
	userResponse := response.UserProfile{
		ID:           trx.User.ID,
		Nama:         trx.User.Nama,
		NoTelp:       trx.User.NoTelp,
		TanggalLahir: trx.User.TanggalLahir.Format("2006-01-02"),
		JenisKelamin: trx.User.JenisKelamin,
		Tentang:      trx.User.Tentang,
		Pekerjaan:    trx.User.Pekerjaan,
		Email:        trx.User.Email,
		IDProvinsi:   trx.User.IDProvinsi,
		IDKota:       trx.User.IDKota,
		IsAdmin:      trx.User.IsAdmin,
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
			ID:            detail.Product.ID,
			NamaProduk:    detail.Product.NamaProduk,
			Slug:          detail.Product.Slug,
			HargaReseller: detail.Product.HargaReseller,
			HargaKonsumen: detail.Product.HargaKonsumen,
			Stok:          detail.Product.Stok,
			Deskripsi:     detail.Product.Deskripsi,
			CreatedAt:     detail.Product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     detail.Product.UpdatedAt.Format("2006-01-02 15:04:05"),
			IDToko:        detail.Product.IDToko,
			IDCategory:    detail.Product.IDCategory,
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
