package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
}

// NewMidtransService membuat instance baru dari MidtransService
func NewMidtransService(serverKey, clientKey string, isProduction bool) MidtransService {
	baseURL := "https://api.sandbox.midtrans.com"
	if isProduction {
		baseURL = "https://api.midtrans.com"
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
	requestBody := map[string]interface{}{
		"payment_type": req.PaymentType,
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": req.GrossAmount,
		},
		"customer_details": req.CustomerDetails,
		"item_details":     req.ItemDetails,
	}

	// Add custom expiry if provided
	if req.CustomExpiry != nil {
		requestBody["custom_expiry"] = req.CustomExpiry
	}

	// Add payment-specific parameters based on payment type
	switch req.PaymentType {
	case "bank_transfer":
		// Bank transfer (BCA, BNI, Mandiri, etc)
		requestBody["bank_transfer"] = map[string]interface{}{
			"bank": "bca", // default, bisa diubah sesuai kebutuhan
		}
	case "virtual_account":
		// Virtual Account
		requestBody["bank_transfer"] = map[string]interface{}{
			"bank": "bca", // default
		}
	case "e_wallet":
		// E-Wallet (GoPay, OVO, DANA, LinkAja)
		requestBody["e_wallet"] = map[string]interface{}{
			"store": "gopay", // default, bisa diubah
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

	// Check for errors
	if response.StatusCode != "201" && response.StatusCode != "200" {
		return nil, fmt.Errorf("midtrans API error: %s - %s", response.StatusCode, response.StatusMessage)
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
	if response.StatusCode != "200" {
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
