package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Задаёт методы бизнес-логики
type OrderService interface {
	CreateOrder(ctx context.Context, pord *orders.PrimaryOrder) (uuid.UUID, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*orders.Order, error)
	GetOrders(ctx context.Context) ([]*orders.Order, error)
	Preschedule(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error
}

// Содержит логику обработчиков для заявок
type OrderHandler struct {
	orderService OrderService
}

// NewOrderHandler создаёт новый экземпляр OrderHandler.
// Использует внедрение зависимостей для OrderService.
func NewOrderHandler(os OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: os,
	}
}

type PrescheduleRequest struct {
	ScheduledFor *time.Time `json:"scheduled_for"`
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.orderService.GetOrders(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id", "details": err.Error()})
		return
	}

	order, err := h.orderService.GetOrder(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CreateNewOrder обрабатывает HTTP POST-запрос на создание нового заказа.
// @Summary Создаёт новый заказ
// @Description Принимает данные первичной заявки и создаёт новый заказ.
// @Accept json
// @Produce json
// @Param order body orders.PrimaryOrder true "Первичная заявка для создания нового заказа"
// @Success 201 {object} orders.Order
// @Failure 400 {object} gin.H "Неверные данные запроса"
// @Router /orders [post]
func (h *OrderHandler) CreateNewOrder(c *gin.Context) {
	var primaryOrder orders.PrimaryOrder

	// Десериализуем JSON-тело запроса в структуру primaryOrder.
	if err := c.ShouldBindJSON(&primaryOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// Вызываем сервис бизнес-логики для создания заказа.
	newOrder, err := h.orderService.CreateOrder(c, &primaryOrder)
	if err != nil {
		// TODO: Check on various errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new order", "details": err.Error()})
		return
	}

	// Отправляем успешный ответ с созданным заказом и статусом 201 Created.
	c.JSON(http.StatusCreated, newOrder)
}

func (h *OrderHandler) PrescheduleOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id", "details": err.Error()})
		return
	}

	var req PrescheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if err := h.orderService.Preschedule(c, id, req.ScheduledFor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
