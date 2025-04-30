package main

import (
	"CarStore/OrderService/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoURI := os.Getenv("ORDER_SERVICE_MONGO_URI")
	port := os.Getenv("ORDER_SERVICE_PORT")
	if port == "" {
		port = "8083"
	}
	if mongoURI == "" {
		log.Fatal("Env variable MONGO_URI not set")
	}

	router := gin.Default()
	if err := routes.Setup(router, mongoURI); err != nil {
		log.Fatalf("Routes setup failed: %v", err)
	}

	log.Printf("Order service running on port %s", port)
	router.Run(":" + port)
}
