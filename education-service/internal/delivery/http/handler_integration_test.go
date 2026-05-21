package http_delivery_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	http_delivery "aitu-superapp/education-service/internal/delivery/http"
	"aitu-superapp/education-service/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type mockEduUseCase struct{}

func (m *mockEduUseCase) CalculateGPA(ctx context.Context, studentID string) (float32, error) {
	return 3.75, nil
}
func (m *mockEduUseCase) SubmitAttendance(ctx context.Context, studentID, scheduleID, status, date string) error {
	return nil
}
func (m *mockEduUseCase) EnrollToCourse(ctx context.Context, studentID, courseID string) error {
	return nil
}
func (m *mockEduUseCase) GetExamScheduleJSON(ctx context.Context, studentID string) (string, error) {
	return `[]`, nil
}
func (m *mockEduUseCase) GetTodaySchedule(ctx context.Context, studentID string) ([]domain.Schedule, error) {
	return []domain.Schedule{{ID: "s-1", SubjectID: "sub-1", Room: "301", StartTime: time.Now()}}, nil
}
func (m *mockEduUseCase) GetWeekSchedule(ctx context.Context, studentID string) ([]domain.Schedule, error) {
	return []domain.Schedule{}, nil
}
func (m *mockEduUseCase) GetGradesList(ctx context.Context, studentID string) ([]domain.Grade, error) {
	return []domain.Grade{{SubjectID: "sub-1", Score: 90, Type: "midterm"}}, nil
}
func (m *mockEduUseCase) GetAttendanceHistory(ctx context.Context, studentID string) ([]domain.Attendance, error) {
	return []domain.Attendance{}, nil
}
func (m *mockEduUseCase) GetCoursesCatalog(ctx context.Context) ([]domain.Subject, error) {
	return []domain.Subject{{ID: "1", Title: "Go Programming", Credits: 5}}, nil
}
func (m *mockEduUseCase) DropCourse(ctx context.Context, studentID, courseID string) error {
	return nil
}
func (m *mockEduUseCase) SearchTeachers(ctx context.Context, query string) ([]string, error) {
	return []string{"Dr. Smith"}, nil
}
func (m *mockEduUseCase) GetTeacherOfficeHours(ctx context.Context, teacherID string) ([]string, error) {
	return []string{"Mon 14:00-16:00"}, nil
}

func setupEduApp() *fiber.App {
	app := fiber.New()
	http_delivery.RegisterRESTHandlers(app, &mockEduUseCase{})
	return app
}

func TestIntegration_GetScheduleToday(t *testing.T) {
	app := setupEduApp()
	req := httptest.NewRequest("GET", "/api/v1/education/schedule/today", nil)
	req.Header.Set("X-User-ID", "student-123")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIntegration_GetGradesList(t *testing.T) {
	app := setupEduApp()
	req := httptest.NewRequest("GET", "/api/v1/education/grades/list", nil)
	req.Header.Set("X-User-ID", "student-123")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIntegration_GetCoursesCatalog(t *testing.T) {
	app := setupEduApp()
	req := httptest.NewRequest("GET", "/api/v1/education/courses/catalog", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIntegration_GetAttendanceHistory(t *testing.T) {
	app := setupEduApp()
	req := httptest.NewRequest("GET", "/api/v1/education/attendance/history", nil)
	req.Header.Set("X-User-ID", "student-123")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
