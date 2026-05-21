package domain

import (
	"context"
	"time"
)

type DocumentRequest struct {
	ID        string `json:"id"`
	StudentID string `json:"student_id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

type FinanceAccount struct {
	ID      string  `json:"id"`
	UserID  string  `json:"user_id"`
	Balance float32 `json:"balance"`
}

type Transaction struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Amount    float32   `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type Dormitory struct {
	ID         string `json:"id"`
	RoomNumber string `json:"room_number"`
	OccupantID string `json:"occupant_id"`
}

type AdminRepository interface {
	GetBalance(ctx context.Context, userID string) (float32, error)
	UpdateBalance(ctx context.Context, userID string, amount float32) error
	CreateTransaction(ctx context.Context, userID string, amount float32, transType string) (string, error)
	InitiatePaymentTx(ctx context.Context, userID string, amount float32) (string, error)
	CreateDocumentRequest(ctx context.Context, studentID, docType string) (string, error)
	UpdateProfile(ctx context.Context, studentID, email, phone string) error
	GetDocumentStatus(ctx context.Context, studentID string) ([]DocumentRequest, error)
	GetDormitoryInfo(ctx context.Context, studentID string) (*Dormitory, error)
	RequestDormRepair(ctx context.Context, studentID, description string) error
	GetFinanceHistory(ctx context.Context, userID string) ([]Transaction, error)
	GetScholarshipInfo(ctx context.Context, studentID string) (map[string]interface{}, error)
	GetLibraryDebts(ctx context.Context, studentID string) ([]string, error)
	ReserveLibraryBook(ctx context.Context, studentID, bookID string) error
}

type AdminUseCase interface {
	GetAccountBalance(ctx context.Context, userID string) (float32, error)
	InitiatePayment(ctx context.Context, userID string, amount float32, purpose string) (string, error)
	RequestCertificate(ctx context.Context, studentID, docType string) (string, error)
	UpdateStudentProfile(ctx context.Context, studentID, email, phone string) error
	GetDocumentStatus(ctx context.Context, studentID string) ([]DocumentRequest, error)
	GetDormitoryInfo(ctx context.Context, studentID string) (*Dormitory, error)
	RequestDormRepair(ctx context.Context, studentID, description string) error
	GetFinanceHistory(ctx context.Context, userID string) ([]Transaction, error)
	GetScholarshipInfo(ctx context.Context, studentID string) (map[string]interface{}, error)
	GetLibraryDebts(ctx context.Context, studentID string) ([]string, error)
	ReserveLibraryBook(ctx context.Context, studentID, bookID string) error
}

type AdminEventPublisher interface {
	PublishPaymentSuccess(userID string, amount float32) error
}
