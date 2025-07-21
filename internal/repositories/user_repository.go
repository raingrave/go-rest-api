package repositories

import (
	"github.com/raingrave/apirest/internal"
	"github.com/raingrave/apirest/internal/models"
)

func CreateUser(user models.User) (int64, error) {
	query := `INSERT INTO users (name, email, created_at) VALUES ($1, $2, NOW()) RETURNING id`
	var id int64
	err := internal.DB.QueryRow(query, user.Name, user.Email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetUser(id int64) (models.User, error) {
	var user models.User
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	err := internal.DB.Get(&user, query, id)
	return user, err
}

func UpdateUser(id int64, user models.User) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err := internal.DB.Exec(query, user.Name, user.Email, id)
	return err
}

func DeleteUser(id int64) error {
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
