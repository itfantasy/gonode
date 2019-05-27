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

func UcFirst(str string) string {
	var upperStr string
	cp := []rune(str)
	for i := 0; i < len(cp); i++ {
		if i == 0 {
			if cp[i] >= 97 && cp[i] <= 122 {
				cp[i] -= 32
				upperStr += string(cp[i])
			} else {
				return str
			}
		} else {
			upperStr += string(cp[i])
		}
	}
	return upperStr
}

func UcWords(str string) string {
	var upperStr string
	strs := strings.Split(str, " ")
	for _, item := range strs {
		upperStr += UcFirst(item)
	}
	return upperStr
}
