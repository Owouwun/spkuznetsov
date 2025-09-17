package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/Owouwun/spkuznetsov/internal/core/api/handlers"
	"github.com/Owouwun/spkuznetsov/internal/core/logic/auth"
	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	repository_auth "github.com/Owouwun/spkuznetsov/internal/core/repository/services/auth"
	repository_orders "github.com/Owouwun/spkuznetsov/internal/core/repository/services/orders"
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
	db := prepareDB()
	router := prepareRouter(db)

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func prepareRouter(db *gorm.DB) *gin.Engine {
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
		apiOrders.POST("/preschedule/:id", orderHandler.Preschedule)
		apiOrders.POST("/assign/:ordID/:empID", orderHandler.Assign)
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

func prepareDB() *gorm.DB {
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
