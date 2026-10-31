package models

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/petaki/support-go/forms"
)

const maxPasswordBytes = 72

// User type.
type User struct {
	ID            int       `json:"id" redis:"id"`
	Username      string    `json:"username" redis:"username"`
	Email         string    `json:"email" redis:"email"`
	Password      []byte    `json:"-" redis:"password"`
	RememberToken string    `json:"-" redis:"remember_token"`
	IsEnabled     bool      `json:"is_enabled" redis:"is_enabled"`
	CreatedAt     Timestamp `json:"created_at" redis:"created_at"`
}

// UserRepository type.
type UserRepository interface {
	Create(*User) error
	Find(int) (*User, error)
	FindAll() ([]*User, error)
	Update(*User, *User) error
	UpdateRememberToken(*User, string) error
	Authenticate(string, string) (*User, error)
	AuthenticateByRememberToken(int, string) (*User, error)
	Delete(*User) error
}

// NewUser function.
func NewUser() *User {
	return &User{
		IsEnabled: true,
	}
}

// UserCreateRules function.
func UserCreateRules(form *forms.Form) {
	form.Required("username", "email", "password", "is_enabled")
	isString(form, "username", "email", "password")
	isBool(form, "is_enabled")
	form.MatchesPattern("username", forms.UsernameRegexp)
	form.Min("username", 3)
	form.Max("username", 20)
	form.MatchesPattern("email", forms.EmailRegexp)
	form.Min("password", 8)
	maxBytes(form, "password", maxPasswordBytes)
}

// UserUpdateRules function.
func UserUpdateRules(form *forms.Form) {
	form.Required("username", "email", "is_enabled")
	isString(form, "username", "email", "password")
	isBool(form, "is_enabled")
	form.MatchesPattern("username", forms.UsernameRegexp)
	form.Min("username", 3)
	form.Max("username", 20)
	form.MatchesPattern("email", forms.EmailRegexp)
	form.Min("password", 8)
	maxBytes(form, "password", maxPasswordBytes)
}

// LoginRules function.
func LoginRules(form *forms.Form) {
	form.Required("username_or_email", "password")
	isString(form, "username_or_email", "password")
	isBool(form, "remember")
}

// Fill function.
func (u *User) Fill(form *forms.Form) *User {
	u.Username, _ = form.Data["username"].(string)
	u.Email, _ = form.Data["email"].(string)
	u.IsEnabled, _ = form.Data["is_enabled"].(bool)
	password, _ := form.Data["password"].(string)
	u.Password = []byte(password)

	return u
}

// Gravatar function.
func (u *User) Gravatar(size int) string {
	return fmt.Sprintf("https://gravatar.com/avatar/%x?s=%d", md5.Sum([]byte(strings.ToLower(strings.TrimSpace(u.Email)))), size)
}

// RememberCookie function.
func (u *User) RememberCookie() []byte {
	return fmt.Appendf(nil, "%v|%s", u.ID, u.RememberToken)
}
