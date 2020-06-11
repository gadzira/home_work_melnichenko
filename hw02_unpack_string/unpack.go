package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString ...
var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	// Check if a string is empty
	if s == "" {
		return "", nil
	}
	// Check if a first simbol not digit
	if unicode.IsDigit([]rune(s)[0]) {
		return "", ErrInvalidString
	}

	// Check if a string contains more than one number
	twoOrMoreNumber := regexp.MustCompile(`[0-9]{2,}`)
	if twoOrMoreNumber.MatchString(s) {
		return "", ErrInvalidString
	}

	// Check if a string haven't some digis
	lessOneNumber := regexp.MustCompile(`[0-9]{1,}`)
	if !lessOneNumber.MatchString(s) {
		return s, nil
	}

	var sb strings.Builder
	for i, v := range s {
		if unicode.IsDigit(v) {
			sv, _ := strconv.Atoi(string(v))
			sb.WriteString(strings.Repeat(string(s[i-1]), sv-1))
		} else {
			sb.WriteString(string(v))
		}
	}
	return sb.String(), nil
}
