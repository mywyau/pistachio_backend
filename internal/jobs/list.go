package jobs

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListJobsHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		rows, err := db.Query(ctx,
			`SELECT 
                j.id,
                j.title,
                j.status,
                j.estimate,
                j.created_at,
                c.name
             FROM jobs j
             JOIN customers c ON c.id = j.customer_id
             ORDER BY j.created_at DESC`,
		)

		if err != nil {
			http.Error(w, "failed to query jobs: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		items := []JobListItem{}

		for rows.Next() {
			var item JobListItem
			var createdAt time.Time

			err := rows.Scan(
				&item.JobID,
				&item.Title,
				&item.Status,
				&item.Estimate,
				&createdAt,
				&item.CustomerName,
			)

			if err != nil {
				http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}

			item.CreatedAt = createdAt.Format(time.RFC3339)
			items = append(items, item)
		}

		json.NewEncoder(w).Encode(items)
	}
}
