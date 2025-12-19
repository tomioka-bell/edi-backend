package utils

import (
	"encoding/json"
	"fmt"
)

func ToStringSlice(v interface{}) ([]string, error) {
	switch t := v.(type) {
	case []string:
		return t, nil
	case []interface{}:
		out := make([]string, 0, len(t))
		for _, it := range t {
			s, ok := it.(string)
			if !ok {
				return nil, fmt.Errorf("requirements contains non-string value")
			}
			out = append(out, s)
		}
		return out, nil
	case string:
		var out []string
		if err := json.Unmarshal([]byte(t), &out); err != nil {
			return nil, fmt.Errorf("requirements string is not valid JSON array: %w", err)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported requirements type: %T", v)
	}
}
