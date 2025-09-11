package repository_postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/seagumineko/spkuznetsov/internal/core/logic/auth"
	"github.com/seagumineko/spkuznetsov/internal/core/logic/requests"
	"github.com/seagumineko/spkuznetsov/internal/core/repository"
	deterrs "github.com/seagumineko/spkuznetsov/internal/errors"
)

type PostgresRequestRepository struct {
	db *sql.DB
}

// TODO: Add error handling
func NewRequestRepository(db *sql.DB) repository.RequestRepository {
	return &PostgresRequestRepository{db: db}
}

func (r *PostgresRequestRepository) CreateRequest(ctx context.Context, req *requests.Request) (int64, error) {
	query := `
		INSERT INTO requests (
			client_name, 
			client_phone, 
			address, 
			client_description, 
			public_link, 
			employee_id, 
			cancel_reason, 
			status, 
			employee_description, 
			scheduled_for
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	var employeeID interface{}
	if req.Employee != nil {
		employeeID = req.Employee.ID
	} else {
		employeeID = nil
	}

	var cancelReason interface{}
	if req.CancelReason != nil {
		cancelReason = *req.CancelReason
	} else {
		cancelReason = nil
	}

	var scheduledFor interface{}
	if req.ScheduledFor != nil {
		scheduledFor = *req.ScheduledFor
	} else {
		scheduledFor = nil
	}

	var id int64
	err := r.db.QueryRowContext(ctx, query,
		req.ClientName,
		req.ClientPhone,
		req.Address,
		req.ClientDescription,
		req.PublicLink,
		employeeID,
		cancelReason,
		req.Status,
		req.EmployeeDescription,
		scheduledFor,
	).Scan(&id)

	if err != nil {
		return 0, deterrs.NewDetErr(
			deterrs.QueryInsertFailed,
			deterrs.WithOriginalError(err),
		)
	}

	return id, nil
}

func (r *PostgresRequestRepository) UpdateRequest(ctx context.Context, id int64, req *requests.Request) error {
	query := `
		UPDATE requests 
		SET 
			client_name = $1,
			client_phone = $2,
			address = $3,
			client_description = $4,
			public_link = $5,
			employee_id = $6,
			cancel_reason = $7,
			status = $8,
			employee_description = $9,
			scheduled_for = $10
		WHERE id = $11
	`

	var employeeID interface{}
	if req.Employee != nil {
		employeeID = req.Employee.ID
	} else {
		employeeID = nil
	}

	var cancelReason interface{}
	if req.CancelReason != nil {
		cancelReason = *req.CancelReason
	} else {
		cancelReason = nil
	}

	var scheduledFor interface{}
	if req.ScheduledFor != nil {
		scheduledFor = *req.ScheduledFor
	} else {
		scheduledFor = nil
	}

	result, err := r.db.ExecContext(ctx, query,
		req.ClientName,
		req.ClientPhone,
		req.Address,
		req.ClientDescription,
		req.PublicLink,
		employeeID,
		cancelReason,
		req.Status,
		req.EmployeeDescription,
		scheduledFor,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("request with id %d not found", id)
	}

	return nil
}

func (r *PostgresRequestRepository) GetRequest(ctx context.Context, id int64) (*requests.Request, error) {
	query := `
		SELECT 
			id,
			client_name,
			client_phone,
			address,
			client_description,
			public_link,
			employee_id,
			cancel_reason,
			status,
			employee_description,
			scheduled_for
		FROM requests 
		WHERE id = $1
	`

	req := &requests.Request{}
	var (
		employeeID   sql.NullInt64
		cancelReason sql.NullString
		scheduledFor sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&req.ID,
		&req.ClientName,
		&req.ClientPhone,
		&req.Address,
		&req.ClientDescription,
		&req.PublicLink,
		&employeeID,
		&cancelReason,
		&req.Status,
		&req.EmployeeDescription,
		&scheduledFor,
	)

	if err != nil {
		return nil, deterrs.NewDetErr(
			deterrs.QuerySelectFailed,
			deterrs.WithOriginalError(err),
		)
	}

	if employeeID.Valid {
		employee, err := r.GetEmployeeByID(ctx, employeeID.Int64)
		if err != nil {
			return nil, fmt.Errorf("failed to get employee: %w", err)
		}
		req.Employee = employee
	}

	if cancelReason.Valid {
		req.CancelReason = &cancelReason.String
	}

	if scheduledFor.Valid {
		req.ScheduledFor = &scheduledFor.Time
	}

	return req, nil
}

func (r *PostgresRequestRepository) GetEmployeeByID(ctx context.Context, id int64) (*auth.Employee, error) {
	query := `SELECT id, name FROM employees WHERE id = $1`

	employee := &auth.Employee{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&employee.ID,
		&employee.Name,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	return employee, nil
}
