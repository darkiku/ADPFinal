package nats_delivery

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type EventPublisher interface {
	PublishLowAttendanceWarning(studentID string, courseID string) error
}

type natsPublisher struct {
	nc *nats.Conn
}

func NewNatsPublisher(nc *nats.Conn) EventPublisher {
	return &natsPublisher{nc: nc}
}

func (p *natsPublisher) PublishLowAttendanceWarning(studentID string, courseID string) error {
	// Формируем payload (тело сообщения)
	payload := map[string]string{
		"student_id": studentID,
		"course_id":  courseID,
		"reason":     "Attendance dropped below 70%",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Публикуем в топик edu.attendance.warning
	err = p.nc.Publish("edu.attendance.warning", data)
	if err != nil {
		log.Printf("Failed to publish to NATS: %v", err)
		return err
	}

	log.Printf("Published warning to NATS for student %s", studentID)
	return nil
}
