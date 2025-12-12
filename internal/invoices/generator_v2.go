package invoices

import (
	"fmt"
	"os"
	"path/filepath"
	// "time"

	"github.com/jung-kurt/gofpdf"
)

const logoPath = "assets/gnome.png"

func GenerateInvoicePDF_v2(data InvoiceData, outputDir string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// Try to add logo
	if _, err := os.Stat(logoPath); err == nil {
		pdf.Image(logoPath, 20, 15, 30, 0, false, "", 0, "")
	}

	// Company Name
	pdf.SetFont("Helvetica", "B", 24)
	pdf.SetXY(20, 20)
	pdf.Cell(0, 10, "Pistachio")
	pdf.Ln(15)

	// Invoice Title
	pdf.SetFont("Helvetica", "B", 18)
	pdf.Cell(0, 10, "Invoice")
	pdf.Ln(12)

	// Invoice Meta Section
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Invoice ID: %s", data.InvoiceID))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Job ID: %s", data.JobID))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Customer: %s", data.Customer))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Date: %s", data.CreatedAt.Format("02 Jan 2006")))
	pdf.Ln(15)

	// Amount Section
	pdf.SetFont("Helvetica", "B", 16)
	pdf.Cell(0, 10, fmt.Sprintf("Amount Due: £%.2f", data.Amount))
	pdf.Ln(20)

	// Horizontal line
	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(10)

	// Footer Note
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.MultiCell(0, 5,
		"Thank you for choosing Pistachio.\nPayments are due within 14 days unless otherwise agreed.",
		"",    // border
		"",    // alignment
		false, // fill
	)

	// Ensure output directory exists
	os.MkdirAll(outputDir, os.ModePerm)

	filename := fmt.Sprintf("%s.pdf", data.InvoiceID)
	outputPath := filepath.Join(outputDir, filename)

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
