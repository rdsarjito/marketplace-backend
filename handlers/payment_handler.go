package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rdsarjito/marketplace-backend/services"
)

type PaymentHandler struct {
	trxService services.TRXService
}

func NewPaymentHandler(trxService services.TRXService) *PaymentHandler {
	return &PaymentHandler{
		trxService: trxService,
	}
}

// HandleWebhook handles webhook notification from Midtrans
// This endpoint should be publicly accessible (no auth middleware)
// Midtrans will send POST request to this endpoint when payment status changes
func (h *PaymentHandler) HandleWebhook(c *fiber.Ctx) error {
	// Parse notification body from Midtrans
	var notification map[string]interface{}
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request body",
		})
	}

	// Handle payment webhook
	if err := h.trxService.HandlePaymentWebhook(notification); err != nil {
		// Log error but still return 200 to Midtrans
		// Midtrans will retry if we return error status
		// In production, you should log this error properly
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	// Return success response to Midtrans
	// Midtrans expects 200 OK status
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Webhook processed successfully",
	})
}

// CheckPaymentStatus manually checks payment status from Midtrans
// This can be used for polling or manual verification
func (h *PaymentHandler) CheckPaymentStatus(c *fiber.Ctx) error {
	// Get order_id from query or body
	orderID := c.Query("order_id")
	if orderID == "" {
		var body map[string]string
		if err := c.BodyParser(&body); err == nil {
			orderID = body["order_id"]
		}
	}

	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "order_id is required",
		})
	}

	// Create notification map for service
	notification := map[string]interface{}{
		"order_id": orderID,
	}

	// Handle payment webhook (which will verify and update status)
	if err := h.trxService.HandlePaymentWebhook(notification); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Payment status checked and updated successfully",
	})
}

