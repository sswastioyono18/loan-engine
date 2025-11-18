package util

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB represents a database connection wrapper that implements the Driver interface
type DB struct {
	DB     *sql.DB
	SqlxDB *sqlx.DB
}

// GetDB returns the underlying database connection
func (d *DB) GetDB() *sql.DB {
	return d.DB
}

// GetUtilDB returns the sqlx.DB instance
func (d *DB) GetUtilDB() *sqlx.DB {
	return d.SqlxDB
}

// Close closes the database connection
func (d *DB) Close() error {
	if d.SqlxDB != nil {
		return d.SqlxDB.Close()
	}
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// InitDB initializes the database connection
func InitDB() (*DB, error) {
	// Get database connection details from environment variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "loan_engine_user"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "loan_engine_password"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "loan_engine_db"
	}

	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	sqlxDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := sqlxDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	sqlxDB.SetMaxOpenConns(25)
	sqlxDB.SetMaxIdleConns(5)
	sqlxDB.SetConnMaxLifetime(30 * time.Minute)

	// Wrap with sqlx
	sqlxDBWrapped := sqlx.NewDb(sqlxDB, "postgres")

	return &DB{DB: sqlxDB, SqlxDB: sqlxDBWrapped}, nil
}
