package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// CustomTime is a custom time type that formats as "YYYY-MM-DD HH:MM:SS" in JSON
type CustomTime struct {
	time.Time
}

// MarshalJSON formats time as "2006-01-02 15:04:05"
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf(`"%s"`, ct.Time.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// UnmarshalJSON parses time from "2006-01-02 15:04:05" format
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		ct.Time = time.Time{}
		return nil
	}

	// Remove quotes
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// Try to parse with multiple formats
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		time.RFC3339,
	}

	var err error
	for _, format := range formats {
		ct.Time, err = time.Parse(format, str)
		if err == nil {
			return nil
		}
	}

	return err
}

// Value implements driver.Valuer for database
func (ct CustomTime) Value() (driver.Value, error) {
	if ct.Time.IsZero() {
		return nil, nil
	}
	return ct.Time, nil
}

// Scan implements sql.Scanner for database
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		ct.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		ct.Time = v
		return nil
	case []byte:
		t, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		ct.Time = t
		return nil
	case string:
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		ct.Time = t
		return nil
	}

	return fmt.Errorf("cannot scan type %T into CustomTime", value)
}
