package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository akses data user untuk modul auth.
type Repository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, u User) (User, error) {
	const q = `
		INSERT INTO users (email, name, password_hash, role)
		VALUES ($1, $2, $3, 'user')
		RETURNING id, email, name, password_hash, role`
	row := r.pool.QueryRow(ctx, q, u.Email, u.Name, u.PasswordHash)
	var out User
	err := row.Scan(&out.ID, &out.Email, &out.Name, &out.PasswordHash, &out.Role)
	if err != nil {
		return User{}, fmt.Errorf("auth repository: create gagal: %w", err)
	}
	return out, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (User, error) {
	const q = `SELECT id, email, name, password_hash, role FROM users WHERE email = $1`
	var out User
	err := r.pool.QueryRow(ctx, q, email).Scan(&out.ID, &out.Email, &out.Name, &out.PasswordHash, &out.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("auth repository: get by email gagal: %w", err)
	}
	return out, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (User, error) {
	const q = `SELECT id, email, name, password_hash, role FROM users WHERE id = $1`
	var out User
	err := r.pool.QueryRow(ctx, q, id).Scan(&out.ID, &out.Email, &out.Name, &out.PasswordHash, &out.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("auth repository: get by id gagal: %w", err)
	}
	return out, nil
}
