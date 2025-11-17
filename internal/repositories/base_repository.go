package repositories

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type BaseRepository struct {
	driver Driver
}

func NewBaseRepository(driver Driver) *BaseRepository {
	return &BaseRepository{
		driver: driver,
	}
}

// GetDB returns the underlying database connection
func (r *BaseRepository) GetDB() *sql.DB {
	return r.driver.GetDB()
}

// GetUtilDB returns the sqlx.DB instance
func (r *BaseRepository) GetUtilDB() *sqlx.DB {
	return r.driver.GetUtilDB()
}

// Common methods for all repositories
func (r *BaseRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.GetDB().BeginTx(ctx, nil)
}

func (r *BaseRepository) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (r *BaseRepository) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}
