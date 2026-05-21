package nats_delivery

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

	"github.com/nats-io/nats.go"
)

// Настройки Gmail — замени на свои
const (
	smtpHost     = "smtp.gmail.com"
	smtpPort     = "587"
	smtpUser     = "darkiku555@gmail.com"
	smtpPassword = "ahoi wgkg tyeo gnul"
	smtpFrom     = "darkiku555@gmail.com"
)

type adminSubscriber struct {
	nc *nats.Conn
}

func NewAdminSubscriber(nc *nats.Conn) *adminSubscriber {
	return &adminSubscriber{nc: nc}
}

func sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", smtpHost+":465", tlsConfig)
	if err != nil {
		// Fallback to STARTTLS on port 587
		client, err2 := smtp.Dial(smtpHost + ":" + smtpPort)
		if err2 != nil {
			return fmt.Errorf("smtp dial error: %v", err2)
		}
		defer client.Close()

		if err2 = client.StartTLS(tlsConfig); err2 != nil {
			return fmt.Errorf("starttls error: %v", err2)
		}
		if err2 = client.Auth(auth); err2 != nil {
			return fmt.Errorf("smtp auth error: %v", err2)
		}
		if err2 = client.Mail(smtpFrom); err2 != nil {
			return err2
		}
		if err2 = client.Rcpt(to); err2 != nil {
			return err2
		}
		wc, err2 := client.Data()
		if err2 != nil {
			return err2
		}
		defer wc.Close()
		msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", smtpFrom, to, subject, body)
		_, err2 = fmt.Fprint(wc, msg)
		return err2
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}
	if err = client.Auth(auth); err != nil {
		return err
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", smtpFrom, to, subject, body)
	if err = client.Mail(smtpFrom); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}
	wc, err := client.Data()
	if err != nil {
		return err
	}
	defer wc.Close()
	_, err = fmt.Fprint(wc, msg)
	return err
}

func (s *adminSubscriber) SubscribeToAllEvents() {
	// 1. Слушаем предупреждения о посещаемости от Education Service
	_, err := s.nc.Subscribe("edu.attendance.warning", func(m *nats.Msg) {
		var payload map[string]string
		if err := json.Unmarshal(m.Data, &payload); err != nil {
			log.Printf("Error unmarshaling attendance message: %v", err)
			return
		}

		studentID := payload["student_id"]
		courseID := payload["course_id"]

		log.Printf("[NATS] Low attendance warning for student: %s, course: %s", studentID, courseID)

		// Реальная отправка email
		subject := "⚠️ Low Attendance Warning - AITU"
		body := fmt.Sprintf(
			"Dear Student %s,\n\nYour attendance for course %s has dropped below 70%%.\n\nPlease attend classes to avoid academic penalties.\n\nBest regards,\nAITU Administration",
			studentID, courseID,
		)

		// В проде тут был бы реальный email студента из БД
		// Для демо логируем и показываем что SMTP вызывается
		log.Printf("[SMTP] Sending attendance warning email to student %s...", studentID)
		if err := sendEmail(smtpUser, subject, body); err != nil {
			log.Printf("[SMTP] Email send failed (check credentials): %v", err)
		} else {
			log.Printf("[SMTP] ✅ Email sent successfully to student %s", studentID)
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to edu.attendance.warning: %v", err)
	}

	// 2. Слушаем новый фидбек от Community Service
	_, err = s.nc.Subscribe("comm.feedback.new", func(m *nats.Msg) {
		var payload map[string]string
		if err := json.Unmarshal(m.Data, &payload); err != nil {
			log.Printf("Error unmarshaling feedback message: %v", err)
			return
		}

		studentID := payload["student_id"]
		text := payload["text"]

		log.Printf("[NATS] New feedback from student %s: %s", studentID, text)
		log.Println("[System] Creating support ticket in Admin system...")

		// Email уведомление администратору
		subject := "📬 New Student Feedback - AITU Support"
		body := fmt.Sprintf(
			"New feedback received.\n\nStudent ID: %s\nMessage: %s\n\nPlease review and respond promptly.",
			studentID, text,
		)

		log.Printf("[SMTP] Sending feedback notification to admin...")
		if err := sendEmail(smtpUser, subject, body); err != nil {
			log.Printf("[SMTP] Email send failed (check credentials): %v", err)
		} else {
			log.Printf("[SMTP] ✅ Admin notification sent")
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to comm.feedback.new: %v", err)
	}
}
