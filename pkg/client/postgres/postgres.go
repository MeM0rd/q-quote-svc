package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB

const driverName = "postgres"

func Open() {
	var err error
	log.Printf("starting postgresql")

	dsn := fmt.Sprintf("host=postgresql user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	DB, err = sql.Open(driverName, dsn)
	if err != nil {
		log.Fatalf("postgres connection error: %s", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Postgres ping failde: %v", err)
	}
}

func Close() {
	log.Printf("closing postgresql")
	err := DB.Close()
	if err != nil {
		log.Fatalf("Closing postgres error: %s", err)
	}
}
