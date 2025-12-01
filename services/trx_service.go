package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	HandlePaymentWebhook(notification map[string]interface{}) error
	CheckPaymentStatus(userID, trxID int) (*response.TRXResponse, error)
}

type trxService struct {
	trxRepo         repositories.TRXRepository
	productRepo     repositories.ProductRepository
	addressRepo     repositories.AddressRepository
	shopRepo        repositories.ShopRepository
	categoryRepo    repositories.CategoryRepository
	userRepo        repositories.UserRepository
	midtransService MidtransService
	emailService    EmailService
	frontendURL     string // Frontend URL for payment redirect
}

func NewTRXService(trxRepo repositories.TRXRepository, productRepo repositories.ProductRepository, addressRepo repositories.AddressRepository, shopRepo repositories.ShopRepository, categoryRepo repositories.CategoryRepository, userRepo repositories.UserRepository, midtransService MidtransService, emailService EmailService, frontendURL string) TRXService {
	return &trxService{
		trxRepo:         trxRepo,
		productRepo:     productRepo,
		addressRepo:     addressRepo,
		shopRepo:        shopRepo,
		categoryRepo:    categoryRepo,
		userRepo:        userRepo,
		midtransService: midtransService,
		emailService:    emailService,
		frontendURL:     frontendURL,
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
	s.attachVANumbersIfNeeded(trx, &trxResponse)
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

	var lastPaymentResp *CreatePaymentResponse
	var lastPaymentVANumbers []response.PaymentVANumber

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

		// Build finish URL for redirect after payment
		finishURL := fmt.Sprintf("%s/payment/%d", s.frontendURL, trx.ID)

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
			FinishURL: finishURL,
		}

		// Call Midtrans service
		paymentResp, err := s.midtransService.CreatePayment(midtransReq)
		if err != nil {
			// If payment creation fails, still return transaction but with error status
			// Use UpdatePaymentStatus to avoid updating created_at
			s.trxRepo.UpdatePaymentStatus(trx.ID, "failed", "", "", "", nil, "")
			// Return more detailed error message
			return nil, fmt.Errorf("failed to create payment with Midtrans: %w. Please check your Midtrans Server Key configuration", err)
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
		// For bank_transfer, RedirectURL might be empty, but we can use actions or va_numbers
		// For e_wallet (gopay, ovo, etc), RedirectURL might be empty, but we can use actions
		paymentURL := paymentResp.RedirectURL
		if paymentURL == "" && len(paymentResp.Actions) > 0 {
			// Try to get URL from actions
			// Priority: deeplink-redirect > generate-qr-code-v2 > generate-qr-code > others
			for _, action := range paymentResp.Actions {
				name, _ := action["name"].(string)
				url, ok := action["url"].(string)
				if ok && url != "" {
					// Prefer deeplink-redirect for e-wallet, or any URL for bank transfer
					if name == "deeplink-redirect" {
						paymentURL = url
						break
					} else if name == "generate-qr-code-v2" && paymentURL == "" {
						paymentURL = url
					} else if name == "generate-qr-code" && paymentURL == "" {
						paymentURL = url
					} else if paymentURL == "" {
						paymentURL = url
					}
				}
			}
		}
		// If still empty and we have va_numbers, we'll show VA info in payment status page
		// PaymentURL can be empty for bank_transfer - frontend will handle displaying VA numbers

		vaNumbers := mapVANumbersFromMidtrans(paymentResp.VaNumbers)
		vaNumbersJSON := serializeVANumbersToJSON(vaNumbers)

		// Use UpdatePaymentStatus to only update payment fields (avoid updating created_at)
		if err := s.trxRepo.UpdatePaymentStatus(
			trx.ID,
			"pending_payment",
			paymentResp.Token,
			paymentURL,
			paymentResp.OrderID,
			expiryTime,
			vaNumbersJSON,
		); err != nil {
			return nil, fmt.Errorf("failed to update transaction with payment info: %w", err)
		}

		lastPaymentResp = paymentResp
		lastPaymentVANumbers = vaNumbers
	}

	// Get created transaction with relations
	createdTRX, err := s.trxRepo.GetByID(trx.ID)
	if err != nil {
		return nil, err
	}

	trxResponse := s.mapTRXToResponse(*createdTRX)
	if lastPaymentResp != nil {
		vaNumbers := lastPaymentVANumbers
		if len(vaNumbers) == 0 {
			vaNumbers = mapVANumbersFromMidtrans(lastPaymentResp.VaNumbers)
		}
		if len(vaNumbers) > 0 {
			trxResponse.PaymentVANumbers = vaNumbers
		} else {
			s.attachVANumbersIfNeeded(createdTRX, &trxResponse)
		}
	} else {
		s.attachVANumbersIfNeeded(createdTRX, &trxResponse)
	}
	return &trxResponse, nil
}

