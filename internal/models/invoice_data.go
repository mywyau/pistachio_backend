package models

import "time"

//
// ─── BUSINESS (SUPPLIER) DETAILS ─────────────────────────────────────────────
//
type BusinessInfo struct {
    Name        string  // e.g. "John Smith Plumbing"
    Address     string  // multiline allowed in PDF
    Email       string
    Phone       string
    Website     string  // optional
    VATNumber   string  // optional
    CompanyReg  string  // optional (for Ltd companies)
    LogoPath    string  // optional custom logo
}

//
// ─── CUSTOMER DETAILS ───────────────────────────────────────────────────────
//
type CustomerInfo struct {
    Name    string
    Email   string
    Address string
}

//
// ─── INVOICE ITEM ───────────────────────────────────────────────────────────
//
type InvoiceItem_v3 struct {
    Description string   // e.g. "Fix leaking tap"
    Quantity    float64  // allow decimals for hours worked
    UnitPrice   float64
    LineTotal   float64  // precomputed or computed in PDF generator
}

//
// ─── TOTALS ─────────────────────────────────────────────────────────────────
//
type InvoiceTotals struct {
    Subtotal float64
    TaxRate  float64 // % e.g. 20.0 for VAT
    TaxAmount float64
    TotalAmount float64
}

//
// ─── PAYMENT DETAILS ─────────────────────────────────────────────────────────
//
type PaymentInfo struct {
    BankName     string
    AccountName  string
    SortCode     string
    AccountNumber string
    IBAN         string // optional
    BIC          string // optional
    PaymentLink  string // Stripe or bank link
    Notes        string // e.g. "Payment due within 14 days"
}

//
// ─── MAIN INVOICE STRUCT ─────────────────────────────────────────────────────
//
type InvoiceData_v3 struct {
    InvoiceID      string      // internal UUID
    InvoiceNumber  string      // visible number (e.g. "INV-0041")
    IssueDate      time.Time
    DueDate        time.Time

    Business       BusinessInfo
    Customer       CustomerInfo

    Items          []InvoiceItem_v3
    Totals         InvoiceTotals
    Payment        PaymentInfo

    FooterNotes    string      // e.g. warranty notes, disclaimers
}
