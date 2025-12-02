package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// MidtransService interface untuk payment gateway operations
type MidtransService interface {
	CreatePayment(req *CreatePaymentRequest) (*CreatePaymentResponse, error)
	VerifyPayment(orderID string) (*PaymentStatusResponse, error)
	HandleWebhook(notification map[string]interface{}) (*PaymentStatusResponse, error)
}

type midtransService struct {
	serverKey    string
	clientKey    string
	isProduction bool
	baseURL      string
	client       *http.Client
}

// CreatePaymentRequest untuk membuat payment request ke Midtrans
type CreatePaymentRequest struct {
	OrderID         string                   `json:"order_id"`
	GrossAmount     int                      `json:"gross_amount"`
	PaymentType     string                   `json:"payment_type"` // virtual_account, e_wallet, bank_transfer, etc
	CustomerDetails map[string]interface{}   `json:"customer_details"`
	ItemDetails     []map[string]interface{} `json:"item_details"`
	CustomExpiry    *CustomExpiry            `json:"custom_expiry,omitempty"`
	FinishURL       string                   `json:"finish_url,omitempty"` // URL untuk redirect setelah payment selesai
}

// CustomExpiry untuk set expiration time
type CustomExpiry struct {
	ExpiryDuration int    `json:"expiry_duration"` // in minutes
	Unit           string `json:"unit"`            // "minute"
}

// CreatePaymentResponse response dari Midtrans setelah create payment
type CreatePaymentResponse struct {
	Token             string                   `json:"token,omitempty"`
	RedirectURL       string                   `json:"redirect_url,omitempty"`
	StatusCode        string                   `json:"status_code"`
	StatusMessage     string                   `json:"status_message"`
	TransactionID     string                   `json:"transaction_id,omitempty"`
	OrderID           string                   `json:"order_id,omitempty"`
	GrossAmount       string                   `json:"gross_amount,omitempty"`
	PaymentType       string                   `json:"payment_type,omitempty"`
	TransactionTime   string                   `json:"transaction_time,omitempty"`
	TransactionStatus string                   `json:"transaction_status,omitempty"`
	VaNumbers         []map[string]interface{} `json:"va_numbers,omitempty"`
	Actions           []map[string]interface{} `json:"actions,omitempty"`
	ExpiryTime        string                   `json:"expiry_time,omitempty"`
	QRString          string                   `json:"qr_string,omitempty"`
}

// PaymentStatusResponse response untuk status payment
type PaymentStatusResponse struct {
	StatusCode        string                   `json:"status_code"`
	StatusMessage     string                   `json:"status_message"`
	TransactionID     string                   `json:"transaction_id"`
	OrderID           string                   `json:"order_id"`
	GrossAmount       string                   `json:"gross_amount"`
	Currency          string                   `json:"currency"`
	PaymentType       string                   `json:"payment_type"`
	TransactionTime   string                   `json:"transaction_time"`
	TransactionStatus string                   `json:"transaction_status"`
	SettlementTime    string                   `json:"settlement_time,omitempty"`
	VaNumbers         []map[string]interface{} `json:"va_numbers,omitempty"`
	Actions           []map[string]interface{} `json:"actions,omitempty"`
	FraudStatus       string                   `json:"fraud_status,omitempty"`
	QRString          string                   `json:"qr_string,omitempty"`
}

