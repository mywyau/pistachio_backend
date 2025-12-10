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

func GetJobDetailHandler(db *pgxpool.Pool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        jobIDParam := chi.URLParam(r, "id")
        jobID, err := uuid.Parse(jobIDParam)
        if err != nil {
            http.Error(w, "invalid job id", http.StatusBadRequest)
            return
        }

        ctx := context.Background()

        // 1️⃣ Query job + customer
        var detail JobDetail
        var cust CustomerInfo
        var createdAt time.Time

        err = db.QueryRow(ctx,
            `SELECT 
                j.id, j.title, j.description, j.status, j.estimate, j.created_at,
                c.id, c.name, c.email, c.phone, c.address
             FROM jobs j
             JOIN customers c ON j.customer_id = c.id
             WHERE j.id = $1`,
            jobID,
        ).Scan(
            &detail.ID, &detail.Title, &detail.Description, &detail.Status, &detail.Estimate, &createdAt,
            &cust.ID, &cust.Name, &cust.Email, &cust.Phone, &cust.Address,
        )

        if err != nil {
            http.Error(w, "job not found: "+err.Error(), http.StatusNotFound)
            return
        }

        detail.CreatedAt = createdAt.Format(time.RFC3339)

        // 2️⃣ Query notes
        notesRows, err := db.Query(ctx,
            `SELECT id, text, created_at FROM job_notes WHERE job_id = $1 ORDER BY created_at DESC`,
            jobID,
        )

        if err != nil {
            http.Error(w, "failed to fetch notes: "+err.Error(), http.StatusInternalServerError)
            return
        }
        defer notesRows.Close()

        notes := []JobNote{}
        for notesRows.Next() {
            var n JobNote
            var noteCreated time.Time
            err := notesRows.Scan(&n.ID, &n.Text, &noteCreated)
            if err != nil {
                http.Error(w, "failed to scan note: "+err.Error(), http.StatusInternalServerError)
                return
            }
            n.CreatedAt = noteCreated.Format(time.RFC3339)
            notes = append(notes, n)
        }

        // 3️⃣ Query photos
        photosRows, err := db.Query(ctx,
            `SELECT id, file_url, created_at FROM job_photos WHERE job_id = $1 ORDER BY created_at DESC`,
            jobID,
        )

        if err != nil {
            http.Error(w, "failed to fetch photos: "+err.Error(), http.StatusInternalServerError)
            return
        }
        defer photosRows.Close()

        photos := []JobPhoto{}
        for photosRows.Next() {
            var p JobPhoto
            var photoCreated time.Time
            err := photosRows.Scan(&p.ID, &p.FileURL, &photoCreated)
            if err != nil {
                http.Error(w, "failed to scan photo: "+err.Error(), http.StatusInternalServerError)
                return
            }
            p.CreatedAt = photoCreated.Format(time.RFC3339)
            photos = append(photos, p)
        }

        // 4️⃣ Assemble full response
        resp := JobDetailResponse{
            Job:      detail,
            Customer: cust,
            Notes:    notes,
            Photos:   photos,
        }

        json.NewEncoder(w).Encode(resp)
    }
}
