package jobs

import (
    "context"
    "encoding/json"
    "net/http"
    "path/filepath"
    "time"

	"pistachio/internal/invoices"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Request format matching the new invoice-only UI
type CreateInvoiceRequest struct {
    CustomerName    string        `json:"customer_name"`
    CustomerEmail   string        `json:"customer_email"`
    CustomerAddress string        `json:"customer_address"`
    Items           []invoices.InvoiceItem_v3 `json:"items"`
}

func CreateInvoiceHandler_v3(db *pgxpool.Pool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        var req CreateInvoiceRequest

        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid JSON body", http.StatusBadRequest)
            return
        }

        // Validation
        if req.CustomerName == "" {
            http.Error(w, "customer_name is required", http.StatusBadRequest)
            return
        }

        if len(req.Items) == 0 {
            http.Error(w, "invoice must include at least one item", http.StatusBadRequest)
            return
        }

        // Calculate totals
        total := 0.0
        for i := range req.Items {
            item := &req.Items[i]
            item.LineTotal = float64(item.Quantity) * item.UnitPrice
            total += item.LineTotal
        }

        invoiceID := uuid.New()
        now := time.Now()
        ctx := context.Background()

        // Convert items to JSON for DB
        itemsJSON, err := json.Marshal(req.Items)
        if err != nil {
            http.Error(w, "failed to encode items", http.StatusInternalServerError)
            return
        }

        // STEP 1 — Insert invoice with placeholder PDF URL
        placeholderPDF := "pending"

        _, err = db.Exec(ctx,
            `INSERT INTO invoices 
            (id, customer_name, customer_email, customer_address, items, total, pdf_url, created_at)
             VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
            invoiceID,
            req.CustomerName,
            req.CustomerEmail,
            req.CustomerAddress,
            itemsJSON,
            total,
            placeholderPDF,
            now,
        )

        if err != nil {
            http.Error(w, "failed to insert invoice: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // STEP 2 — Generate PDF using your v3 generator
        pdfPath, err := invoices.GenerateInvoicePDF_v3(invoices.InvoiceData_v3{
            InvoiceID:       invoiceID.String(),
            CustomerName:    req.CustomerName,
            CustomerEmail:   req.CustomerEmail,
            CustomerAddress: req.CustomerAddress,
            Items:           req.Items,
            TotalAmount:     total,
            CreatedAt:       now,
        }, "uploads/invoices")

        if err != nil {
            http.Error(w, "failed to generate PDF: "+err.Error(), http.StatusInternalServerError)
            return
        }

        pdfURL := "/uploads/invoices/" + filepath.Base(pdfPath)

        // STEP 3 — Update PDF URL in DB
        _, err = db.Exec(ctx,
            `UPDATE invoices SET pdf_url=$1 WHERE id=$2`,
            pdfURL, invoiceID,
        )

        if err != nil {
            http.Error(w, "failed to store PDF URL", http.StatusInternalServerError)
            return
        }

        // STEP 4 — Respond with invoice info
        json.NewEncoder(w).Encode(map[string]any{
            "invoice_id": invoiceID,
            "customer_name": req.CustomerName,
            "total": total,
            "pdf_url": pdfURL,
            "created_at": now.Format(time.RFC3339),
        })
    }
}
