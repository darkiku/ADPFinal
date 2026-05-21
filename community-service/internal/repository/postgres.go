package repository

import (
	"context"
	"database/sql"

	"aitu-superapp/community-service/internal/domain"
	_ "github.com/lib/pq"
)

type postgresRepo struct{ db *sql.DB }

func NewPostgresRepository(db *sql.DB) domain.PostgresRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) SaveBooking(ctx context.Context, b *domain.CoworkingBooking) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO coworking_bookings (space_id, student_id, start_time, end_time) VALUES ($1, $2, $3, $4)`,
		b.SpaceID, b.StudentID, b.StartTime, b.EndTime)
	return err
}

func (r *postgresRepo) GetEventsList(ctx context.Context) ([]domain.Event, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, date, location FROM events ORDER BY date ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Event
	for rows.Next() {
		var e domain.Event
		rows.Scan(&e.ID, &e.Title, &e.Date, &e.Location)
		out = append(out, e)
	}
	return out, nil
}

func (r *postgresRepo) RegisterForEvent(ctx context.Context, eventID, studentID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO event_registrations (event_id, student_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		eventID, studentID)
	return err
}

func (r *postgresRepo) CancelEventRegistration(ctx context.Context, eventID, studentID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM event_registrations WHERE event_id=$1 AND student_id=$2`, eventID, studentID)
	return err
}

func (r *postgresRepo) CancelEvent(ctx context.Context, eventID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, eventID)
	return err
}

func (r *postgresRepo) GetCoworkingList(ctx context.Context) ([]string, error) {
	return []string{"Room 101", "Room 102", "Quiet Zone A", "Conference Room B"}, nil
}

func (r *postgresRepo) CancelBooking(ctx context.Context, bookingID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM coworking_bookings WHERE id = $1`, bookingID)
	return err
}

func (r *postgresRepo) GetMarketplaceItems(ctx context.Context) ([]domain.Ad, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, student_id, title, description, price FROM marketplace_ads ORDER BY created_at DESC`)
	if err != nil {
		return []domain.Ad{}, nil
	}
	defer rows.Close()
	var out []domain.Ad
	for rows.Next() {
		var a domain.Ad
		rows.Scan(&a.ID, &a.StudentID, &a.Title, &a.Description, &a.Price)
		out = append(out, a)
	}
	return out, nil
}

func (r *postgresRepo) CreateAd(ctx context.Context, ad *domain.Ad) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO marketplace_ads (student_id, title, description, price) VALUES ($1, $2, $3, $4) RETURNING id`,
		ad.StudentID, ad.Title, ad.Description, ad.Price).Scan(&id)
	return id, err
}

func (r *postgresRepo) ReportLostFound(ctx context.Context, item *domain.LostFound) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO lost_found (student_id, item, description, status) VALUES ($1, $2, $3, $4)`,
		item.StudentID, item.Item, item.Description, item.Status)
	return err
}

func (r *postgresRepo) GetLostFoundList(ctx context.Context) ([]domain.LostFound, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, student_id, item, description, status FROM lost_found ORDER BY id DESC`)
	if err != nil {
		return []domain.LostFound{}, nil
	}
	defer rows.Close()
	var out []domain.LostFound
	for rows.Next() {
		var lf domain.LostFound
		rows.Scan(&lf.ID, &lf.StudentID, &lf.Item, &lf.Description, &lf.Status)
		out = append(out, lf)
	}
	return out, nil
}

func (r *postgresRepo) SubmitFeedback(ctx context.Context, feedback *domain.Feedback) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO feedback (student_id, text, category) VALUES ($1, $2, $3)`,
		feedback.StudentID, feedback.Text, feedback.Category)
	return err
}

func (r *postgresRepo) GetNewsFeed(ctx context.Context) ([]domain.NewsItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, body, created_at FROM news ORDER BY created_at DESC LIMIT 20`)
	if err != nil {
		return []domain.NewsItem{}, nil
	}
	defer rows.Close()
	var out []domain.NewsItem
	for rows.Next() {
		var n domain.NewsItem
		rows.Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt)
		out = append(out, n)
	}
	return out, nil
}
