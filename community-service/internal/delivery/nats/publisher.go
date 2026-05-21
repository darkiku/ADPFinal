package nats_delivery

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type communityPublisher struct {
	nc *nats.Conn
}

func NewCommunityPublisher(nc *nats.Conn) *communityPublisher {
	return &communityPublisher{nc: nc}
}

func (p *communityPublisher) PublishNewFeedback(studentID, feedbackText string) error {
	payload := map[string]string{
		"student_id": studentID,
		"text":       feedbackText,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Топик comm.feedback.new (из README)
	err = p.nc.Publish("comm.feedback.new", data)
	if err != nil {
		log.Printf("Failed to publish feedback to NATS: %v", err)
		return err
	}

	log.Println("Feedback sent to NATS queue")
	return nil
}
