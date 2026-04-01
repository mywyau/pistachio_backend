package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"pistachio/internal/database"
	"pistachio/internal/jobs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// --- Config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // local default
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://pistachio:pistachio_pwd@localhost:5432/pistachio_db"
	}

	// --- Database
	db := database.ConnectDatabase(dbURL)
	defer db.Close()

	// --- Router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- CORS (prod + local)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://your-frontend-domain.vercel.app", // add later
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	// --- Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	r.Post("/jobs", jobs.CreateJobHandler(db))
	r.Get("/jobs", jobs.ListJobsHandler(db))
	r.Get("/jobs/{id}", jobs.GetJobDetailHandler(db))
	r.Post("/jobs/{id}/notes", jobs.CreateNoteHandler(db))
	r.Post("/invoices", jobs.CreateInvoiceHandler(db))
	r.Put("/jobs/{id}/status", jobs.UpdateJobStatusHandler(db))

	// Local file serving (V1 only)
	r.Handle("/uploads/*",
		http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))),
	)
	r.Post("/jobs/{id}/photos", jobs.UploadPhotoHandler(db, "uploads/photos"))

	// --- Server
	log.Printf("🚀 API running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
