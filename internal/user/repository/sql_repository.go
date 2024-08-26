package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"


	"github.com/MichaelGenchev/smart-home-system/pkg/models"
	"github.com/lib/pq"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUserExists   = errors.New("user already exists")
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		// Check for unique constraint violation on email
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrUserExists
		}
		return err
	}
	return nil
}

func (r *postgresRepository) GetUser(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, updated_at = $4
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, time.Now())
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *postgresRepository) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}