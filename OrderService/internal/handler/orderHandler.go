package handler

import (
	"CarStore/OrderService/internal/entity"
	"CarStore/OrderService/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type OrderHandler struct {
	uc *usecase.OrderUsecase
}

func NewOrderHandler(rg *gin.RouterGroup, orderUc *usecase.OrderUsecase) {
	h := &OrderHandler{uc: orderUc}
	rg.POST("/", h.Create)
	rg.GET("/", h.List)
	rg.GET("/:id", h.GetByID)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)

}

func (o *OrderHandler) Create(c *gin.Context) {
	var req struct {
		UserID     string  `json:"userId"`
		CarID      string  `json:"cartId"`
		Quantity   int     `json:"quantity"`
		TotalPrice float64 `json:"totalPrice"`
		Status     string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUUID, _ := uuid.Parse(req.UserID)
	carUUID, _ := uuid.Parse(req.CarID)
	order := &entity.Order{
		ID:         uuid.New(),
		UserID:     userUUID,
		CarID:      carUUID,
		Quantity:   req.Quantity,
		TotalPrice: req.TotalPrice,
		CreatedAt:  time.Now(),
	}
	if req.Status == "" {
		order.Status = "pending"
	} else {
		order.Status = req.Status
	}
	if err := o.uc.Create(c.Request.Context(), order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (o *OrderHandler) List(c *gin.Context) {
	orders, err := o.uc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.uc.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (o *OrderHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order entity.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order.ID = uid
	if err := o.uc.Update(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (o *OrderHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := o.uc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
