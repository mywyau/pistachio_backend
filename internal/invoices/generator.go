package invoices

import (
	"fmt"
	"os"
	"path/filepath"

	// "time"
	"pistachio/internal/models"

	"github.com/jung-kurt/gofpdf"
)

const logoPath = "assets/gnome.png"


func drawAddress(pdf *gofpdf.Fpdf, addr models.CustomerAddress, lineHeight float64) {
	lines := []string{
		addr.Line1,
		addr.Line2,
		addr.City,
		addr.Postcode,
		addr.Country,
	}

	for _, line := range lines {
		if line != "" {
			pdf.MultiCell(0, lineHeight, line, "", "", false)
		}
	}
}

func drawAddressBusiness(pdf *gofpdf.Fpdf, addr models.BusinessAddress, lineHeight float64) {
	lines := []string{
		addr.Line1,
		addr.Line2,
		addr.City,
		addr.Postcode,
		addr.Country,
	}

	for _, line := range lines {
		if line != "" {
			pdf.MultiCell(0, lineHeight, line, "", "", false)
		}
	}
}

func drawItemsTableHeader(pdf *gofpdf.Fpdf, colDesc, colQty, colUnit, colTotal float64) {
	pdf.SetFont("Roboto", "B", 12)
	pdf.SetFillColor(230, 230, 230)

	pdf.CellFormat(colDesc, 8, "Description", "1", 0, "", true, 0, "")
	pdf.CellFormat(colQty, 8, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colUnit, 8, "Unit Price", "1", 0, "R", true, 0, "")
	pdf.CellFormat(colTotal, 8, "Total", "1", 0, "R", true, 0, "")
	pdf.Ln(8)

	pdf.SetFont("Roboto", "", 12)
}

