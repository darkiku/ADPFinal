package usecase

import (
	"context"
	"errors"

	"aitu-superapp/admin-finance-service/internal/domain"
)

type adminUseCase struct {
	repo      domain.AdminRepository
	publisher domain.AdminEventPublisher
}

func NewAdminUseCase(repo domain.AdminRepository, pub domain.AdminEventPublisher) domain.AdminUseCase {
	return &adminUseCase{repo: repo, publisher: pub}
}

func (uc *adminUseCase) GetAccountBalance(ctx context.Context, userID string) (float32, error) {
	return uc.repo.GetBalance(ctx, userID)
}

func (uc *adminUseCase) InitiatePayment(ctx context.Context, userID string, amount float32, purpose string) (string, error) {
	if amount <= 0 {
		return "", errors.New("amount must be > 0")
	}

	// Атомарная транзакция БД: BEGIN/COMMIT/ROLLBACK внутри репозитория
	txID, err := uc.repo.InitiatePaymentTx(ctx, userID, amount)
	if err != nil {
		return "", err
	}

	if uc.publisher != nil {
		_ = uc.publisher.PublishPaymentSuccess(userID, amount)
	}

	return txID, nil
}

func (uc *adminUseCase) RequestCertificate(ctx context.Context, studentID, docType string) (string, error) {
	return uc.repo.CreateDocumentRequest(ctx, studentID, docType)
}

func (uc *adminUseCase) UpdateStudentProfile(ctx context.Context, studentID, email, phone string) error {
	return uc.repo.UpdateProfile(ctx, studentID, email, phone)
}

func (uc *adminUseCase) GetDocumentStatus(ctx context.Context, studentID string) ([]domain.DocumentRequest, error) {
	return uc.repo.GetDocumentStatus(ctx, studentID)
}
func (uc *adminUseCase) GetDormitoryInfo(ctx context.Context, studentID string) (*domain.Dormitory, error) {
	return uc.repo.GetDormitoryInfo(ctx, studentID)
}
func (uc *adminUseCase) RequestDormRepair(ctx context.Context, studentID, desc string) error {
	return uc.repo.RequestDormRepair(ctx, studentID, desc)
}
func (uc *adminUseCase) GetFinanceHistory(ctx context.Context, userID string) ([]domain.Transaction, error) {
	return uc.repo.GetFinanceHistory(ctx, userID)
}
func (uc *adminUseCase) GetScholarshipInfo(ctx context.Context, studentID string) (map[string]interface{}, error) {
	return uc.repo.GetScholarshipInfo(ctx, studentID)
}
func (uc *adminUseCase) GetLibraryDebts(ctx context.Context, studentID string) ([]string, error) {
	return uc.repo.GetLibraryDebts(ctx, studentID)
}
func (uc *adminUseCase) ReserveLibraryBook(ctx context.Context, studentID, bookID string) error {
	return uc.repo.ReserveLibraryBook(ctx, studentID, bookID)
}
