package usecase

import (
	"context"
	"encoding/json"
	"time"

	"aitu-superapp/education-service/internal/domain"
)

type educationUseCase struct {
	repo      domain.EducationRepository
	publisher domain.EventPublisher
}

func NewEducationUseCase(repo domain.EducationRepository, pub domain.EventPublisher) domain.EducationUseCase {
	return &educationUseCase{repo: repo, publisher: pub}
}

func (uc *educationUseCase) CalculateGPA(ctx context.Context, studentID string) (float32, error) {
	return uc.repo.GetStudentGPA(ctx, studentID)
}

func (uc *educationUseCase) SubmitAttendance(ctx context.Context, studentID, scheduleID, status, dateStr string) error {
	parsedDate, _ := time.Parse("2006-01-02", dateStr)
	att := &domain.Attendance{StudentID: studentID, ScheduleID: scheduleID, Status: status, Date: parsedDate}
	if err := uc.repo.SaveAttendance(ctx, att); err != nil {
		return err
	}
	if status == "absent" {
		_ = uc.publisher.PublishLowAttendanceWarning(studentID, scheduleID)
	}
	return nil
}

func (uc *educationUseCase) EnrollToCourse(ctx context.Context, studentID, courseID string) error {
	return uc.repo.EnrollStudent(ctx, studentID, courseID)
}

func (uc *educationUseCase) GetExamScheduleJSON(ctx context.Context, studentID string) (string, error) {
	exams, err := uc.repo.GetExamSchedule(ctx, studentID)
	if err != nil {
		return "[]", err
	}
	b, _ := json.Marshal(exams)
	return string(b), nil
}

func (uc *educationUseCase) GetAttendanceStats(ctx context.Context, studentID string) (float32, int, int, error) {
	return uc.repo.GetAttendanceStats(ctx, studentID)
}

func (uc *educationUseCase) GetTodaySchedule(ctx context.Context, studentID string) ([]domain.Schedule, error) {
	return uc.repo.GetTodaySchedule(ctx, studentID)
}
func (uc *educationUseCase) GetWeekSchedule(ctx context.Context, studentID string) ([]domain.Schedule, error) {
	return uc.repo.GetWeekSchedule(ctx, studentID)
}
func (uc *educationUseCase) GetGradesList(ctx context.Context, studentID string) ([]domain.Grade, error) {
	return uc.repo.GetGradesList(ctx, studentID)
}
func (uc *educationUseCase) GetAttendanceHistory(ctx context.Context, studentID string) ([]domain.Attendance, error) {
	return uc.repo.GetAttendanceHistory(ctx, studentID)
}
func (uc *educationUseCase) GetCoursesCatalog(ctx context.Context) ([]domain.Subject, error) {
	return uc.repo.GetCoursesCatalog(ctx)
}
func (uc *educationUseCase) DropCourse(ctx context.Context, studentID, courseID string) error {
	return uc.repo.DropCourse(ctx, studentID, courseID)
}
func (uc *educationUseCase) SearchTeachers(ctx context.Context, query string) ([]domain.Teacher, error) {
	return uc.repo.SearchTeachers(ctx, query)
}
func (uc *educationUseCase) GetTeacherOfficeHours(ctx context.Context, teacherID string) ([]string, error) {
	return uc.repo.GetTeacherOfficeHours(ctx, teacherID)
}
