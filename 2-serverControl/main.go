package main

import (
    "fmt"
    "2-serverControl/config"
    "2-serverControl/db"
    "2-serverControl/logger"
    "2-serverControl/utils"
    "time"
    _ "github.com/lib/pq"
	"runtime" // горутины в runtime
    "strings"
    "sync"
)


func main() {
    logger.Init()

    maxCores := 2
    runtime.GOMAXPROCS(maxCores)
    logger.Info.Println("🍀 Максимальное количество ядер: 🍀", runtime.GOMAXPROCS(0))


	cfg := config.LoadConfig()

	database, err := db.NewDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		logger.Error.Fatalf("Не удалось подключиться к базе данных: %v", err)
	} else {
        logger.Info.Println("Подключение к базе данных успешно установлено!")
    }
	defer database.Close()

	sites := []string{
		"https://telegram.org",
		"https://youtube.com",
		"https://github.com",
		"https://google.com",
		"https://yandex.ru",
		"https://stepik.org",
		"https://vk.com",
	}

    // канал для передачи пингов
    ch := make(chan string)
    var wg sync.WaitGroup

    // запуск горутины для каждого сайта в одной горутине
    for _, site := range sites {
        wg.Add(1)
        // горутина
        go func(site string) {
            defer wg.Done()
            ticker := time.NewTicker(30 * time.Second)
            defer ticker.Stop()

            for {
                select {
                case <- ticker.C:
                    // когда таймер сработает, то измеряем пинг
                    result, err := utils.PingSite(site)
                    if err != nil {
                        result = fmt.Sprintf("%s: %v", site, err)
                    }

                    // отправляем результат в канал
                    ch <- result
                }
            }
        }(site)
    }

    // всего ждем 3 минуты
    timeout := time.After(3 * time.Minute)

    // создадим еще одну горутину для сохранения резутатов
    // как только взят пинг, то происходит сохранение результатов в БД
    // но пинги могут копиться в result
    for {
        select {
        case result := <-ch:
            logger.Info.Println("result: ", result)
            index := strings.LastIndex(result, ":")
            url := result[:index]
            timeStr := strings.TrimSpace(result[index+1:])
            err := database.SavePingResult(url, timeStr)
            if err != nil {
                logger.Info.Printf("Ошибка при сохранении результата: %v", err)
            } else {
                logger.Info.Printf("Сохранение результата: %s", result)
            }
        case <-timeout: // если время закончилось
            logger.Info.Println("Мониторинг завершен.")
            return
        }
    }
    wg.Wait()
    close(ch) // закрываем канал
}
