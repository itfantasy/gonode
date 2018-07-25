package strs

import (
	"strings"
)

func StartsWith(str string, substr string) bool {
	return strings.HasPrefix(str, substr)
}

func EndsWith(str string, substr string) bool {
	return strings.HasSuffix(str, substr)
}
