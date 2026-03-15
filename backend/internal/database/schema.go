package database

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func Migrate(db *sql.DB) error {
	// Wrap enum creation in DO blocks so they don't fail on re-runs.
	idempotentSchema := `
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'appointment_status') THEN
        CREATE TYPE appointment_status AS ENUM ('scheduled', 'confirmed', 'cancelled', 'completed');
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_type') THEN
        CREATE TYPE notification_type AS ENUM ('reminder', 'confirmation_request');
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_status') THEN
        CREATE TYPE notification_status AS ENUM ('pending', 'sent', 'failed');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    customer_id UUID NOT NULL REFERENCES customers(id),
    appointment_time TIMESTAMPTZ NOT NULL,
    status appointment_status NOT NULL DEFAULT 'scheduled',
    service_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL REFERENCES appointments(id),
    type notification_type NOT NULL,
    status notification_status NOT NULL DEFAULT 'pending',
    scheduled_for TIMESTAMPTZ NOT NULL,
    sent_at TIMESTAMPTZ
);
`

	if _, err := db.Exec(idempotentSchema); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	slog.Info("database schema migrated")

	if err := seed(db); err != nil {
		return fmt.Errorf("seed: %w", err)
	}

	return nil
}

func seed(db *sql.DB) error {
	// Only seed if the tenants table is empty.
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM tenants").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		slog.Info("seed skipped, data already exists")
		return nil
	}

	seedSQL := `
WITH t AS (
    INSERT INTO tenants (id, name, email, phone)
    VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Bella Beauty Studio', 'contact@bellastudio.com', '+48500100200')
    RETURNING id
),
c1 AS (
    INSERT INTO customers (id, tenant_id, first_name, last_name, phone)
    VALUES ('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', (SELECT id FROM t), 'Anna', 'Kowalska', '+48501111222')
    RETURNING id
),
c2 AS (
    INSERT INTO customers (id, tenant_id, first_name, last_name, phone)
    VALUES ('b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', (SELECT id FROM t), 'Maria', 'Nowak', '+48502333444')
    RETURNING id
)
INSERT INTO appointments (tenant_id, customer_id, appointment_time, status, service_name) VALUES
    ((SELECT id FROM t), (SELECT id FROM c1), now() + interval '1 day',   'scheduled',  'Haircut & Styling'),
    ((SELECT id FROM t), (SELECT id FROM c2), now() + interval '2 days',  'confirmed',  'Manicure'),
    ((SELECT id FROM t), (SELECT id FROM c1), now() + interval '3 days',  'scheduled',  'Color & Highlights'),
    ((SELECT id FROM t), (SELECT id FROM c2), now() - interval '1 day',   'completed',  'Facial Treatment');
`

	if _, err := db.Exec(seedSQL); err != nil {
		return err
	}

	slog.Info("seed data inserted")
	return nil
}
