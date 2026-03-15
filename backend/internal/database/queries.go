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

func CreateAppointment(db *sql.DB, appt models.Appointment) error {
	query := `
		INSERT INTO appointments (tenant_id, customer_id, appointment_time, service_name)
		VALUES ($1, $2, $3, $4)
	`

	_, err := db.Exec(query, appt.TenantID, appt.CustomerID, appt.AppointmentTime, appt.ServiceName)
	if err != nil {
		return fmt.Errorf("CreateAppointment: %w", err)
	}

	return nil
}
