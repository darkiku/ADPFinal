package domain

import (
	"context"
	"time"
)

type Event struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
}

type CoworkingBooking struct {
	ID        string    `json:"id"`
	SpaceID   string    `json:"space_id"`
	StudentID string    `json:"student_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type Ad struct {
	ID          string  `json:"id"`
	StudentID   string  `json:"student_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type LostFound struct {
	ID          string `json:"id"`
	StudentID   string `json:"student_id"`
	Item        string `json:"item"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type Feedback struct {
	ID        string `json:"id"`
	StudentID string `json:"student_id"`
	Text      string `json:"text"`
	Category  string `json:"category"`
}

type NewsItem struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// Контракт для Postgres
type PostgresRepository interface {
	SaveBooking(ctx context.Context, booking *CoworkingBooking) error
	GetEventsList(ctx context.Context) ([]Event, error)
	RegisterForEvent(ctx context.Context, eventID, studentID string) error
	CancelEventRegistration(ctx context.Context, eventID, studentID string) error
	CancelEvent(ctx context.Context, eventID string) error
	GetCoworkingList(ctx context.Context) ([]string, error)
	CancelBooking(ctx context.Context, bookingID string) error
	GetMarketplaceItems(ctx context.Context) ([]Ad, error)
	CreateAd(ctx context.Context, ad *Ad) (string, error)
	ReportLostFound(ctx context.Context, item *LostFound) error
	GetLostFoundList(ctx context.Context) ([]LostFound, error)
	SubmitFeedback(ctx context.Context, feedback *Feedback) error
	GetNewsFeed(ctx context.Context) ([]NewsItem, error)
}

// Контракт для UseCase
type CommunityUseCase interface {
	// gRPC
	BookSpace(ctx context.Context, spaceID, studentID, startTimeStr, endTimeStr string) error
	RegisterForEvent(ctx context.Context, eventID, studentID string) error
	CreateMarketplaceAd(ctx context.Context, studentID, title, description string, price float32) (string, error)
	SubmitFeedback(ctx context.Context, studentID, text, category string) error

	// REST
	GetEventsList(ctx context.Context) ([]Event, error)
	CancelEvent(ctx context.Context, eventID string) error
	CancelEventRegistration(ctx context.Context, eventID, studentID string) error
	GetCoworkingList(ctx context.Context) ([]string, error)
	CancelBooking(ctx context.Context, bookingID string) error
	GetMarketplaceItems(ctx context.Context) ([]Ad, error)
	ReportLostFound(ctx context.Context, studentID, item, description, status string) error
	GetLostFoundList(ctx context.Context) ([]LostFound, error)
	GetNewsFeed(ctx context.Context) ([]NewsItem, error)
}
