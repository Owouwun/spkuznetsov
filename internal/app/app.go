package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Owouwun/spkuznetsov/internal/core/api/handlers"
	"github.com/Owouwun/spkuznetsov/internal/core/logic/auth"
	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	repository_auth "github.com/Owouwun/spkuznetsov/internal/core/repository/services/auth"
	repository_orders "github.com/Owouwun/spkuznetsov/internal/core/repository/services/orders"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbConnectionTimeout = 30 * time.Second
)

func PrepareRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	prepareOrders(router, db)
	prepareEmployees(router, db)

	return router
}

func prepareOrders(router *gin.Engine, db *gorm.DB) {
	orderRepo := repository_orders.NewOrderRepository(db)
	orderService := orders.NewOrderService(orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	apiOrders := router.Group("/api/v1/orders")
	{
		apiOrders.GET("/", orderHandler.GetAll)
		apiOrders.GET("/:id", orderHandler.GetByID)
		apiOrders.POST("", orderHandler.Create)
		apiOrders.PATCH("/:id/preschedule", orderHandler.Preschedule)
		apiOrders.PATCH("/:ordID/assign/:empID", orderHandler.Assign)
		apiOrders.PATCH("/:id/schedule", orderHandler.Schedule)
		apiOrders.PATCH("/:id/progress", orderHandler.Progress)
		apiOrders.PATCH("/:id/complete", orderHandler.Complete)
		apiOrders.PATCH("/:id/close", orderHandler.Close)
		apiOrders.PATCH("/:id/cancel", orderHandler.Cancel)
	}
}

func prepareEmployees(router *gin.Engine, db *gorm.DB) {
	authRepo := repository_auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	apiEmployees := router.Group("/api/v1/employees")
	{
		apiEmployees.GET("/", authHandler.GetEmployees)
		apiEmployees.GET("/:id", authHandler.GetEmployeeByID)
		apiEmployees.POST("", authHandler.Create)
	}
}

func PrepareDB() *gorm.DB {
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

	return db
}
