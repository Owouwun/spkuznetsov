package handlers

import (
	"net/http"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	"github.com/gin-gonic/gin"
)

// Задаёт методы бизнес-логики
type OrderService interface {
	CreateNewOrder(pord *orders.PrimaryOrder) (*orders.Order, error)
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

type CreateNewOrderResponseBody struct {
	Order *orders.Order
	Err   error
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
	newOrder, err := h.orderService.CreateNewOrder(&primaryOrder)
	if err != nil {
		// TODO: Check on various errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new order"})
		return
	}

	// Отправляем успешный ответ с созданным заказом и статусом 201 Created.
	c.JSON(http.StatusCreated, &CreateNewOrderResponseBody{newOrder, nil})
}
