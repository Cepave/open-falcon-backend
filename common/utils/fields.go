package utils

import (
	"strings"
)

func DictedFieldstring(s string) map[string]interface{} {
	if s == "" {
		return map[string]interface{}{}
	}

	field_dict := make(map[string]interface{})
	fields := strings.Split(s, ",")
	for _, field := range fields {
		field_pair := strings.SplitN(field, "=", 2)
		if len(field_pair) == 2 {
			key := strings.TrimSpace(field_pair[0])
			val := strings.TrimSpace(field_pair[1])
			field_dict[key] = val
		}
	}
	return field_dict
}
