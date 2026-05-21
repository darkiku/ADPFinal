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

	grpc_delivery "aitu-superapp/admin-finance-service/internal/delivery/grpc"
	http_delivery "aitu-superapp/admin-finance-service/internal/delivery/http"
	nats_delivery "aitu-superapp/admin-finance-service/internal/delivery/nats"
	"aitu-superapp/admin-finance-service/internal/repository"
	"aitu-superapp/admin-finance-service/internal/usecase"
	pb "aitu-superapp/proto"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5443/admin_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Инициализация NATS Subscriber (Слушаем ивенты)
	subscriber := nats_delivery.NewAdminSubscriber(nc)
	subscriber.SubscribeToAllEvents()

	repo := repository.NewPostgresRepository(db)
	publisher := nats_delivery.NewAdminPublisher(nc)
	uc := usecase.NewAdminUseCase(repo, publisher)

	// REST сервер
	app := fiber.New()
	http_delivery.RegisterRESTHandlers(app, uc)
	go func() {
		log.Println("🌐 Admin REST Server running on port 8002...")
		log.Fatal(app.Listen(":8002"))
	}()

	// gRPC сервер
	handler := grpc_delivery.NewAdminHandler(uc)
	lis, _ := net.Listen("tcp", ":50052")
	grpcServer := grpc.NewServer()
	pb.RegisterAdminServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Println("✅ Admin & Finance Service is running on gRPC port 50052...")
	log.Fatal(grpcServer.Serve(lis))
}
