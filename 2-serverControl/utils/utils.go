package utils

import (
	"fmt"
	"net/http"
	"time"
)

func PingSite(url string) (string, error) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	return fmt.Sprintf("%s: %s", url, duration), nil
}