package worker

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"project_abc/backend/internal/database"
)

type SMSConfig struct {
	APIToken    string
	Sender      string
	FrontendURL string
}

func StartSMSReminderJob(db *sql.DB, cfg SMSConfig) {
	if cfg.Sender == "" {
		cfg.Sender = "Test"
	}

	mode := "live"
	if cfg.APIToken == "" {
		mode = "simulated (no SMSAPI_TOKEN)"
	}

	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		slog.Info("sms worker started", "interval", "1m", "mode", mode)

		run(db, cfg)

		for range ticker.C {
			run(db, cfg)
		}
	}()
}

func run(db *sql.DB, cfg SMSConfig) {
	appointments, err := database.GetScheduledAppointments(db)
	if err != nil {
		slog.Error("sms worker: failed to query appointments", "error", err)
		return
	}

	if len(appointments) == 0 {
		slog.Info("sms worker: no pending appointments to notify")
		return
	}

	for _, appt := range appointments {
		magicLink := fmt.Sprintf("%s/#/p/%s", cfg.FrontendURL, appt.ID)
		message := fmt.Sprintf(
			"Hello %s, your appointment for %s is coming up. Confirm or cancel here: %s",
			appt.CustomerFirstName, appt.ServiceName, magicLink,
		)

		err := sendSMS(cfg, appt.CustomerPhone, message)
		if err != nil {
			slog.Error("sms worker: send failed",
				"appointment_id", appt.ID,
				"phone", appt.CustomerPhone,
				"error", err,
			)

			if dbErr := database.CreateNotification(db, appt.ID, "reminder", "failed"); dbErr != nil {
				slog.Error("sms worker: failed to record notification", "error", dbErr)
			}
			continue
		}

		if dbErr := database.CreateNotification(db, appt.ID, "reminder", "sent"); dbErr != nil {
			slog.Error("sms worker: failed to record notification", "error", dbErr)
		}

		slog.Info("sms worker: notification sent",
			"appointment_id", appt.ID,
			"phone", appt.CustomerPhone,
		)
	}

	slog.Info("sms worker: batch complete", "count", len(appointments))
}

func sendSMS(cfg SMSConfig, phone string, message string) error {
	// Fallback to simulation if no token is configured.
	if cfg.APIToken == "" {
		slog.Info("sms worker: SMS (simulated)",
			"to", phone,
			"message", message,
		)
		return nil
	}

	data := url.Values{
		"to":      {phone},
		"from":    {cfg.Sender},
		"message": {message},
		"format":  {"json"},
	}

	req, err := http.NewRequest("POST", "https://api.smsapi.pl/sms.do", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+cfg.APIToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("SMSAPI error %d: %s", resp.StatusCode, string(body))
	}

	slog.Info("sms worker: SMSAPI response",
		"status", resp.StatusCode,
		"body", string(body),
	)

	return nil
}
