// +build linux

package logger

import (
	"os"
)

type brush func(string) string

func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []brush{
	newBrush("1;32"), // Debug
	newBrush("1;36"), // Info
	newBrush("1;33"), // Warn
	newBrush("1;31"), // Error
	newBrush("1;41"), // Fatal
}

func (l *LogInfo) Println() {
	msg := colors[l.Level](l.FormatString())
	os.Stdout.Write(append([]byte(msg), '\n'))
}
