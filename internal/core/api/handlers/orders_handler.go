package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Задаёт методы бизнес-логики
type OrderService interface {
	Create(ctx context.Context, pord *orders.PrimaryOrder) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*orders.Order, error)
	GetAll(ctx context.Context) ([]*orders.Order, error)
	Preschedule(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error
	Assign(ctx context.Context, id uuid.UUID, empID uint) error
	Schedule(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error
	Progress(ctx context.Context, id uuid.UUID, empDescr string) error
	Complete(ctx context.Context, id uuid.UUID) error
	Close(ctx context.Context, id uuid.UUID) error
	Cancel(ctx context.Context, id uuid.UUID, reason string) error
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

func (h *OrderHandler) GetAll(c *gin.Context) {
	orders, err := h.orderService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id", "details": err.Error()})
		return
	}

	order, err := h.orderService.GetByID(c, id)
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
func (h *OrderHandler) Create(c *gin.Context) {
	var primaryOrder orders.PrimaryOrder

	// Десериализуем JSON-тело запроса в структуру primaryOrder.
	if err := c.ShouldBindJSON(&primaryOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// Вызываем сервис бизнес-логики для создания заказа.
	newOrderID, err := h.orderService.Create(c, &primaryOrder)
	if err != nil {
		// TODO: Check on various errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new order", "details": err.Error()})
		return
	}

	// Отправляем успешный ответ с ID созданного заказа и статусом 201 Created.
	c.JSON(http.StatusCreated, newOrderID)
}

func (h *OrderHandler) Preschedule(c *gin.Context) {
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

func (h *OrderHandler) Assign(c *gin.Context) {
	ordID, err := uuid.Parse(c.Param("ordID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id", "details": err.Error()})
		return
	}

	empID, err := strconv.ParseUint(c.Param("empID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id", "details": err.Error()})
		return
	}

	if err := h.orderService.Assign(c, ordID, uint(empID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *OrderHandler) Schedule(c *gin.Context) {
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

	if err := h.orderService.Schedule(c, id, req.ScheduledFor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

type ProgressRequest struct {
	EmployeeDescription string `json:"employee_description"`
}

func (h *OrderHandler) Progress(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id", "details": err.Error()})
		return
	}

	var req ProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if err := h.orderService.Progress(c, id, req.EmployeeDescription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *OrderHandler) Complete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id", "details": err.Error()})
		return
	}

	if err := h.orderService.Complete(c, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *OrderHandler) Close(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id", "details": err.Error()})
		return
	}

	if err := h.orderService.Close(c, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

type CancelRequest struct {
	CancelReason string `json:"cancel_reason"`
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id", "details": err.Error()})
		return
	}

	var req CancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if err := h.orderService.Cancel(c, id, req.CancelReason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
