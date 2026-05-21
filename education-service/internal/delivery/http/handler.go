package http_delivery

import (
	"aitu-superapp/education-service/internal/domain"
	"github.com/gofiber/fiber/v2"
)

func RegisterRESTHandlers(app *fiber.App, uc domain.EducationUseCase) {
	edu := app.Group("/api/v1/education")

	// 1. GET /schedule/today
	edu.Get("/schedule/today", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID") // Достаем из хедера, прокинутого Gateway
		schedule, err := uc.GetTodaySchedule(c.Context(), studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(schedule)
	})

	// 2. GET /schedule/week
	edu.Get("/schedule/week", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		schedule, err := uc.GetWeekSchedule(c.Context(), studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(schedule)
	})

	// 3. GET /grades/list
	edu.Get("/grades/list", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		grades, err := uc.GetGradesList(c.Context(), studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(grades)
	})

	// 4. GET /attendance/history
	edu.Get("/attendance/history", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		atts, err := uc.GetAttendanceHistory(c.Context(), studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(atts)
	})

	// 5. GET /courses/catalog
	edu.Get("/courses/catalog", func(c *fiber.Ctx) error {
		catalog, err := uc.GetCoursesCatalog(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(catalog)
	})

	// 6. DELETE /courses/drop
	edu.Delete("/courses/drop/:course_id", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		courseID := c.Params("course_id")
		err := uc.DropCourse(c.Context(), studentID, courseID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(204)
	})

	// 7. GET /teachers/search
	edu.Get("/teachers/search", func(c *fiber.Ctx) error {
		query := c.Query("q")
		teachers, err := uc.SearchTeachers(c.Context(), query)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(teachers)
	})

	// 8. GET /teachers/office-hours
	edu.Get("/teachers/office-hours/:teacher_id", func(c *fiber.Ctx) error {
		teacherID := c.Params("teacher_id")
		hours, err := uc.GetTeacherOfficeHours(c.Context(), teacherID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(hours)
	})
}
