package api

import (
	"encoding/json"
	"net/http"

	"github.com/buniekbua/service-health-dashboard/internal/storage"
)

type StatusResponse struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}

func StartServer(s *storage.Storage) error {
	http.HandleFunc("/status", StatusHandler(s))
	return http.ListenAndServe(":8080", nil)
}

func StatusHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		statuses := s.GetAllStatuses()
		json.NewEncoder(w).Encode(statuses)
	}
}
