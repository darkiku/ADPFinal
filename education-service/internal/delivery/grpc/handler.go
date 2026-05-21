package grpc_delivery

import (
	"context"
	"strings"

	"aitu-superapp/education-service/internal/domain"
	pb "aitu-superapp/proto"
)

type EducationHandler struct {
	pb.UnimplementedEducationServiceServer
	useCase domain.EducationUseCase
}

func NewEducationHandler(uc domain.EducationUseCase) *EducationHandler {
	return &EducationHandler{useCase: uc}
}

func (h *EducationHandler) GetStudentGPA(ctx context.Context, req *pb.GPARequest) (*pb.GPAResponse, error) {
	gpa, err := h.useCase.CalculateGPA(ctx, req.StudentId)
	if err != nil {
		return nil, err
	}
	return &pb.GPAResponse{Gpa: gpa}, nil
}

func (h *EducationHandler) SubmitAttendance(ctx context.Context, req *pb.AttendanceRequest) (*pb.AttendanceResponse, error) {
	err := h.useCase.SubmitAttendance(ctx, req.StudentId, req.ScheduleId, req.Status, req.Date)
	if err != nil {
		return &pb.AttendanceResponse{Success: false}, err
	}
	return &pb.AttendanceResponse{Success: true}, nil
}

func (h *EducationHandler) GetAttendanceStats(ctx context.Context, req *pb.AttendanceStatsRequest) (*pb.AttendanceStatsResponse, error) {
	percent, total, present, err := h.useCase.GetAttendanceStats(ctx, req.StudentId)
	if err != nil {
		return nil, err
	}
	return &pb.AttendanceStatsResponse{Percent: percent, Total: int32(total), Present: int32(present)}, nil
}

func (h *EducationHandler) GetScheduleToday(ctx context.Context, req *pb.ScheduleTodayRequest) (*pb.ScheduleTodayResponse, error) {
	schedules, err := h.useCase.GetTodaySchedule(ctx, req.StudentId)
	if err != nil {
		return nil, err
	}
	var items []*pb.ScheduleItem
	for _, s := range schedules {
		items = append(items, &pb.ScheduleItem{
			Subject: s.SubjectTitle,
			Room:    s.Room,
			Time:    s.StartTime.Format("15:04"),
			Type:    "Lecture",
		})
	}
	return &pb.ScheduleTodayResponse{Items: items}, nil
}

func (h *EducationHandler) GetGrades(ctx context.Context, req *pb.GetGradesRequest) (*pb.GetGradesResponse, error) {
	grades, err := h.useCase.GetGradesList(ctx, req.StudentId)
	if err != nil {
		return nil, err
	}
	var items []*pb.GradeItem
	for _, g := range grades {
		items = append(items, &pb.GradeItem{Subject: g.SubjectTitle, Score: g.Score, Type: g.Type})
	}
	return &pb.GetGradesResponse{Grades: items}, nil
}

func (h *EducationHandler) GetCoursesCatalog(ctx context.Context, req *pb.CoursesRequest) (*pb.CoursesResponse, error) {
	courses, err := h.useCase.GetCoursesCatalog(ctx)
	if err != nil {
		return nil, err
	}
	var items []*pb.CourseItem
	for _, c := range courses {
		items = append(items, &pb.CourseItem{Id: c.ID, Title: c.Title, Credits: int32(c.Credits)})
	}
	return &pb.CoursesResponse{Courses: items}, nil
}

func (h *EducationHandler) EnrollToCourse(ctx context.Context, req *pb.EnrollRequest) (*pb.EnrollResponse, error) {
	err := h.useCase.EnrollToCourse(ctx, req.StudentId, req.CourseId)
	if err != nil {
		return &pb.EnrollResponse{Success: false, Message: err.Error()}, err
	}
	return &pb.EnrollResponse{Success: true, Message: "Enrolled successfully"}, nil
}

func (h *EducationHandler) DropCourse(ctx context.Context, req *pb.DropCourseRequest) (*pb.DropCourseResponse, error) {
	err := h.useCase.DropCourse(ctx, req.StudentId, req.CourseId)
	if err != nil {
		return &pb.DropCourseResponse{Success: false}, err
	}
	return &pb.DropCourseResponse{Success: true}, nil
}

func (h *EducationHandler) SearchTeacher(ctx context.Context, req *pb.SearchTeacherRequest) (*pb.SearchTeacherResponse, error) {
	teachers, err := h.useCase.SearchTeachers(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	var items []*pb.TeacherItem
	for _, t := range teachers {
		items = append(items, &pb.TeacherItem{
			Name:       t.Name,
			Department: t.Department,
			Email:      "",
		})
	}
	return &pb.SearchTeacherResponse{Teachers: items}, nil
}

func (h *EducationHandler) GetTeacherOfficeHours(ctx context.Context, req *pb.OfficeHoursRequest) (*pb.OfficeHoursResponse, error) {
	hours, _ := h.useCase.GetTeacherOfficeHours(ctx, req.TeacherId)
	return &pb.OfficeHoursResponse{Hours: strings.Join(hours, ", "), Room: ""}, nil
}

func (h *EducationHandler) GetExamSchedule(ctx context.Context, req *pb.ExamScheduleRequest) (*pb.ExamScheduleResponse, error) {
	j, _ := h.useCase.GetExamScheduleJSON(ctx, req.StudentId)
	return &pb.ExamScheduleResponse{ScheduleJson: j}, nil
}

func (h *EducationHandler) GetStudentProfile(ctx context.Context, req *pb.StudentProfileRequest) (*pb.StudentProfileResponse, error) {
	return &pb.StudentProfileResponse{Name: "", Email: "", Specialty: "", Year: 1}, nil
}
