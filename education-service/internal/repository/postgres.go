package repository

import (
	"context"
	"database/sql"
	"time"

	"aitu-superapp/education-service/internal/domain"

	_ "github.com/lib/pq"
)

type postgresRepo struct{ db *sql.DB }

func NewPostgresRepository(db *sql.DB) domain.EducationRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetStudentGPA(ctx context.Context, studentID string) (float32, error) {
	var gpa float32
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(AVG(
			CASE
				WHEN score >= 95 THEN 4.33 WHEN score >= 90 THEN 4.00
				WHEN score >= 85 THEN 3.67 WHEN score >= 80 THEN 3.33
				WHEN score >= 75 THEN 3.00 WHEN score >= 70 THEN 2.67
				WHEN score >= 65 THEN 2.33 WHEN score >= 60 THEN 2.00
				WHEN score >= 55 THEN 1.67 WHEN score >= 50 THEN 1.00
				ELSE 0.00
			END
		), 0) FROM grades WHERE student_id = $1`, studentID).Scan(&gpa)
	return gpa, err
}

func (r *postgresRepo) SaveAttendance(ctx context.Context, att *domain.Attendance) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO attendance (student_id, schedule_id, status, date) VALUES ($1, $2, $3, $4)`,
		att.StudentID, att.ScheduleID, att.Status, att.Date)
	return err
}

func (r *postgresRepo) EnrollStudent(ctx context.Context, studentID, courseID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO student_courses (student_id, subject_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		studentID, courseID)
	return err
}

func (r *postgresRepo) GetExamSchedule(ctx context.Context, studentID string) ([]domain.Schedule, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.subject_id, sub.title, s.room, s.start_time, s.day_of_week
		FROM schedule s
		JOIN subjects sub ON sub.id = s.subject_id
		JOIN student_courses sc ON sc.subject_id = s.subject_id
		WHERE sc.student_id = $1`, studentID)
	if err != nil {
		return []domain.Schedule{}, nil
	}
	defer rows.Close()
	var out []domain.Schedule
	for rows.Next() {
		var s domain.Schedule
		rows.Scan(&s.ID, &s.SubjectID, &s.SubjectTitle, &s.Room, &s.StartTime, &s.DayOfWeek)
		out = append(out, s)
	}
	return out, nil
}

func (r *postgresRepo) GetAttendanceStats(ctx context.Context, studentID string) (float32, int, int, error) {
	var total, present int
	r.db.QueryRowContext(ctx,
		`SELECT COUNT(*), COALESCE(SUM(CASE WHEN status='present' THEN 1 ELSE 0 END),0) FROM attendance WHERE student_id=$1`,
		studentID).Scan(&total, &present)
	if total == 0 {
		return 0, 0, 0, nil
	}
	return float32(present) / float32(total) * 100, total, present, nil
}

func (r *postgresRepo) GetTodaySchedule(ctx context.Context, studentID string) ([]domain.Schedule, error) {
	today := int(time.Now().Weekday())
	if today == 0 {
		today = 7
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.subject_id, sub.title, s.room, s.start_time, s.day_of_week
		FROM schedule s JOIN subjects sub ON sub.id = s.subject_id
		WHERE s.day_of_week = $1`, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Schedule
	for rows.Next() {
		var s domain.Schedule
		rows.Scan(&s.ID, &s.SubjectID, &s.SubjectTitle, &s.Room, &s.StartTime, &s.DayOfWeek)
		out = append(out, s)
	}
	return out, nil
}

func (r *postgresRepo) GetWeekSchedule(ctx context.Context, studentID string) ([]domain.Schedule, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.subject_id, sub.title, s.room, s.start_time, s.day_of_week
		FROM schedule s JOIN subjects sub ON sub.id = s.subject_id
		ORDER BY s.day_of_week, s.start_time`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Schedule
	for rows.Next() {
		var s domain.Schedule
		rows.Scan(&s.ID, &s.SubjectID, &s.SubjectTitle, &s.Room, &s.StartTime, &s.DayOfWeek)
		out = append(out, s)
	}
	return out, nil
}

func (r *postgresRepo) GetGradesList(ctx context.Context, studentID string) ([]domain.Grade, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT g.id, g.student_id, g.subject_id, sub.title, g.score, g.type
		FROM grades g JOIN subjects sub ON sub.id = g.subject_id
		WHERE g.student_id = $1`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Grade
	for rows.Next() {
		var g domain.Grade
		rows.Scan(&g.ID, &g.StudentID, &g.SubjectID, &g.SubjectTitle, &g.Score, &g.Type)
		out = append(out, g)
	}
	return out, nil
}

func (r *postgresRepo) GetAttendanceHistory(ctx context.Context, studentID string) ([]domain.Attendance, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, student_id, schedule_id, status, date FROM attendance WHERE student_id = $1`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Attendance
	for rows.Next() {
		var a domain.Attendance
		rows.Scan(&a.ID, &a.StudentID, &a.ScheduleID, &a.Status, &a.Date)
		out = append(out, a)
	}
	return out, nil
}

func (r *postgresRepo) GetCoursesCatalog(ctx context.Context) ([]domain.Subject, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, credits, teacher_id FROM subjects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Subject
	for rows.Next() {
		var s domain.Subject
		rows.Scan(&s.ID, &s.Title, &s.Credits, &s.TeacherID)
		out = append(out, s)
	}
	return out, nil
}

func (r *postgresRepo) DropCourse(ctx context.Context, studentID, courseID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM student_courses WHERE student_id=$1 AND subject_id=$2`, studentID, courseID)
	return err
}

func (r *postgresRepo) SearchTeachers(ctx context.Context, query string) ([]domain.Teacher, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, department FROM teachers WHERE name ILIKE $1 OR department ILIKE $1`,
		"%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Teacher
	for rows.Next() {
		var t domain.Teacher
		rows.Scan(&t.ID, &t.Name, &t.Department)
		out = append(out, t)
	}
	return out, nil
}

func (r *postgresRepo) GetTeacherOfficeHours(ctx context.Context, teacherID string) ([]string, error) {
	return []string{"Mon 10:00-12:00", "Wed 14:00-16:00", "Fri 10:00-11:00"}, nil
}
