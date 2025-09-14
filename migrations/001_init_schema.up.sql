CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.employees (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.orders (
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

CREATE INDEX idx_orders_status ON public.orders(status);
CREATE INDEX idx_orders_employee_id ON public.orders(employee_id);
CREATE INDEX idx_orders_scheduled_for ON public.orders(scheduled_for);