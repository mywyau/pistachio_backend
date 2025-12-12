package invoices

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

const logoPath_v3 = "assets/gnome.png"

type InvoiceData_v3 struct {
	InvoiceID       string
	CustomerName    string
	CustomerEmail   string
	CustomerAddress string

	Items []InvoiceItem_v3

	TotalAmount float64
	CreatedAt   time.Time
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

	// Register fonts
	pdf.AddUTF8Font("Roboto", "", "assets/font/Roboto-Regular.ttf")
	pdf.AddUTF8Font("Roboto", "B", "assets/font/Roboto-Bold.ttf")

	pdf.SetFont("Roboto", "", 12)

	// ================================================
	// HEADER: LOGO + COMPANY INFO + INVOICE META
	// ================================================
	if _, err := os.Stat(logoPath_v3); err == nil {
		pdf.Image(logoPath_v3, 20, 20, 30, 0, false, "", 0, "")
	}

	// Company Info (next to logo)
	pdf.SetXY(60, 20)
	pdf.SetFont("Roboto", "B", 20)
	pdf.Cell(0, 10, "Pistachio")

	pdf.Ln(5)
	pdf.SetY(pdf.GetY() + 10) // push the title further down

	pdf.SetFont("Roboto", "", 12)
	pdf.SetX(60)
	pdf.Cell(0, 6, "Pistachio Services")
	pdf.Ln(6)

	pdf.SetX(60)
	pdf.Cell(0, 6, "pistachio.example")
	pdf.Ln(6)

	pdf.SetX(60)
	pdf.Cell(0, 6, "support@pistachio.com")

	pdf.Ln(10)
	pdf.SetY(pdf.GetY() + 10) // push the title further down

	pdf.SetFont("Roboto", "", 12)
	pdf.Cell(0, 6, fmt.Sprintf("Invoice ID: %s", data.InvoiceID))
	pdf.Ln(6)

	pdf.Cell(0, 6, fmt.Sprintf("Date: %s", data.CreatedAt.Format("02 Jan 2006")))

	// Spacing before main content
	pdf.Ln(5)
	pdf.SetY(pdf.GetY() + 10) // push the title further down

	// ================================================
	// INVOICE TITLE
	// ================================================
	pdf.SetFont("Roboto", "B", 22)
	pdf.Cell(0, 12, "INVOICE")
	pdf.Ln(18)

	// ================================================
	// BILL TO SECTION
	// ================================================
	pdf.SetFont("Roboto", "B", 14)
	pdf.Cell(0, 8, "Bill To:")
	pdf.Ln(10)

	pdf.SetFont("Roboto", "", 12)
	pdf.Cell(0, 6, data.CustomerName)
	pdf.Ln(6)

	if data.CustomerAddress != "" {
		pdf.MultiCell(0, 6, data.CustomerAddress, "", "", false)
		pdf.Ln(2)
	}

	if data.CustomerEmail != "" {
		pdf.Cell(0, 6, data.CustomerEmail)
		pdf.Ln(10)
	}

	// ================================================
	// ITEMS TABLE
	// ================================================
	pdf.SetFont("Roboto", "B", 12)

	// Column widths (must total <= 170)
	colDesc := 80.0
	colQty := 20.0
	colUnit := 35.0
	colTotal := 35.0

	// Header
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(colDesc, 8, "Description", "1", 0, "", true, 0, "")
	pdf.CellFormat(colQty, 8, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colUnit, 8, "Unit Price", "1", 0, "R", true, 0, "")
	pdf.CellFormat(colTotal, 8, "Total", "1", 0, "R", true, 0, "")
	pdf.Ln(8)

	// Rows
	pdf.SetFont("Roboto", "", 12)
	for _, item := range data.Items {

		// Page break check
		if pdf.GetY() > 260 {
			pdf.AddPage()

			// Re-draw header
			pdf.SetFont("Roboto", "B", 12)
			pdf.SetFillColor(230, 230, 230)
			pdf.CellFormat(colDesc, 8, "Description", "1", 0, "", true, 0, "")
			pdf.CellFormat(colQty, 8, "Qty", "1", 0, "C", true, 0, "")
			pdf.CellFormat(colUnit, 8, "Unit Price", "1", 0, "R", true, 0, "")
			pdf.CellFormat(colTotal, 8, "Total", "1", 0, "R", true, 0, "")
			pdf.Ln(8)
			pdf.SetFont("Roboto", "", 12)
		}

		pdf.CellFormat(colDesc, 8, item.Description, "1", 0, "", false, 0, "")
		pdf.CellFormat(colQty, 8, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colUnit, 8, fmt.Sprintf("£%.2f", item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(colTotal, 8, fmt.Sprintf("£%.2f", item.LineTotal), "1", 0, "R", false, 0, "")
		pdf.Ln(8)
	}

	// // ================================================
	// // ITEMS TABLE
	// // ================================================
	// pdf.SetFont("Roboto", "B", 12)

	// // Header background
	// pdf.SetFillColor(230, 230, 230)
	// pdf.CellFormat(90, 8, "Description", "1", 0, "", true, 0, "")
	// pdf.CellFormat(20, 8, "Qty", "1", 0, "C", true, 0, "")
	// pdf.CellFormat(40, 8, "Unit Price", "1", 0, "R", true, 0, "")
	// pdf.CellFormat(40, 8, "Total", "1", 0, "R", true, 0, "")
	// pdf.Ln(8)

	// // Rows
	// pdf.SetFont("Roboto", "", 12)
	// for _, item := range data.Items {
	// 	if pdf.GetY() > 260 { // A4 is ~297mm high, margins = 20 top/bottom
	// 		pdf.AddPage()

	// 		// --- Redraw table header on new page ---
	// 		pdf.SetFont("Roboto", "B", 12)
	// 		pdf.SetFillColor(230, 230, 230)
	// 		pdf.CellFormat(90, 8, "Description", "1", 0, "", true, 0, "")
	// 		pdf.CellFormat(20, 8, "Qty", "1", 0, "C", true, 0, "")
	// 		pdf.CellFormat(40, 8, "Unit Price", "1", 0, "R", true, 0, "")
	// 		pdf.CellFormat(40, 8, "Total", "1", 0, "R", true, 0, "")
	// 		pdf.Ln(8)

	// 		pdf.SetFont("Roboto", "", 12)
	// 	}
	// 	// --- Write row ---
	// 	pdf.CellFormat(90, 8, item.Description, "1", 0, "", false, 0, "")
	// 	pdf.CellFormat(20, 8, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
	// 	pdf.CellFormat(40, 8, fmt.Sprintf("£%.2f", item.UnitPrice), "1", 0, "R", false, 0, "")
	// 	pdf.CellFormat(40, 8, fmt.Sprintf("£%.2f", item.LineTotal), "1", 0, "R", false, 0, "")
	// 	pdf.Ln(8)
	// }

	// ================================================
	// TOTALS SECTION
	// ================================================
	pdf.Ln(8)
	pdf.SetFont("Roboto", "", 12)

	subtotal := data.TotalAmount // no tax yet

	rightCol := 50.0  // width of value column
	labelCol := 120.0 // width of label column (120 + 50 = 170)

	pdf.CellFormat(labelCol, 8, "Subtotal:", "", 0, "R", false, 0, "")
	pdf.CellFormat(rightCol, 8, fmt.Sprintf("£%.2f", subtotal), "", 0, "R", false, 0, "")
	pdf.Ln(6)

	pdf.CellFormat(labelCol, 8, "Tax (0%):", "", 0, "R", false, 0, "")
	pdf.CellFormat(rightCol, 8, "£0.00", "", 0, "R", false, 0, "")
	pdf.Ln(6)

	pdf.SetFont("Roboto", "B", 14)
	pdf.CellFormat(labelCol, 10, "Amount Due:", "", 0, "R", false, 0, "")
	pdf.CellFormat(rightCol, 10, fmt.Sprintf("£%.2f", data.TotalAmount), "", 0, "R", false, 0, "")
	pdf.Ln(16)

	// ================================================
	// FOOTER
	// ================================================
	pdf.SetFont("Roboto", "", 10)
	pdf.SetTextColor(120, 120, 120)
	pdf.MultiCell(0, 5,
		"Thank you for your business.\nPayments due within 14 days unless otherwise agreed.",
		"", "", false)

	pdf.Ln(5)
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(4)

	pdf.SetFont("Roboto", "", 9)
	pdf.Cell(0, 5, "Pistachio Ltd • pistachio.example • support@pistachio.com")

	// Save PDF
	os.MkdirAll(outputDir, os.ModePerm)
	filename := fmt.Sprintf("%s.pdf", data.InvoiceID)
	path := filepath.Join(outputDir, filename)

	err := pdf.OutputFileAndClose(path)
	if err != nil {
		return "", err
	}

	return path, nil
}
