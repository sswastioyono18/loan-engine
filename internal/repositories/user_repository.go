package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kitabisa/loan-engine/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUserID(ctx context.Context, userID string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	UpdatePassword(ctx context.Context, id int, hashedPassword string) error
}

type userRepositoryImpl struct {
	base *BaseRepository
}

func NewUserRepository(driver Driver) UserRepository {
	return &userRepositoryImpl{
		base: NewBaseRepository(driver),
	}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			user_id, email, password_hash, user_type, full_name, is_active
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	db := r.base.GetUtilDB()
	err := db.QueryRowContext(
		ctx, query,
		user.UserID, user.Email, user.PasswordHash, user.UserType,
		user.FullName, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, user_id, email, password_hash, user_type, full_name,
		       is_active, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user models.User
	db := r.base.GetUtilDB()
	err := db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, user_id, email, password_hash, user_type, full_name,
		       is_active, created_at, updated_at
		FROM users WHERE email = $1
	`

	var user models.User
	db := r.base.GetUtilDB()
	err := db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) GetByUserID(ctx context.Context, userID string) (*models.User, error) {
	query := `
		SELECT id, user_id, email, password_hash, user_type, full_name,
		       is_active, created_at, updated_at
		FROM users WHERE user_id = $1
	`

	var user models.User
	db := r.base.GetUtilDB()
	err := db.GetContext(ctx, &user, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET
			user_id = $1, email = $2, user_type = $3, full_name = $4,
			is_active = $5, updated_at = NOW()
		WHERE id = $6
	`

	db := r.base.GetUtilDB()
	result, err := db.ExecContext(
		ctx, query,
		user.UserID, user.Email, user.UserType, user.FullName,
		user.IsActive, user.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"
	db := r.base.GetUtilDB()
	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepositoryImpl) UpdatePassword(ctx context.Context, id int, hashedPassword string) error {
	query := "UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2"

	db := r.base.GetUtilDB()
	result, err := db.ExecContext(ctx, query, hashedPassword, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}