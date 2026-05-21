package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"

	httpdelivery "aitu-superapp/auth-service/internal/delivery/http"
	"aitu-superapp/auth-service/internal/repository"
	"aitu-superapp/auth-service/internal/usecase"
)

func main() {
	authDB, err := sql.Open("postgres", "host=localhost port=5445 user=postgres password=postgres dbname=auth_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer authDB.Close()

	adminDB, err := sql.Open("postgres", "host=localhost port=5443 user=postgres password=postgres dbname=admin_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer adminDB.Close()

	repo := repository.NewAuthRepository(authDB, adminDB)
	uc := usecase.NewAuthUseCase(repo)
	handler := httpdelivery.NewAuthHandler(uc)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool { return true },
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, OPTIONS",
	}))

	handler.RegisterRoutes(app)

	log.Println("✅ Auth service: http://localhost:4000")
	log.Fatal(app.Listen(":4000"))
}
