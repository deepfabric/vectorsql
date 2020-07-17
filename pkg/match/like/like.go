package like

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

func Compile(expr string) (*like, error) {
	switch {
	case len(expr) == 0:
		return &like{fn: func(s string) bool { return s == "" }}, nil
	case len(expr) == 1 && expr[0] == '%':
		return &like{fn: func(s string) bool { return true }}, nil
	case len(expr) == 1 && expr[0] == '_':
		return &like{fn: func(s string) bool { return len(s) != 0 }}, nil
	}
	reg, err := regexp.Compile(convert(expr))
	if err != nil {
		return nil, err
	}
	return &like{reg: reg}, nil
}

func (l *like) MatchString(s string) bool {
	if l.fn != nil {
		return l.fn(s)
	}
	return l.reg.MatchString(s)
}

func convert(expr string) string {
	return fmt.Sprintf("^(?s:%s)$", replace(expr))
}

func replace(s string) string {
	var oc rune

	r := make([]byte, len(s)+strings.Count(s, `%`))
	w := 0
	start := 0
	for len(s) > start {
		c, wid := utf8.DecodeRuneInString(s[start:])
		if oc == '\\' {
			w += copy(r[w:], s[start:start+wid])
			start += wid
			oc = 0
			continue
		}
		switch c {
		case '_':
			w += copy(r[w:], []byte{'*'})
		case '%':
			w += copy(r[w:], []byte{'.', '*'})
		case '\\':
		default:
			w += copy(r[w:], s[start:start+wid])
		}
		start += wid
		oc = c
	}
	return string(r)
}
