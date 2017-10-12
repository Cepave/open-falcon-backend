package utils

import (
	"fmt"
	"math"

	"github.com/juju/errors"
)

// Shorten string to "<heading chars> <more> <tailing chars>(total size)" if the length of input is
// greater than "maxSize".
//
// For example of 4 characters of maximum size(more is "..."):
//
// 	"abcd" -> "abcd"
// 	"abcde" -> "ab ...(4) de"
// 	"This is hello world!" -> "Th ...(4) d!"
//
// If "more" is empty, the result string would be "<heading chars> <tailing chars>"
func ShortenStringToSize(s string, more string, maxSize int) string {
	if maxSize <= 0 {
		panic(errors.Details(errors.Errorf("Maximum size is not viable")))
	}

	if len(s) <= maxSize {
		return s
	}

	srcRune := []rune(s)
	srcSize := len(srcRune)

	headingSize := int(math.Ceil(float64(maxSize) / 2))
	tailingSize := maxSize / 2
	if tailingSize == 0 {
		tailingSize = 1
	}

	headingText := srcRune[:headingSize]
	tailingText := srcRune[srcSize-tailingSize:]

	if len(more) > 0 {
		return fmt.Sprintf("%s %s %s", string(headingText), more, string(tailingText))
	}

	return fmt.Sprintf("%s %s", string(headingText), string(tailingText))
}
