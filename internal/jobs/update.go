package jobs

import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type UpdateStatusRequest struct {
    Status string `json:"status"`
}

type UpdateStatusResponse struct {
    JobID  uuid.UUID `json:"job_id"`
    Status string    `json:"status"`
}

// Valid allowed statuses for jobs
var ValidStatuses = map[string]bool{
    "new":           true,
    "in_progress":   true,
    "waiting_parts": true,
    "completed":     true,
    "invoiced":      true,
}

func UpdateJobStatusHandler(db *pgxpool.Pool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract job ID from URL
        jobIDParam := chi.URLParam(r, "id")
        jobID, err := uuid.Parse(jobIDParam)
        if err != nil {
            http.Error(w, "invalid job id", http.StatusBadRequest)
            return
        }

        // Parse request body
        var req UpdateStatusRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid JSON", http.StatusBadRequest)
            return
        }

        // Validate requested status
        if !ValidStatuses[req.Status] {
            http.Error(w, "invalid job status", http.StatusBadRequest)
            return
        }

        ctx := context.Background()

        // Update job status
        result, err := db.Exec(ctx,
            `UPDATE jobs 
             SET status = $1, updated_at = $2
             WHERE id = $3`,
            req.Status,
            time.Now(),
            jobID,
        )

        if err != nil {
            http.Error(w, "failed to update job: "+err.Error(), http.StatusInternalServerError)
            return
        }

        if result.RowsAffected() == 0 {
            http.Error(w, "job not found", http.StatusNotFound)
            return
        }

        // Return success response
        resp := UpdateStatusResponse{
            JobID:  jobID,
            Status: req.Status,
        }

        json.NewEncoder(w).Encode(resp)
    }
}
