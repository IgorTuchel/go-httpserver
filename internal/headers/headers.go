package headers

import (
	"errors"
	"slices"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := strings.Index(string(data), "\r\n")
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	headerline := string(data)[:idx]
	header := strings.SplitN(headerline, ":", 2)
	if len(header) != 2 {
		return 0, false, errors.New("malformed header")
	}
	if header[0][len(header[0])-1:] == " " {
		return 0, false, errors.New("malformed header")
	}
	key, value := strings.TrimSpace(header[0]), strings.TrimSpace(header[1])
	if !validToken([]byte(key)) {
		return 0, false, errors.New("malformed header")
	}

	key = strings.ToLower(key)
	h.Set(key, value)

	return idx + 2, false, nil
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{v, value}, ", ")
	}
	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	v, ok := h[key]
	return v, ok
}

func validToken(b []byte) bool {
	for _, c := range b {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	if c >= 'A' && c <= 'Z' {
		return true
	}
	if c >= '0' && c <= '9' {
		return true
	}

	return slices.Contains(tokenChars, c)
}
