package main

import (
	"log"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	http_delivery "aitu-superapp/api-gateway/internal/delivery/http"
	pb "aitu-superapp/proto"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-User-ID, X-User-Role",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, OPTIONS",
		AllowCredentials: false,
	}))

	prometheus := fiberprometheus.New("aitu-gateway")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Static("/", "./frontend")

	connEdu, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer connEdu.Close()
	eduClient := pb.NewEducationServiceClient(connEdu)

	connAdmin, _ := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer connAdmin.Close()
	adminClient := pb.NewAdminServiceClient(connAdmin)

	connComm, _ := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer connComm.Close()
	commClient := pb.NewCommunityServiceClient(connComm)

	http_delivery.SetupRoutes(app, eduClient, adminClient, commClient)

	log.Println("Gateway: http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
