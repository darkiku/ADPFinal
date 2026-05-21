package grpc_delivery

import (
	"context"

	"aitu-superapp/admin-finance-service/internal/domain"
	pb "aitu-superapp/proto"
)

type AdminHandler struct {
	pb.UnimplementedAdminServiceServer
	useCase domain.AdminUseCase
}

func NewAdminHandler(uc domain.AdminUseCase) *AdminHandler {
	return &AdminHandler{useCase: uc}
}

func (h *AdminHandler) GetAccountBalance(ctx context.Context, req *pb.BalanceRequest) (*pb.BalanceResponse, error) {
	balance, err := h.useCase.GetAccountBalance(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.BalanceResponse{Balance: balance}, nil
}

func (h *AdminHandler) InitiatePayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	txID, err := h.useCase.InitiatePayment(ctx, req.UserId, req.Amount, req.Purpose)
	if err != nil {
		return &pb.PaymentResponse{Success: false}, err
	}
	return &pb.PaymentResponse{Success: true, TransactionId: txID}, nil
}

func (h *AdminHandler) RequestCertificate(ctx context.Context, req *pb.CertificateRequest) (*pb.CertificateResponse, error) {
	requestID, err := h.useCase.RequestCertificate(ctx, req.StudentId, req.DocType)
	if err != nil {
		return nil, err
	}
	return &pb.CertificateResponse{RequestId: requestID, Status: "pending"}, nil
}

func (h *AdminHandler) UpdateStudentProfile(ctx context.Context, req *pb.ProfileUpdateRequest) (*pb.ProfileUpdateResponse, error) {
	err := h.useCase.UpdateStudentProfile(ctx, req.StudentId, req.Email, req.Phone)
	if err != nil {
		return &pb.ProfileUpdateResponse{Success: false}, err
	}
	return &pb.ProfileUpdateResponse{Success: true}, nil
}

func (h *AdminHandler) GetFinanceHistory(ctx context.Context, req *pb.FinanceHistoryRequest) (*pb.FinanceHistoryResponse, error) {
	transactions, err := h.useCase.GetFinanceHistory(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var items []*pb.TransactionItem
	for _, t := range transactions {
		items = append(items, &pb.TransactionItem{
			Amount:  t.Amount,
			Purpose: t.Type,
			Date:    t.CreatedAt.Format("2006-01-02 15:04"),
			Type:    t.Type,
		})
	}
	if items == nil {
		items = []*pb.TransactionItem{}
	}
	return &pb.FinanceHistoryResponse{Transactions: items}, nil
}

func (h *AdminHandler) GetDormitoryInfo(ctx context.Context, req *pb.DormitoryRequest) (*pb.DormitoryResponse, error) {
	info, err := h.useCase.GetDormitoryInfo(ctx, req.StudentId)
	if err != nil {
		return nil, err
	}
	// студент не заселён — возвращаем пустой ответ без паники
	if info == nil {
		return &pb.DormitoryResponse{
			Room:     "Not assigned",
			Building: "Not assigned",
			Floor:    "-",
		}, nil
	}
	return &pb.DormitoryResponse{
		Room:     info.RoomNumber,
		Building: info.ID,
		Floor:    "-",
	}, nil
}

func (h *AdminHandler) GetScholarshipInfo(ctx context.Context, req *pb.ScholarshipRequest) (*pb.ScholarshipResponse, error) {
	info, err := h.useCase.GetScholarshipInfo(ctx, req.StudentId)
	if err != nil {
		return nil, err
	}
	scholarshipType, _ := info["type"].(string)
	if scholarshipType == "" {
		scholarshipType, _ = info["status"].(string)
	}
	amount, _ := info["amount"].(float64)
	currency, _ := info["currency"].(string)
	if currency == "" {
		currency = "KZT"
	}
	nextPayment, _ := info["next_payment"].(string)
	if nextPayment == "" {
		nextPayment = "2025-06-01"
	}
	return &pb.ScholarshipResponse{
		Type:        scholarshipType,
		Amount:      float32(amount),
		Currency:    currency,
		NextPayment: nextPayment,
	}, nil
}

func (h *AdminHandler) GetLibraryDebts(ctx context.Context, req *pb.LibraryDebtsRequest) (*pb.LibraryDebtsResponse, error) {
	books, err := h.useCase.GetLibraryDebts(ctx, req.StudentId)
	if err != nil {
		return nil, err
	}
	if books == nil {
		books = []string{}
	}
	return &pb.LibraryDebtsResponse{Books: books, Fine: 0}, nil
}

func (h *AdminHandler) GetDocumentStatus(ctx context.Context, req *pb.DocStatusRequest) (*pb.DocStatusResponse, error) {
	docs, err := h.useCase.GetDocumentStatus(ctx, req.StudentId)
	if err != nil {
		return nil, err
	}
	var items []*pb.DocItem
	for _, d := range docs {
		items = append(items, &pb.DocItem{
			Id:     d.ID,
			Type:   d.Type,
			Status: d.Status,
		})
	}
	if items == nil {
		items = []*pb.DocItem{}
	}
	return &pb.DocStatusResponse{Documents: items}, nil
}

func (h *AdminHandler) GetLaundryStatus(ctx context.Context, req *pb.LaundryRequest) (*pb.LaundryResponse, error) {
	return &pb.LaundryResponse{
		Available: 3,
		Total:     8,
		NextFree:  "14:30",
	}, nil
}
