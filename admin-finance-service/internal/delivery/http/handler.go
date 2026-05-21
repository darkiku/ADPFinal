package http_delivery

import (
	"aitu-superapp/admin-finance-service/internal/domain"
	"github.com/gofiber/fiber/v2"
)

func RegisterRESTHandlers(app *fiber.App, uc domain.AdminUseCase) {
	admin := app.Group("/api/v1/admin")

	// 1. GET /docs/status
	admin.Get("/docs/status", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		docs, err := uc.GetDocumentStatus(c.Context(), studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(docs)
	})

	// 2. GET /dormitory/info
	admin.Get("/dormitory/info", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		info, err := uc.GetDormitoryInfo(c.Context(), studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(info)
	})

	// 3. POST /dormitory/repair
	admin.Post("/dormitory/repair", func(c *fiber.Ctx) error {
		type Req struct {
			Description string `json:"description"`
		}
		var body Req
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "bad request"})
		}

		err := uc.RequestDormRepair(c.Context(), c.Get("X-User-ID"), body.Description)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(201)
	})

	// 4. GET /laundry/status
	admin.Get("/laundry/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"machine_1": "available", "machine_2": "in_use"}) // Заглушка
	})

	// 5. GET /finance/history
	admin.Get("/finance/history", func(c *fiber.Ctx) error {
		history, err := uc.GetFinanceHistory(c.Context(), c.Get("X-User-ID"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(history)
	})

	// 6. GET /scholarship/info
	admin.Get("/scholarship/info", func(c *fiber.Ctx) error {
		info, err := uc.GetScholarshipInfo(c.Context(), c.Get("X-User-ID"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(info)
	})

	// 7. GET /library/debts
	admin.Get("/library/debts", func(c *fiber.Ctx) error {
		debts, err := uc.GetLibraryDebts(c.Context(), c.Get("X-User-ID"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(debts)
	})

	// 8. POST /library/reserve
	admin.Post("/library/reserve", func(c *fiber.Ctx) error {
		type Req struct {
			BookID string `json:"book_id"`
		}
		var body Req
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "bad request"})
		}

		err := uc.ReserveLibraryBook(c.Context(), c.Get("X-User-ID"), body.BookID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(201)
	})
}
