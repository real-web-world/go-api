package cache

import "time"

// todo wait complete
type Cache interface {
	Get() (interface{}, error)
	Set(k string, v interface{}) error
	SetEx(k string, v interface{}, ex time.Duration) error
	Del(k string) error
	Inc(k string) error
	Dec(k string) error
}
