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

// Log is a simple wrapper around Logf, that automatically builds a format string, output is separated by spaces
func Log(a ...any) {
	format := strings.Repeat("%v ", len(a))
	fmt.Println(fmt.Sprintf(format, a...))
}

// Log is a simple wrapper around Logf, that automatically builds a format string, output is separated by new lines
func Logn(a ...any) {
	format := strings.Repeat("%v\n", len(a))
	fmt.Println(fmt.Sprintf(format, a...))
}

// LogMatrix is a quick means to output a 2d grid
func LogMatrix(a [][]string) {
	out := make([]string, len(a))
	for i := range a {
		out[i] = strings.Join(a[i], "")
	}
	fmt.Println(fmt.Sprintf("%s", strings.Join(out, "\n")))
}