func (s *trxService) generateInvoiceCode() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("INV-%d", timestamp)
}

// mapPaymentMethodToMidtransType maps our payment method to Midtrans payment type
// For e_wallet, returns the specific wallet name (gopay, ovo, dana, linkaja)
// For other methods, returns the payment type
func (s *trxService) mapPaymentMethodToMidtransType(methodBayar string) string {
	switch methodBayar {
	case "virtual_account", "va":
		return "virtual_account" // Will be configured as VA in the service
	case "gopay":
		return "gopay"
	case "ovo":
		return "ovo"
	case "dana":
		return "dana"
	case "linkaja":
		return "linkaja"
	case "e_wallet", "ewallet":
		return "gopay" // Default to gopay for generic e_wallet
	case "bank_transfer", "bank_transfer_bca", "bank_transfer_bni", "bank_transfer_mandiri":
		return "bank_transfer"
	case "credit_card", "cc":
		return "credit_card"
	default:
		return "bank_transfer" // Default fallback
	}
}

// HandlePaymentWebhook handles webhook notification from Midtrans
func (s *trxService) HandlePaymentWebhook(notification map[string]interface{}) error {
	// Get order_id from notification
	orderID, ok := notification["order_id"].(string)
	if !ok {
		return fmt.Errorf("invalid notification: missing order_id")
	}

	// Verify payment status from Midtrans
	paymentStatus, err := s.midtransService.VerifyPayment(orderID)
	if err != nil {
		return fmt.Errorf("failed to verify payment: %w", err)
	}

	// Find transaction by invoice code (order_id)
	trx, err := s.trxRepo.GetByInvoiceCode(orderID)
	if err != nil {
		return fmt.Errorf("transaction not found: %w", err)
	}

	// Store old status for email notification
	oldStatus := trx.PaymentStatus

	// Map Midtrans transaction status to our payment status
	paymentStatusStr := s.mapMidtransStatusToPaymentStatus(paymentStatus.TransactionStatus)
	vaNumbers := mapVANumbersFromMidtrans(paymentStatus.VaNumbers)
	vaNumbersJSON := serializeVANumbersToJSON(vaNumbers)

	// Update transaction payment status (use UpdatePaymentStatus to avoid updating created_at)
	if err := s.trxRepo.UpdatePaymentStatus(trx.ID, paymentStatusStr, "", "", "", nil, vaNumbersJSON); err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// Send email notification if status changed to paid or expired
	if oldStatus != paymentStatusStr {
		// Get user email from transaction
		user, err := s.userRepo.GetByID(trx.IDUser)
		if err == nil {
			// Send email notification asynchronously (don't block webhook)
			go func() {
				if paymentStatusStr == "paid" {
					_ = s.emailService.SendPaymentSuccessEmail(user.Email, trx.KodeInvoice, trx.HargaTotal)
				} else if paymentStatusStr == "expired" {
					_ = s.emailService.SendPaymentExpiredEmail(user.Email, trx.KodeInvoice, trx.HargaTotal)
				}
			}()
		}
	}

	return nil
}

// mapMidtransStatusToPaymentStatus maps Midtrans transaction status to our payment status
func (s *trxService) mapMidtransStatusToPaymentStatus(midtransStatus string) string {
	switch midtransStatus {
	case "settlement":
		return "paid"
	case "pending":
		return "pending_payment"
	case "expire":
		return "expired"
	case "cancel", "deny":
		return "cancelled"
	case "failure":
		return "failed"
	default:
		return "pending_payment"
	}
}

// CheckPaymentStatus manually checks payment status for a transaction
func (s *trxService) CheckPaymentStatus(userID, trxID int) (*response.TRXResponse, error) {
	// Get transaction and validate ownership
	trx, err := s.trxRepo.GetByID(trxID)
	if err != nil {
		return nil, errors.New("Transaction not found")
	}

	if trx.IDUser != userID {
		return nil, errors.New(constants.ErrForbidden)
	}

	// Only check payment for non-COD transactions
	if trx.MethodBayar == "COD" {
		return nil, errors.New("COD transactions do not require payment verification")
	}

	// Check if transaction has invoice code
	if trx.KodeInvoice == "" {
		return nil, errors.New("Transaction does not have invoice code")
	}

	// Verify payment status from Midtrans
	paymentStatus, err := s.midtransService.VerifyPayment(trx.KodeInvoice)
	if err != nil {
		return nil, fmt.Errorf("failed to verify payment: %w", err)
	}

	// Map Midtrans status to our payment status
	paymentStatusStr := s.mapMidtransStatusToPaymentStatus(paymentStatus.TransactionStatus)
	vaNumbers := mapVANumbersFromMidtrans(paymentStatus.VaNumbers)
	vaNumbersJSON := serializeVANumbersToJSON(vaNumbers)

	// Update transaction payment status (use UpdatePaymentStatus to avoid updating created_at)
	if err := s.trxRepo.UpdatePaymentStatus(trxID, paymentStatusStr, "", "", "", nil, vaNumbersJSON); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// Get updated transaction with relations
	updatedTRX, err := s.trxRepo.GetByID(trxID)
	if err != nil {
		return nil, err
	}

	trxResponse := s.mapTRXToResponse(*updatedTRX)
	vaNumbersUpdated := mapVANumbersFromMidtrans(paymentStatus.VaNumbers)
	if len(vaNumbersUpdated) > 0 {
		trxResponse.PaymentVANumbers = vaNumbersUpdated
	} else {
		s.attachVANumbersIfNeeded(updatedTRX, &trxResponse)
	}
	return &trxResponse, nil
}

