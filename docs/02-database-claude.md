# Database Schema (PostgreSQL)

Use UUIDs as primary keys. We need migration scripts (e.g., using `golang-migrate`).

1. **tenants**: `id` (UUID), `name` (String), `email` (String), `phone` (String), `created_at` (Timestamp).
2. **customers**: `id` (UUID), `tenant_id` (FK), `first_name` (String), `last_name` (String), `phone` (String), `created_at` (Timestamp).
3. **appointments**: `id` (UUID), `tenant_id` (FK), `customer_id` (FK), `appointment_time` (Timestamp), `status` (Enum: 'scheduled', 'confirmed', 'cancelled', 'completed'), `service_name` (String), `created_at` (Timestamp).
4. **notifications**: `id` (UUID), `appointment_id` (FK), `type` (Enum: 'reminder', 'confirmation_request'), `status` (Enum: 'pending', 'sent', 'failed'), `scheduled_for` (Timestamp), `sent_at` (Timestamp).