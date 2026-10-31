package models

import (
	"testing"

	"github.com/petaki/support-go/forms"
)

func TestRulesRejectWrongTypes(t *testing.T) {
	tests := []struct {
		name  string
		rules func(*forms.Form)
		data  map[string]any
		field string
	}{
		{"command name", CommandCreateRules, map[string]any{"name": 1.0, "image": "light", "value": "id"}, "name"},
		{"command image", CommandUpdateRules, map[string]any{"name": "a", "image": "bogus", "value": "id"}, "image"},
		{"user is_enabled", UserCreateRules, map[string]any{"username": "bob", "email": "bob@example.com", "password": "password123", "is_enabled": "yes"}, "is_enabled"},
		{"user password", UserUpdateRules, map[string]any{"username": "bob", "email": "bob@example.com", "password": 12345678.0, "is_enabled": true}, "password"},
		{"login remember", LoginRules, map[string]any{"username_or_email": "bob", "password": "secret", "remember": "yes"}, "remember"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := forms.New(tt.data)
			tt.rules(form)

			if len(form.Errors[tt.field]) == 0 {
				t.Fatalf("expected an error for %q, got %v", tt.field, form.Errors)
			}
		})
	}
}

func TestFillWithoutOptionalFields(t *testing.T) {
	form := forms.New(map[string]any{"username": "bob", "email": "bob@example.com", "is_enabled": true})
	UserUpdateRules(form)

	if !form.IsValid() {
		t.Fatalf("expected a valid form, got %v", form.Errors)
	}

	user := NewUser().Fill(form)
	if len(user.Password) != 0 {
		t.Fatalf("expected an empty password, got %q", user.Password)
	}
}
