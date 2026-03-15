package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"project_abc/backend/internal/database"
	"project_abc/backend/internal/models"
	"project_abc/backend/internal/worker"
)

func authMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token := strings.TrimPrefix(header, "Bearer ")

			if token == "" || token == header || token != secret {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	projectName := os.Getenv("PROJECT_NAME")
	if projectName == "" {
		projectName = "project_abc"
	}

	adminSecret := os.Getenv("ADMIN_SECRET")
	if adminSecret == "" {
		slog.Warn("ADMIN_SECRET is not set, admin routes will reject all requests")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database
	db, err := database.Connect()
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		slog.Error("database migration failed", "error", err)
		os.Exit(1)
	}

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:         300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"project": projectName,
		})
	})

	r.Route("/api", func(r chi.Router) {
		// Public routes — used by Magic Link (no auth)
		r.Get("/appointments/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			appt, err := database.GetAppointmentByID(db, id)
			if err != nil {
				slog.Error("failed to get appointment", "error", err)
				http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(appt)
		})

		r.Patch("/appointments/{id}/status", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")

			var body struct {
				Status string `json:"status"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
				return
			}

			if body.Status != "confirmed" && body.Status != "cancelled" {
				http.Error(w, `{"error":"status must be confirmed or cancelled"}`, http.StatusBadRequest)
				return
			}

			if err := database.UpdateAppointmentStatus(db, id, body.Status); err != nil {
				slog.Error("failed to update appointment status", "error", err)
				http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": body.Status})
		})

		// Protected routes — Dashboard (require ADMIN_SECRET)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware(adminSecret))

			r.Get("/appointments", func(w http.ResponseWriter, r *http.Request) {
				appointments, err := database.GetAllAppointments(db)
				if err != nil {
					slog.Error("failed to get appointments", "error", err)
					http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(appointments)
			})

			r.Post("/appointments", func(w http.ResponseWriter, r *http.Request) {
				var input models.CreateAppointmentInput
				if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
					http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
					return
				}

				if input.FirstName == "" || input.LastName == "" || input.Phone == "" ||
					input.ServiceName == "" || input.AppointmentTime.IsZero() {
					http.Error(w, `{"error":"all fields are required"}`, http.StatusBadRequest)
					return
				}

				if err := database.CreateAppointmentWithCustomer(db, input); err != nil {
					slog.Error("failed to create appointment", "error", err)
					http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]string{"status": "created"})
			})
		})
	})

	// SMS Worker
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	worker.StartSMSReminderJob(db, worker.SMSConfig{
		APIToken:    os.Getenv("SMSAPI_TOKEN"),
		Sender:      os.Getenv("SMS_SENDER"),
		FrontendURL: frontendURL,
	})

	slog.Info("starting server", "port", port, "project", projectName)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
