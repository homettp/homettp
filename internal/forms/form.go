package forms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var UsernameRegex = regexp.MustCompile("^[a-zA-Z0-9\\.\\-_]+$")
var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Form struct {
	Data   map[string]interface{}
	Errors Bag
}

func New(w http.ResponseWriter, r *http.Request) (*Form, error) {
	var data map[string]interface{}

	if r.Header.Get("Content-Type") != "" {
		value := r.Header.Get("Content-Type")

		if !strings.Contains(value, "/json") && !strings.Contains(value, "+json") {
			return nil, &Error{
				Status: http.StatusUnsupportedMediaType,
				Msg:    "Content-Type header contains an invalid value",
			}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return nil, &Error{
				Status: http.StatusBadRequest,
				Msg:    fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset),
			}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return nil, &Error{
				Status: http.StatusBadRequest,
				Msg:    fmt.Sprintf("Request body contains badly-formed JSON"),
			}

		case errors.As(err, &unmarshalTypeError):
			return nil, &Error{
				Status: http.StatusBadRequest,
				Msg: fmt.Sprintf(
					"Request body contains an invalid value for the %q field (at position %d)",
					unmarshalTypeError.Field,
					unmarshalTypeError.Offset,
				),
			}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")

			return nil, &Error{
				Status: http.StatusBadRequest,
				Msg:    fmt.Sprintf("Request body contains unknown field %s", fieldName),
			}

		case errors.Is(err, io.EOF):
			return nil, &Error{
				Status: http.StatusBadRequest,
				Msg:    "Request body must not be empty",
			}

		case err.Error() == "http: request body too large":
			return nil, &Error{
				Status: http.StatusRequestEntityTooLarge,
				Msg:    "Request body must not be larger than 1MB",
			}

		default:
			return nil, err
		}
	}

	return &Form{
		data,
		map[string][]string{},
	}, nil
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		if f.Data[field] == nil || strings.TrimSpace(fmt.Sprintf("%v", f.Data[field])) == "" {
			f.Errors.Add(field, "The field is required.")
		}
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	if f.Data[field] == nil {
		return
	}

	switch value := f.Data[field].(type) {
	case string:
		if value == "" {
			return
		}

		if !pattern.MatchString(value) {
			f.Errors.Add(field, "The field format is invalid.")
		}
	}
}

func (f *Form) Min(field string, min float64) {
	if f.Data[field] == nil {
		return
	}

	switch value := f.Data[field].(type) {
	case float64:
		if value < min {
			f.Errors.Add(field, fmt.Sprintf("The field must be at least %v.", min))
		}
	}
}

func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}
