package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"pistachio/internal/invoices"
	"pistachio/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Request format matching the new invoice-only UI
type CreateInvoiceRequest struct {
	CustomerName    string               `json:"customer_name"`
	CustomerEmail   string               `json:"customer_email"`
	CustomerAddress string               `json:"customer_address"`
	Items           []models.InvoiceItem `json:"items"`
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

		// --- Calculate totals ---
		subtotal := 0.0
		for i := range req.Items {
			item := &req.Items[i]
			item.LineTotal = float64(item.Quantity) * item.UnitPrice
			subtotal += item.LineTotal
		}

		taxRate := 0.0 // no tax for now
		taxAmount := 0.0
		total := subtotal + taxAmount

		invoiceID := uuid.New()
		now := time.Now()
		dueDate := now.Add(14 * 24 * time.Hour) // default: 14 days

		// --- Build full invoice model ---
		invoiceData := models.InvoiceData{
			InvoiceID:     invoiceID.String(),
			InvoiceNumber: fmt.Sprintf("INV-%s", invoiceID.String()[0:8]),
			IssueDate:     now,
			DueDate:       dueDate,

			Business: models.BusinessInfo{
				Name:      "Pistachio Ltd",
				Address:   "123 Example Street\nLondon, UK",
				Email:     "support@pistachio.com",
				Phone:     "+44 0000 000000",
				Website:   "https://pistachio.example",
				VATNumber: "",
				LogoPath:  "assets/gnome.png",
			},

			Customer: models.CustomerInfo{
				Name:    req.CustomerName,
				Email:   req.CustomerEmail,
				Address: req.CustomerAddress,
			},

			Items: req.Items,

			Totals: models.InvoiceTotals{
				Subtotal:    subtotal,
				TaxRate:     taxRate,
				TaxAmount:   taxAmount,
				TotalAmount: total,
			},

			Payment: models.PaymentInfo{
				BankName:      "Barclays",
				AccountName:   "Pistachio Ltd",
				SortCode:      "00-00-00",
				AccountNumber: "00000000",
				Notes:         "Payment due in 30 days.",
			},

			FooterNotes: "some placeholder footer note",
		}

		// STEP 1 — Insert into DB (JSON items)
		ctx := context.Background()

		itemsJSON, err := json.Marshal(req.Items)
		if err != nil {
			http.Error(w, "cannot encode items json", 500)
			return
		}

		placeholderPDF := "pending"

		_, err = db.Exec(ctx, `
            INSERT INTO invoices 
            (id, customer_name, customer_email, customer_address, items, total, pdf_url, created_at)
            VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
        `,
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
			http.Error(w, "failed to insert invoice: "+err.Error(), 500)
			return
		}

		// STEP 2 — Generate PDF
		pdfPath, err := invoices.GenerateInvoicePDF(invoiceData, "uploads/invoices")
		if err != nil {
			http.Error(w, "failed to generate PDF: "+err.Error(), 500)
			return
		}

		pdfURL := "/uploads/invoices/" + filepath.Base(pdfPath)

		// STEP 3 — Update DB with final PDF URL
		_, err = db.Exec(ctx, `UPDATE invoices SET pdf_url=$1 WHERE id=$2`, pdfURL, invoiceID)
		if err != nil {
			http.Error(w, "failed to update pdf url: "+err.Error(), 500)
			return
		}

		// STEP 4 — Respond
		json.NewEncoder(w).Encode(map[string]any{
			"invoice_id":     invoiceID.String(),
			"invoice_number": invoiceData.InvoiceNumber,
			"total":          invoiceData.Totals.TotalAmount,
			"pdf_url":        pdfURL,
			"issue_date":     now,
			"due_date":       dueDate,
		})
	}
}
