// +build windows

package logger

import (
	"syscall"
)

var (
	kernel32     *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
	proc         *syscall.LazyProc = kernel32.NewProc(`SetConsoleTextAttribute`)
	closeHandle  *syscall.LazyProc = kernel32.NewProc(`CloseHandle`)
	defaultColor int               = 15
)

var colors = []int{
	10, // Debug
	11, // Info
	14, // Warn
	12, // Error
	4,  // Fatal
}

func (l *LogInfo) Println() {
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(colors[l.Level]))
	print(l.FormatString() + "\r\n")
	closeHandle.Call(handle)

	handle2, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(defaultColor))
	closeHandle.Call(handle2)
}