// NewMidtransService membuat instance baru dari MidtransService
func NewMidtransService(serverKey, clientKey string, isProduction bool) MidtransService {
	baseURL := "https://api.sandbox.midtrans.com"
	if isProduction {
		baseURL = "https://api.midtrans.com"
	}

	// Validate server key (should not be empty)
	if serverKey == "" {
		// Log warning but don't fail - will fail when trying to create payment
		fmt.Println("WARNING: Midtrans Server Key is empty. Payment creation will fail.")
	}

	return &midtransService{
		serverKey:    serverKey,
		clientKey:    clientKey,
		isProduction: isProduction,
		baseURL:      baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreatePayment membuat payment request ke Midtrans
func (s *midtransService) CreatePayment(req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	// Build request body
	// Note: For e_wallet, payment_type should be the specific wallet (gopay, ovo, dana, linkaja)
	// not "e_wallet" itself
	paymentType := req.PaymentType

	requestBody := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": req.GrossAmount,
		},
		"customer_details": req.CustomerDetails,
		"item_details":     req.ItemDetails,
	}

	// Add finish_url if provided (for redirect after payment)
	if req.FinishURL != "" {
		requestBody["finish_url"] = req.FinishURL
	}

	// Add custom expiry if provided
	if req.CustomExpiry != nil {
		requestBody["custom_expiry"] = req.CustomExpiry
	}

	// Add payment-specific parameters based on payment type
	switch paymentType {
	case "bank_transfer", "virtual_account":
		// Bank transfer (BCA, BNI, Mandiri, etc)
		requestBody["payment_type"] = "bank_transfer"
		requestBody["bank_transfer"] = map[string]interface{}{
			"bank": "bca", // default, bisa diubah sesuai kebutuhan
		}
	case "e_wallet", "gopay", "ovo", "dana", "linkaja":
		// E-Wallet - payment_type should be the specific wallet name
		// Map to specific wallet or default to gopay
		walletType := "gopay" // default
		if paymentType == "ovo" {
			walletType = "ovo"
		} else if paymentType == "dana" {
			walletType = "dana"
		} else if paymentType == "linkaja" {
			walletType = "linkaja"
		}
		requestBody["payment_type"] = walletType
	case "credit_card", "cc":
		requestBody["payment_type"] = "credit_card"
	default:
		// Default to bank_transfer
		requestBody["payment_type"] = "bank_transfer"
		requestBody["bank_transfer"] = map[string]interface{}{
			"bank": "bca",
		}
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/v2/charge", s.baseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	auth := base64.StdEncoding.EncodeToString([]byte(s.serverKey + ":"))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+auth)

	// Execute request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var response CreatePaymentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("[Midtrans] CreatePayment Response for OrderID %s:", req.OrderID)
	log.Printf("[Midtrans] Full Response Body: %s", string(body))
	log.Printf("[Midtrans] Payment Type: %s", response.PaymentType)
	log.Printf("[Midtrans] QR String: %s", response.QRString)
	log.Printf("[Midtrans] Actions: %+v", response.Actions)
	log.Printf("[Midtrans] Status Code: %s", response.StatusCode)
	log.Printf("[Midtrans] Status Message: %s", response.StatusMessage)

	// Check for errors
	if response.StatusCode != "201" && response.StatusCode != "200" {
		// Return more detailed error message including response body for debugging
		return nil, fmt.Errorf("midtrans API error [%s]: %s. Response: %s", response.StatusCode, response.StatusMessage, string(body))
	}

	// For bank_transfer/virtual_account, Midtrans doesn't always return redirect_url
	// Instead, it returns VA numbers or transfer instructions
	// We need to check if we have either redirect_url OR va_numbers
	hasRedirectURL := response.RedirectURL != ""
	hasToken := response.Token != ""
	hasVaNumbers := len(response.VaNumbers) > 0
	hasActions := len(response.Actions) > 0

	// For bank_transfer/virtual_account, va_numbers or actions are valid
	// For e_wallet (gopay, ovo, etc), actions with deeplink-redirect or qr-code are valid
	// For credit_card, redirect_url or token is required
	isBankTransfer := paymentType == "bank_transfer" || paymentType == "virtual_account"
	isEWallet := paymentType == "gopay" || paymentType == "ovo" || paymentType == "dana" || paymentType == "linkaja"

	if isBankTransfer {
		// Bank transfer is valid if we have va_numbers or actions
		if !hasVaNumbers && !hasActions && !hasRedirectURL {
			return nil, fmt.Errorf("midtrans API returned empty payment data for bank transfer. Status: %s - %s. Response: %s", response.StatusCode, response.StatusMessage, string(body))
		}
	} else if isEWallet {
		// For e_wallet, actions with deeplink-redirect or qr-code are valid
		if !hasActions && !hasRedirectURL && !hasToken {
			return nil, fmt.Errorf("midtrans API returned empty payment data for e-wallet. Status: %s - %s. Response: %s", response.StatusCode, response.StatusMessage, string(body))
		}
	} else {
		// For credit_card, we need redirect_url or token
		if !hasRedirectURL && !hasToken {
			return nil, fmt.Errorf("midtrans API returned empty payment URL. Status: %s - %s. Response: %s", response.StatusCode, response.StatusMessage, string(body))
		}
	}

	return &response, nil
}

// VerifyPayment mengecek status payment dari Midtrans
func (s *midtransService) VerifyPayment(orderID string) (*PaymentStatusResponse, error) {
	url := fmt.Sprintf("%s/v2/%s/status", s.baseURL, orderID)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	auth := base64.StdEncoding.EncodeToString([]byte(s.serverKey + ":"))
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+auth)

	// Execute request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var response PaymentStatusResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	// In sandbox, Midtrans sometimes returns 201 with message "Success, transaction is found"
	// Treat both 200 and 201 as successful status responses.
	if response.StatusCode != "200" && response.StatusCode != "201" {
		return nil, fmt.Errorf("midtrans API error: %s - %s", response.StatusCode, response.StatusMessage)
	}

	return &response, nil
}

// HandleWebhook memproses webhook notification dari Midtrans
func (s *midtransService) HandleWebhook(notification map[string]interface{}) (*PaymentStatusResponse, error) {
	// Extract order_id from notification
	orderID, ok := notification["order_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid notification: missing order_id")
	}

	// Verify payment status using order_id
	return s.VerifyPayment(orderID)
}
