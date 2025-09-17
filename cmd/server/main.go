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
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Postgres migration
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Migrations from file
	_ "github.com/lib/pq"                                      // Register Postgres driver
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

	runMigrations(dbConn)

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
		apiGroup.GET("/orders", orderHandler.GetOrders)
		apiGroup.GET("/orders/:id", orderHandler.GetOrder)
		apiGroup.POST("/orders", orderHandler.CreateNewOrder)
		apiGroup.POST("/orders/:id/preschedule", orderHandler.PrescheduleOrder)
	}

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func runMigrations(dbConn string) {
	log.Println("Running database migrations...")

	m, err := migrate.New(
		"file://migrations", // Путь к папке с миграциями
		dbConn,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations applied successfully!")
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
			defer func() {
				err := db.Close()
				if err != nil {
					log.Fatal(err)
				}
			}()

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
