package repository

import (
	"context"
	"database/sql"
	"errors"

	"aitu-superapp/auth-service/internal/domain"
)

type authRepo struct {
	db      *sql.DB
	adminDB *sql.DB
}

func NewAuthRepository(db *sql.DB, adminDB *sql.DB) domain.AuthRepository {
	return &authRepo{db: db, adminDB: adminDB}
}

func (r *authRepo) CreateUser(ctx context.Context, u *domain.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, password_hash, name, role, verified, verify_token)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		u.ID, u.Email, u.PasswordHash, u.Name, u.Role, u.Verified, u.VerifyToken,
	)
	return err
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, role, verified FROM users WHERE email=$1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.Verified)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *authRepo) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, name, role, verified FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.Verified)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *authRepo) VerifyUser(ctx context.Context, token string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE users SET verified=true, verify_token='' WHERE verify_token=$1`, token,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("invalid verification token")
	}
	return nil
}

func (r *authRepo) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, email).Scan(&exists)
	return exists, err
}

func (r *authRepo) CreateFinanceAccount(ctx context.Context, userID string) error {
	if r.adminDB == nil {
		return nil
	}
	_, err := r.adminDB.ExecContext(ctx,
		`INSERT INTO finance_accounts (user_id, balance) VALUES ($1, 0) ON CONFLICT (user_id) DO NOTHING`,
		userID,
	)
	return err
}
