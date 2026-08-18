package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONB json.RawMessage

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	if !json.Valid(j) {
		return nil, fmt.Errorf("invalid jsonb value")
	}
	return string(j), nil
}

func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = JSONB(append([]byte(nil), v...))
	case string:
		*j = JSONB([]byte(v))
	default:
		return fmt.Errorf("cannot scan %T into JSONB", value)
	}

	return nil
}

func (j JSONB) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return json.RawMessage(j).MarshalJSON()
}

func (j *JSONB) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*j = nil
		return nil
	}
	if !json.Valid(data) {
		return fmt.Errorf("invalid jsonb value")
	}
	*j = JSONB(append([]byte(nil), data...))
	return nil
}

func (j JSONB) ToRawMessage() json.RawMessage {
	if len(j) == 0 {
		return nil
	}
	return json.RawMessage(j)
}
