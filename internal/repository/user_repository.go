package repository

import (
	"database/sql"
	"fmt"

	"github.com/majidkabir/TinyGo/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.CreateUserRequest) (*models.User, error) {
	query := `INSERT INTO users (name, email, age) 
	VALUES ($1, $2, $3) 
	RETURNING id, name, email, age, created_at, updated_at`

	var result models.User
	err := r.db.QueryRow(query, user.Name, user.Email, user.Age).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Age,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &result, nil
}

// GetByID retrieves a user by their ID.
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `SELECT id, name, email, age, created_at, updated_at 
	FROM users WHERE id = $1`

	var result models.User
	err := r.db.QueryRow(query, id).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Age,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &result, nil
}

// GetByEmail retrieves a user by their email address.
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, age, created_at, updated_at 
	FROM users WHERE email = $1`

	var result models.User
	err := r.db.QueryRow(query, email).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Age,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &result, nil
}

func (r *UserRepository) GetAll(limit, offset int) ([]*models.User, error) {
	query := `SELECT id, name, email, age, created_at, updated_at FROM users LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Age,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}

func (r *UserRepository) Update(id int, user *models.UpdateUserRequest) (*models.User, error) {
	query := `UPDATE users SET name = COALESCE($1, name), 
	email = COALESCE($2, email), age = COALESCE($3, age), 
	updated_at = NOW() WHERE id = $4 
	RETURNING id, name, email, age, created_at, updated_at`

	var result models.User
	err := r.db.QueryRow(query, user.Name, user.Email, user.Age, id).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Age,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &result, nil
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := r.db.QueryRow(query).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}
