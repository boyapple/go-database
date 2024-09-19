package filter

import (
	"database/sql"
	"time"
)

func NewTimeFieldFilter(fields []string) Filter {
	m := make(TimeField)
	for _, field := range fields {
		m[field] = struct{}{}
	}
	return m
}

// TimeField 时间字段map
type TimeField map[string]struct{}

// Type 时间字段类型过滤
func (t TimeField) Type(name string) (interface{}, bool) {
	_, ok := t[name]
	if !ok {
		return nil, false
	}
	return &sql.NullTime{}, true
}

// Value 时间字段值过滤
func (t TimeField) Value(name string, value interface{}) (interface{}, bool) {
	_, ok := t[name]
	if !ok {
		return nil, false
	}
	switch v := value.(type) {
	case time.Time:
		return v.Unix(), true
	case *sql.NullTime:
		return v.Time.Unix(), true
	case int64:
		return time.Unix(v, 0), true
	default:
		return nil, false
	}
}
