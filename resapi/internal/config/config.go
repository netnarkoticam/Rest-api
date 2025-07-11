package config

import (
	"flag"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

// Config хранит конфигурацию сервера
type Config struct {
	Address string
}

// ReadConfig читает конфигурацию из флагов командной строки
func ReadConfig() Config {
	var address string
	flag.StringVar(&address, "address", ":8080", "Адрес сервера")
	flag.Parse()
	return Config{Address: address}
}

// ConnectDB устанавливает соединение с базой данных PostgreSQL
func ConnectDB() (*sql.DB, error) {
	// Жёстко прописанные параметры из docker-compose.yml
	host := "db"
	port := 5432
	user := "user"
	password := "password"
	dbname := "usersdb"

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение
	return db, db.Ping()
}