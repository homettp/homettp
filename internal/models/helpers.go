package models

import (
	"fmt"
	"sync"

	"github.com/petaki/support-go/forms"
	"golang.org/x/crypto/bcrypt"
)

var dummyPasswordHash = sync.OnceValue(func() []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte("homettp"), 12)

	return hash
})

func isString(form *forms.Form, fields ...string) {
	for _, field := range fields {
		if _, ok := form.Data[field].(string); !ok && form.Data[field] != nil {
			form.Errors.Add(field, "The field must be a string.")
		}
	}
}

func isBool(form *forms.Form, fields ...string) {
	for _, field := range fields {
		if _, ok := form.Data[field].(bool); !ok && form.Data[field] != nil {
			form.Errors.Add(field, "The field must be a boolean.")
		}
	}
}

func isCommandImage(form *forms.Form, field string) {
	if image, ok := form.Data[field].(string); ok && image != "" && !CommandImage(image).IsValid() {
		form.Errors.Add(field, "The selected image is invalid.")
	}
}

func rejectCredentials(password string) error {
	_ = bcrypt.CompareHashAndPassword(dummyPasswordHash(), []byte(password))

	return ErrInvalidCredentials
}

func maxBytes(form *forms.Form, field string, max int) {
	if value, ok := form.Data[field].(string); ok && len(value) > max {
		form.Errors.Add(field, fmt.Sprintf("The field may not be greater than %d bytes.", max))
	}
}
