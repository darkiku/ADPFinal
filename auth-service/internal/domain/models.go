package domain

import "context"

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	Verified     bool   `json:"verified"`
	VerifyToken  string `json:"-"`
}

type AuthRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	VerifyUser(ctx context.Context, token string) error
	UserExists(ctx context.Context, email string) (bool, error)
	CreateFinanceAccount(ctx context.Context, userID string) error
}

type AuthUseCase interface {
	Register(ctx context.Context, email, password, name string) error
	Login(ctx context.Context, email, password string) (string, *User, error)
	VerifyEmail(ctx context.Context, token string) error
	GetProfile(ctx context.Context, userID string) (*User, error)
}
