package usecase_test

import (
	"context"
	"errors"
	"testing"

	"aitu-superapp/community-service/internal/domain"
	"aitu-superapp/community-service/internal/repository"
	"aitu-superapp/community-service/internal/usecase"
)

// --- Mock PostgresRepository ---

type mockCommRepo struct {
	err error
}

func (m *mockCommRepo) SaveBooking(ctx context.Context, b *domain.CoworkingBooking) error {
	return m.err
}
func (m *mockCommRepo) GetEventsList(ctx context.Context) ([]domain.Event, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []domain.Event{{ID: "evt-1", Title: "Hackathon"}}, nil
}
func (m *mockCommRepo) CancelEvent(ctx context.Context, eventID string) error { return m.err }
func (m *mockCommRepo) GetCoworkingList(ctx context.Context) ([]string, error) {
	return nil, m.err
}
func (m *mockCommRepo) CancelBooking(ctx context.Context, bookingID string) error { return m.err }
func (m *mockCommRepo) GetMarketplaceItems(ctx context.Context) ([]domain.Ad, error) {
	return nil, m.err
}
func (m *mockCommRepo) CreateAd(ctx context.Context, ad *domain.Ad) (string, error) {
	return "ad-001", m.err
}
func (m *mockCommRepo) ReportLostFound(ctx context.Context, item *domain.LostFound) error {
	return m.err
}
func (m *mockCommRepo) GetLostFoundList(ctx context.Context) ([]domain.LostFound, error) {
	return nil, m.err
}
func (m *mockCommRepo) SubmitFeedback(ctx context.Context, f *domain.Feedback) error { return m.err }

var _ domain.PostgresRepository = (*mockCommRepo)(nil)

// --- Mock RedisRepository ---

type mockRedisRepo struct {
	locked bool
}

func (m *mockRedisRepo) AcquireLock(ctx context.Context, spaceID, timeSlot string) (bool, error) {
	return !m.locked, nil
}
func (m *mockRedisRepo) ReleaseLock(ctx context.Context, spaceID, timeSlot string) error {
	return nil
}

var _ repository.RedisRepository = (*mockRedisRepo)(nil)

// --- Tests ---

func TestBookSpace_Success(t *testing.T) {
	uc := usecase.NewCommunityUseCase(&mockCommRepo{}, &mockRedisRepo{locked: false})

	err := uc.BookSpace(context.Background(), "space-1", "student-1", "2026-05-18T10:00:00Z", "2026-05-18T12:00:00Z")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestBookSpace_AlreadyLocked(t *testing.T) {
	uc := usecase.NewCommunityUseCase(&mockCommRepo{}, &mockRedisRepo{locked: true})

	err := uc.BookSpace(context.Background(), "space-1", "student-1", "2026-05-18T10:00:00Z", "2026-05-18T12:00:00Z")
	if err == nil {
		t.Fatal("expected error when space is locked, got nil")
	}
}

func TestRegisterForEvent_Success(t *testing.T) {
	uc := usecase.NewCommunityUseCase(&mockCommRepo{}, &mockRedisRepo{})

	err := uc.RegisterForEvent(context.Background(), "event-1", "student-1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestRegisterForEvent_DBError(t *testing.T) {
	uc := usecase.NewCommunityUseCase(&mockCommRepo{err: errors.New("db down")}, &mockRedisRepo{})

	err := uc.RegisterForEvent(context.Background(), "event-1", "student-1")
	if err == nil {
		t.Fatal("expected error when db fails, got nil")
	}
}

func TestGetEventsList_Success(t *testing.T) {
	uc := usecase.NewCommunityUseCase(&mockCommRepo{}, &mockRedisRepo{})

	events, err := uc.GetEventsList(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(events) == 0 {
		t.Error("expected at least 1 event")
	}
}

func TestSubmitFeedback_Success(t *testing.T) {
	uc := usecase.NewCommunityUseCase(&mockCommRepo{}, &mockRedisRepo{})

	err := uc.SubmitFeedback(context.Background(), "student-1", "Great coworking!")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCreateMarketplaceAd_Success(t *testing.T) {
	uc := usecase.NewCommunityUseCase(&mockCommRepo{}, &mockRedisRepo{})

	id, err := uc.CreateMarketplaceAd(context.Background(), "student-1", "Calculus Book", "Used textbook", 3000)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if id == "" {
		t.Error("expected ad ID, got empty")
	}
}
