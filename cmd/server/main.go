package main

import (
	"log"

	_ "github.com/Owouwun/spkuznetsov/cmd/docs"
	"github.com/Owouwun/spkuznetsov/internal/app"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Orders Management Service
// @version         1.0
// @description     A web service for managing and tracking service requests.
// @contact.name    Ivan Kuznetsov
// @contact.email   kuznetsovivangio@gmail.com
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	db := app.PrepareDB()
	router := app.PrepareRouter(db)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
