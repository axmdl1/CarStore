package main

import (
	orderpb "CarStore/OrderService/api/pb/order"
	"CarStore/OrderService/internal/handler"
	"CarStore/OrderService/internal/repository"
	"CarStore/OrderService/internal/usecase"
	"CarStore/OrderService/pkg/mongo"
	"CarStore/UserService/pkg/auth"
	"CarStore/UserService/pkg/jwt"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	_ = godotenv.Load()
	uri := os.Getenv("ORDER_SERVICE_MONGO_URI")
	dbName := os.Getenv("ORDER_SERVICE_DB_NAME")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("ORDER_SERVICE_PORT")
	if port == "" {
		port = "50054"
	}

	client, err := mongo.NewMongoClient(uri + dbName)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(dbName)

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("NATS connect failed: %v", err)
	}

	repo := repository.NewOrderRepo(db)
	uc := usecase.NewOrderUsecase(repo, nc)
	jwtSvc := jwt.NewJWTService(jwtSecret, "OrderService")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(auth.UnaryAuthInterceptor(*jwtSvc)))

	orderpb.RegisterOrderServiceServer(grpcServer, handler.NewOrderHandler(uc))

	log.Printf("gRPC OrderService listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
