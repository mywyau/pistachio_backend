package jobs

import (
    "context"
    "encoding/json"
    "net/http"
    // "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

func CreateJobHandler(db *pgxpool.Pool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req CreateJobRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid JSON", http.StatusBadRequest)
            return
        }

        ctx := context.Background()

        // 1️⃣ Create customer
        customerID := uuid.New()

        _, err := db.Exec(ctx,
            `INSERT INTO customers (id, name, email, phone, address, created_at)
             VALUES ($1, $2, $3, $4, $5, NOW())`,
            customerID,
            req.Customer.Name,
            req.Customer.Email,
            req.Customer.Phone,
            req.Customer.Address,
        )

        if err != nil {
            http.Error(w, "failed to create customer: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // 2️⃣ Create job
        jobID := uuid.New()

        _, err = db.Exec(ctx,
            `INSERT INTO jobs (id, customer_id, title, description, estimate, status, created_at, updated_at)
             VALUES ($1, $2, $3, $4, $5, 'new', NOW(), NOW())`,
            jobID,
            customerID,
            req.Title,
            req.Description,
            req.Estimate,
        )

        if err != nil {
            http.Error(w, "failed to create job: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // 3️⃣ Return response
        resp := CreateJobResponse{
            JobID:      jobID,
            CustomerID: customerID,
            Title:      req.Title,
            Status:     "new",
        }

        json.NewEncoder(w).Encode(resp)
    }
}