func GenerateInvoicePDF(data models.InvoiceData, outputDir string) (string, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	pdf.AliasNbPages("")
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Roboto", "", 9)
		pdf.SetTextColor(120, 120, 120)
		pdf.CellFormat(
			0,
			10,
			fmt.Sprintf("Page %d of {nb}", pdf.PageNo()),
			"",
			0,
			"C",
			false,
			0,
			"",
		)
	})

	// Register fonts
	pdf.AddUTF8Font("Roboto", "", "assets/font/Roboto-Regular.ttf")
	pdf.AddUTF8Font("Roboto", "B", "assets/font/Roboto-Bold.ttf")
	pdf.AddUTF8Font("Roboto", "I", "assets/font/Roboto-Italic.ttf")

	pdf.SetFont("Roboto", "", 12)

	// ============================
	// HEADER SECTION
	// ============================
	// Logo

	pageWidth := 210.0
	rightMargin := 20.0
	leftMargin := 20.0

	logoWidth := 30.0
	logoX := pageWidth - rightMargin - logoWidth
	logoY := 20.0

	headerTop := 20.0

	pageHeight := 297.0
	bottomMargin := 20.0
	footerHeight := 35.0
	totalsHeight := 30.0
	// rowHeight := 8.0

	usableBottomY := pageHeight - bottomMargin - footerHeight - totalsHeight

	if data.Business.LogoPath != "" {
		if _, err := os.Stat(data.Business.LogoPath); err == nil {
			pdf.Image(data.Business.LogoPath, logoX, logoY, logoWidth, 0, false, "", 0, "")
		}
	}

	// Business info (right of logo)
	// pdf.SetXY(logoX - 60, 60)
	pdf.SetXY(leftMargin, headerTop)
	pdf.SetFont("Roboto", "B", 20)
	pdf.Cell(0, 10, data.Business.Name)
	pdf.Ln(10)

	pdf.SetFont("Roboto", "", 12)

	leftColX := leftMargin
	leftColWidth := 90.0

	// Business Address
	if data.Business.BusinessAddress.Line1 != "" {
		drawAddressBusiness(pdf, data.Business.BusinessAddress, 6)
		pdf.Ln(2)
	}

	// Business Email
	if data.Business.Email != "" {
		pdf.SetX(leftColX)
		pdf.Cell(leftColWidth, 6, data.Business.Email)
		pdf.Ln(6)
	}

	// Business Phone
	if data.Business.Phone != "" {
		pdf.SetX(leftColX)
		pdf.Cell(leftColWidth, 6, data.Business.Phone)
		pdf.Ln(6)
	}

	// Website
	if data.Business.Website != "" {
		pdf.SetX(leftColX)
		pdf.Cell(leftColWidth, 6, data.Business.Website)
		pdf.Ln(6)
	}

	// VAT number / Company registration
	if data.Business.VATNumber != "" {
		pdf.SetX(leftColX)
		pdf.Cell(0, 6, "VAT: "+data.Business.VATNumber)
		pdf.Ln(6)
	}

	if data.Business.CompanyReg != "" {
		pdf.SetX(leftColX)
		pdf.Cell(0, 6, "Company Reg: "+data.Business.CompanyReg)
		pdf.Ln(6)
	}

	pdf.Ln(5)

	// Invoice metadata on right side
	pdf.SetFont("Roboto", "B", 12)
	pdf.SetX(130)
	pdf.Cell(40, 6, "Invoice No:")
	pdf.SetFont("Roboto", "", 12)
	pdf.Cell(40, 6, data.InvoiceNumber)
	pdf.Ln(6)

	pdf.SetX(130)
	pdf.SetFont("Roboto", "B", 12)
	pdf.Cell(40, 6, "Issue Date:")
	pdf.SetFont("Roboto", "", 12)
	pdf.Cell(40, 6, data.IssueDate.Format("02 Jan 2006"))
	pdf.Ln(6)

	pdf.SetX(130)
	pdf.SetFont("Roboto", "B", 12)
	pdf.Cell(40, 6, "Due Date:")
	pdf.SetFont("Roboto", "", 12)
	pdf.Cell(40, 6, data.DueDate.Format("02 Jan 2006"))

	pdf.Ln(15)

	// =====================================
	// TITLE
	// =====================================
	pdf.SetFont("Roboto", "B", 22)
	pdf.Cell(0, 12, "INVOICE")
	pdf.Ln(18)

	// =====================================
	// BILL TO SECTION
	// =====================================
	pdf.SetFont("Roboto", "B", 14)
	pdf.Cell(0, 8, "Bill To:")
	pdf.Ln(10)

	pdf.SetFont("Roboto", "", 12)
	pdf.Cell(0, 6, data.Customer.Name)
	pdf.Ln(6)

	if data.Customer.CustomerAddress.Line1 != "" {
		drawAddress(pdf, data.Customer.CustomerAddress, 6)
		pdf.Ln(2)
	}

	if data.Customer.Email != "" {
		pdf.Cell(0, 6, data.Customer.Email)
		pdf.Ln(10)
	}

	// =====================================
	// THANK YOU MESSAGE
	// =====================================
	pdf.Ln(4)
	pdf.SetFont("Roboto", "", 11)
	pdf.SetTextColor(80, 80, 80)

	pdf.MultiCell(
		0,
		6,
		"Thank you for your business!",
		"",
		"",
		false,
	)

	pdf.Ln(6)

	// Reset for table
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Roboto", "B", 12)

	// =====================================
	// ITEMS TABLE
	// =====================================
	pdf.SetFont("Roboto", "B", 12)

	colDesc := 80.0
	colQty := 20.0
	colUnit := 35.0
	colTotal := 35.0

	drawItemsTableHeader(pdf, colDesc, colQty, colUnit, colTotal)

	for i, item := range data.Items {

		// Estimate row height (safe even if description grows later)
		rowHeight := 8.0

		// BEFORE writing the row, check space
		if pdf.GetY()+rowHeight > usableBottomY {

			// Optional continuation note
			pdf.SetFont("Roboto", "I", 9)
			pdf.SetTextColor(120, 120, 120)

			// Move cursor to right margin and align text right
			pdf.CellFormat(
				0, // full width
				5,
				"Items continued on next page…",
				"",
				1,   // line break
				"R", // RIGHT aligned
				false,
				0,
				"",
			)
			pdf.SetTextColor(0, 0, 0)

			pdf.Ln(6)

			pdf.AddPage()
			drawItemsTableHeader(pdf, colDesc, colQty, colUnit, colTotal)
		}

		pdf.CellFormat(colDesc, rowHeight, item.Description, "1", 0, "", false, 0, "")
		pdf.CellFormat(colQty, rowHeight, fmt.Sprintf("%.2f", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colUnit, rowHeight, fmt.Sprintf("£%.2f", item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(colTotal, rowHeight, fmt.Sprintf("£%.2f", item.LineTotal), "1", 0, "R", false, 0, "")
		pdf.Ln(rowHeight)

		_ = i // (useful later for zebra striping)
	}

	// =====================================
	// TOTALS
	// =====================================

	if pdf.GetY()+totalsHeight > pageHeight-bottomMargin {
		pdf.AddPage()
	}

	pdf.Ln(6)

	labelCol := 120.0
	rightCol := 50.0

	pdf.SetFont("Roboto", "", 12)
	pdf.CellFormat(labelCol, 8, "Subtotal:", "", 0, "R", false, 0, "")
	pdf.CellFormat(rightCol, 8, fmt.Sprintf("£%.2f", data.Totals.Subtotal), "", 0, "R", false, 0, "")
	pdf.Ln(6)

	pdf.CellFormat(labelCol, 8, fmt.Sprintf("Tax (%.1f%%):", data.Totals.TaxRate), "", 0, "R", false, 0, "")
	pdf.CellFormat(rightCol, 8, fmt.Sprintf("£%.2f", data.Totals.TaxAmount), "", 0, "R", false, 0, "")
	pdf.Ln(6)

	pdf.SetFont("Roboto", "B", 14)
	pdf.CellFormat(labelCol, 10, "Amount Due:", "", 0, "R", false, 0, "")
	pdf.CellFormat(rightCol, 10, fmt.Sprintf("£%.2f", data.Totals.TotalAmount), "", 0, "R", false, 0, "")
	pdf.Ln(15)

	// =====================================
	// PAYMENT INFORMATION
	// =====================================
	pdf.SetFont("Roboto", "B", 14)
	pdf.Cell(0, 8, "Payment Details")
	pdf.Ln(10)

	pdf.SetFont("Roboto", "", 12)

	if data.Payment.AccountName != "" {
		pdf.Cell(0, 6, "Account Name: "+data.Payment.AccountName)
		pdf.Ln(6)
	}
	if data.Payment.BankName != "" {
		pdf.Cell(0, 6, "Bank: "+data.Payment.BankName)
		pdf.Ln(6)
	}
	if data.Payment.SortCode != "" {
		pdf.Cell(0, 6, "Sort Code: "+data.Payment.SortCode)
		pdf.Ln(6)
	}
	if data.Payment.AccountNumber != "" {
		pdf.Cell(0, 6, "Account Number: "+data.Payment.AccountNumber)
		pdf.Ln(6)
	}
	if data.Payment.IBAN != "" {
		pdf.Cell(0, 6, "IBAN: "+data.Payment.IBAN)
		pdf.Ln(6)
	}
	if data.Payment.BIC != "" {
		pdf.Cell(0, 6, "BIC: "+data.Payment.BIC)
		pdf.Ln(6)
	}

	if data.Payment.Notes != "" {
		pdf.Ln(4)
		pdf.MultiCell(0, 6, data.Payment.Notes, "", "", false)
	}

	pdf.Ln(10)

	// =====================================
	// FOOTER NOTES (CONTENT)
	// =====================================
	if data.FooterNotes != "" {

		// Ensure notes don’t collide with footer
		if pdf.GetY()+20 > pageHeight-bottomMargin {
			pdf.AddPage()
		}

		pdf.SetFont("Roboto", "", 10)
		pdf.SetTextColor(100, 100, 100)
		pdf.MultiCell(0, 5, data.FooterNotes, "", "", false)
	}

	// Save file
	os.MkdirAll(outputDir, os.ModePerm)
	filename := fmt.Sprintf("%s.pdf", data.InvoiceID)
	path := filepath.Join(outputDir, filename)

	err := pdf.OutputFileAndClose(path)
	if err != nil {
		return "", err
	}

	return path, nil
}
