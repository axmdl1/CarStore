package main

import (
	"CarStore/UserService/internal/handler"
	"CarStore/UserService/internal/repository"
	"CarStore/UserService/internal/usecase"
	"CarStore/UserService/pkg/jwt"
	"CarStore/UserService/pkg/mongo"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoURI := os.Getenv("MONGO_URI")
	jwtSecret := os.Getenv("JWT_SECRET")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if mongoURI == "" || jwtSecret == "" {
		log.Fatal("Environment variables MONGO_URI and JWT_SECRET must be set")
	}

	// Initialize MongoDB client
	client, err := mongo.NewMongoClient(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db := client.Database(dbName)

	// Setup repository, JWT service, and usecase
	userRepo := repository.NewUserRepository(db)
	jwtSvc := jwt.NewJWTService(jwtSecret, "UserService")
	userUC := usecase.NewUserUsecase(userRepo, jwtSvc)

	// Initialize Gin router and register routes
	router := gin.Default()
	api := router.Group("/user")
	handler.NewAuthHandler(api, userUC)

	// Start HTTP server
	log.Printf("UserService listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
