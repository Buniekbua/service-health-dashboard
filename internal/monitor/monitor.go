package monitor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/buniekbua/service-health-dashboard/internal/storage"
)

func CheckStatus(client *http.Client, url string) (int, error) {
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func monitorURL(client *http.Client, url string, s *storage.Storage) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		status, err := CheckStatus(client, url)
		if err != nil {
			fmt.Println("Request error for ", url, ":", err)
			s.UpdateStatus(url, 0)
		} else {
			s.UpdateStatus(url, status)
			fmt.Println("Status for ", url, " is", status)
		}
		<-ticker.C
	}
}

func StartMonitoring(urls []string, s *storage.Storage) {
	client := &http.Client{}
	for _, url := range urls {
		go monitorURL(client, url, s)
	}
}
