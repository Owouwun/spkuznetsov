CREATE TABLE IF NOT EXISTS employees (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);
INSERT INTO employees (name, created_at, updated_at) VALUES ('Петр Петров', NOW(), NOW());

CREATE TABLE requests (
    id BIGSERIAL PRIMARY KEY,
    client_name TEXT NOT NULL,
    client_phone TEXT NOT NULL,
    address TEXT NOT NULL,
    client_description TEXT,
    public_link TEXT,
    employee_id BIGINT REFERENCES employees(id) ON DELETE SET NULL,
    cancel_reason TEXT,
    status INTEGER NOT NULL,
    employee_description TEXT,
    scheduled_for TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);
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
    scheduled_for,
    created_at,
    updated_at
) VALUES (
    'Иван Иванов',
    '+71112223344',
    'ул. Тестовая, д. 1',
    'Что-то сломалось',
    'abracadabra',
    1,
    NULL,
    1,
    NULL,
    NULL,
    NOW(),
    NOW()
);

CREATE INDEX idx_requests_status ON requests(status);
CREATE INDEX idx_requests_employee_id ON requests(employee_id);
CREATE INDEX idx_requests_scheduled_for ON requests(scheduled_for);