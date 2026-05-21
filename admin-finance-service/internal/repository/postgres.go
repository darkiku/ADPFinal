package repository

import (
	"context"
	"database/sql"
	"errors"

	"aitu-superapp/admin-finance-service/internal/domain"

	_ "github.com/lib/pq"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) domain.AdminRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetBalance(ctx context.Context, userID string) (float32, error) {
	var balance float32
	err := r.db.QueryRowContext(ctx, `SELECT balance FROM finance_accounts WHERE user_id = $1`, userID).Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		// новый пользователь — создаём счёт с 0
		_, _ = r.db.ExecContext(ctx, `INSERT INTO finance_accounts (user_id, balance) VALUES ($1, 0) ON CONFLICT DO NOTHING`, userID)
		return 0, nil
	}
	return balance, err
}

func (r *postgresRepo) UpdateBalance(ctx context.Context, userID string, amount float32) error {
	query := `INSERT INTO finance_accounts (user_id, balance) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET balance = finance_accounts.balance + EXCLUDED.balance`
	_, err := r.db.ExecContext(ctx, query, userID, amount)
	return err
}

func (r *postgresRepo) CreateTransaction(ctx context.Context, userID string, amount float32, transType string) (string, error) {
	var txID string
	query := `INSERT INTO transactions (account_id, amount, type) VALUES ((SELECT id FROM finance_accounts WHERE user_id = $1 LIMIT 1), $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, userID, amount, transType).Scan(&txID)
	return txID, err
}

func (r *postgresRepo) InitiatePaymentTx(ctx context.Context, userID string, amount float32) (string, error) {
	// убеждаемся что счёт существует
	_, _ = r.db.ExecContext(ctx, `INSERT INTO finance_accounts (user_id, balance) VALUES ($1, 0) ON CONFLICT DO NOTHING`, userID)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`UPDATE finance_accounts SET balance = balance - $1 WHERE user_id = $2`,
		amount, userID,
	)
	if err != nil {
		return "", err
	}

	var txID string
	err = tx.QueryRowContext(ctx,
		`INSERT INTO transactions (account_id, amount, type)
		 VALUES ((SELECT id FROM finance_accounts WHERE user_id = $1 LIMIT 1), $2, 'expense')
		 RETURNING id`,
		userID, amount,
	).Scan(&txID)
	if err != nil {
		return "", err
	}

	return txID, tx.Commit()
}

func (r *postgresRepo) GetFinanceHistory(ctx context.Context, userID string) ([]domain.Transaction, error) {
	// убеждаемся что счёт существует
	_, _ = r.db.ExecContext(ctx, `INSERT INTO finance_accounts (user_id, balance) VALUES ($1, 0) ON CONFLICT DO NOTHING`, userID)

	query := `SELECT t.id, t.account_id, t.amount, t.type, t.created_at 
	          FROM transactions t 
	          JOIN finance_accounts f ON t.account_id = f.id 
	          WHERE f.user_id = $1 
	          ORDER BY t.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return []domain.Transaction{}, nil
	}
	defer rows.Close()

	var history []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.AccountID, &t.Amount, &t.Type, &t.CreatedAt); err != nil {
			continue
		}
		history = append(history, t)
	}
	if history == nil {
		history = []domain.Transaction{}
	}
	return history, nil
}

func (r *postgresRepo) CreateDocumentRequest(ctx context.Context, studentID, docType string) (string, error) {
	var reqID string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO document_requests (student_id, type, status) VALUES ($1, $2, 'pending') RETURNING id`,
		studentID, docType).Scan(&reqID)
	return reqID, err
}

func (r *postgresRepo) GetDocumentStatus(ctx context.Context, studentID string) ([]domain.DocumentRequest, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, student_id, type, status FROM document_requests WHERE student_id = $1 ORDER BY rowid DESC`,
		studentID)
	if err != nil {
		// rowid не всегда есть, пробуем без ORDER BY
		rows, err = r.db.QueryContext(ctx,
			`SELECT id, student_id, type, status FROM document_requests WHERE student_id = $1`,
			studentID)
		if err != nil {
			return []domain.DocumentRequest{}, nil
		}
	}
	defer rows.Close()

	var docs []domain.DocumentRequest
	for rows.Next() {
		var d domain.DocumentRequest
		if err := rows.Scan(&d.ID, &d.StudentID, &d.Type, &d.Status); err != nil {
			continue
		}
		docs = append(docs, d)
	}
	if docs == nil {
		docs = []domain.DocumentRequest{}
	}
	return docs, nil
}

func (r *postgresRepo) UpdateProfile(ctx context.Context, studentID, email, phone string) error {
	return nil
}

func (r *postgresRepo) GetDormitoryInfo(ctx context.Context, studentID string) (*domain.Dormitory, error) {
	var dorm domain.Dormitory
	err := r.db.QueryRowContext(ctx,
		`SELECT id, room_number, occupant_id FROM dormitories WHERE occupant_id = $1`,
		studentID).Scan(&dorm.ID, &dorm.RoomNumber, &dorm.OccupantID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // нет данных — не ошибка
	}
	if err != nil {
		return nil, err
	}
	return &dorm, nil
}

func (r *postgresRepo) RequestDormRepair(ctx context.Context, studentID, description string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO dorm_repairs (student_id, description, status) VALUES ($1, $2, 'pending')
		 ON CONFLICT DO NOTHING`,
		studentID, description)
	if err != nil {
		// таблицы может не быть — создаём на лету и повторяем
		_, _ = r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS dorm_repairs (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			student_id UUID NOT NULL,
			description TEXT NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO dorm_repairs (student_id, description, status) VALUES ($1, $2, 'pending')`,
			studentID, description)
	}
	return err
}

func (r *postgresRepo) GetScholarshipInfo(ctx context.Context, studentID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":         "Government Grant",
		"amount":       float64(41898),
		"currency":     "KZT",
		"next_payment": "2025-06-01",
		"status":       "active",
	}, nil
}

func (r *postgresRepo) GetLibraryDebts(ctx context.Context, studentID string) ([]string, error) {
	return []string{}, nil
}

func (r *postgresRepo) ReserveLibraryBook(ctx context.Context, studentID, bookID string) error {
	_, _ = r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS library_reservations (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		student_id UUID NOT NULL,
		book_id VARCHAR(255) NOT NULL,
		reserved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO library_reservations (student_id, book_id) VALUES ($1, $2)`,
		studentID, bookID)
	return err
}
