package repositories

import (
	"database/sql"
	"fmt"
	"time"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Driver interface defines the database driver contract
type Driver interface {
	// GetDB returns the underlying database connection
	GetDB() *sql.DB
	GetUtilDB() *sqlx.DB
	// Close closes the database connection
	Close() error
}

// DBUtil wraps sqlx.DB to provide additional functionality
type DBUtil struct {
	*sqlx.DB
}

// GetDB returns the underlying sql.DB
func (db *DBUtil) GetDB() *sql.DB {
	return db.DB.DB
}

// postgresDriver implements the Driver interface using PostgreSQL and sqlx
type postgresDriver struct {
	db *DBUtil
}

// NewPostgreSQLDriver creates a new PostgreSQL driver instance
func NewPostgreSQLDriver(connectionString string) (Driver, error) {
	// Connect to PostgreSQL
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Wrap with sqlx
	sqlxDB := sqlx.NewDb(db, "postgres")

	return &postgresDriver{db: &DBUtil{sqlxDB}}, nil
}

// GetDB returns the underlying database connection
func (d *postgresDriver) GetDB() *sql.DB {
	return d.db.GetDB()
}

// GetUtilDB returns the sqlx.DB instance
func (d *postgresDriver) GetUtilDB() *sqlx.DB {
	return d.db.DB
}

// Close closes the database connection
func (d *postgresDriver) Close() error {
	return d.db.Close()
}