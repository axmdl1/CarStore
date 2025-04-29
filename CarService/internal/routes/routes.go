package routes

import (
	"CarStore/CarService/internal/handler"
	"CarStore/CarService/internal/repository"
	"CarStore/CarService/internal/usecase"
	"CarStore/CarService/pkg/mongo"
	"github.com/gin-gonic/gin"
	"os"
)

func Setup(r *gin.Engine, mongoURI string) error {
	client, err := mongo.NewMongoClient(mongoURI)
	if err != nil {
		return err
	}
	db := client.Database(os.Getenv("CAR_SERVICE_DB_NAME"))

	carRepo := repository.NewCarRepo(db)
	carUC := usecase.NewCarUsecase(carRepo)

	api := r.Group("/cars")
	handler.NewCarHandler(api, carUC)
	return nil
}
