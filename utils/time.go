package utils

import (
	"database/sql/driver"
	"time"
)

type Time struct {
}

func (t *Time) Match(v driver.Value) bool {
	_, ok := v.(time.Time)

	return ok
}
