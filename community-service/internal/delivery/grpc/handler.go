package grpc

import (
	"context"

	"aitu-superapp/community-service/internal/domain"
	pb "aitu-superapp/proto"
)

type CommunityHandler struct {
	pb.UnimplementedCommunityServiceServer
	useCase domain.CommunityUseCase
}

func NewCommunityHandler(uc domain.CommunityUseCase) *CommunityHandler {
	return &CommunityHandler{useCase: uc}
}

func (h *CommunityHandler) BookCoworkingSpace(ctx context.Context, req *pb.BookSpaceRequest) (*pb.BookSpaceResponse, error) {
	err := h.useCase.BookSpace(ctx, req.SpaceId, req.StudentId, req.StartTime, req.EndTime)
	if err != nil {
		return &pb.BookSpaceResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.BookSpaceResponse{Success: true, Message: "Successfully booked"}, nil
}

func (h *CommunityHandler) RegisterForEvent(ctx context.Context, req *pb.EventRegistrationRequest) (*pb.EventRegistrationResponse, error) {
	err := h.useCase.RegisterForEvent(ctx, req.EventId, req.StudentId)
	if err != nil {
		return &pb.EventRegistrationResponse{Success: false}, err
	}
	return &pb.EventRegistrationResponse{Success: true}, nil
}

func (h *CommunityHandler) CancelEventRegistration(ctx context.Context, req *pb.CancelEventRequest) (*pb.CancelEventResponse, error) {
	err := h.useCase.CancelEventRegistration(ctx, req.EventId, req.StudentId)
	if err != nil {
		return &pb.CancelEventResponse{Success: false}, err
	}
	return &pb.CancelEventResponse{Success: true}, nil
}

func (h *CommunityHandler) CreateMarketplaceAd(ctx context.Context, req *pb.AdRequest) (*pb.AdResponse, error) {
	id, err := h.useCase.CreateMarketplaceAd(ctx, req.StudentId, req.Title, req.Description, req.Price)
	if err != nil {
		return &pb.AdResponse{Success: false}, err
	}
	return &pb.AdResponse{AdId: id, Success: true}, nil
}

func (h *CommunityHandler) SubmitFeedback(ctx context.Context, req *pb.FeedbackRequest) (*pb.FeedbackResponse, error) {
	err := h.useCase.SubmitFeedback(ctx, req.StudentId, req.Text, req.Category)
	if err != nil {
		return &pb.FeedbackResponse{Success: false}, err
	}
	return &pb.FeedbackResponse{Success: true}, nil
}

func (h *CommunityHandler) GetEventsList(ctx context.Context, req *pb.EventsListRequest) (*pb.EventsListResponse, error) {
	events, err := h.useCase.GetEventsList(ctx)
	if err != nil {
		return nil, err
	}
	var items []*pb.EventItem
	for _, e := range events {
		items = append(items, &pb.EventItem{
			Id:       e.ID,
			Title:    e.Title,
			Date:     e.Date.Format("2006-01-02 15:04"),
			Location: e.Location,
		})
	}
	return &pb.EventsListResponse{Events: items}, nil
}

func (h *CommunityHandler) GetCoworkingSpaces(ctx context.Context, req *pb.CoworkingSpacesRequest) (*pb.CoworkingSpacesResponse, error) {
	spaces, err := h.useCase.GetCoworkingList(ctx)
	if err != nil {
		return nil, err
	}
	var items []*pb.SpaceItem
	for _, s := range spaces {
		items = append(items, &pb.SpaceItem{
			Name:      s,
			Available: true,
			Capacity:  10,
		})
	}
	return &pb.CoworkingSpacesResponse{Spaces: items}, nil
}

func (h *CommunityHandler) CancelBooking(ctx context.Context, req *pb.CancelBookingRequest) (*pb.CancelBookingResponse, error) {
	err := h.useCase.CancelBooking(ctx, req.BookingId)
	if err != nil {
		return &pb.CancelBookingResponse{Success: false}, err
	}
	return &pb.CancelBookingResponse{Success: true}, nil
}

func (h *CommunityHandler) GetMarketplaceItems(ctx context.Context, req *pb.MarketplaceRequest) (*pb.MarketplaceResponse, error) {
	items, err := h.useCase.GetMarketplaceItems(ctx)
	if err != nil {
		return nil, err
	}
	var out []*pb.MarketplaceItem
	for _, a := range items {
		out = append(out, &pb.MarketplaceItem{
			Id:          a.ID,
			Title:       a.Title,
			Description: a.Description,
			Price:       a.Price,
			Seller:      a.StudentID,
		})
	}
	return &pb.MarketplaceResponse{Items: out}, nil
}

func (h *CommunityHandler) ReportLostFound(ctx context.Context, req *pb.LostFoundRequest) (*pb.LostFoundResponse, error) {
	err := h.useCase.ReportLostFound(ctx, req.StudentId, req.Item, req.Description, req.Status)
	if err != nil {
		return &pb.LostFoundResponse{Success: false}, err
	}
	return &pb.LostFoundResponse{Success: true}, nil
}

func (h *CommunityHandler) GetLostFoundList(ctx context.Context, req *pb.GetLostFoundRequest) (*pb.GetLostFoundResponse, error) {
	items, err := h.useCase.GetLostFoundList(ctx)
	if err != nil {
		return nil, err
	}
	var out []*pb.LostFoundItem
	for _, lf := range items {
		out = append(out, &pb.LostFoundItem{Id: lf.ID, Item: lf.Item, Description: lf.Description, Status: lf.Status})
	}
	return &pb.GetLostFoundResponse{Items: out}, nil
}

func (h *CommunityHandler) GetNewsFeed(ctx context.Context, req *pb.NewsFeedRequest) (*pb.NewsFeedResponse, error) {
	news, err := h.useCase.GetNewsFeed(ctx)
	if err != nil {
		return nil, err
	}
	var items []*pb.NewsItem
	for _, n := range news {
		items = append(items, &pb.NewsItem{
			Id:       n.ID,
			Title:    n.Title,
			Content:  n.Body,
			Date:     n.CreatedAt.Format("2006-01-02"),
			Category: "news",
		})
	}
	return &pb.NewsFeedResponse{News: items}, nil
}
