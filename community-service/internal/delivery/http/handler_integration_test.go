package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	http_delivery "aitu-superapp/community-service/internal/delivery/http"
	"aitu-superapp/community-service/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type mockCommUseCase struct{}

func (m *mockCommUseCase) BookSpace(ctx context.Context, spaceID, studentID, start, end string) error {
	return nil
}
func (m *mockCommUseCase) RegisterForEvent(ctx context.Context, eventID, studentID string) error {
	return nil
}
func (m *mockCommUseCase) CreateMarketplaceAd(ctx context.Context, studentID, title, desc string, price float32) (string, error) {
	return "ad-001", nil
}
func (m *mockCommUseCase) SubmitFeedback(ctx context.Context, studentID, text string) error {
	return nil
}
func (m *mockCommUseCase) GetEventsList(ctx context.Context) ([]domain.Event, error) {
	return []domain.Event{{ID: "evt-1", Title: "Hackathon 2026", Location: "Main Hall"}}, nil
}
func (m *mockCommUseCase) CancelEvent(ctx context.Context, eventID string) error { return nil }
func (m *mockCommUseCase) GetCoworkingList(ctx context.Context) ([]string, error) {
	return []string{"Room 101", "Room 102"}, nil
}
func (m *mockCommUseCase) CancelBooking(ctx context.Context, bookingID string) error { return nil }
func (m *mockCommUseCase) GetMarketplaceItems(ctx context.Context) ([]domain.Ad, error) {
	return []domain.Ad{{ID: "ad-1", Title: "Calculus Book", Price: 3000}}, nil
}
func (m *mockCommUseCase) ReportLostFound(ctx context.Context, studentID, item, desc, status string) error {
	return nil
}
func (m *mockCommUseCase) GetLostFoundList(ctx context.Context) ([]domain.LostFound, error) {
	return []domain.LostFound{}, nil
}
func (m *mockCommUseCase) GetNewsFeed(ctx context.Context) ([]string, error) {
	return []string{"AITU wins ranking!", "New library opens"}, nil
}

func setupCommApp() *fiber.App {
	app := fiber.New()
	h := http_delivery.NewCommunityHandler(&mockCommUseCase{})
	h.RegisterRoutes(app)
	return app
}

func TestIntegration_GetEventsList(t *testing.T) {
	app := setupCommApp()
	req := httptest.NewRequest("GET", "/api/v1/community/events/list", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIntegration_GetCoworkingList(t *testing.T) {
	app := setupCommApp()
	req := httptest.NewRequest("GET", "/api/v1/community/coworking/list", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIntegration_GetMarketplaceItems(t *testing.T) {
	app := setupCommApp()
	req := httptest.NewRequest("GET", "/api/v1/community/marketplace/items", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIntegration_GetNewsFeed(t *testing.T) {
	app := setupCommApp()
	req := httptest.NewRequest("GET", "/api/v1/community/news/feed", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIntegration_PostLostFoundReport(t *testing.T) {
	app := setupCommApp()
	body, _ := json.Marshal(map[string]string{
		"item": "Laptop", "description": "Black Dell", "status": "lost",
	})
	req := httptest.NewRequest("POST", "/api/v1/community/lost-found/report", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "student-123")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		t.Errorf("expected 200 or 201, got %d", resp.StatusCode)
	}
}
