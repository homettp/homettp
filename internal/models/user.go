package models

import (
	"crypto/md5"
	"fmt"
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
	UpdateRememberToken(*User, string) error
	Authenticate(string, string) (*User, error)
	AuthenticateByRememberToken(int, string) (*User, error)
}

func (u *User) Gravatar(size int) string {
	return fmt.Sprintf("https://gravatar.com/avatar/%x?s=%d", md5.Sum([]byte(u.Email)), size)
}

func (u *User) RememberCookie() []byte {
	return []byte(fmt.Sprintf("%v|%s", u.Id, u.RememberToken))
}
