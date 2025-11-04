package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/buniekbua/service-health-dashboard/internal/monitor"
	"github.com/buniekbua/service-health-dashboard/internal/storage"
)

func TestStatusEndpoint_TableDriven(t *testing.T) {
	cases := []struct {
		name     string
		setup    func() *storage.Storage
		expected map[string]int
	}{
		{
			name: "empty storage",
			setup: func() *storage.Storage {
				return storage.NewStorage()
			},
			expected: map[string]int{},
		},
		{
			name: "one url 200",
			setup: func() *storage.Storage {
				s := storage.NewStorage()
				s.UpdateStatus("https://test.url", 200)
				return s
			},
			expected: map[string]int{"https://test.url": 200},
		},
		{
			name: "two urls mixed",
			setup: func() *storage.Storage {
				s := storage.NewStorage()
				s.UpdateStatus("https://up.url", 200)
				s.UpdateStatus("https://fail.url", 0)
				return s
			},
			expected: map[string]int{
				"https://up.url":   200,
				"https://fail.url": 0,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.setup()
			req := httptest.NewRequest(http.MethodGet, "/status", nil)
			w := httptest.NewRecorder()

			handler := StatusHandler(s)

			handler.ServeHTTP(w, req)

			res := w.Result()
			if res.StatusCode != http.StatusOK {
				t.Fatalf("Expected status 200, got %d", res.StatusCode)
			}

			var data map[string]int
			err := json.NewDecoder(res.Body).Decode(&data)
			if err != nil {
				t.Fatalf("JSON decode error: %v", err)
			}
			if len(data) != len(tc.expected) {
				t.Errorf("Expected %d items, got %d", len(tc.expected), len(data))
			}

			for url, status := range tc.expected {
				if data[url] != status {
					t.Errorf("URL: %s - expected %d, got %d", url, status, data[url])
				}
			}
		})
	}
}

func TestStatusEnpoint(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{
			name:       "valid GET",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid POST",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	s := storage.NewStorage()
	s.UpdateStatus("https://www.google.com", 200)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/status", nil)
			w := httptest.NewRecorder()

			handler := StatusHandler(s)
			handler.ServeHTTP(w, req)

			res := w.Result()
			if res.StatusCode != tc.wantStatus {
				t.Errorf("Expected status %d, got %d", tc.wantStatus, res.StatusCode)
			}
			if tc.wantStatus == http.StatusOK {
				var data map[string]int
				if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
					t.Errorf("Failed to decode JSON: %v", err)
				}
				if data["https://www.google.com"] != 200 {
					t.Errorf("Expected status 200 for https://www.google.com, got %d", data["https://www.google.com"])
				}
			}
		})
	}
}

func TestIntegration_MonitorAndAPI(t *testing.T) {
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testSrv.Close()

	s := storage.NewStorage()

	client := &http.Client{}

	status, err := monitor.CheckStatus(client, testSrv.URL)
	if err != nil {
		t.Fatalf("Monitor error: %v", err)
	}
	s.UpdateStatus(testSrv.URL, status)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	handler := StatusHandler(s)
	handler.ServeHTTP(w, req)

	// Walidacja odpowiedzi
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Status code: %d", res.StatusCode)
	}
	var data map[string]int
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if data[testSrv.URL] != 200 {
		t.Errorf("Expected status 200, got %d", data[testSrv.URL])
	}
}

func TestAddURLHandler(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantInURLs []string
	}{
		{
			name:       "valid POST",
			method:     http.MethodPost,
			body:       `{"url":"https://example.com"}`,
			wantStatus: http.StatusCreated,
			wantInURLs: []string{"https://example.com"},
		},
		{
			name:       "empty body",
			method:     http.MethodPost,
			body:       "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			method:     http.MethodPost,
			body:       `{"url":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "wrong method GET",
			method:     http.MethodGet,
			body:       ``,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := storage.NewStorage()

			req := httptest.NewRequest(tc.method, "/urls", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler := AddURLHandler(s)
			handler.ServeHTTP(w, req)

			res := w.Result()
			if res.StatusCode != tc.wantStatus {
				t.Errorf("Expected status: %d, got: %d", tc.wantStatus, res.StatusCode)
			}

			if tc.wantStatus == http.StatusCreated {
				urls := s.GetURLs()
				for _, expectedURL := range tc.wantInURLs {
					found := false
					for _, u := range urls {
						if u == expectedURL {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected URL %q to be added, but it was not present", expectedURL)
					}
				}
			}
		})
	}
}
