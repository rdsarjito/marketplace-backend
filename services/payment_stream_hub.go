package services

import (
	"log"
	"sync"
)

// PaymentClient represents a single SSE client channel
type PaymentClient chan string

// PaymentHub manages SSE subscribers per transaction ID
type PaymentHub struct {
	mu      sync.RWMutex
	clients map[int]map[PaymentClient]struct{}
}

// NewPaymentHub creates a new PaymentHub instance
func NewPaymentHub() *PaymentHub {
	return &PaymentHub{
		clients: make(map[int]map[PaymentClient]struct{}),
	}
}

// Subscribe registers a new client for a given transaction ID
func (h *PaymentHub) Subscribe(trxID int) PaymentClient {
	ch := make(PaymentClient, 10) // Buffered channel to prevent blocking

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[trxID] == nil {
		h.clients[trxID] = make(map[PaymentClient]struct{})
	}
	h.clients[trxID][ch] = struct{}{}

	return ch
}

// GetClientCount returns the number of clients subscribed to a transaction ID
func (h *PaymentHub) GetClientCount(trxID int) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[trxID]; ok {
		return len(clients)
	}
	return 0
}

// Unsubscribe removes a client from a given transaction ID
func (h *PaymentHub) Unsubscribe(trxID int, ch PaymentClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.clients[trxID]; ok {
		if _, exists := clients[ch]; exists {
			delete(clients, ch)
			close(ch)
		}
		if len(clients) == 0 {
			delete(h.clients, trxID)
		}
	}
}

// Publish sends a message to all clients subscribed to a transaction ID
func (h *PaymentHub) Publish(trxID int, msg string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[trxID]; ok {
		clientCount := len(clients)
		if clientCount > 0 {
			log.Printf("[SSE] Publishing to %d client(s) for transaction %d: %s", clientCount, trxID, msg)
			sentCount := 0
			for ch := range clients {
				select {
				case ch <- msg:
					sentCount++
				default:
					// Drop message if client is not ready to receive (channel buffer full)
					log.Printf("[SSE] Warning: Failed to send message to client for transaction %d (channel full)", trxID)
				}
			}
			log.Printf("[SSE] Sent message to %d/%d client(s) for transaction %d", sentCount, clientCount, trxID)
		} else {
			log.Printf("[SSE] No clients connected for transaction %d, message will be dropped", trxID)
		}
	} else {
		log.Printf("[SSE] No clients subscribed for transaction %d", trxID)
	}
}

// PaymentStatusHub is a global hub instance used by transaction service and handlers
var PaymentStatusHub = NewPaymentHub()
