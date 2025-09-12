package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Owouwun/spkuznetsov/internal/core/api/handlers"
	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
)

type orderServiceImpl struct{}

func (s *orderServiceImpl) CreateNewOrder(pord *orders.PrimaryOrder) (*orders.Order, error) {
	return pord.CreateNewOrder()
}

func main() {
	router := gin.Default()

	// DI
	orderService := &orderServiceImpl{}
	orderHandler := handlers.NewOrderHandler(orderService)

	// Routing
	router.POST("/orders", orderHandler.CreateNewOrder)

	// Launch
	if err := router.Run(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}
}
