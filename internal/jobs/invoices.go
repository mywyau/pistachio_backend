package jobs

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
    "path/filepath"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "pistachio/internal/invoices"
)

func CreateInvoiceHandler(db *pgxpool.Pool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        jobIDParam := chi.URLParam(r, "id")
        jobID, err := uuid.Parse(jobIDParam)
        if err != nil {
            http.Error(w, "invalid job id", http.StatusBadRequest)
            return
        }

        var req CreateInvoiceRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid JSON", http.StatusBadRequest)
            return
        }

        if req.Amount <= 0 {
            http.Error(w, "invoice amount must be > 0", http.StatusBadRequest)
            return
        }

        ctx := context.Background()

        // 1️⃣ Check job exists and status
        var currentStatus string
        err = db.QueryRow(ctx,
            `SELECT status FROM jobs WHERE id = $1`,
            jobID,
        ).Scan(&currentStatus)

        if err != nil {
            http.Error(w, "job not found", http.StatusNotFound)
            return
        }

        if currentStatus != "completed" {
            http.Error(w, "job must be 'completed' before invoicing", http.StatusBadRequest)
            return
        }

        // Fetch customer name
        var custName string
        err = db.QueryRow(ctx,
            `SELECT c.name 
             FROM jobs j 
             JOIN customers c ON j.customer_id = c.id
             WHERE j.id = $1`,
            jobID,
        ).Scan(&custName)

        if err != nil {
            http.Error(w, "failed to fetch customer name: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // 2️⃣ Insert invoice with placeholder PDF URL
        invoiceID := uuid.New()
        placeholderPDF := "pending"

        _, err = db.Exec(ctx,
            `INSERT INTO invoices (id, job_id, amount, pdf_url, created_at)
             VALUES ($1, $2, $3, $4, NOW())`,
            invoiceID, jobID, req.Amount, placeholderPDF,
        )

        if err != nil {
            http.Error(w, "failed to create invoice: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // 3️⃣ Generate PDF file
        pdfPath, err := invoices.GenerateInvoicePDF(invoices.InvoiceData{
            InvoiceID: invoiceID.String(),
            JobID:     jobID.String(),
            Customer:  custName,
            Amount:    req.Amount,
            CreatedAt: time.Now(),
        }, "uploads/invoices")

        if err != nil {
            http.Error(w, "failed to generate PDF: "+err.Error(), http.StatusInternalServerError)
            return
        }

        pdfURL := "/uploads/invoices/" + filepath.Base(pdfPath)

        // 4️⃣ Update invoice with actual PDF URL
        _, err = db.Exec(ctx,
            `UPDATE invoices SET pdf_url=$1 WHERE id=$2`,
            pdfURL, invoiceID,
        )

        if err != nil {
            http.Error(w, "failed to update PDF URL: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // 5️⃣ Update job status to invoiced
        _, err = db.Exec(ctx,
            `UPDATE jobs SET status='invoiced', updated_at=NOW() WHERE id=$1`,
            jobID,
        )

        if err != nil {
            http.Error(w, "failed to update job status: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // 6️⃣ Return response
        resp := InvoiceResponse{
            InvoiceID: invoiceID,
            JobID:     jobID,
            Amount:    req.Amount,
            PDFURL:    pdfURL,
            CreatedAt: time.Now().Format(time.RFC3339),
        }

        json.NewEncoder(w).Encode(resp)
    }
}
