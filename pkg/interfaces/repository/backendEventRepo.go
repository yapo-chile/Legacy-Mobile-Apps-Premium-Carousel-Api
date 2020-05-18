package repository

import (
	"encoding/json"
	"fmt"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

// producer allows to push events to queue
type producer struct {
	handler              KafkaProducer
	premiumProductsTopic string
}

// MakeBackendEventsProducer creates new instance of Producer for backend events
func MakeBackendEventsProducer(handler KafkaProducer, premiumProductsTopic string) usecases.BackendEventsRepository {
	return &producer{
		handler:              handler,
		premiumProductsTopic: premiumProductsTopic,
	}
}

// kafkaMessage defines the valid input supported by backend events
type kafkaMessage struct {
	Type      EventType   `json:"type"`
	Date      string      `json:"date"`
	Timestamp string      `json:"ts_utc"`
	Content   interface{} `json:"content"`
}

// EventType defines valid events to push to backend events
type EventType string

const (
	// PremiumCarouselPurchase represents a premium carousel purchase event type
	PremiumCarouselPurchase EventType = "premium_carousel_purchase"
)

// Push pushes given product to backend events through purchases topic
func (p *producer) PushSoldProduct(product domain.Product) error {
	switch product.Type {
	case domain.PremiumCarousel:
		content := map[string]interface{}{
			"id":              product.ID,
			"type":            product.Type,
			"user_id":         product.UserID,
			"email":           product.Email,
			"purchase_id":     product.Purchase.ID,
			"purchase_number": product.Purchase.Number,
			"purchase_price":  product.Purchase.Price,
			"purchase_status": product.Purchase.Status,
			"purchase_type":   product.Purchase.Type,
			"status":          product.Status,
			"expired_at":      product.ExpiredAt.String(),
			"created_at":      product.CreatedAt.String(),
		}
		message := kafkaMessage{
			Type:      PremiumCarouselPurchase,
			Date:      product.CreatedAt.Format("2006-01-02 15:04:05"),
			Timestamp: fmt.Sprintf("%d", product.CreatedAt.Unix()),
			Content:   content,
		}
		bytes, _ := json.Marshal(message) // nolint
		return p.handler.SendMessage(p.premiumProductsTopic, bytes)
	default:
		return fmt.Errorf("Product not supported")
	}
}
