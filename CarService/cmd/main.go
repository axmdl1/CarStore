package main

import (
	"CarStore/UserService/pkg/auth"
	"CarStore/UserService/pkg/jwt"
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
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
	mongoURI := os.Getenv("MONGO_URI_ATLAS")
	dbName := os.Getenv("CAR_SERVICE_DB_NAME")
	jwtSecret := os.Getenv("JWT_SECRET")
	grpcPort := os.Getenv("CAR_SERVICE_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}

	if mongoURI == "" || dbName == "" {
		log.Fatal("MONGO_URI and DB_NAME must be set")
	}

	// connect to Mongo
	client, err := mongo.NewMongoClient(mongoURI + dbName)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}
	db := client.Database(dbName)

	// wire layers
	carRepo := repository.NewCarRepo(db)
	carUC := usecase.NewCarUsecase(carRepo)
	jwtSvc := jwt.NewJWTService(jwtSecret, "CarService")

	//nats connection
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("NATS connect: %v", err)
	}
	_, err = nc.Subscribe("order.created", func(m *nats.Msg) {
		var evt struct {
			CarID    string `json:"car_id"`
			Quantity int    `json:"quantity"`
		}
		if err := json.Unmarshal(m.Data, &evt); err != nil {
			log.Printf("bad event: %v", err)
			return
		}
		newStock, err := carUC.DecreaseStock(context.Background(), evt.CarID, evt.Quantity)
		if err != nil {
			log.Printf("decrease stock failed: %v", err)
		} else {
			log.Printf("stock for %s decreased by %d, now %d", evt.CarID, evt.Quantity, newStock)
		}
	})
	if err != nil {
		log.Fatalf("NATS subscribe: %v", err)
	}

	// start gRPC server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("listen on %s failed: %v", grpcPort, err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(auth.UnaryAuthInterceptor(*jwtSvc)))

	// register gRPC handler
	carpetpb.RegisterCarServiceServer(grpcServer, handler.NewCarHandler(carUC))

	log.Printf("gRPC CarService listening on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC serve error: %v", err)
	}
}
