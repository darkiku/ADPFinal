package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpc_delivery "aitu-superapp/education-service/internal/delivery/grpc"
	http_delivery "aitu-superapp/education-service/internal/delivery/http"
	nats_delivery "aitu-superapp/education-service/internal/delivery/nats"
	"aitu-superapp/education-service/internal/repository"
	"aitu-superapp/education-service/internal/usecase"
	pb "aitu-superapp/proto"
)

func main() {
	// 1. Инфраструктура
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5442/education_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// 2. Слои (Clean Architecture DI)
	repo := repository.NewPostgresRepository(db)
	publisher := nats_delivery.NewNatsPublisher(nc)
	uc := usecase.NewEducationUseCase(repo, publisher)

	// 3. Запуск REST (Fiber) в фоне
	app := fiber.New()
	http_delivery.RegisterRESTHandlers(app, uc)
	go func() {
		log.Println("🌐 Education REST Server running on port 8001...")
		log.Fatal(app.Listen(":8001"))
	}()

	// 4. Запуск gRPC
	handler := grpc_delivery.NewEducationHandler(uc)
	lis, _ := net.Listen("tcp", ":50051")
	grpcServer := grpc.NewServer()
	pb.RegisterEducationServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Println("✅ Education Service is running on gRPC port 50051...")
	log.Fatal(grpcServer.Serve(lis))
}
