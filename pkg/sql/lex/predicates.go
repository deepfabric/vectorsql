// Copyright 2017 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package lex

import (
	"unicode"
)

// isASCII returns true if all the characters in s are ASCII.
func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// IsDigit returns true if the character is between 0 and 9.
func IsDigit(ch int) bool {
	return ch >= '0' && ch <= '9'
}

// IsHexDigit returns true if the character is a valid hexadecimal digit.
func IsHexDigit(ch int) bool {
	return (ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}

// IsIdentStart returns true if the character is valid at the start of an identifier.
func IsIdentStart(ch int) bool {
	return (ch >= 'A' && ch <= 'Z') ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= 128 && ch <= 255) ||
		(ch == '_')
}

// IsIdentMiddle returns true if the character is valid inside an identifier.
func IsIdentMiddle(ch int) bool {
	return IsIdentStart(ch) || IsDigit(ch) || ch == '$'
}
