package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/Owouwun/spkuznetsov/internal/core/api/handlers"
	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	repository_orders "github.com/Owouwun/spkuznetsov/internal/core/repository/orders"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbConnectionTimeout = 30 * time.Second
)

func main() {
	dbConn := os.Getenv("DATABASE_CONN")
	if dbConn == "" {
		log.Fatal("DATABASE_CONN environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbConnectionTimeout)
	defer cancel()

	if err := waitForDBReady(ctx, dbConn); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Connecting to the PostgreSQL database...")
	db, err := gorm.Open(postgres.Open(dbConn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	orderRepo := repository_orders.NewOrderRepository(db)

	orderService := orders.NewOrderService(orderRepo)

	orderHandler := handlers.NewOrderHandler(orderService)

	router := gin.Default()

	apiGroup := router.Group("/api/v1")
	{
		apiGroup.POST("/orders", orderHandler.CreateNewOrder)
		apiGroup.GET("/orders/:id", orderHandler.GetOrder)
	}

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func waitForDBReady(ctx context.Context, dbConn string) error {
	log.Println("Waiting for database to be ready...")

	done := make(chan error)

	go func() {
		for {
			db, err := sql.Open("postgres", dbConn)
			if err != nil {
				done <- err
				return
			}
			defer db.Close()

			if err := db.Ping(); err == nil {
				done <- nil
				return
			}

			// Wait till the next try
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done(): // Timeout
		return ctx.Err()
	}
}
