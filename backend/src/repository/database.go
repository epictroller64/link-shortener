package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Db *pgxpool.Pool

func InitDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v\n", err)
	}

	Db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	fmt.Println("Connected to database")
}

// Call this function when shutting down your application
func CloseDatabase() {
	if Db != nil {
		Db.Close()
	}
}
