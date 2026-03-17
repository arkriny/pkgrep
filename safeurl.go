package main

import (
	"fmt"
	"unicode"
)

// safeURLSegment checks whether a string can be safely placed in URL segment.
func safeURLSegment(s string) error {
	for _, r := range s {
		if !unicode.IsLetter(r) &&
			!unicode.IsDigit(r) &&
			r != '-' &&
			r != '.' &&
			r != '_' &&
			r != '~' {
			return fmt.Errorf("disallowed character in URL %q", r)
		}
	}
	return nil
}
