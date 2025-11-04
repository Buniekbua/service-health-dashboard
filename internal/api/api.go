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
	http.HandleFunc("/urls", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			AddURLHandler(s).ServeHTTP(w, r)
		case http.MethodDelete:
			RemoveURLHandler(s).ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
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

func AddURLHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			URL string `json:"url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		s.AddURL(req.URL)
		w.WriteHeader(http.StatusCreated)
	}
}

func RemoveURLHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			URL string `json:"url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		s.RemoveURL(req.URL)
		w.WriteHeader(http.StatusOK)
	}
}
