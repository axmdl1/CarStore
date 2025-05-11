// cmd/main.go
package main

import (
	"CarStore/UserService/pkg/auth"
	"CarStore/UserService/pkg/email"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	userpb "CarStore/UserService/api/pb/user"
	"CarStore/UserService/internal/handler"
	"CarStore/UserService/internal/repository"
	"CarStore/UserService/internal/usecase"
	"CarStore/UserService/pkg/jwt"
	"CarStore/UserService/pkg/mongo"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, reading environment directly")
	}

	// env vars
	mongoURI := os.Getenv("MONGO_URI_ATLAS")
	dbName := os.Getenv("USER_SERVICE_DB_NAME")
	jwtSecret := os.Getenv("JWT_SECRET")
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpFrom := os.Getenv("SMTP_FROM")

	// sanity check
	if mongoURI == "" || dbName == "" || jwtSecret == "" {
		log.Fatal("MONGO_URI, DB_NAME and JWT_SECRET must be set")
	}

	// connect Mongo
	client, err := mongo.NewMongoClient(mongoURI + dbName)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	db := client.Database(dbName)

	// wiring
	userRepo := repository.NewUserRepository(db)
	jwtSvc := jwt.NewJWTService(jwtSecret, "UserService")
	emailSvc := email.NewSMTPSender(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom)
	userUC := usecase.NewUserUsecase(userRepo, jwtSvc, emailSvc)

	// gRPC server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("listen on %s: %v", grpcPort, err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(auth.UnaryAuthInterceptor(*jwtSvc)))

	// register your service implementation
	userpb.RegisterUserServiceServer(grpcServer, handler.NewAuthHandler(userUC))

	log.Printf("gRPC UserService listening on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC serve error: %v", err)
	}
}
