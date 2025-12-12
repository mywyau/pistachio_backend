package main

import (
	"fmt"
	"net/http"
	"pistachio/internal/database"
	"pistachio/internal/jobs"

	"github.com/go-chi/chi/v5"
)

func main() {
	dbUrl := "postgres://pistachio:pistachio_pwd@localhost:5432/pistachio_db"

	db := database.ConnectDatabase(dbUrl)
	defer db.Close()

	r := chi.NewRouter()

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	r.Post("/jobs", jobs.CreateJobHandler(db))

	r.Get("/jobs", jobs.ListJobsHandler(db))

	r.Get("/jobs/{id}", jobs.GetJobDetailHandler(db))

	r.Post("/jobs/{id}/notes", jobs.CreateNoteHandler(db))

	// r.Post("/jobs/{id}/invoice", jobs.CreateInvoiceHandler(db))

	r.Post("/invoices", jobs.CreateInvoiceHandler_v3(db))

	r.Put("/jobs/{id}/status", jobs.UpdateJobStatusHandler(db))

	// Serve static uploaded files (for local dev)
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Upload endpoint
	r.Post("/jobs/{id}/photos", jobs.UploadPhotoHandler(db, "uploads/photos"))

	fmt.Println("API running on :8080")
	http.ListenAndServe(":8080", r)
}
