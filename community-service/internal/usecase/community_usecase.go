package usecase

import (
	"context"
	"errors"
	"time"

	"aitu-superapp/community-service/internal/domain"
	"aitu-superapp/community-service/internal/repository"
)

type communityUseCase struct {
	pgRepo    domain.PostgresRepository
	redisRepo repository.RedisRepository
}

func NewCommunityUseCase(pgRepo domain.PostgresRepository, redisRepo repository.RedisRepository) domain.CommunityUseCase {
	return &communityUseCase{pgRepo: pgRepo, redisRepo: redisRepo}
}

func (uc *communityUseCase) BookSpace(ctx context.Context, spaceID, studentID, startTimeStr, endTimeStr string) error {
	timeSlot := startTimeStr + "-" + endTimeStr
	locked, err := uc.redisRepo.AcquireLock(ctx, spaceID, timeSlot)
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("space is currently being booked by someone else")
	}
	defer func() {
		_ = uc.redisRepo.ReleaseLock(ctx, spaceID, timeSlot)
	}()

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return err
	}
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return err
	}

	booking := &domain.CoworkingBooking{
		SpaceID: spaceID, StudentID: studentID, StartTime: startTime, EndTime: endTime,
	}
	return uc.pgRepo.SaveBooking(ctx, booking)
}

func (uc *communityUseCase) RegisterForEvent(ctx context.Context, eventID, studentID string) error {
	return uc.pgRepo.RegisterForEvent(ctx, eventID, studentID)
}

func (uc *communityUseCase) CancelEventRegistration(ctx context.Context, eventID, studentID string) error {
	return uc.pgRepo.CancelEventRegistration(ctx, eventID, studentID)
}

func (uc *communityUseCase) CreateMarketplaceAd(ctx context.Context, studentID, title, desc string, price float32) (string, error) {
	ad := &domain.Ad{StudentID: studentID, Title: title, Description: desc, Price: price}
	return uc.pgRepo.CreateAd(ctx, ad)
}

func (uc *communityUseCase) SubmitFeedback(ctx context.Context, studentID, text, category string) error {
	return uc.pgRepo.SubmitFeedback(ctx, &domain.Feedback{StudentID: studentID, Text: text, Category: category})
}

func (uc *communityUseCase) GetEventsList(ctx context.Context) ([]domain.Event, error) {
	return uc.pgRepo.GetEventsList(ctx)
}

func (uc *communityUseCase) CancelEvent(ctx context.Context, eventID string) error {
	return uc.pgRepo.CancelEvent(ctx, eventID)
}

func (uc *communityUseCase) GetCoworkingList(ctx context.Context) ([]string, error) {
	return uc.pgRepo.GetCoworkingList(ctx)
}

func (uc *communityUseCase) CancelBooking(ctx context.Context, bookingID string) error {
	return uc.pgRepo.CancelBooking(ctx, bookingID)
}

func (uc *communityUseCase) GetMarketplaceItems(ctx context.Context) ([]domain.Ad, error) {
	return uc.pgRepo.GetMarketplaceItems(ctx)
}

func (uc *communityUseCase) ReportLostFound(ctx context.Context, studentID, item, desc, status string) error {
	return uc.pgRepo.ReportLostFound(ctx, &domain.LostFound{StudentID: studentID, Item: item, Description: desc, Status: status})
}

func (uc *communityUseCase) GetLostFoundList(ctx context.Context) ([]domain.LostFound, error) {
	return uc.pgRepo.GetLostFoundList(ctx)
}

func (uc *communityUseCase) GetNewsFeed(ctx context.Context) ([]domain.NewsItem, error) {
	return uc.pgRepo.GetNewsFeed(ctx)
}
