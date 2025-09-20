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

// OrderHandler содержит зависимости и логику HTTP-обработчиков.
type OrderHandler struct {
	orderService OrderService
}

func NewOrderHandler(os OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: os,
	}
}

// DTO

// PrescheduleRequest represents a request to set or update a scheduled time for an order.
// swagger:model PrescheduleRequest
type PrescheduleRequest struct {
	ScheduledFor *time.Time `json:"scheduled_for"`
}

// ProgressRequest represents data to report progress on an order.
// swagger:model ProgressRequest
type ProgressRequest struct {
	EmployeeDescription string `json:"employee_description"`
}

// CancelRequest represents reason for cancelling an order.
// swagger:model CancelRequest
type CancelRequest struct {
	CancelReason string `json:"cancel_reason"`
}

// Handlers

// GetAll godoc
// @Summary Get all orders
// @Description Returns list of all orders
// @Tags orders
// @Produce json
// @Success 200 {array} orders.Order
// @Failure 500 {object} map[string]interface{}
// @Router /orders [get]
func (h *OrderHandler) GetAll(c *gin.Context) {
	orders, err := h.orderService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetByID godoc
// @Summary Get order by ID
// @Description Get single order by UUID
// @Tags orders
// @Produce json
// @Param id path string true "Order ID" Format(uuid)
// @Success 200 {object} orders.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id} [get]
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

// Create godoc
// @Summary Create a new order
// @Description Create order with PrimaryOrder payload
// @Tags orders
// @Accept json
// @Produce json
// @Param order body orders.PrimaryOrder true "Primary order payload"
// @Success 201 {string} string "new order UUID"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var primaryOrder orders.PrimaryOrder

	if err := c.ShouldBindJSON(&primaryOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	newOrderID, err := h.orderService.Create(c, &primaryOrder)
	if err != nil {
		// TODO: Check on various errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newOrderID)
}

// Preschedule godoc
// @Summary Preschedule an order (provisional scheduling)
// @Description Set or update a provisional scheduled time for the order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID" Format(uuid)
// @Param body body PrescheduleRequest true "Preschedule payload"
// @Success 200 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/preschedule [patch]
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

// Assign godoc
// @Summary Assign employee to order
// @Description Assign employee by numeric ID to an order
// @Tags orders
// @Produce json
// @Param id path string true "Order ID" Format(uuid)
// @Param empID path int true "Employee ID"
// @Success 200 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/assign/{empID} [patch]
func (h *OrderHandler) Assign(c *gin.Context) {
	ordID, err := uuid.Parse(c.Param("id"))
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

// Schedule godoc
// @Summary Schedule an order (final scheduling)
// @Description Set the final scheduled time for an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID" Format(uuid)
// @Param body body PrescheduleRequest true "Schedule payload"
// @Success 200 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/schedule [patch]
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

// Progress godoc
// @Summary Report progress for an order
// @Description Attach employee progress/notes to an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID" Format(uuid)
// @Param body body ProgressRequest true "Progress payload"
// @Success 200 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/progress [patch]
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

// Complete godoc
// @Summary Mark order as completed
// @Tags orders
// @Produce json
// @Param id path string true "Order ID" Format(uuid)
// @Success 200 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/complete [patch]
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

// Close godoc
// @Summary Close an order
// @Tags orders
// @Produce json
// @Param id path string true "Order ID" Format(uuid)
// @Success 200 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/close [patch]
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

// Cancel godoc
// @Summary Cancel an order
// @Description Cancel with a reason
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID" Format(uuid)
// @Param body body CancelRequest true "Cancel payload"
// @Success 200 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/cancel [patch]
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
