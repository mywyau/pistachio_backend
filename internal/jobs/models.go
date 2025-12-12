package jobs

import "github.com/google/uuid"

type CreateJobRequest struct {
    Customer struct {
        Name    string `json:"name"`
        Email   string `json:"email"`
        Phone   string `json:"phone"`
        Address string `json:"address"`
    } `json:"customer"`

    Title       string  `json:"title"`
    Description string  `json:"description"`
    Estimate    float64 `json:"estimate"`
}

type CreateJobResponse struct {
    JobID      uuid.UUID `json:"job_id"`
    CustomerID uuid.UUID `json:"customer_id"`
    Title      string    `json:"title"`
    Status     string    `json:"status"`
}

type JobListItem struct {
    JobID        uuid.UUID `json:"job_id"`
    Title        string    `json:"title"`
    Status       string    `json:"status"`
    Estimate     float64   `json:"estimate"`
    CreatedAt    string    `json:"created_at"`
    CustomerName string    `json:"customer_name"`
}

type JobDetail struct {
    ID          uuid.UUID `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    Estimate    float64   `json:"estimate"`
    CreatedAt   string    `json:"created_at"`
}

type CustomerInfo struct {
    ID      uuid.UUID `json:"id"`
    Name    string    `json:"name"`
    Email   string    `json:"email"`
    Phone   string    `json:"phone"`
    Address string    `json:"address"`
}

type JobNote struct {
    ID        uuid.UUID `json:"id"`
    Text      string    `json:"text"`
    CreatedAt string    `json:"created_at"`
}

type JobPhoto struct {
    ID        uuid.UUID `json:"id"`
    FileURL   string    `json:"file_url"`
    CreatedAt string    `json:"created_at"`
}

type JobDetailResponse struct {
    Job      JobDetail     `json:"job"`
    Customer CustomerInfo  `json:"customer"`
    Notes    []JobNote     `json:"notes"`
    Photos   []JobPhoto    `json:"photos"`
}

type CreateNoteRequest struct {
    Text string `json:"text"`
}

type NoteResponse struct {
    ID        uuid.UUID `json:"id"`
    JobID     uuid.UUID `json:"job_id"`
    Text      string    `json:"text"`
    CreatedAt string    `json:"created_at"`
}

type PhotoResponse struct {
    ID        uuid.UUID `json:"id"`
    JobID     uuid.UUID `json:"job_id"`
    FileURL   string    `json:"file_url"`
    CreatedAt string    `json:"created_at"`
}

// type CreateInvoiceRequest struct {
//     Amount float64 `json:"amount"`
// }

type InvoiceResponse struct {
    InvoiceID uuid.UUID `json:"invoice_id"`
    JobID     uuid.UUID `json:"job_id"`
    Amount    float64   `json:"amount"`
    PDFURL    string    `json:"pdf_url"`
    CreatedAt string    `json:"created_at"`
}
