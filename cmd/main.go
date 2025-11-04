package main

import (
	"fmt"

	"github.com/buniekbua/service-health-dashboard/internal/api"
	"github.com/buniekbua/service-health-dashboard/internal/monitor"
	"github.com/buniekbua/service-health-dashboard/internal/storage"
)

func main() {
	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://nonexistent.example.com",
	}

	s := storage.NewStorage()
	monitor.StartMonitoring(urls, s)

	fmt.Println("Starting HTTP Server on :8080")
	if err := api.StartServer(s); err != nil {
		fmt.Println("Server error: ", err)
	}

}
