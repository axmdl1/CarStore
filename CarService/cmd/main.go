package main

import (
	"CarStore/CarService/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoURI := os.Getenv("CAR_SERVICE_MONGO_URI")
	port := os.Getenv("CAR_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}
	if mongoURI == "" {
		log.Fatal("CAR_SERVICE_MONGO_URI environment variable not set")
	}

	router := gin.Default()
	if err := routes.Setup(router, mongoURI); err != nil {
		log.Fatalf("Routes setup failed: %v", err)
	}

	log.Printf("Car service running on port %s", port)
	router.Run(":" + port)
}
