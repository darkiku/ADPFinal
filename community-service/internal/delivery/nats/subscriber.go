package nats_delivery

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type communitySubscriber struct {
	nc *nats.Conn
}

func NewCommunitySubscriber(nc *nats.Conn) *communitySubscriber {
	return &communitySubscriber{nc: nc}
}

func (s *communitySubscriber) SubscribeToFinanceEvents() {
	_, err := s.nc.Subscribe("fin.payment.success", func(m *nats.Msg) {
		var payload map[string]interface{}
		if err := json.Unmarshal(m.Data, &payload); err != nil {
			log.Printf("Error unmarshaling NATS message: %v", err)
			return
		}

		userID := payload["user_id"].(string)
		log.Printf("Community Service received payment success for user: %s. Unlocking premium events...", userID)

		// Здесь должна быть логика обновления статуса в Postgres
		// Например: uc.UnlockPremiumEventAccess(userID)
	})

	if err != nil {
		log.Fatalf("Failed to subscribe to fin.payment.success: %v", err)
	}
}
