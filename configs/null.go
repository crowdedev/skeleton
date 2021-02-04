package configs

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"time"
)

type NullString struct {
	sql.NullString
	GetValue string
	Valid    bool // Valid is true if String is not NULL
}

func NewNullString(value string, valid bool) NullString {
	return NullString{GetValue: value, Valid: valid}
}

func (n *NullString) Scan(value interface{}) error {
	if isNil(value) {
		n.GetValue, n.Valid = "", false
		return nil
	}
	n.Valid = true
	if err := n.NullString.Scan(value); err != nil {
		return err
	}
	n.GetValue = n.NullString.String
	return nil
}

func (n NullString) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.GetValue, nil
}

func (n NullString) Get() string {
	return n.GetValue
}

func (n *NullString) Set(value string, valid bool) {
	n.GetValue, n.Valid = value, valid
}

type NullInt32 struct {
	sql.NullInt32
	GetValue int32
	Valid    bool // Valid is true if Int32 is not NULL
}

func NewNullInt32(value int32, valid bool) NullInt32 {
	return NullInt32{GetValue: value, Valid: valid}
}

func (n *NullInt32) Scan(value interface{}) error {
	if isNil(value) {
		n.GetValue, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	if err := n.NullInt32.Scan(value); err != nil {
		return err
	}
	n.GetValue = n.NullInt32.Int32
	return nil
}

func (n NullInt32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.GetValue), nil
}

func (n NullInt32) Get() int32 {
	return n.GetValue
}

func (n *NullInt32) Set(value int32, valid bool) {
	n.GetValue, n.Valid = value, valid
}

type NullTime struct {
	sql.NullTime
	Seconds int64
	Nanos   int32
	Valid   bool // Valid is true if Time is not NULL
}

func NewNullTime(t time.Time, valid bool) NullTime {
	return NullTime{Seconds: int64(t.Unix()), Nanos: int32(t.Nanosecond()), Valid: valid}
}

func (n *NullTime) Scan(value interface{}) error {
	if isNil(value) {
		n.Seconds, n.Nanos, n.Valid = 0, 0, false
		return nil
	}
	n.Valid = true
	if err := n.NullTime.Scan(value); err != nil {
		return err
	}
	if n.NullTime.Valid {
		n.Seconds = n.NullTime.Time.Unix()
		n.Nanos = int32(n.NullTime.Time.Nanosecond())
	}
	return nil
}

func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return time.Unix(int64(n.Seconds), int64(n.Nanos)).UTC(), nil
}

func (n NullTime) Get() time.Time {
	return time.Unix(int64(n.Seconds), int64(n.Nanos)).UTC()
}

func (n *NullTime) Set(t time.Time, valid bool) {
	n.Seconds = int64(t.Unix())
	n.Nanos = int32(t.Nanosecond())
	n.Valid = valid
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	// detect (*T)(nil)
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
