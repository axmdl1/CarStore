package routes

import (
	"CarStore/OrderService/internal/handler"
	"CarStore/OrderService/internal/repository"
	"CarStore/OrderService/internal/usecase"
	"CarStore/OrderService/pkg/mongo"
	"github.com/gin-gonic/gin"
	"os"
)

func Setup(r *gin.Engine, mongoURI string) error {
	client, err := mongo.NewMongoClient(mongoURI)
	if err != nil {
		return err
	}
	db := client.Database(os.Getenv("ORDER_SERVICE_DB_NAME"))

	orderRepo := repository.NewOrderRepo(db)
	orderUC := usecase.NewOrderUsecase(orderRepo)

	api := r.Group("/order")
	handler.NewOrderHandler(api, orderUC)
	return nil
}
