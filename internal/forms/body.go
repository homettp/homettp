package forms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DecodeBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value := r.Header.Get("Content-Type")

		if !strings.Contains(value, "/json") && !strings.Contains(value, "+json") {
			return &Error{
				Status: http.StatusUnsupportedMediaType,
				Msg:    "Content-Type header contains an invalid value",
			}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return &Error{
				Status: http.StatusBadRequest,
				Msg:    fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset),
			}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return &Error{
				Status: http.StatusBadRequest,
				Msg:    "Request body contains badly-formed JSON",
			}

		case errors.As(err, &unmarshalTypeError):
			return &Error{
				Status: http.StatusBadRequest,
				Msg: fmt.Sprintf(
					"Request body contains an invalid value for the %q field (at position %d)",
					unmarshalTypeError.Field,
					unmarshalTypeError.Offset,
				),
			}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")

			return &Error{
				Status: http.StatusBadRequest,
				Msg:    fmt.Sprintf("Request body contains unknown field %s", fieldName),
			}

		case errors.Is(err, io.EOF):
			return &Error{
				Status: http.StatusBadRequest,
				Msg:    "Request body must not be empty",
			}

		case err.Error() == "http: request body too large":
			return &Error{
				Status: http.StatusRequestEntityTooLarge,
				Msg:    "Request body must not be larger than 1MB",
			}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return &Error{
			Status: http.StatusBadRequest,
			Msg:    "Request body must only contain a single JSON object",
		}
	}

	return nil
}
