package http

import (
	"strings"

	"aitu-superapp/auth-service/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "aitu-superapp-secret-2026"

type AuthHandler struct{ uc domain.AuthUseCase }

func NewAuthHandler(uc domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Get("/verify", h.VerifyEmail)
	auth.Get("/profile", h.GetProfile)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if body.Email == "" || body.Password == "" || body.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "email, password and name are required"})
	}
	if err := h.uc.Register(c.Context(), body.Email, body.Password, body.Name); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Registered! Check your email to verify your account."})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	token, user, err := h.uc.Login(c.Context(), body.Email, body.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"token": token,
		"user":  fiber.Map{"id": user.ID, "email": user.Email, "name": user.Name, "role": user.Role},
	})
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "token required"})
	}
	if err := h.uc.VerifyEmail(c.Context(), token); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendString("Email verified! You can now login at http://localhost:5173")
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
	}
	userID := claims["user_id"].(string)
	user, err := h.uc.GetProfile(c.Context(), userID)
	if err != nil || user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(user)
}
