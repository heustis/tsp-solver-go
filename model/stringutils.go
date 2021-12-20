package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ToString converts an object to a string.
func ToString(value interface{}) string {
	if p, okay := value.(Printable); okay {
		return p.ToString()
	} else if jsonBytes, err := json.Marshal(value); err == nil && strings.Compare("null", string(jsonBytes)) != 0 {
		return string(jsonBytes)
	} else {
		return fmt.Sprintf("%v", value)
	}
}
