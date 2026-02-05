package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"ecommerce-backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context, page, limit, rangeDays int) ([]models.User, int, error)
	Update(ctx context.Context, user *models.User) error
	UpdateRole(ctx context.Context, id uuid.UUID, role string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (email, password_hash, first_name, last_name, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		log.Printf("❌ Error creating user: %v", err)
		return err
	}

	log.Printf("✅ User created with ID: %v, Email: %v", user.ID, user.Email)
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
        SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users
        SET email = $1, first_name = $2, last_name = $3, updated_at = NOW()
        WHERE id = $4
    `

	_, err := r.db.Exec(ctx, query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.ID,
	)

	return err
}

func (r *userRepository) UpdateRole(ctx context.Context, id uuid.UUID, role string) error {
	query := `
        UPDATE users
        SET role = $1, updated_at = NOW()
        WHERE id = $2
    `

	_, err := r.db.Exec(ctx, query, role, id)
	return err
}

func (r *userRepository) GetAll(ctx context.Context, page, limit, rangeDays int) ([]models.User, int, error) {
	offset := (page - 1) * limit

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if rangeDays > 0 {
		whereClause += fmt.Sprintf(" AND created_at >= NOW() - $%d * INTERVAL '1 day'", argCount)
		args = append(args, rangeDays)
		argCount++
	}

	countQuery := "SELECT COUNT(*) FROM users " + whereClause
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
        SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at
        FROM users
    ` + whereClause + ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argCount) + ` OFFSET $` + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
