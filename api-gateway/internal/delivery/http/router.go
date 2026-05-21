package http

import (
	"bytes"
	"io"
	"net/http"

	"aitu-superapp/api-gateway/internal/middleware"
	pb "aitu-superapp/proto"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, eduClient pb.EducationServiceClient, adminClient pb.AdminServiceClient, commClient pb.CommunityServiceClient) {

	// ========== AUTH (proxy to auth-service :4000) ==========
	authProxy := func(c *fiber.Ctx) error {
		url := "http://localhost:4000" + c.OriginalURL()
		req, _ := http.NewRequest(c.Method(), url, bytes.NewReader(c.Body()))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", c.Get("Authorization"))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return c.Status(502).JSON(fiber.Map{"error": "auth service unavailable"})
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		c.Set("Content-Type", "application/json")
		return c.Status(resp.StatusCode).Send(body)
	}

	app.Post("/auth/register", authProxy)
	app.Post("/auth/login", authProxy)
	app.Get("/auth/verify", authProxy)
	app.Get("/auth/profile", authProxy)

	// Все /api/v1 роуты защищены JWT middleware
	api := app.Group("/api/v1", middleware.JWTProtected())

	// ========== EDUCATION ==========
	edu := api.Group("/education")

	edu.Get("/gpa/:student_id", func(c *fiber.Ctx) error {
		resp, err := eduClient.GetStudentGPA(c.Context(), &pb.GPARequest{StudentId: c.Params("student_id")})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"student_id": c.Params("student_id"), "gpa": resp.Gpa})
	})

	edu.Get("/schedule/today", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		resp, err := eduClient.GetScheduleToday(c.Context(), &pb.ScheduleTodayRequest{StudentId: studentID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"schedule": resp.Items})
	})

	edu.Get("/grades/list", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		resp, err := eduClient.GetGrades(c.Context(), &pb.GetGradesRequest{StudentId: studentID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"grades": resp.Grades})
	})

	edu.Get("/attendance/history", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		resp, err := eduClient.GetAttendanceStats(c.Context(), &pb.AttendanceStatsRequest{StudentId: studentID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"percent": resp.Percent, "total": resp.Total, "present": resp.Present})
	})

	edu.Get("/courses/catalog", func(c *fiber.Ctx) error {
		resp, err := eduClient.GetCoursesCatalog(c.Context(), &pb.CoursesRequest{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"courses": resp.Courses})
	})

	edu.Post("/courses/enroll", func(c *fiber.Ctx) error {
		type Req struct {
			StudentID string `json:"student_id"`
			CourseID  string `json:"course_id"`
		}
		var body Req
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}
		resp, err := eduClient.EnrollToCourse(c.Context(), &pb.EnrollRequest{StudentId: body.StudentID, CourseId: body.CourseID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	edu.Delete("/courses/drop", func(c *fiber.Ctx) error {
		type Req struct {
			StudentID string `json:"student_id"`
			CourseID  string `json:"course_id"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := eduClient.DropCourse(c.Context(), &pb.DropCourseRequest{StudentId: body.StudentID, CourseId: body.CourseID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	edu.Get("/teachers/search", func(c *fiber.Ctx) error {
		query := c.Query("q", "")
		resp, err := eduClient.SearchTeacher(c.Context(), &pb.SearchTeacherRequest{Query: query})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"teachers": resp.Teachers})
	})

	edu.Get("/teachers/office-hours", func(c *fiber.Ctx) error {
		teacherID := c.Query("teacher_id", "")
		resp, err := eduClient.GetTeacherOfficeHours(c.Context(), &pb.OfficeHoursRequest{TeacherId: teacherID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	edu.Get("/exam/schedule", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		resp, err := eduClient.GetExamSchedule(c.Context(), &pb.ExamScheduleRequest{StudentId: studentID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"schedule": resp.ScheduleJson})
	})

	edu.Get("/profile", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		resp, err := eduClient.GetStudentProfile(c.Context(), &pb.StudentProfileRequest{StudentId: studentID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	edu.Post("/attendance/submit", func(c *fiber.Ctx) error {
		type Req struct {
			StudentID  string `json:"student_id"`
			ScheduleID string `json:"schedule_id"`
			Status     string `json:"status"`
			Date       string `json:"date"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := eduClient.SubmitAttendance(c.Context(), &pb.AttendanceRequest{
			StudentId: body.StudentID, ScheduleId: body.ScheduleID, Status: body.Status, Date: body.Date,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	// ========== FINANCE ==========
	admin := api.Group("/finance")

	admin.Get("/balance/:user_id", func(c *fiber.Ctx) error {
		resp, err := adminClient.GetAccountBalance(c.Context(), &pb.BalanceRequest{UserId: c.Params("user_id")})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"user_id": c.Params("user_id"), "balance": resp.Balance, "currency": resp.Currency})
	})

	admin.Post("/pay", func(c *fiber.Ctx) error {
		type Req struct {
			UserID  string  `json:"user_id"`
			Amount  float32 `json:"amount"`
			Purpose string  `json:"purpose"`
		}
		var body Req
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}
		resp, err := adminClient.InitiatePayment(c.Context(), &pb.PaymentRequest{
			UserId: body.UserID, Amount: body.Amount, Purpose: body.Purpose,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	admin.Get("/history/:user_id", func(c *fiber.Ctx) error {
		resp, err := adminClient.GetFinanceHistory(c.Context(), &pb.FinanceHistoryRequest{UserId: c.Params("user_id")})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"transactions": resp.Transactions})
	})

	admin.Get("/docs/status", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		resp, err := adminClient.GetDocumentStatus(c.Context(), &pb.DocStatusRequest{StudentId: studentID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"documents": resp.Documents})
	})

	admin.Post("/docs/request", func(c *fiber.Ctx) error {
		type Req struct {
			StudentID string `json:"student_id"`
			DocType   string `json:"doc_type"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := adminClient.RequestCertificate(c.Context(), &pb.CertificateRequest{
			StudentId: body.StudentID, DocType: body.DocType,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	admin.Get("/dormitory/info", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		resp, err := adminClient.GetDormitoryInfo(c.Context(), &pb.DormitoryRequest{StudentId: studentID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	admin.Post("/dormitory/repair", func(c *fiber.Ctx) error {
		type Req struct {
			StudentID   string `json:"student_id"`
			Description string `json:"description"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := adminClient.RequestDormRepair(c.Context(), &pb.RepairRequest{
			StudentId: body.StudentID, Description: body.Description,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	admin.Get("/scholarship/info", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		resp, err := adminClient.GetScholarshipInfo(c.Context(), &pb.ScholarshipRequest{StudentId: studentID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	admin.Get("/library/debts", func(c *fiber.Ctx) error {
		studentID := c.Get("X-User-ID")
		resp, err := adminClient.GetLibraryDebts(c.Context(), &pb.LibraryDebtsRequest{StudentId: studentID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"debts": resp.Books, "fine": resp.Fine})
	})

	admin.Post("/library/reserve", func(c *fiber.Ctx) error {
		type Req struct {
			StudentID string `json:"student_id"`
			BookID    string `json:"book_id"`
			Title     string `json:"title"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := adminClient.ReserveLibraryBook(c.Context(), &pb.ReserveBookRequest{
			StudentId: body.StudentID, BookId: body.BookID, Title: body.Title,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	admin.Get("/laundry/status", func(c *fiber.Ctx) error {
		resp, err := adminClient.GetLaundryStatus(c.Context(), &pb.LaundryRequest{Building: "A"})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	admin.Put("/profile", func(c *fiber.Ctx) error {
		type Req struct {
			StudentID string `json:"student_id"`
			Email     string `json:"email"`
			Phone     string `json:"phone"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := adminClient.UpdateStudentProfile(c.Context(), &pb.ProfileUpdateRequest{
			StudentId: body.StudentID, Email: body.Email, Phone: body.Phone,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	// ========== COMMUNITY ==========
	comm := api.Group("/community")

	comm.Post("/coworking/book", func(c *fiber.Ctx) error {
		type Req struct {
			SpaceID   string `json:"space_id"`
			StudentID string `json:"student_id"`
			StartTime string `json:"start_time"`
			EndTime   string `json:"end_time"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := commClient.BookCoworkingSpace(c.Context(), &pb.BookSpaceRequest{
			SpaceId: body.SpaceID, StudentId: body.StudentID, StartTime: body.StartTime, EndTime: body.EndTime,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	comm.Get("/coworking/list", func(c *fiber.Ctx) error {
		resp, err := commClient.GetCoworkingSpaces(c.Context(), &pb.CoworkingSpacesRequest{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"spaces": resp.Spaces})
	})

	comm.Delete("/coworking/cancel", func(c *fiber.Ctx) error {
		type Req struct {
			BookingID string `json:"booking_id"`
			StudentID string `json:"student_id"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := commClient.CancelBooking(c.Context(), &pb.CancelBookingRequest{
			BookingId: body.BookingID, StudentId: body.StudentID,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	comm.Get("/events/list", func(c *fiber.Ctx) error {
		resp, err := commClient.GetEventsList(c.Context(), &pb.EventsListRequest{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"events": resp.Events})
	})

	comm.Post("/events/register", func(c *fiber.Ctx) error {
		type Req struct {
			EventID   string `json:"event_id"`
			StudentID string `json:"student_id"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := commClient.RegisterForEvent(c.Context(), &pb.EventRegistrationRequest{
			EventId: body.EventID, StudentId: body.StudentID,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	comm.Delete("/events/cancel", func(c *fiber.Ctx) error {
		type Req struct {
			EventID   string `json:"event_id"`
			StudentID string `json:"student_id"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := commClient.CancelEventRegistration(c.Context(), &pb.CancelEventRequest{
			EventId: body.EventID, StudentId: body.StudentID,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	comm.Get("/marketplace/items", func(c *fiber.Ctx) error {
		resp, err := commClient.GetMarketplaceItems(c.Context(), &pb.MarketplaceRequest{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"items": resp.Items})
	})

	comm.Post("/marketplace/create", func(c *fiber.Ctx) error {
		type Req struct {
			StudentID   string  `json:"student_id"`
			Title       string  `json:"title"`
			Description string  `json:"description"`
			Price       float32 `json:"price"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := commClient.CreateMarketplaceAd(c.Context(), &pb.AdRequest{
			StudentId: body.StudentID, Title: body.Title, Description: body.Description, Price: body.Price,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	comm.Post("/lost-found/report", func(c *fiber.Ctx) error {
		type Req struct {
			StudentID   string `json:"student_id"`
			Item        string `json:"item"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := commClient.ReportLostFound(c.Context(), &pb.LostFoundRequest{
			StudentId: body.StudentID, Item: body.Item, Description: body.Description, Status: body.Status,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})

	comm.Get("/lost-found/list", func(c *fiber.Ctx) error {
		resp, err := commClient.GetLostFoundList(c.Context(), &pb.GetLostFoundRequest{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"items": resp.Items})
	})

	comm.Get("/news/feed", func(c *fiber.Ctx) error {
		resp, err := commClient.GetNewsFeed(c.Context(), &pb.NewsFeedRequest{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"news": resp.News})
	})

	comm.Post("/feedback", func(c *fiber.Ctx) error {
		type Req struct {
			StudentID string `json:"student_id"`
			Text      string `json:"text"`
			Category  string `json:"category"`
		}
		var body Req
		c.BodyParser(&body)
		resp, err := commClient.SubmitFeedback(c.Context(), &pb.FeedbackRequest{
			StudentId: body.StudentID, Text: body.Text, Category: body.Category,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	})
}
