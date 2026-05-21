package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpc_delivery "aitu-superapp/community-service/internal/delivery/grpc"
	http_delivery "aitu-superapp/community-service/internal/delivery/http"
	nats_delivery "aitu-superapp/community-service/internal/delivery/nats"
	"aitu-superapp/community-service/internal/repository"
	"aitu-superapp/community-service/internal/usecase"
	pb "aitu-superapp/proto"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5444/community_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer redisClient.Close()

	subscriber := nats_delivery.NewCommunitySubscriber(nc)
	subscriber.SubscribeToFinanceEvents()

	pgRepo := repository.NewPostgresRepository(db)
	redisRepo := repository.NewRedisRepository(redisClient)
	uc := usecase.NewCommunityUseCase(pgRepo, redisRepo)

	app := fiber.New()
	httpHandler := http_delivery.NewCommunityHandler(uc)
	httpHandler.RegisterRoutes(app)
	go func() {
		log.Println("🌐 Community REST Server running on port 8003...")
		log.Fatal(app.Listen(":8003"))
	}()

	handler := grpc_delivery.NewCommunityHandler(uc)
	lis, _ := net.Listen("tcp", ":50053")
	grpcServer := grpc.NewServer()
	pb.RegisterCommunityServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Println("Community Service is running on gRPC port 50053")
	log.Fatal(grpcServer.Serve(lis))
}
