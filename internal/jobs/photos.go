package jobs

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

func UploadPhotoHandler(db *pgxpool.Pool, uploadDir string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        jobIDParam := chi.URLParam(r, "id")
        jobID, err := uuid.Parse(jobIDParam)
        if err != nil {
            http.Error(w, "invalid job id", http.StatusBadRequest)
            return
        }

        // 1️⃣ Ensure job exists
        ctx := context.Background()
        var exists bool
        err = db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM jobs WHERE id = $1)`, jobID).Scan(&exists)
        if err != nil || !exists {
            http.Error(w, "job not found", http.StatusNotFound)
            return
        }

        // 2️⃣ Parse multipart form
        r.ParseMultipartForm(10 << 20) // 10MB

        file, header, err := r.FormFile("file")
        if err != nil {
            http.Error(w, "file upload error: "+err.Error(), http.StatusBadRequest)
            return
        }
        defer file.Close()

        // 3️⃣ Create unique file path
        photoID := uuid.New()
        ext := filepath.Ext(header.Filename)
        fileName := fmt.Sprintf("%s%s", photoID.String(), ext)

        savePath := filepath.Join(uploadDir, fileName)

        // Ensure directory exists
        os.MkdirAll(uploadDir, os.ModePerm)

        // 4️⃣ Save file locally
        dst, err := os.Create(savePath)
        if err != nil {
            http.Error(w, "failed to save file: "+err.Error(), http.StatusInternalServerError)
            return
        }
        defer dst.Close()

        _, err = io.Copy(dst, file)
        if err != nil {
            http.Error(w, "failed to write file: "+err.Error(), http.StatusInternalServerError)
            return
        }

        fileURL := "/uploads/photos/" + fileName

        // 5️⃣ Insert DB entry
        _, err = db.Exec(ctx,
            `INSERT INTO job_photos (id, job_id, file_url, created_at)
             VALUES ($1, $2, $3, NOW())`,
            photoID,
            jobID,
            fileURL,
        )

        if err != nil {
            http.Error(w, "db insert failed: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // 6️⃣ Respond
        resp := PhotoResponse{
            ID:        photoID,
            JobID:     jobID,
            FileURL:   fileURL,
            CreatedAt: time.Now().Format(time.RFC3339),
        }

        json.NewEncoder(w).Encode(resp)
    }
}
