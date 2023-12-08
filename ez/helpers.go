package ez

import (
	"fmt"
	"strconv"
	"strings"
)

// Atoi ignores the errors in strconv.Atoi and returns the response
func Atoi(in string) int {
	out, _ := strconv.Atoi(in)
	return out
}

// Logf is a simple wrapper around Println + Sprintf
func Logf(format string, a ...any) {
	fmt.Println(fmt.Sprintf(format, a...))
}

// Log is a simple wrapper around Logf, that automatically builds a format string
func Log(a ...any) {
	format := strings.Repeat("%v ", len(a))
	fmt.Println(fmt.Sprintf(format, a...))
}
