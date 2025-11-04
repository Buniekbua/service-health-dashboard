package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := &http.Client{}

	status, err := CheckStatus(client, ts.URL)
	if err != nil {
		t.Fatalf("Expected: no error, given: %v", err)
	}
	if status != 200 {
		t.Errorf("Expected: 200, given: %d", status)
	}
}

func TestCheckStatusError(t *testing.T) {
	client := &http.Client{}

	_, err := CheckStatus(client, "http://nonexistent.invalid")
	if err == nil {
		t.Error("Error expected for non existent URL")
	}
}
