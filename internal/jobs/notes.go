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

func CreateNoteHandler(db *pgxpool.Pool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        jobIDParam := chi.URLParam(r, "id")
        jobID, err := uuid.Parse(jobIDParam)
        if err != nil {
            http.Error(w, "invalid job id", http.StatusBadRequest)
            return
        }

        var req CreateNoteRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid JSON", http.StatusBadRequest)
            return
        }

        if req.Text == "" {
            http.Error(w, "note text cannot be empty", http.StatusBadRequest)
            return
        }

        ctx := context.Background()

        // 1️⃣ Ensure job exists
        var exists bool
        err = db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM jobs WHERE id = $1)`, jobID).Scan(&exists)
        if err != nil || !exists {
            http.Error(w, "job not found", http.StatusNotFound)
            return
        }

        // 2️⃣ Insert note
        noteID := uuid.New()

        _, err = db.Exec(ctx,
            `INSERT INTO job_notes (id, job_id, text, created_at)
             VALUES ($1, $2, $3, NOW())`,
            noteID,
            jobID,
            req.Text,
        )

        if err != nil {
            http.Error(w, "failed to insert note: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // 3️⃣ Return response
        resp := NoteResponse{
            ID:        noteID,
            JobID:     jobID,
            Text:      req.Text,
            CreatedAt: time.Now().Format(time.RFC3339),
        }

        json.NewEncoder(w).Encode(resp)
    }
}
