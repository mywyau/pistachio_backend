package models

import "time"

// Business (Supplier)
type BusinessInfo struct {
	Name       string
	BusinessAddress     BusinessAddress
	Email      string
	Phone      string
	Website    string
	VATNumber  string
	CompanyReg string
	LogoPath   string
}

// Customer
type CustomerInfo struct {
	Name    string
	Email   string
	CustomerAddress CustomerAddress `json:"address"`
}

// Customer
type BusinessAddress struct {
	Line1    string `json:"line1"`
	Line2    string `json:"line2"`
	City     string `json:"city"`
	Postcode string `json:"postcode"`
	Country  string `json:"country"`
}


// Customer
type CustomerAddress struct {
	Line1    string `json:"line1"`
	Line2    string `json:"line2"`
	City     string `json:"city"`
	Postcode string `json:"postcode"`
	Country  string `json:"country"`
}

// Item
type InvoiceItem struct {
	Description string
	Quantity    float64
	UnitPrice   float64
	LineTotal   float64
}

// Totals
type InvoiceTotals struct {
	Subtotal    float64
	TaxRate     float64
	TaxAmount   float64
	TotalAmount float64
}

// Payment
type PaymentInfo struct {
	BankName      string
	AccountName   string
	SortCode      string
	AccountNumber string
	IBAN          string
	BIC           string
	PaymentLink   string
	Notes         string
}

// Full Invoice Structure
type InvoiceData struct {
	InvoiceID     string // internal UUID
	InvoiceNumber string // visible invoice number e.g. INV-0041

	IssueDate time.Time
	DueDate   time.Time

	Business BusinessInfo
	Customer CustomerInfo

	Items   []InvoiceItem
	Totals  InvoiceTotals
	Payment PaymentInfo

	FooterNotes string // optional footer or custom text
}
