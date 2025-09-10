CREATE TYPE status_enum AS ENUM (
    'New',
    'Prescheduled',
    'Assigned',
    'Scheduled',
    'InProgress',
    'Done',
    'Paid',
    'Canceled'
);

CREATE TABLE IF NOT EXISTS employees (
    id BIGSERIAL PRIMARY KEY
);

CREATE TABLE requests (
    id BIGSERIAL PRIMARY KEY,
    client_name TEXT NOT NULL,
    client_phone TEXT NOT NULL,
    address TEXT NOT NULL,
    client_description TEXT,
    public_link TEXT,
    employee_id BIGINT REFERENCES employees(id) ON DELETE SET NULL,
    cancel_reason TEXT,
    status status_enum NOT NULL DEFAULT 'New',
    employee_description TEXT,
    scheduled_for TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_requests_status ON requests(status);
CREATE INDEX idx_requests_employee_id ON requests(employee_id);
CREATE INDEX idx_requests_scheduled_for ON requests(scheduled_for);