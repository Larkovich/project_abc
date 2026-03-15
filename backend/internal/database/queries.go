package database

import (
	"database/sql"
	"fmt"

	"project_abc/backend/internal/models"
)

func GetAllAppointments(db *sql.DB) ([]models.Appointment, error) {
	query := `
		SELECT
			a.id, a.tenant_id, a.customer_id,
			a.appointment_time, a.status, a.service_name, a.created_at,
			c.first_name, c.last_name
		FROM appointments a
		JOIN customers c ON c.id = a.customer_id
		ORDER BY a.appointment_time ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllAppointments: %w", err)
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var a models.Appointment
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.CustomerID,
			&a.AppointmentTime, &a.Status, &a.ServiceName, &a.CreatedAt,
			&a.CustomerFirstName, &a.CustomerLastName,
		); err != nil {
			return nil, fmt.Errorf("GetAllAppointments scan: %w", err)
		}
		appointments = append(appointments, a)
	}

	return appointments, rows.Err()
}

func GetAppointmentByID(db *sql.DB, id string) (models.Appointment, error) {
	query := `
		SELECT
			a.id, a.tenant_id, a.customer_id,
			a.appointment_time, a.status, a.service_name, a.created_at,
			c.first_name, c.last_name
		FROM appointments a
		JOIN customers c ON c.id = a.customer_id
		WHERE a.id = $1
	`

	var a models.Appointment
	err := db.QueryRow(query, id).Scan(
		&a.ID, &a.TenantID, &a.CustomerID,
		&a.AppointmentTime, &a.Status, &a.ServiceName, &a.CreatedAt,
		&a.CustomerFirstName, &a.CustomerLastName,
	)
	if err != nil {
		return a, fmt.Errorf("GetAppointmentByID: %w", err)
	}

	return a, nil
}

func UpdateAppointmentStatus(db *sql.DB, id string, status string) error {
	query := `UPDATE appointments SET status = $1 WHERE id = $2`

	res, err := db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("UpdateAppointmentStatus: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("UpdateAppointmentStatus: appointment %s not found", id)
	}

	return nil
}

func GetScheduledAppointments(db *sql.DB) ([]models.Appointment, error) {
	query := `
		SELECT
			a.id, a.tenant_id, a.customer_id,
			a.appointment_time, a.status, a.service_name, a.created_at,
			c.first_name, c.last_name
		FROM appointments a
		JOIN customers c ON c.id = a.customer_id
		WHERE a.status = 'scheduled'
		  AND a.appointment_time > now()
		ORDER BY a.appointment_time ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetScheduledAppointments: %w", err)
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var a models.Appointment
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.CustomerID,
			&a.AppointmentTime, &a.Status, &a.ServiceName, &a.CreatedAt,
			&a.CustomerFirstName, &a.CustomerLastName,
		); err != nil {
			return nil, fmt.Errorf("GetScheduledAppointments scan: %w", err)
		}
		appointments = append(appointments, a)
	}

	return appointments, rows.Err()
}

const defaultTenantID = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

func CreateAppointmentWithCustomer(db *sql.DB, input models.CreateAppointmentInput) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("CreateAppointmentWithCustomer begin: %w", err)
	}
	defer tx.Rollback()

	// Find existing customer by phone for this tenant, or create a new one.
	var customerID string
	err = tx.QueryRow(
		`SELECT id FROM customers WHERE phone = $1 AND tenant_id = $2`,
		input.Phone, defaultTenantID,
	).Scan(&customerID)

	if err == sql.ErrNoRows {
		err = tx.QueryRow(
			`INSERT INTO customers (tenant_id, first_name, last_name, phone)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			defaultTenantID, input.FirstName, input.LastName, input.Phone,
		).Scan(&customerID)
	}
	if err != nil {
		return fmt.Errorf("CreateAppointmentWithCustomer customer: %w", err)
	}

	// Insert the appointment.
	_, err = tx.Exec(
		`INSERT INTO appointments (tenant_id, customer_id, appointment_time, service_name)
		 VALUES ($1, $2, $3, $4)`,
		defaultTenantID, customerID, input.AppointmentTime, input.ServiceName,
	)
	if err != nil {
		return fmt.Errorf("CreateAppointmentWithCustomer appointment: %w", err)
	}

	return tx.Commit()
}
