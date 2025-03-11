package main

import (
    "1-dataAPI/logger"
    "1-dataAPI/config"
    "1-dataAPI/db"
    "1-dataAPI/utils"
    _ "github.com/lib/pq" // сюда тоже нужно импортировать драйвера postgresql
    "sync"
    "runtime"
    "fmt"
)

func main() {
    logger.Init()


    maxCores := 2
    runtime.GOMAXPROCS(maxCores)
    logger.Info.Println("🍀 Максимальное количество ядер: 🍀", runtime.GOMAXPROCS(0))

    // загружаем конфигурации
    cfg := config.LoadConfig()

    // формируем строку подключения
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

    // подключаемся к БД
    dbConn, err := db.ConnectDB(psqlInfo)
    if err != nil {
        logger.Error.Printf("Ошибка подключения к базе данных: %v", err)
    }
    defer dbConn.Close()

    logger.Info.Println("Подключение к базе данных успешно установлено!")


    currencies := [4]string{"EUR", "USD", "CNY", "BTC"}
    ratesChan := make(chan utils.CurrencyRate, len(currencies))  // буферизированный канал

    // WaitGroup для ожидания завершения всех горутин
    // это необходимо, чтобы основная программа не завершилась раньше, чем все горутины закончат свою работу
    var wg sync.WaitGroup

    // запускаем горутины
    for _, currency := range currencies {
        wg.Add(1) // увеличиваем счетчик на 1
        // анонимная
        // цикл не ждет завершения горутины и переходит к следующей итерации
        go func(currency string) {
            defer wg.Done() // уменьшаем счетчик при завершении горутины
            rate, err := utils.GetCurrencyRate(currency)
            if err != nil {
                logger.Error.Printf("Ошибка при получении курса для %s: %v\n", currency, err)
                return
            }
            ratesChan <- utils.CurrencyRate{Currency: currency, Rate: rate}
        }(currency) // передаем название валюты в явном виде
    }

    // ожидаем завершения всех горутин
    wg.Wait() // блокирует основной поток, пока счетчик не станет равным нулю.
    close(ratesChan) // закрываем канал после завершения всех горутин

    // сохранение курсов в бд
    for rate := range ratesChan {
        err := utils.SaveRateToDB(dbConn, rate.Currency, rate.Rate)
        if err != nil {
            logger.Error.Printf("Ошибка при сохранении курса для %s: %v\n", rate.Currency, err)
            continue
        }
        logger.Info.Printf("Курс %s: %.4f\n", rate.Currency, rate.Rate)
    }
}