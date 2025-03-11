package utils

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "1-dataAPI/logger"
)

type CurrencyRate struct {
    Currency string
    Rate     float64
}


func GetCurrencyRate(currency string) (float64, error) {
    var rate float64
    var url string

    switch currency {
    case "BTC":
        url = "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"
    default:
        url = fmt.Sprintf("https://www.cbr-xml-daily.ru/daily_json.js") // остальное возьмем у ЦБ PФ
    }

    resp, err := http.Get(url)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close() // закрываем http-запрос

    var result map[string]interface{} // отображение, где значение может быть любым типом данных
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return 0, err
    }
    logger.Info.Println("JSON:", result)

    if currency == "BTC" {
        logger.Info.Println("bitcoin", result["bitcoin"])
        // приводим bitcon в тип string а его значение в тип interface{} с помощью .(map[string]interface{})
        // и забираем значение "usd" в типе float64
        rate = result["bitcoin"].(map[string]interface{})["usd"].(float64) // result["bitcoin"] = еще одна map
    } else {
        valute, ok := result["Valute"].(map[string]interface{})
        if !ok {
            return 0, fmt.Errorf("неверный формат данных для Valute")
        }

        currencyData, ok := valute[currency].(map[string]interface{})
        if !ok {
            return 0, fmt.Errorf("валюта %s не найдена", currency)
        }

        rate, ok = currencyData["Value"].(float64)
        if !ok {
            return 0, fmt.Errorf("неверный формат курса для %s", currency)
        }
    }
    

    return rate, nil
}

// сохраняем в бд
func SaveRateToDB(db *sql.DB, currency string, rate float64) error {
    query := `INSERT INTO currency_rates (valute, rate) VALUES ($1, $2)`
    _, err := db.Exec(query, currency, rate)
    return err
}