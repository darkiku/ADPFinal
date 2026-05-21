package usecase_test

import (
	"context"
	"errors"
	"testing"

	"aitu-superapp/education-service/internal/domain"
	"aitu-superapp/education-service/internal/usecase"
)

type mockEduRepo struct {
	gpa float32
	err error
}

func (m *mockEduRepo) GetStudentGPA(ctx context.Context, studentID string) (float32, error) {
	return m.gpa, m.err
}
func (m *mockEduRepo) SaveAttendance(ctx context.Context, att *domain.Attendance) error {
	return m.err
}
func (m *mockEduRepo) EnrollStudent(ctx context.Context, studentID, courseID string) error {
	return m.err
}
func (m *mockEduRepo) GetExamSchedule(ctx context.Context, studentID string) ([]domain.Schedule, error) {
	return nil, m.err
}
func (m *mockEduRepo) GetTodaySchedule(ctx context.Context, studentID string) ([]domain.Schedule, error) {
	return nil, m.err
}
func (m *mockEduRepo) GetWeekSchedule(ctx context.Context, studentID string) ([]domain.Schedule, error) {
	return nil, m.err
}
func (m *mockEduRepo) GetGradesList(ctx context.Context, studentID string) ([]domain.Grade, error) {
	return nil, m.err
}
func (m *mockEduRepo) GetAttendanceHistory(ctx context.Context, studentID string) ([]domain.Attendance, error) {
	return nil, m.err
}
func (m *mockEduRepo) GetCoursesCatalog(ctx context.Context) ([]domain.Subject, error) {
	return nil, m.err
}
func (m *mockEduRepo) DropCourse(ctx context.Context, studentID, courseID string) error {
	return m.err
}
func (m *mockEduRepo) SearchTeachers(ctx context.Context, query string) ([]string, error) {
	return nil, m.err
}
func (m *mockEduRepo) GetTeacherOfficeHours(ctx context.Context, teacherID string) ([]string, error) {
	return nil, m.err
}

type mockPublisher struct{}

func (m *mockPublisher) PublishLowAttendanceWarning(studentID, courseID string) error {
	return nil
}

func TestCalculateGPA_Success(t *testing.T) {
	repo := &mockEduRepo{gpa: 3.75}
	uc := usecase.NewEducationUseCase(repo, &mockPublisher{})

	gpa, err := uc.CalculateGPA(context.Background(), "student-123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if gpa != 3.75 {
		t.Errorf("expected GPA 3.75, got: %v", gpa)
	}
}

func TestCalculateGPA_RepoError(t *testing.T) {
	repo := &mockEduRepo{err: errors.New("db connection failed")}
	uc := usecase.NewEducationUseCase(repo, &mockPublisher{})

	_, err := uc.CalculateGPA(context.Background(), "student-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSubmitAttendance_AbsentPublishesNATS(t *testing.T) {
	repo := &mockEduRepo{}
	uc := usecase.NewEducationUseCase(repo, &mockPublisher{})

	err := uc.SubmitAttendance(context.Background(), "student-123", "schedule-456", "absent", "2026-05-18")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestEnrollToCourse_Success(t *testing.T) {
	repo := &mockEduRepo{}
	uc := usecase.NewEducationUseCase(repo, &mockPublisher{})

	err := uc.EnrollToCourse(context.Background(), "student-123", "course-456")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestEnrollToCourse_Error(t *testing.T) {
	repo := &mockEduRepo{err: errors.New("course full")}
	uc := usecase.NewEducationUseCase(repo, &mockPublisher{})

	err := uc.EnrollToCourse(context.Background(), "student-123", "course-456")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
