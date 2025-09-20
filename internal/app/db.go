package app

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Postgres migration
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Migrations from file
	_ "github.com/lib/pq"                                      // Register Postgres driver
)

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
