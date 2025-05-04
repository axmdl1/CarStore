package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	carpetpb "CarStore/CarService/api/pb/car"
	"CarStore/CarService/internal/handler"
	"CarStore/CarService/internal/repository"
	"CarStore/CarService/internal/usecase"
	"CarStore/CarService/pkg/mongo"
)

func main() {
	// load .env if present
	_ = godotenv.Load()

	// read env
	mongoURI := os.Getenv("CAR_SERVICE_MONGO_URI")
	dbName := os.Getenv("CAR_SERVICE_DB_NAME")
	grpcPort := os.Getenv("CAR_SERVICE_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}

	if mongoURI == "" || dbName == "" {
		log.Fatal("MONGO_URI and DB_NAME must be set")
	}

	// connect to Mongo
	client, err := mongo.NewMongoClient(mongoURI)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}
	db := client.Database(dbName)

	// wire layers
	carRepo := repository.NewCarRepo(db)
	carUC := usecase.NewCarUsecase(carRepo)

	// start gRPC server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("listen on %s failed: %v", grpcPort, err)
	}
	grpcServer := grpc.NewServer()

	// register gRPC handler
	carpetpb.RegisterCarServiceServer(grpcServer, handler.NewCarHandler(carUC))

	log.Printf("gRPC CarService listening on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC serve error: %v", err)
	}
}
