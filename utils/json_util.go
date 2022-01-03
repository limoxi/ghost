package utils

import (
	"encoding/json"
)

func Decode(jsonStr string, container interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), container)
	return err
}
