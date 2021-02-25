package models

import (
	"fmt"
	"time"
)

// Timestamp type.
type Timestamp time.Time

// MarshalJSON function.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	value := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))

	return []byte(value), nil
}

// RedisArg function.
func (t Timestamp) RedisArg() interface{} {
	return time.Time(t).Format(time.RFC3339)
}

// RedisScan function.
func (t *Timestamp) RedisScan(src interface{}) error {
	bs, ok := src.([]byte)
	if !ok {
		return ErrTimestamp
	}

	parsed, err := time.Parse(time.RFC3339, string(bs))
	if err != nil {
		return err
	}

	*t = Timestamp(parsed)

	return nil
}
