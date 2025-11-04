package storage

import "testing"

func TestStorage(t *testing.T) {
	s := NewStorage()

	s.UpdateStatus("http://test.url", 200)
	statuses := s.GetAllStatuses()

	if status, ok := statuses["http://test.url"]; !ok {
		t.Errorf("No status for expected URL")
	} else if status != 200 {
		t.Errorf("Expected: 200, given: %d", status)
	}
}
