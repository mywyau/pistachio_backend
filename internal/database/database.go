package database

import (
    "context"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDatabase(url string) *pgxpool.Pool {
    cfg, err := pgxpool.ParseConfig(url)
    if err != nil {
        log.Fatalf("Unable to parse db config: %v\n", err)
    }

    pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
    if err != nil {
        log.Fatalf("Unable to connect to db: %v\n", err)
    }

    return pool
}
