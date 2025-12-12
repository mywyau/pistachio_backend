package invoices

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/jung-kurt/gofpdf"
	"time"
)

const logoPath_v3 = "assets/gnome.png"

type InvoiceData_v3 struct {
	InvoiceID       string
	CustomerName    string
	CustomerEmail   string
	CustomerAddress string

	Items           []InvoiceItem_v3

	TotalAmount     float64
	CreatedAt       time.Time
}

type InvoiceItem_v3 struct {
	Description string
	Quantity    int
	UnitPrice   float64
	LineTotal   float64
}

func GenerateInvoicePDF_v3(data InvoiceData_v3, outputDir string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// Register UTF-8 Font BEFORE using it
	pdf.AddUTF8Font("Roboto", "", "assets/font/Roboto-Regular.ttf")
	pdf.AddUTF8Font("Roboto", "B", "assets/font/Roboto-Bold.ttf")
	
	pdf.SetFont("Roboto", "", 12)

	// Optional: Title font sizes
	headerFont := 26.00
	sectionTitleFont := 14.00

	// Add logo if it exists
	if _, err := os.Stat(logoPath_v3); err == nil {
		pdf.Image(logoPath_v3, 20, 15, 30, 0, false, "", 0, "")
	}

	//
	// COMPANY NAME
	//
	pdf.SetFont("Roboto", "B", headerFont)
	pdf.SetXY(20, 20)
	pdf.Cell(0, 10, "Pistachio")
	pdf.Ln(18)

	//
	// INVOICE TITLE
	//
	pdf.SetFont("Roboto", "B", 20)
	pdf.Cell(0, 12, "Invoice")
	pdf.Ln(14)

	//
	// INVOICE META INFORMATION
	//
	pdf.SetFont("Roboto", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Invoice ID: %s", data.InvoiceID))
	pdf.Ln(6)

	pdf.Cell(0, 8, fmt.Sprintf("Date: %s", data.CreatedAt.Format("02 Jan 2006")))
	pdf.Ln(12)

	//
	// CUSTOMER DETAILS
	//
	pdf.SetFont("Roboto", "B", sectionTitleFont)
	pdf.Cell(0, 10, "Bill To:")
	pdf.Ln(8)

	pdf.SetFont("Roboto", "", 12)
	pdf.Cell(0, 7, data.CustomerName)
	pdf.Ln(5)

	if data.CustomerAddress != "" {
		pdf.MultiCell(0, 6, data.CustomerAddress, "", "", false)
		pdf.Ln(2)
	}

	if data.CustomerEmail != "" {
		pdf.Cell(0, 7, data.CustomerEmail)
		pdf.Ln(10)
	}

	//
	// TABLE HEADER
	//
	pdf.SetFont("Roboto", "B", 12)

	pdf.CellFormat(90, 8, "Description", "1", 0, "", false, 0, "")
	pdf.CellFormat(20, 8, "Qty", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, "Unit Price", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 8, "Total", "1", 0, "R", false, 0, "")
	pdf.Ln(8)

	//
	// TABLE ROWS
	//
	pdf.SetFont("Roboto", "", 12)

	for _, item := range data.Items {
		pdf.CellFormat(90, 8, item.Description, "1", 0, "", false, 0, "")
		pdf.CellFormat(20, 8, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("£%.2f", item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("£%.2f", item.LineTotal), "1", 0, "R", false, 0, "")
		pdf.Ln(8)
	}

	pdf.Ln(6)

	//
	// TOTAL SECTION
	//
	pdf.SetFont("Roboto", "B", 16)
	pdf.Cell(0, 10, fmt.Sprintf("Amount Due: £%.2f", data.TotalAmount))
	pdf.Ln(14)

	//
	// FOOTER
	//
	pdf.SetFont("Roboto", "", 10)
	pdf.SetTextColor(120, 120, 120)

	pdf.MultiCell(0, 5,
		"Thank you for your business.\nPayments are due within 14 days unless otherwise agreed.",
		"", "", false,
	)

	//
	// SAVE FILE
	//
	os.MkdirAll(outputDir, os.ModePerm)

	filename := fmt.Sprintf("%s.pdf", data.InvoiceID)
	outputPath := filepath.Join(outputDir, filename)

	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		return "", err
	}

	return outputPath, nil
}
