package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"

	carpetpb "CarStore/CarService/api/pb/car"
	"CarStore/CarService/internal/entity"
	"CarStore/CarService/internal/usecase"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CarHandler struct {
	carpetpb.UnimplementedCarServiceServer
	uc *usecase.CarUsecase
}

func NewCarHandler(uc *usecase.CarUsecase) carpetpb.CarServiceServer {
	return &CarHandler{uc: uc}
}

func (h *CarHandler) CreateCar(ctx context.Context, req *carpetpb.CreateCarRequest) (*carpetpb.CreateCarResponse, error) {
	log.Printf("CreateCar request: %+v", req)
	if req.Car == nil {
		return nil, status.Error(codes.InvalidArgument, "car payload is required")
	}
	e := &entity.Car{
		Brand:          req.Car.Brand,
		Model:          req.Car.Model,
		Year:           int(req.Car.Year),
		Price:          req.Car.Price,
		Description:    req.Car.Description,
		EngineCapacity: req.Car.EngineCapacity,
		Mileage:        int(req.Car.Mileage),
		Gearbox:        req.Car.Gearbox,
		EngineType:     req.Car.EngineType,
		Stock:          int(req.Car.Stock),
	}

	if err := h.uc.Create(ctx, e); err != nil {
		return nil, err
	}

	return &carpetpb.CreateCarResponse{
		Car: &carpetpb.Car{
			Id:             e.ID.String(),
			Brand:          e.Brand,
			Model:          e.Model,
			Year:           int32(e.Year),
			Price:          e.Price,
			Description:    e.Description,
			EngineCapacity: e.EngineCapacity,
			Mileage:        int32(e.Mileage),
			Gearbox:        e.Gearbox,
			EngineType:     e.EngineType,
			Stock:          int32(e.Stock),
			CreatedAt:      timestamppb.New(e.CreatedAt),
		},
	}, nil
}

func (h *CarHandler) GetCar(ctx context.Context, req *carpetpb.GetCarRequest) (*carpetpb.GetCarResponse, error) {
	log.Printf("GetCar request: %+v", req)
	e, err := h.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &carpetpb.GetCarResponse{
		Car: &carpetpb.Car{
			Id:             e.ID.String(),
			Brand:          e.Brand,
			Model:          e.Model,
			Year:           int32(e.Year),
			Price:          e.Price,
			Description:    e.Description,
			EngineCapacity: e.EngineCapacity,
			Mileage:        int32(e.Mileage),
			Gearbox:        e.Gearbox,
			EngineType:     e.EngineType,
			CreatedAt:      timestamppb.New(e.CreatedAt),
		},
	}, nil
}

func (h *CarHandler) UpdateCar(ctx context.Context, req *carpetpb.UpdateCarRequest) (*carpetpb.UpdateCarResponse, error) {
	log.Printf("UpdateCar request: %+v", req)
	uid, err := uuid.Parse(req.Car.Id)
	if err != nil {
		return nil, err
	}
	e := &entity.Car{
		ID:             uid,
		Brand:          req.Car.Brand,
		Model:          req.Car.Model,
		Year:           int(req.Car.Year),
		Price:          req.Car.Price,
		Description:    req.Car.Description,
		EngineCapacity: req.Car.EngineCapacity,
		Mileage:        int(req.Car.Mileage),
		Gearbox:        req.Car.Gearbox,
		EngineType:     req.Car.EngineType,
		CreatedAt:      req.Car.CreatedAt.AsTime(),
	}

	if err := h.uc.Update(ctx, e); err != nil {
		return nil, err
	}
	return &carpetpb.UpdateCarResponse{Car: req.Car}, nil
}

func (h *CarHandler) DeleteCar(ctx context.Context, req *carpetpb.DeleteCarRequest) (*carpetpb.DeleteCarResponse, error) {
	log.Printf("DeleteCar request: %+v", req)
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &carpetpb.DeleteCarResponse{Success: true}, nil
}

func (h *CarHandler) ListCars(ctx context.Context, req *carpetpb.ListCarsRequest) (*carpetpb.ListCarsResponse, error) {
	log.Printf("ListCars request: %+v", req)
	es, err := h.uc.List(ctx)
	if err != nil {
		return nil, err
	}
	resp := &carpetpb.ListCarsResponse{}
	for _, e := range es {
		resp.Cars = append(resp.Cars, &carpetpb.Car{
			Id:             e.ID.String(),
			Brand:          e.Brand,
			Model:          e.Model,
			Year:           int32(e.Year),
			Price:          e.Price,
			Description:    e.Description,
			EngineCapacity: e.EngineCapacity,
			Mileage:        int32(e.Mileage),
			Gearbox:        e.Gearbox,
			EngineType:     e.EngineType,
			Stock:          int32(e.Stock),
			CreatedAt:      timestamppb.New(e.CreatedAt),
		})
	}
	return resp, nil
}

func (h *CarHandler) DecreaseStock(ctx context.Context, req *carpetpb.DecreaseStockRequest) (*carpetpb.DecreaseStockResponse, error) {
	newStock, err := h.uc.DecreaseStock(ctx, req.CarId, int(req.Quantity))
	if err != nil {
		return nil, err
	}
	return &carpetpb.DecreaseStockResponse{NewStock: int32(newStock)}, nil
}
