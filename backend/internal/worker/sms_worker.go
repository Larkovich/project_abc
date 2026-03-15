package worker

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"project_abc/backend/internal/database"
)

func StartSMSReminderJob(db *sql.DB, frontendURL string) {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		slog.Info("sms worker started", "interval", "1m")

		// Run immediately on startup, then on every tick.
		run(db, frontendURL)

		for range ticker.C {
			run(db, frontendURL)
		}
	}()
}

func run(db *sql.DB, frontendURL string) {
	appointments, err := database.GetScheduledAppointments(db)
	if err != nil {
		slog.Error("sms worker: failed to query appointments", "error", err)
		return
	}

	if len(appointments) == 0 {
		slog.Info("sms worker: no scheduled appointments to notify")
		return
	}

	for _, appt := range appointments {
		magicLink := fmt.Sprintf("%s/#/p/%s", frontendURL, appt.ID)

		slog.Info("sms worker: sending SMS (simulated)",
			"to", appt.CustomerFirstName+" "+appt.CustomerLastName,
			"appointment_id", appt.ID,
			"service", appt.ServiceName,
			"message", fmt.Sprintf(
				"Hello %s, your appointment for %s is coming up. Confirm or cancel here: %s",
				appt.CustomerFirstName, appt.ServiceName, magicLink,
			),
		)
	}

	slog.Info("sms worker: batch complete", "count", len(appointments))
}
