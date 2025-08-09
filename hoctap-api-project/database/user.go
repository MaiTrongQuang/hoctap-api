package database

import (
	"database/sql"
	"fmt"
	"time"
)

// User represents a user in the database
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository handles user database operations
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{db: DB}
}

// GetAllUsers retrieves all users from the database
func (ur *UserRepository) GetAllUsers() ([]User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC`

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return users, nil
}

// GetUserByID retrieves a user by ID
func (ur *UserRepository) GetUserByID(id int) (*User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?`

	var user User
	err := ur.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

// CreateUser creates a new user in the database
func (ur *UserRepository) CreateUser(name, email string) (*User, error) {
	// Check if email already exists
	if exists, err := ur.emailExists(email); err != nil {
		return nil, fmt.Errorf("failed to check email existence: %v", err)
	} else if exists {
		return nil, fmt.Errorf("user with email '%s' already exists", email)
	}

	query := `INSERT INTO users (name, email) VALUES (?, ?)`

	result, err := ur.db.Exec(query, name, email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	// Retrieve the created user
	return ur.GetUserByID(int(id))
}

// UpdateUser updates an existing user
func (ur *UserRepository) UpdateUser(id int, name, email string) (*User, error) {
	// Check if user exists
	if _, err := ur.GetUserByID(id); err != nil {
		return nil, err
	}

	// Check if email already exists for another user
	if exists, err := ur.emailExistsForOtherUser(email, id); err != nil {
		return nil, fmt.Errorf("failed to check email existence: %v", err)
	} else if exists {
		return nil, fmt.Errorf("user with email '%s' already exists", email)
	}

	query := `UPDATE users SET name = ?, email = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	_, err := ur.db.Exec(query, name, email, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// Retrieve the updated user
	return ur.GetUserByID(id)
}

// DeleteUser deletes a user by ID
func (ur *UserRepository) DeleteUser(id int) error {
	// Check if user exists
	if _, err := ur.GetUserByID(id); err != nil {
		return err
	}

	query := `DELETE FROM users WHERE id = ?`

	result, err := ur.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

// GetUsersCount returns the total number of users
func (ur *UserRepository) GetUsersCount() (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := ur.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %v", err)
	}

	return count, nil
}

// Helper function to check if email exists
func (ur *UserRepository) emailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	var count int
	err := ur.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Helper function to check if email exists for another user
func (ur *UserRepository) emailExistsForOtherUser(email string, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ? AND id != ?`

	var count int
	err := ur.db.QueryRow(query, email, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// SeedUsers creates some initial users for testing
func (ur *UserRepository) SeedUsers() error {
	// Check if users already exist
	count, err := ur.GetUsersCount()
	if err != nil {
		return err
	}

	// Only seed if no users exist
	if count > 0 {
		return nil
	}

	initialUsers := []struct {
		Name  string
		Email string
	}{
		{"John Doe", "john@example.com"},
		{"Jane Smith", "jane@example.com"},
		{"Alice Johnson", "alice@example.com"},
	}

	for _, user := range initialUsers {
		_, err := ur.CreateUser(user.Name, user.Email)
		if err != nil {
			return fmt.Errorf("failed to seed user %s: %v", user.Name, err)
		}
	}

	return nil
}
