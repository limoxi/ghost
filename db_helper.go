package ghost

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type Json string

func (jsonStr Json) Value() (driver.Value, error) {
	content := []byte(jsonStr)
	return content, nil
}

func (jsonStr *Json) Scan(value interface{}) error {
	if value == nil {
		*jsonStr = ""
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	*jsonStr = Json(bytes)
	return nil
}
