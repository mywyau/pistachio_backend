package invoices

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type InvoiceData struct {
	InvoiceID string
	JobID     string
	Customer  string
	Amount    float64
	CreatedAt time.Time
}

func GenerateInvoicePDF(data InvoiceData, outputDir string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")


	pdf.AddUTF8Font("Roboto", "", "assets/fonts/Roboto-Regular.ttf")
	pdf.SetFont("Roboto", "", 12)
	// pdf.SetFont("Arial", "B", 20)
	pdf.AddPage()

	// Invoice Header
	pdf.Cell(40, 10, "Pistachio — Invoice")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)

	// Invoice metadata
	pdf.Cell(40, 10, fmt.Sprintf("Invoice ID: %s", data.InvoiceID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Job ID: %s", data.JobID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Customer: %s", data.Customer))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Created: %s", data.CreatedAt.Format("02 Jan 2006")))
	pdf.Ln(12)

	// Amount
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Amount Due: £%.2f", data.Amount))

	// Ensure output directory exists
	os.MkdirAll(outputDir, os.ModePerm)

	// File path
	filename := fmt.Sprintf("%s.pdf", data.InvoiceID)
	outputPath := filepath.Join(outputDir, filename)

	// Save PDF
	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
