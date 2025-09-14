package repository_postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"testing"

	"github.com/Owouwun/spkuznetsov/internal/testutils"
	"github.com/Owouwun/spkuznetsov/pkg/logger"
	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	migpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	containerPort = "5432"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	ctx := context.Background()

	dsnProvider := func(host string, port nat.Port) string {
		return fmt.Sprintf("host=%s port=%d user=user password=password dbname=testdb sslmode=disable", host, port.Int())
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{containerPort + "/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(1),
			wait.ForListeningPort(containerPort),
			wait.ForSQL(nat.Port(containerPort), "postgres", dsnProvider),
		),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	port, err := postgresContainer.MappedPort(ctx, containerPort)
	if err != nil {
		t.Fatal(err)
	}

	finalDSN := dsnProvider(host, port)
	sqlDB, err := sql.Open("postgres", finalDSN)
	if err != nil {
		t.Fatal(err)
	}

	err = runTestMigrations(sqlDB)
	if err != nil {
		t.Fatal(err)
		err = sqlDB.Close()
		if err != nil {
			t.Log(err)
		}
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "public.",
			SingularTable: false,
		},
	})
	if err != nil {
		t.Fatal(err)
		err = sqlDB.Close()
		if err != nil {
			t.Log(err)
		}
	}

	cleanup := func() {
		err = sqlDB.Close()
		if err != nil {
			t.Log(err)
		}
		logger.LogIfErr(t, "error while terminating container: %v",
			postgresContainer.Terminate, ctx,
		)
	}

	return gormDB, cleanup
}

func runTestMigrations(db *sql.DB) error {
	driver, err := migpostgres.WithInstance(db, &migpostgres.Config{})
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

func TestOrderRepository_CreateOrder(t *testing.T) {
	gormDB, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(gormDB)
	newOrder := testutils.NewTestOrder()

	ordID, err := repo.CreateOrder(context.Background(), newOrder)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	request, err := repo.GetOrder(context.Background(), ordID)
	if err != nil {
		t.Fatalf("Failed to get request: %v", err)
	}

	testutils.ValidateOrder(t, newOrder, request)
}

func TestOrderRepository_UpdateOrder(t *testing.T) {
	gormDB, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(gormDB)

	createdOrder := testutils.NewTestOrder()
	ordID, err := repo.CreateOrder(context.Background(), createdOrder)
	if err != nil {
		t.Fatalf("Failed to create request for update: %v", err)
	}

	updatedOrder := createdOrder
	updatedOrder.ID = ordID
	updatedOrder.Address = "Updated GORM Test Address"

	err = repo.UpdateOrder(context.Background(), updatedOrder)
	if err != nil {
		t.Fatalf("Failed to update request: %v", err)
	}

	resultOrder, err := repo.GetOrder(context.Background(), ordID)
	if err != nil {
		t.Fatalf("Failed to get updated request: %v", err)
	}

	testutils.ValidateOrder(t, updatedOrder, resultOrder)
}
