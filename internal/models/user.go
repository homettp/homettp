package models

import (
	"crypto/md5"
	"fmt"

	"github.com/homettp/homettp/internal/forms"
)

type User struct {
	Id            int       `json:"id" redis:"id"`
	Username      string    `json:"username" redis:"username"`
	Email         string    `json:"email" redis:"email"`
	Password      []byte    `json:"-" redis:"password"`
	RememberToken string    `json:"-" redis:"remember_token"`
	IsEnabled     bool      `json:"is_enabled" redis:"is_enabled"`
	CreatedAt     Timestamp `json:"created_at" redis:"created_at"`
}

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

func NewUser() *User {
	return &User{
		IsEnabled: true,
	}
}

func UserCreateRules(form *forms.Form) {
	form.Required("username", "email", "password", "is_enabled")
	form.MatchesPattern("username", forms.UsernameRegex)
	form.Min("username", 3)
	form.Max("username", 20)
	form.MatchesPattern("email", forms.EmailRegex)
	form.Min("password", 8)
}

func UserUpdateRules(form *forms.Form) {
	form.Required("username", "email", "is_enabled")
	form.MatchesPattern("username", forms.UsernameRegex)
	form.Min("username", 3)
	form.Max("username", 20)
	form.MatchesPattern("email", forms.EmailRegex)
	form.Min("password", 8)
}

func (u *User) Fill(form *forms.Form) *User {
	u.Username = form.Data["username"].(string)
	u.Email = form.Data["email"].(string)
	u.IsEnabled = form.Data["is_enabled"].(bool)
	u.Password = []byte(form.Data["password"].(string))

	return u
}

func (u *User) Gravatar(size int) string {
	return fmt.Sprintf("https://gravatar.com/avatar/%x?s=%d", md5.Sum([]byte(u.Email)), size)
}

func (u *User) RememberCookie() []byte {
	return []byte(fmt.Sprintf("%v|%s", u.Id, u.RememberToken))
}
