package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	dsn := "postgres://quorbit:quorbit@localhost:15432/go_test?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v\n", err)
	}

	DB = pool

	log.Println("Connected to postgres")
}
