package main

import (
	"log"

	"github.com/Owouwun/spkuznetsov/internal/app"
)

func main() {
	db := app.PrepareDB()
	router := app.PrepareRouter(db)

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
