package db

import (
    "database/sql"
    "fmt"
    "time"
	"2-serverControl/logger"
)

type DB struct {
    connection *sql.DB
}

func NewDB(host, port, user, password, dbname string) (*DB, error) {
	connectionStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{connection: db}, nil
}

// функция SavePingResult принадлежит типу указатель на структуру
func (db *DB) SavePingResult(url, status string) error {
    query := `
        INSERT INTO server_status (server_url, status, created_at)
        VALUES ($1, $2, $3)`

    _, err := db.connection.Exec(query, url, status, time.Now())
    if err != nil {
        logger.Error.Printf("ошибка при выполнении запроса: %v", err)
    }
    return nil
}

func (db *DB) Close() error {
	return db.connection.Close()
}