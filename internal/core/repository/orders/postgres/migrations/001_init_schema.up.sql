CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS employees (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
INSERT INTO employees (name) VALUES ('Петр Петров');

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_name TEXT NOT NULL,
    client_phone TEXT NOT NULL,
    address TEXT NOT NULL,
    client_description TEXT,
    employee_id BIGINT REFERENCES employees(id) ON DELETE SET NULL,
    cancel_reason TEXT,
    status INTEGER NOT NULL,
    employee_description TEXT,
    scheduled_for TIMESTAMP WITH TIME ZONE
);
INSERT INTO orders (
    client_name,
    client_phone,
    address,
    client_description,
    employee_id,
    cancel_reason,
    status,
    employee_description,
    scheduled_for
) VALUES (
    'Тесть Тестя',
    '+71234567890',
    'ул. Тестовая, д. 1',
    'Что-то сломалось',
    1,
    NULL,
    1,
    NULL,
    NULL
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_employee_id ON orders(employee_id);
CREATE INDEX idx_orders_scheduled_for ON orders(scheduled_for);