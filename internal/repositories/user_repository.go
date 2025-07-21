package repositories

import (
	"github.com/google/uuid"
	"github.com/raingrave/apirest/internal"
	"github.com/raingrave/apirest/internal/models"
)

func CreateUser(user models.User) (uuid.UUID, error) {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id uuid.UUID
	err := internal.DB.QueryRow(query, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func GetUser(id uuid.UUID) (models.User, error) {
	var user models.User
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	err := internal.DB.Get(&user, query, id)
	return user, err
}

func UpdateUser(id uuid.UUID, user models.User) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err := internal.DB.Exec(query, user.Name, user.Email, id)
	return err
}

func DeleteUser(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := internal.DB.Exec(query, id)
	return err
}

func ListUsers() ([]models.User, error) {
	var users []models.User
	query := `SELECT id, name, email, created_at FROM users`
	err := internal.DB.Select(&users, query)
	return users, err
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1`
	err := internal.DB.Get(&user, query, email)
	return user, err
}
