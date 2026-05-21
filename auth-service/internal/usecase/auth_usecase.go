package usecase

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"aitu-superapp/auth-service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtSecret    = "aitu-superapp-secret-2026"
	smtpHost     = "smtp.gmail.com"
	smtpPort     = "587"
	smtpUser     = "darkiku555@gmail.com"
	smtpPassword = "cfrb xryz fzeq znmv"
)

type authUseCase struct{ repo domain.AuthRepository }

func NewAuthUseCase(repo domain.AuthRepository) domain.AuthUseCase {
	return &authUseCase{repo: repo}
}

func (uc *authUseCase) Register(ctx context.Context, email, password, name string) error {
	exists, err := uc.repo.UserExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	verifyToken := uuid.New().String()
	userID := uuid.New().String()
	user := &domain.User{
		ID:           userID,
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		Role:         "student",
		Verified:     false,
		VerifyToken:  verifyToken,
	}

	if err := uc.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	// Автоматически создаём финансовый аккаунт с начальным балансом 0
	if err := uc.repo.CreateFinanceAccount(ctx, userID); err != nil {
		fmt.Println("[Register] failed to create finance account:", err)
		// Не прерываем регистрацию из-за этой ошибки
	}

	go uc.sendVerificationEmail(email, name, verifyToken)
	return nil
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", nil, errors.New("user not found")
	}
	if !user.Verified {
		return "", nil, errors.New("email not verified, check your inbox")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", nil, err
	}
	return tokenStr, user, nil
}

func (uc *authUseCase) VerifyEmail(ctx context.Context, token string) error {
	return uc.repo.VerifyUser(ctx, token)
}

func (uc *authUseCase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return uc.repo.GetUserByID(ctx, userID)
}

func (uc *authUseCase) sendVerificationEmail(to, name, token string) {
	verifyURL := fmt.Sprintf("http://localhost:4000/auth/verify?token=%s", token)
	subject := "Verify your AITU SuperApp email"
	body := fmt.Sprintf(
		"Hi %s!\n\nWelcome to AITU SuperApp.\n\nClick to verify:\n%s\n\nAITU Team", name, verifyURL,
	)
	msg := fmt.Sprintf(
		"From: AITU SuperApp <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		smtpUser, to, subject, body,
	)

	tlsConfig := &tls.Config{ServerName: smtpHost}
	client, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		fmt.Println("[SMTP] dial error:", err)
		return
	}
	defer client.Quit()

	if err = client.StartTLS(tlsConfig); err != nil {
		fmt.Println("[SMTP] starttls error:", err)
		return
	}
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	if err = client.Auth(auth); err != nil {
		fmt.Println("[SMTP] auth error:", err)
		return
	}
	if err = client.Mail(smtpUser); err != nil {
		fmt.Println("[SMTP] MAIL FROM error:", err)
		return
	}
	if err = client.Rcpt(to); err != nil {
		fmt.Println("[SMTP] RCPT TO error:", err)
		return
	}
	wc, err := client.Data()
	if err != nil {
		fmt.Println("[SMTP] DATA error:", err)
		return
	}
	_, err = fmt.Fprint(wc, msg)
	if err != nil {
		fmt.Println("[SMTP] write error:", err)
	}
	wc.Close()
	fmt.Printf("[SMTP] Email sent to %s\n", to)
}
