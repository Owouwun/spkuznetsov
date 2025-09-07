package repository_postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"testing"

	"github.com/Owouwun/ipkuznetsov/internal/core/logic/requests"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	containerPort = "5432"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunPostgresContainer() (context.Context, testcontainers.Container, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{containerPort + "/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(containerPort),
		),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	return ctx, postgresContainer, nil
}

func runTestMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func TestRequestRepository_CreateRequest(t *testing.T) {
	ctx, postgresContainer, err := RunPostgresContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer postgresContainer.Terminate(ctx)

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	port, err := postgresContainer.MappedPort(ctx, containerPort)
	if err != nil {
		t.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=user password=password dbname=testdb sslmode=disable", host, port.Int())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = runTestMigrations(db)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewRequestRepository(db)
	newRequest := requests.NewTestRequest()

	requestID, err := repo.CreateRequest(ctx, newRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	request, err := repo.GetRequest(ctx, requestID)
	if err != nil {
		t.Fatalf("Failed to get request: %v", err)
	}

	requests.ValidateRequest(t, newRequest, request)
}
