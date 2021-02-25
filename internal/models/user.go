package models

import (
	"crypto/md5"
	"fmt"

	"github.com/petaki/support-go/forms"
)

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
	form.MatchesPattern("username", forms.UsernameRegexp)
	form.Min("username", 3)
	form.Max("username", 20)
	form.MatchesPattern("email", forms.EmailRegexp)
	form.Min("password", 8)
}

// UserUpdateRules function.
func UserUpdateRules(form *forms.Form) {
	form.Required("username", "email", "is_enabled")
	form.MatchesPattern("username", forms.UsernameRegexp)
	form.Min("username", 3)
	form.Max("username", 20)
	form.MatchesPattern("email", forms.EmailRegexp)
	form.Min("password", 8)
}

// Fill function.
func (u *User) Fill(form *forms.Form) *User {
	u.Username = form.Data["username"].(string)
	u.Email = form.Data["email"].(string)
	u.IsEnabled = form.Data["is_enabled"].(bool)
	u.Password = []byte(form.Data["password"].(string))

	return u
}

// Gravatar function.
func (u *User) Gravatar(size int) string {
	return fmt.Sprintf("https://gravatar.com/avatar/%x?s=%d", md5.Sum([]byte(u.Email)), size)
}

// RememberCookie function.
func (u *User) RememberCookie() []byte {
	return []byte(fmt.Sprintf("%v|%s", u.ID, u.RememberToken))
}
