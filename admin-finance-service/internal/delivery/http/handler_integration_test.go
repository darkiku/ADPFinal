package http_delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	http_delivery "aitu-superapp/admin-finance-service/internal/delivery/http"
	"aitu-superapp/admin-finance-service/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type mockAdminUseCase struct{}

func (m *mockAdminUseCase) GetAccountBalance(ctx context.Context, userID string) (float32, error) {
	return 75000, nil
}
func (m *mockAdminUseCase) InitiatePayment(ctx context.Context, userID string, amount float32, purpose string) (string, error) {
	return "tx-test-001", nil
}
func (m *mockAdminUseCase) RequestCertificate(ctx context.Context, studentID, docType string) (string, error) {
	return "doc-test-001", nil
}
func (m *mockAdminUseCase) UpdateStudentProfile(ctx context.Context, studentID, email, phone string) error {
	return nil
}
func (m *mockAdminUseCase) GetDocumentStatus(ctx context.Context, studentID string) ([]domain.DocumentRequest, error) {
	return []domain.DocumentRequest{{ID: "doc-1", StudentID: studentID, Type: "transcript", Status: "ready"}}, nil
}
func (m *mockAdminUseCase) GetDormitoryInfo(ctx context.Context, studentID string) (*domain.Dormitory, error) {
	return &domain.Dormitory{ID: "dorm-1", RoomNumber: "304", OccupantID: studentID}, nil
}
func (m *mockAdminUseCase) RequestDormRepair(ctx context.Context, studentID, description string) error {
	return nil
}
func (m *mockAdminUseCase) GetFinanceHistory(ctx context.Context, userID string) ([]domain.Transaction, error) {
	return []domain.Transaction{}, nil
}
func (m *mockAdminUseCase) GetScholarshipInfo(ctx context.Context, studentID string) (map[string]interface{}, error) {
	return map[string]interface{}{"type": "State Grant", "amount": 50000}, nil
}
func (m *mockAdminUseCase) GetLibraryDebts(ctx context.Context, studentID string) ([]string, error) {
	return []string{}, nil
}
func (m *mockAdminUseCase) ReserveLibraryBook(ctx context.Context, studentID, bookID string) error {
	return nil
}

func setupAdminApp() *fiber.App {
	app := fiber.New()
	http_delivery.RegisterRESTHandlers(app, &mockAdminUseCase{})
	return app
}

func TestIntegration_GetDocumentStatus(t *testing.T) {
	app := setupAdminApp()
	req := httptest.NewRequest("GET", "/api/v1/admin/docs/status", nil)
	req.Header.Set("X-User-ID", "student-123")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIntegration_GetDormitoryInfo(t *testing.T) {
	app := setupAdminApp()
	req := httptest.NewRequest("GET", "/api/v1/admin/dormitory/info", nil)
	req.Header.Set("X-User-ID", "student-123")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIntegration_PostDormRepair(t *testing.T) {
	app := setupAdminApp()
	body, _ := json.Marshal(map[string]string{"description": "Broken window"})
	req := httptest.NewRequest("POST", "/api/v1/admin/dormitory/repair", bytes.NewReader(body))
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

func TestIntegration_GetScholarshipInfo(t *testing.T) {
	app := setupAdminApp()
	req := httptest.NewRequest("GET", "/api/v1/admin/scholarship/info", nil)
	req.Header.Set("X-User-ID", "student-123")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIntegration_GetLibraryDebts(t *testing.T) {
	app := setupAdminApp()
	req := httptest.NewRequest("GET", "/api/v1/admin/library/debts", nil)
	req.Header.Set("X-User-ID", "student-123")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		t.Errorf("expected 200 or 201, got %d", resp.StatusCode)
	}
}
