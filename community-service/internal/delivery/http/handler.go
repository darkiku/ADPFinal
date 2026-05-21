package http

import (
	"aitu-superapp/community-service/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type CommunityHandler struct {
	useCase domain.CommunityUseCase
}

func NewCommunityHandler(uc domain.CommunityUseCase) *CommunityHandler {
	return &CommunityHandler{useCase: uc}
}

func (h *CommunityHandler) RegisterRoutes(app *fiber.App) {
	comm := app.Group("/api/v1/community")

	comm.Get("/events/list", h.GetEventsList)
	comm.Delete("/events/cancel/:id", h.CancelEvent)
	comm.Get("/coworking/list", h.GetCoworkingList)
	comm.Delete("/coworking/cancel/:id", h.CancelBooking)
	comm.Get("/marketplace/items", h.GetMarketplaceItems)
	comm.Post("/lost-found/report", h.ReportLostFound)
	comm.Get("/lost-found/list", h.GetLostFoundList)
	comm.Get("/news/feed", h.GetNewsFeed)
}

func (h *CommunityHandler) GetEventsList(c *fiber.Ctx) error {
	events, err := h.useCase.GetEventsList(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(events)
}

func (h *CommunityHandler) CancelEvent(c *fiber.Ctx) error {
	err := h.useCase.CancelEvent(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

func (h *CommunityHandler) GetCoworkingList(c *fiber.Ctx) error {
	list, err := h.useCase.GetCoworkingList(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

func (h *CommunityHandler) CancelBooking(c *fiber.Ctx) error {
	err := h.useCase.CancelBooking(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

func (h *CommunityHandler) GetMarketplaceItems(c *fiber.Ctx) error {
	items, err := h.useCase.GetMarketplaceItems(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func (h *CommunityHandler) ReportLostFound(c *fiber.Ctx) error {
	type Req struct {
		Item        string `json:"item"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad request"})
	}

	err := h.useCase.ReportLostFound(c.Context(), c.Get("X-User-ID"), body.Item, body.Description, body.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(201)
}

func (h *CommunityHandler) GetLostFoundList(c *fiber.Ctx) error {
	items, err := h.useCase.GetLostFoundList(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func (h *CommunityHandler) GetNewsFeed(c *fiber.Ctx) error {
	news, err := h.useCase.GetNewsFeed(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(news)
}
