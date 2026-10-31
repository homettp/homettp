//go:build windows

package web

import "testing"

func TestIsSafePayload(t *testing.T) {
	tests := []struct {
		payload  string
		expected bool
	}{
		{"", true},
		{`C:\data\file.txt`, true},
		{`C:\data\`, false},
		{`a\\`, false},
		{"a&b", false},
		{`a"b`, false},
	}

	for _, tt := range tests {
		t.Run(tt.payload, func(t *testing.T) {
			if got := isSafePayload(tt.payload); got != tt.expected {
				t.Fatalf("isSafePayload(%q) = %v, expected %v", tt.payload, got, tt.expected)
			}
		})
	}
}
