package ui

import (
	"fmt"

	"github.com/mattn/go-colorable"
)

const escape = "\x1b"
const (
	black int = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

var (
	stdout = colorable.NewColorableStdout()
	stderr = colorable.NewColorableStderr()
)

// Printf formats according to a format specifier and writes to colorable stdout(green color)
// It returns the number of bytes written and any write error encountered.
func Printf(format string, args ...interface{}) {
	fmt.Fprintf(stdout, color(green, format), args...)
}

// Errorf formats according to a format specifier and writes to colorable stdout(red color)
// It returns the number of bytes written and any write error encountered.
func Errorf(format string, args ...interface{}) {
	fmt.Fprintf(stderr, color(red, format), args...)
}

func color(color int, format string) string {
	return fmt.Sprintf("%s[%dm%s%s[0m", escape, color, format, escape)
}
