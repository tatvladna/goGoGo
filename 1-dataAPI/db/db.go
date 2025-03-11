package db

import (
    "database/sql"
    "fmt"
    "1-dataAPI/logger"
    _ "github.com/lib/pq"
)

func ConnectDB(psqlInfo string) (*sql.DB, error) {

    logger.Info.Println(psqlInfo)
    // подключение к бд
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        logger.Error.Printf("Ошибка при подключении к базе данных: %v\n", err)
        return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
    }

    // проверка подключения
    err = db.Ping()
    if err != nil {
        logger.Error.Printf("Ошибка при проверке подключения к базе данных: %v\n", err)
        return nil, fmt.Errorf("ошибка при проверке подключения к базе данных: %w", err)
    }

    logger.Info.Println("Успешное подключение к базе данных!")
    return db, nil
}