func (s *trxService) mapTRXToResponse(trx model.TRX) response.TRXResponse {
	paymentVANumbers := deserializeVANumbersFromString(trx.PaymentVANumbers)
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

	// Format payment expired at
	var paymentExpiredAtStr string
	if trx.PaymentExpiredAt != nil {
		paymentExpiredAtStr = trx.PaymentExpiredAt.Format("2006-01-02 15:04:05")
	}

	return response.TRXResponse{
		ID:               trx.ID,
		HargaTotal:       trx.HargaTotal,
		KodeInvoice:      trx.KodeInvoice,
		MethodBayar:      trx.MethodBayar,
		PaymentStatus:    trx.PaymentStatus,
		PaymentURL:       trx.PaymentURL,
		PaymentExpiredAt: paymentExpiredAtStr,
		PaymentVANumbers: paymentVANumbers,
		CreatedAt:        trx.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        trx.UpdatedAt.Format("2006-01-02 15:04:05"),
		IDUser:           trx.IDUser,
		IDAlamat:         trx.IDAlamat,
		User:             userResponse,
		Address:          addressResponse,
		DetailTRX:        detailResponses,
	}
}

func (s *trxService) attachVANumbersIfNeeded(trx *model.TRX, trxResponse *response.TRXResponse) {
	if trx == nil || trxResponse == nil {
		return
	}
	if len(trxResponse.PaymentVANumbers) > 0 {
		return
	}
	if !isVirtualAccountMethod(trx.MethodBayar) || trx.KodeInvoice == "" {
		return
	}
	vaNumbers, err := s.fetchVANumbersFromMidtrans(trx.KodeInvoice)
	if err != nil || len(vaNumbers) == 0 {
		return
	}
	trxResponse.PaymentVANumbers = vaNumbers
	_ = s.trxRepo.UpdatePaymentStatus(trx.ID, trx.PaymentStatus, "", "", "", nil, serializeVANumbersToJSON(vaNumbers))
}

func (s *trxService) fetchVANumbersFromMidtrans(orderID string) ([]response.PaymentVANumber, error) {
	status, err := s.midtransService.VerifyPayment(orderID)
	if err != nil {
		return nil, err
	}
	return mapVANumbersFromMidtrans(status.VaNumbers), nil
}

func mapVANumbersFromMidtrans(vaData []map[string]interface{}) []response.PaymentVANumber {
	var result []response.PaymentVANumber
	for _, entry := range vaData {
		bank, _ := entry["bank"].(string)
		vaNumber, _ := entry["va_number"].(string)
		if vaNumber == "" {
			vaNumber, _ = entry["virtual_account_number"].(string)
		}
		if vaNumber == "" {
			continue
		}
		result = append(result, response.PaymentVANumber{
			Bank:     strings.ToUpper(bank),
			VANumber: vaNumber,
		})
	}
	return result
}

func serializeVANumbersToJSON(numbers []response.PaymentVANumber) string {
	if len(numbers) == 0 {
		return ""
	}
	bytes, err := json.Marshal(numbers)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func deserializeVANumbersFromString(data string) []response.PaymentVANumber {
	if strings.TrimSpace(data) == "" {
		return nil
	}
	var numbers []response.PaymentVANumber
	if err := json.Unmarshal([]byte(data), &numbers); err != nil {
		return nil
	}
	return numbers
}

func normalizePaymentMethod(method string) string {
	lower := strings.ToLower(method)
	if normalized, ok := constants.PaymentMethodAliases[lower]; ok {
		return normalized
	}
	return lower
}

func isVirtualAccountMethod(method string) bool {
	switch normalizePaymentMethod(method) {
	case strings.ToLower(constants.PaymentMethodVirtualAccount), strings.ToLower(constants.PaymentMethodBankTransfer):
		return true
	default:
		return false
	}
}
