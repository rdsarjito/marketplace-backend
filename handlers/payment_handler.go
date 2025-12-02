package handlers

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rdsarjito/marketplace-backend/services"
	"github.com/rdsarjito/marketplace-backend/utils"
)

type PaymentHandler struct {
	trxService  services.TRXService
	userService services.UserService
}

func NewPaymentHandler(trxService services.TRXService, userService services.UserService) *PaymentHandler {
	return &PaymentHandler{
		trxService:  trxService,
		userService: userService,
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

// StreamPaymentStatus sends payment status updates via Server-Sent Events (SSE)
// This endpoint expects a JWT token in the query parameter (?token=...)
// and validates that the authenticated user owns the requested transaction.
func (h *PaymentHandler) StreamPaymentStatus(c *fiber.Ctx) error {
	// Validate token from query parameter
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Missing token",
		})
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid token",
		})
	}

	// Optional: ensure user still exists
	if _, err := h.userService.GetMyProfile(claims.UserID); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "User not found",
		})
	}

	// Parse transaction ID from path
	trxIDStr := c.Params("id")
	trxID, err := strconv.Atoi(trxIDStr)
	if err != nil || trxID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid transaction ID",
		})
	}

	// Ensure transaction belongs to the authenticated user
	if _, err := h.trxService.GetDetailTRX(claims.UserID, trxID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no") // Disable nginx buffering

	client := services.PaymentStatusHub.Subscribe(trxID)
	log.Printf("[SSE] Client subscribed for transaction %d (total clients: %d)", trxID, services.PaymentStatusHub.GetClientCount(trxID))

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// Unsubscribe when connection ends (defer inside SetBodyStreamWriter)
		defer func() {
			services.PaymentStatusHub.Unsubscribe(trxID, client)
			log.Printf("[SSE] Client unsubscribed for transaction %d", trxID)
		}()

		// Send initial connection message to establish connection
		// Use proper SSE format with event and data
		fmt.Fprintf(w, "event: connected\n")
		fmt.Fprintf(w, "data: {\"trx_id\": %d, \"status\": \"connected\"}\n\n", trxID)
		w.Flush()
		log.Printf("[SSE] Sent initial connection message for transaction %d", trxID)

		// Create ticker for keepalive (ping every 10 seconds to keep connection alive)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		// Handle messages from hub and keepalive in the same select
		for {
			select {
			case msg, ok := <-client:
				if !ok {
					// Channel closed, connection ended
					log.Printf("[SSE] Client channel closed for transaction %d", trxID)
					return
				}
				// msg is expected to be a JSON string
				fmt.Fprintf(w, "event: payment_updated\n")
				fmt.Fprintf(w, "data: %s\n\n", msg)
				w.Flush()
				log.Printf("[SSE] Sent message to client for transaction %d", trxID)
			case <-ticker.C:
				// Send keepalive comment to keep connection alive
				fmt.Fprintf(w, ": keepalive\n\n")
				w.Flush()
				log.Printf("[SSE] Sent keepalive for transaction %d", trxID)
			}
		}
	})

	return nil
}
