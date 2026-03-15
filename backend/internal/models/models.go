package models

import "time"

type Customer struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type Appointment struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	CustomerID      string    `json:"customer_id"`
	AppointmentTime time.Time `json:"appointment_time"`
	Status          string    `json:"status"`
	ServiceName     string    `json:"service_name"`
	CreatedAt       time.Time `json:"created_at"`

	// Joined from customers table
	CustomerFirstName string `json:"customer_first_name,omitempty"`
	CustomerLastName  string `json:"customer_last_name,omitempty"`
	CustomerPhone     string `json:"customer_phone,omitempty"`
}

type CreateAppointmentInput struct {
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Phone           string    `json:"phone"`
	ServiceName     string    `json:"service_name"`
	AppointmentTime time.Time `json:"appointment_time"`
}
