package main

import (
	orderpb "CarStore/OrderService/api/pb/order"
	"CarStore/OrderService/internal/handler"
	"CarStore/OrderService/internal/repository"
	"CarStore/OrderService/internal/usecase"
	"CarStore/OrderService/pkg/mongo"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	_ = godotenv.Load()
	uri := os.Getenv("ORDER_SERVICE_MONGO_URI")
	dbName := os.Getenv("ORDER_SERVICE_DB_NAME")
	port := os.Getenv("ORDER_SERVICE_PORT")
	if port == "" {
		port = "50054"
	}

	client, err := mongo.NewMongoClient(uri)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(dbName)

	repo := repository.NewOrderRepo(db)
	uc := usecase.NewOrderUsecase(repo)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	orderpb.RegisterOrderServiceServer(grpcServer, handler.NewOrderHandler(uc))

	log.Printf("gRPC OrderService listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
