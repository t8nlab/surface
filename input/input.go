package sfInput

import (
	"errors"
	"fmt"
)

func GetString(m map[string]any, key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", errors.New("missing " + key)
	}

	s, ok := v.(string)
	if !ok {
		return "", errors.New("invalid type for " + key)
	}

	return s, nil
}

func GetBool(m map[string]any, key string, def bool) bool {
	v, ok := m[key]
	if !ok {
		return def
	}
	b, ok := v.(bool)
	if !ok {
		return def
	}
	return b
}

func GetInt(m map[string]any, key string) (int, error) {
	v, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("missing %s", key)
	}

	switch val := v.(type) {
	case float64:
		return int(val), nil
	case int:
		return val, nil
	default:
		return 0, fmt.Errorf("invalid type for %s", key)
	}
}