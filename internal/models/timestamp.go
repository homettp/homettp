package models

import "time"

type Timestamp time.Time

func (t Timestamp) RedisArg() interface{} {
	return time.Time(t).Format(time.RFC3339)
}

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
