package models

import (
	"strings"
	"testing"

	"github.com/petaki/support-go/forms"
)

func TestUserGravatarNormalizesEmail(t *testing.T) {
	expected := (&User{Email: "user@example.com"}).Gravatar(64)
	got := (&User{Email: "  User@Example.COM "}).Gravatar(64)

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}

	if expected != "https://gravatar.com/avatar/b58996c504c5638798eb6b511e6f49af?s=64" {
		t.Fatalf("unexpected gravatar url %q", expected)
	}
}

func TestUserRulesRejectLongPassword(t *testing.T) {
	for _, rules := range []func(*forms.Form){UserCreateRules, UserUpdateRules} {
		form := forms.New(map[string]any{
			"username":   "bob",
			"email":      "bob@example.com",
			"password":   strings.Repeat("a", 73),
			"is_enabled": true,
		})

		rules(form)

		if form.IsValid() {
			t.Fatal("expected a password over 72 bytes to be rejected")
		}
	}
}

func TestUserRulesCountPasswordBytes(t *testing.T) {
	form := forms.New(map[string]any{
		"username":   "bob",
		"email":      "bob@example.com",
		"password":   strings.Repeat("é", 37),
		"is_enabled": true,
	})

	UserCreateRules(form)

	errors := form.Errors["password"]
	if len(errors) != 1 || errors[0] != "The field may not be greater than 72 bytes." {
		t.Fatalf("expected the byte limit error, got %v", form.Errors)
	}
}
