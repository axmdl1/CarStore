package handler

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	orderpb "CarStore/OrderService/api/pb/order"
	"CarStore/OrderService/internal/entity"
	"CarStore/OrderService/internal/usecase"
)

type OrderHandler struct {
	orderpb.UnimplementedOrderServiceServer
	uc *usecase.OrderUsecase
}

func NewOrderHandler(uc *usecase.OrderUsecase) orderpb.OrderServiceServer {
	return &OrderHandler{uc: uc}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	log.Printf("CreateOrder request: %+v", req)
	e := &entity.Order{
		UserID:     uuid.MustParse(req.UserId),
		CarID:      uuid.MustParse(req.CarId),
		Quantity:   int(req.Quantity),
		TotalPrice: req.TotalPrice,
	}
	if err := h.uc.Create(ctx, e); err != nil {
		return nil, err
	}
	return &orderpb.CreateOrderResponse{Order: &orderpb.Order{
		Id:         e.ID.String(),
		UserId:     e.UserID.String(),
		CarId:      e.CarID.String(),
		Quantity:   int32(e.Quantity),
		TotalPrice: e.TotalPrice,
		Status:     e.Status,
		CreatedAt:  timestamppb.New(e.CreatedAt),
	}}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	log.Printf("GetOrder request: %+v", req)
	e, err := h.uc.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &orderpb.GetOrderResponse{Order: &orderpb.Order{
		Id:         e.ID.String(),
		UserId:     e.UserID.String(),
		CarId:      e.CarID.String(),
		Quantity:   int32(e.Quantity),
		TotalPrice: e.TotalPrice,
		Status:     e.Status,
		CreatedAt:  timestamppb.New(e.CreatedAt),
	}}, nil
}

func (h *OrderHandler) UpdateOrder(ctx context.Context, req *orderpb.UpdateOrderRequest) (*orderpb.UpdateOrderResponse, error) {
	log.Printf("UpdateOrder request: %+v", req)
	e := &entity.Order{
		ID:         uuid.MustParse(req.Order.Id),
		UserID:     uuid.MustParse(req.Order.UserId),
		CarID:      uuid.MustParse(req.Order.CarId),
		Quantity:   int(req.Order.Quantity),
		TotalPrice: req.Order.TotalPrice,
		Status:     req.Order.Status,
		CreatedAt:  req.Order.CreatedAt.AsTime(),
	}
	if err := h.uc.Update(ctx, e); err != nil {
		return nil, err
	}
	return &orderpb.UpdateOrderResponse{Order: req.Order}, nil
}

func (h *OrderHandler) DeleteOrder(ctx context.Context, req *orderpb.DeleteOrderRequest) (*orderpb.DeleteOrderResponse, error) {
	log.Printf("DeleteOrder request: %+v", req)
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &orderpb.DeleteOrderResponse{Success: true}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	log.Printf("ListOrders request: %+v", req)
	es, err := h.uc.List(ctx)
	if err != nil {
		return nil, err
	}
	res := &orderpb.ListOrdersResponse{}
	for _, e := range es {
		res.Orders = append(res.Orders, &orderpb.Order{
			Id:         e.ID.String(),
			UserId:     e.UserID.String(),
			CarId:      e.CarID.String(),
			Quantity:   int32(e.Quantity),
			TotalPrice: e.TotalPrice,
			Status:     e.Status,
			CreatedAt:  timestamppb.New(e.CreatedAt),
		})
	}
	return res, nil
}
