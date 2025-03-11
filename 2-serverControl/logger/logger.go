package logger

import (
    "log"
    "os"
)

var (
    Info  *log.Logger
    Error *log.Logger
)

func Init() {
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("Ошибка при открытии файла логов:", err)
    }
    Info = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    Error = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}