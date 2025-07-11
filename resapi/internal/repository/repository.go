package repository

import (
	"fmt"
	"resapi/internal/domain/models"
	"database/sql"
	"strings"
)

// Repository отвечает за взаимодействие с БД
type Repository struct {
	db *sql.DB
}

// New создаёт новый репозиторий
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetAllUser возвращает всех пользователей
func (repo *Repository) GetAllUser() ([]models.User, error) {
	rows, err := repo.db.Query("SELECT id, name, login, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUser возвращает пользователя по ID
func (repo *Repository) GetUser(id int) (models.User, error) {
	var user models.User
	err := repo.db.QueryRow("SELECT id, name, login, password FROM users WHERE id=$1", id).
		Scan(&user.ID, &user.Name, &user.Login, &user.Password)

	if err == sql.ErrNoRows {
		return models.User{}, fmt.Errorf("user not found")
	}

	return user, err
}

// InsertUser добавляет нового пользователя
func (repo *Repository) InsertUser(user models.User) error {
	_, err := repo.db.Exec(
		"INSERT INTO users (name, login, password) VALUES ($1, $2, $3)",
		user.Name, user.Login, user.Password,
	)
	return err
}

// DeleteUser удаляет пользователя по ID
func (repo *Repository) DeleteUser(id int) error {
	res, err := repo.db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// UpdateUser обновляет данные пользователя
func (repo *Repository) UpdateUser(user models.User) error {
	res, err := repo.db.Exec(
		"UPDATE users SET name=$1, login=$2, password=$3 WHERE id=$4",
		user.Name, user.Login, user.Password, user.ID,
	)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// PatchUser частично обновляет данные пользователя
func (repo *Repository) PatchUser(id int, patch map[string]interface{}) (models.User, error) {
	fields := []string{}
	args := []interface{}{}
	i := 1

	if name, ok := patch["name"].(string); ok {
		fields = append(fields, fmt.Sprintf("name=$%d", i))
		args = append(args, name)
		i++
	}
	if login, ok := patch["login"].(string); ok {
		fields = append(fields, fmt.Sprintf("login=$%d", i))
		args = append(args, login)
		i++
	}
	if password, ok := patch["password"].(string); ok {
		fields = append(fields, fmt.Sprintf("password=$%d", i))
		args = append(args, password)
		i++
	}

	if len(fields) == 0 {
		return models.User{}, fmt.Errorf("no fields to update")
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id=$%d RETURNING id, name, login, password",
		strings.Join(fields, ", "), i,
	)

	var user models.User
	err := repo.db.QueryRow(query, args...).Scan(&user.ID, &user.Name, &user.Login, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}