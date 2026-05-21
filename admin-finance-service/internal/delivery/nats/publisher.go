package nats_delivery

import (
	"encoding/json"
	"log"

	"aitu-superapp/admin-finance-service/internal/domain"
	"github.com/nats-io/nats.go"
)

type adminPublisher struct {
	nc *nats.Conn
}

func NewAdminPublisher(nc *nats.Conn) domain.AdminEventPublisher {
	return &adminPublisher{nc: nc}
}

func (p *adminPublisher) PublishPaymentSuccess(userID string, amount float32) error {
	payload := map[string]interface{}{
		"user_id": userID,
		"amount":  amount,
		"status":  "success",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Публикуем событие "fin.payment.success" (как было в таблице требований README)
	err = p.nc.Publish("fin.payment.success", data)
	if err != nil {
		log.Printf("Failed to publish payment success to NATS: %v", err)
		return err
	}

	log.Printf("Published payment success for user %s", userID)
	return nil
}
