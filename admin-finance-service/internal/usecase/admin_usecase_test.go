package usecase_test

import (
	"context"
	"errors"
	"testing"

	"aitu-superapp/admin-finance-service/internal/domain"
	"aitu-superapp/admin-finance-service/internal/usecase"
)

// --- Моки ---

type mockAdminRepo struct {
	balance float32
	txID    string
	err     error
}

func (m *mockAdminRepo) GetBalance(ctx context.Context, userID string) (float32, error) {
	return m.balance, m.err
}
func (m *mockAdminRepo) UpdateBalance(ctx context.Context, userID string, amount float32) error {
	return m.err
}
func (m *mockAdminRepo) CreateTransaction(ctx context.Context, userID string, amount float32, t string) (string, error) {
	return m.txID, m.err
}
func (m *mockAdminRepo) InitiatePaymentTx(ctx context.Context, userID string, amount float32) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return "tx-uuid-001", nil
}
func (m *mockAdminRepo) CreateDocumentRequest(ctx context.Context, studentID, docType string) (string, error) {
	return "doc-001", m.err
}
func (m *mockAdminRepo) UpdateProfile(ctx context.Context, studentID, email, phone string) error {
	return m.err
}
func (m *mockAdminRepo) GetDocumentStatus(ctx context.Context, studentID string) ([]domain.DocumentRequest, error) {
	return nil, m.err
}
func (m *mockAdminRepo) GetDormitoryInfo(ctx context.Context, studentID string) (*domain.Dormitory, error) {
	return nil, m.err
}
func (m *mockAdminRepo) RequestDormRepair(ctx context.Context, studentID, desc string) error {
	return m.err
}
func (m *mockAdminRepo) GetFinanceHistory(ctx context.Context, userID string) ([]domain.Transaction, error) {
	return nil, m.err
}
func (m *mockAdminRepo) GetScholarshipInfo(ctx context.Context, studentID string) (map[string]interface{}, error) {
	return nil, m.err
}
func (m *mockAdminRepo) GetLibraryDebts(ctx context.Context, studentID string) ([]string, error) {
	return nil, m.err
}
func (m *mockAdminRepo) ReserveLibraryBook(ctx context.Context, studentID, bookID string) error {
	return m.err
}

type mockAdminPublisher struct{}

func (m *mockAdminPublisher) PublishPaymentSuccess(userID string, amount float32) error {
	return nil
}

// --- Тесты ---

func TestGetAccountBalance_Success(t *testing.T) {
	repo := &mockAdminRepo{balance: 75000}
	uc := usecase.NewAdminUseCase(repo, &mockAdminPublisher{})

	balance, err := uc.GetAccountBalance(context.Background(), "user-123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if balance != 75000 {
		t.Errorf("expected 75000, got: %v", balance)
	}
}

func TestInitiatePayment_Success(t *testing.T) {
	repo := &mockAdminRepo{}
	uc := usecase.NewAdminUseCase(repo, &mockAdminPublisher{})

	txID, err := uc.InitiatePayment(context.Background(), "user-123", 1500, "Coworking")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if txID == "" {
		t.Error("expected transaction ID, got empty string")
	}
}

func TestInitiatePayment_NegativeAmount(t *testing.T) {
	repo := &mockAdminRepo{}
	uc := usecase.NewAdminUseCase(repo, &mockAdminPublisher{})

	_, err := uc.InitiatePayment(context.Background(), "user-123", -100, "Test")
	if err == nil {
		t.Fatal("expected error for negative amount, got nil")
	}
}

func TestInitiatePayment_DBError(t *testing.T) {
	repo := &mockAdminRepo{err: errors.New("db error")}
	uc := usecase.NewAdminUseCase(repo, &mockAdminPublisher{})

	_, err := uc.InitiatePayment(context.Background(), "user-123", 500, "Test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRequestCertificate_Success(t *testing.T) {
	repo := &mockAdminRepo{}
	uc := usecase.NewAdminUseCase(repo, &mockAdminPublisher{})

	id, err := uc.RequestCertificate(context.Background(), "student-123", "transcript")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if id == "" {
		t.Error("expected document ID, got empty")
	}
}
