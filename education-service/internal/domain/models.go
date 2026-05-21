package domain

import (
	"context"
	"time"
)

type Subject struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Credits   int    `json:"credits"`
	TeacherID string `json:"teacher_id"`
}

type Schedule struct {
	ID           string    `json:"id"`
	SubjectID    string    `json:"subject_id"`
	SubjectTitle string    `json:"subject"`
	Room         string    `json:"room"`
	StartTime    time.Time `json:"start_time"`
	DayOfWeek    int       `json:"day_of_week"`
}

type Grade struct {
	ID           string  `json:"id"`
	StudentID    string  `json:"student_id"`
	SubjectID    string  `json:"subject_id"`
	SubjectTitle string  `json:"subject"`
	Score        float32 `json:"score"`
	Type         string  `json:"type"`
}

type Attendance struct {
	ID         string    `json:"id"`
	StudentID  string    `json:"student_id"`
	ScheduleID string    `json:"schedule_id"`
	Status     string    `json:"status"`
	Date       time.Time `json:"date"`
}

type Teacher struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
}

type EducationRepository interface {
	GetStudentGPA(ctx context.Context, studentID string) (float32, error)
	SaveAttendance(ctx context.Context, att *Attendance) error
	EnrollStudent(ctx context.Context, studentID, courseID string) error
	GetExamSchedule(ctx context.Context, studentID string) ([]Schedule, error)
	GetAttendanceStats(ctx context.Context, studentID string) (percent float32, total int, present int, err error)
	GetTodaySchedule(ctx context.Context, studentID string) ([]Schedule, error)
	GetWeekSchedule(ctx context.Context, studentID string) ([]Schedule, error)
	GetGradesList(ctx context.Context, studentID string) ([]Grade, error)
	GetAttendanceHistory(ctx context.Context, studentID string) ([]Attendance, error)
	GetCoursesCatalog(ctx context.Context) ([]Subject, error)
	DropCourse(ctx context.Context, studentID, courseID string) error
	SearchTeachers(ctx context.Context, query string) ([]Teacher, error)
	GetTeacherOfficeHours(ctx context.Context, teacherID string) ([]string, error)
}

type EducationUseCase interface {
	CalculateGPA(ctx context.Context, studentID string) (float32, error)
	SubmitAttendance(ctx context.Context, studentID, scheduleID, status, date string) error
	EnrollToCourse(ctx context.Context, studentID, courseID string) error
	GetExamScheduleJSON(ctx context.Context, studentID string) (string, error)
	GetAttendanceStats(ctx context.Context, studentID string) (float32, int, int, error)
	GetTodaySchedule(ctx context.Context, studentID string) ([]Schedule, error)
	GetWeekSchedule(ctx context.Context, studentID string) ([]Schedule, error)
	GetGradesList(ctx context.Context, studentID string) ([]Grade, error)
	GetAttendanceHistory(ctx context.Context, studentID string) ([]Attendance, error)
	GetCoursesCatalog(ctx context.Context) ([]Subject, error)
	DropCourse(ctx context.Context, studentID, courseID string) error
	SearchTeachers(ctx context.Context, query string) ([]Teacher, error)
	GetTeacherOfficeHours(ctx context.Context, teacherID string) ([]string, error)
}

type EventPublisher interface {
	PublishLowAttendanceWarning(studentID string, courseID string) error
}
