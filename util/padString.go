package util

import "strings"

func PadStringLeft(s string, length int, char rune) string {
	padding := length - len(s)
	if padding <= 0 {
		return s
	}
	return strings.Repeat(string(char), padding) + s
}

func PadStringRight(s string, length int, char rune) string {
	padding := length - len(s)
	if padding <= 0 {
		return s
	}
	return s + strings.Repeat(string(char), padding)
}
