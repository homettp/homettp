package models

import "testing"

func TestCommandStartsWithPayload(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"%p", true},
		{" %p", true},
		{"%p arg", true},
		{"\t%p.sh", true},
		{"echo %p", false},
		{"ls", false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			command := &Command{Value: tt.value}

			if command.StartsWithPayload() != tt.expected {
				t.Fatalf("expected %v for %q", tt.expected, tt.value)
			}
		})
	}
}
