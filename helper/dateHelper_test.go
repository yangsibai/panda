package helper

import "testing"

func TestGetFullDate(t *testing.T) {
	var s string
	s = GetFullDate()
	if s == "" {
		t.Error("Expected a date string, got", s)
	}
}
