package models

// User represents a user in the system.
type User struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	Age       int    `json:"age" db:"age"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest represents the payload for creating a new user.
type CreateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Email string `json:"email" validate:"omitempty,email"`
	Age   int    `json:"age" validate:"omitempty,gte=0,lte=130"`
}

// UpdateUserRequest represents the payload for updating an existing user.
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Email string `json:"email" validate:"omitempty,email"`
	Age   int    `json:"age" validate:"omitempty,gte=0,lte=130"`
}
