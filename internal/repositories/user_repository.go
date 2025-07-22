package repositories

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/raingrave/apirest/internal"
	"github.com/raingrave/apirest/internal/models"
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	CreateUser(user models.User) (uuid.UUID, error)
	GetUser(id uuid.UUID) (models.User, error)
	UpdateUser(id uuid.UUID, user models.User) error
	DeleteUser(id uuid.UUID) error
	ListUsers() ([]models.User, error)
	GetUserByEmail(email string) (models.User, error)
}

// userRepository is the concrete implementation for PostgreSQL.
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository() UserRepository {
	return &userRepository{db: internal.DB}
}

func (r *userRepository) CreateUser(user models.User) (uuid.UUID, error) {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id uuid.UUID
	err := r.db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *userRepository) GetUser(id uuid.UUID) (models.User, error) {
	var user models.User
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	return user, err
}

func (r *userRepository) UpdateUser(id uuid.UUID, user models.User) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err := r.db.Exec(query, user.Name, user.Email, id)
	return err
}

func (r *userRepository) DeleteUser(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *userRepository) ListUsers() ([]models.User, error) {
	var users []models.User
	query := `SELECT id, name, email, created_at FROM users`
	err := r.db.Select(&users, query)
	return users, err
}

func (r *userRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1`
	err := r.db.Get(&user, query, email)
	return user, err
}